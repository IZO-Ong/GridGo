// Package handlers contains the HTTP logic for the GridGo API.
// This file manages Maze lifecycle of generation, persistence, solving and rendering.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/maze"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/IZO-Ong/gridgo/internal/utils"
	"gorm.io/gorm"
)

// HandleGenerateMaze creates a new maze grid based on the requested algorithm.
func HandleGenerateMaze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit form size to 10MB
	r.ParseMultipartForm(10 << 20)
	rows, _ := strconv.Atoi(r.FormValue("rows"))
	cols, _ := strconv.Atoi(r.FormValue("cols"))
	genType := r.FormValue("type")

	myMaze := maze.NewMaze(rows, cols)
	var originalWeights map[string]int

	switch genType {
	case "image":
		// Extract weights from image luminance/contrast
		file, _, _ := r.FormFile("image")
		if file != nil {
			defer file.Close()
			weights, _ := maze.GetEdgeWeights(file, rows, cols)
			originalWeights = weights
			myMaze.GenerateImageMaze(weights)
		}
	case "kruskal":
		myMaze.GenerateKruskal()
	case "recursive":
		myMaze.GenerateRecursive(0, 0)
	}

	// Finalize maze state and calculate difficulty metrics
	myMaze.SyncGridToWeights(originalWeights)
	myMaze.SetRandomStartEnd()
	
	weightsBytes, _ := json.Marshal(myMaze.Weights)
	stats := myMaze.CalculateStats()
	mazeID := utils.GenerateMazeID()

	userID := middleware.GetUserID(r)
	
	// Prepare DB model
	dbMaze := models.Maze{
		ID:          mazeID, 
		WeightsJSON: string(weightsBytes),
		Rows:        rows, 
		Cols:        cols,
		StartRow:    myMaze.Start[0], 
		StartCol:    myMaze.Start[1],
		EndRow:      myMaze.End[0], 
		EndCol:      myMaze.End[1],
		Complexity:  stats.Complexity, 
	}

	// Associate with user if authenticated
	if userID != "" { dbMaze.CreatorID = &userID }

	if err := db.DB.Create(&dbMaze).Error; err == nil && userID != "" {
        db.RDB.Del(db.Ctx, "user:mazes:"+userID)

        var username string
        if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Pluck("username", &username).Error; err == nil {
            db.RDB.Del(db.Ctx, "profile:public:"+username)
        }
    }

	myMaze.ID = mazeID
	json.NewEncoder(w).Encode(myMaze)
}

// HandleRenderMaze takes a maze structure and returns a binary PNG image.
func HandleRenderMaze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions { return }

	var m maze.Maze
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	m.RenderToWriter(w, 10)
}

// HandleSolveMaze accepts a maze and an algorithm name, then returns 
// the computed path and the list of cells visited during the search.
func HandleSolveMaze(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions { return }
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Maze      maze.Maze `json:"maze"`
		Algorithm string    `json:"algorithm"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var visited [][2]int
	var path [][2]int

	// Route to algorithm in maze package
	switch payload.Algorithm {
	case "astar":
		visited, path = payload.Maze.SolveAStar()
	case "bfs":
		visited, path = payload.Maze.SolveBFS()
	case "greedy":
		visited, path = payload.Maze.SolveGreedy()
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"visited": visited,
		"path":    path,
	})
}

// HandleGetMaze retrieves a maze by ID and reconstructs its full grid 
// object from the stored JSON weights.
func HandleGetMaze(w http.ResponseWriter, r *http.Request) {
    mazeID := r.URL.Query().Get("id")
    cacheKey := "maze:reconstructed:" + mazeID

    // Cache a map[string]interface{} to match response format
    data, err := db.GetOrSet(db.Ctx, cacheKey, 24*time.Hour, func() (*map[string]interface{}, error) {
        var m models.Maze
        if err := db.DB.First(&m, "id = ?", mazeID).Error; err != nil {
            return nil, err
        }

        var savedWeights map[string]int
        json.Unmarshal([]byte(m.WeightsJSON), &savedWeights)

        // Reinflate the grid structure
        reconstructed := maze.NewMaze(m.Rows, m.Cols)
        reconstructed.GenerateImageMaze(savedWeights)

        // Mapping of weights back into cell structs
        for r := range m.Rows {
            for c := range m.Cols {
                if v, ok := savedWeights[fmt.Sprintf("%d-%d-top", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[0] = v }
                if v, ok := savedWeights[fmt.Sprintf("%d-%d-right", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[1] = v }
                if v, ok := savedWeights[fmt.Sprintf("%d-%d-bottom", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[2] = v }
                if v, ok := savedWeights[fmt.Sprintf("%d-%d-left", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[3] = v }
            }
        }

        reconstructed.SetManualStartEnd(m.StartRow, m.StartCol, m.EndRow, m.EndCol)
        
        // Return formatted response map
        result := map[string]interface{}{
            "id": m.ID, "rows": m.Rows, "cols": m.Cols, "grid": reconstructed.Grid,
            "start": [2]int{m.StartRow, m.StartCol}, "end": [2]int{m.EndRow, m.EndCol},
        }
        return &result, nil
    })

    if err != nil {
        http.Error(w, "Maze not found", 404)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// HandleGetMyMazes returns a list of all mazes created by the current user.
func HandleGetMyMazes(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    if userID == "" {
        http.Error(w, "AUTH_REQUIRED", 401)
        return
    }

    cacheKey := "user:mazes:" + userID

    mazes, err := db.GetOrSet(db.Ctx, cacheKey, 10*time.Minute, func() (*[]models.Maze, error) {
        var m []models.Maze
        err := db.DB.Where("creator_id = ?", userID).Order("created_at desc").Find(&m).Error
        return &m, err
    })

    if err != nil {
        http.Error(w, "DB_ERROR", 500)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mazes)
}

// HandleDeleteMaze removes a maze, with authorization check for ownership.
func HandleDeleteMaze(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete { return }
    
    mazeID := r.URL.Query().Get("id")
    userID := middleware.GetUserID(r)

    var m models.Maze
    if err := db.DB.Preload("Creator").Where("id = ? AND creator_id = ?", mazeID, userID).First(&m).Error; err != nil {
        http.Error(w, "UNAUTHORIZED_OR_NOT_FOUND", 403)
        return
    }

    if err := db.DB.Delete(&m).Error; err != nil {
        http.Error(w, "DB_ERROR", 500)
        return
    }

    // 3. INVALIDATION BATCH
    db.RDB.Del(db.Ctx, "user:mazes:"+userID)
    db.RDB.Del(db.Ctx, "maze:reconstructed:"+mazeID)
    
    // Clear the public profile so the stats (maze count) refresh
    if m.Creator != nil {
        db.RDB.Del(db.Ctx, "profile:public:"+m.Creator.Username)
    }

    w.WriteHeader(http.StatusOK)
}

// HandleUpdateThumbnail saves a base64 or URL thumbnail for the maze gallery.
func HandleUpdateThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "METHOD_NOT_ALLOWED", 405); return
	}

	userID := middleware.GetUserID(r)
	var payload struct { ID string; Thumbnail string }
	json.NewDecoder(r.Body).Decode(&payload)

	// Ownership & Immutability Check
	var result *gorm.DB
	if userID != "" {
		result = db.DB.Model(&models.Maze{}).
			Where("id = ? AND creator_id = ? AND (thumbnail IS NULL OR thumbnail = '')", payload.ID, userID).
			Update("thumbnail", payload.Thumbnail)
	} else {
		result = db.DB.Model(&models.Maze{}).
			Where("id = ? AND creator_id IS NULL AND (thumbnail IS NULL OR thumbnail = '')", payload.ID).
			Update("thumbnail", payload.Thumbnail)
	}

	if result.RowsAffected > 0 {
		if userID != "" {
			// Clear Private Gallery
			db.RDB.Del(db.Ctx, "user:mazes:"+userID)
			var username string
			db.DB.Model(&models.User{}).Where("id = ?", userID).Pluck("username", &username)
			db.RDB.Del(db.Ctx, "profile:public:"+username)
		}
		db.RDB.Del(db.Ctx, "maze:reconstructed:"+payload.ID)
	}
	w.WriteHeader(http.StatusOK)
}
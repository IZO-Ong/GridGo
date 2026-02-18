package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/maze"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
)

func HandleGenerateMaze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)
	rows, _ := strconv.Atoi(r.FormValue("rows"))
	cols, _ := strconv.Atoi(r.FormValue("cols"))
	genType := r.FormValue("type")

	myMaze := maze.NewMaze(rows, cols)
	var originalWeights map[string]int

	switch genType {
	case "image":
		file, _, _ := r.FormFile("image")
		defer file.Close()
		weights, _ := maze.GetEdgeWeights(file, rows, cols)
		originalWeights = weights
		myMaze.GenerateImageMaze(weights)
	case "kruskal":
		myMaze.GenerateKruskal()
	case "recursive":
		myMaze.GenerateRecursive(0, 0)
	}

	myMaze.SyncGridToWeights(originalWeights)
	myMaze.SetRandomStartEnd()
	
	weightsBytes, _ := json.Marshal(myMaze.Weights)
	stats := myMaze.CalculateStats()
	mazeID := fmt.Sprintf("M-%d-X", rand.Intn(9000)+1000)

	userID := middleware.GetUserID(r)
	
	dbMaze := models.Maze{
		ID: mazeID, WeightsJSON: string(weightsBytes),
		Rows: rows, Cols: cols,
		StartRow: myMaze.Start[0], StartCol: myMaze.Start[1],
		EndRow: myMaze.End[0], EndCol: myMaze.End[1],
		Complexity: stats.Complexity, DeadEnds: stats.DeadEnds,
	}

	if userID != "" { dbMaze.CreatorID = &userID }

	db.DB.Create(&dbMaze)
	myMaze.ID = mazeID
	json.NewEncoder(w).Encode(myMaze)
}

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

func HandleGetMaze(w http.ResponseWriter, r *http.Request) {
	mazeID := r.URL.Query().Get("id")
	var m models.Maze
	if err := db.DB.First(&m, "id = ?", mazeID).Error; err != nil {
		http.Error(w, "Maze not found", 404)
		return
	}

	var savedWeights map[string]int
	json.Unmarshal([]byte(m.WeightsJSON), &savedWeights)

	reconstructed := maze.NewMaze(m.Rows, m.Cols)
	reconstructed.GenerateImageMaze(savedWeights)

	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if v, ok := savedWeights[fmt.Sprintf("%d-%d-top", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[0] = v }
			if v, ok := savedWeights[fmt.Sprintf("%d-%d-right", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[1] = v }
			if v, ok := savedWeights[fmt.Sprintf("%d-%d-bottom", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[2] = v }
			if v, ok := savedWeights[fmt.Sprintf("%d-%d-left", r, c)]; ok { reconstructed.Grid[r][c].WallWeights[3] = v }
		}
	}

	reconstructed.SetManualStartEnd(m.StartRow, m.StartCol, m.EndRow, m.EndCol)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": m.ID, "rows": m.Rows, "cols": m.Cols, "grid": reconstructed.Grid,
		"start": [2]int{m.StartRow, m.StartCol}, "end": [2]int{m.EndRow, m.EndCol},
	})
}

func HandleGetMyMazes(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    if userID == "" {
        http.Error(w, "AUTH_REQUIRED", 401)
        return
    }

    var mazes []models.Maze
    db.DB.Where("creator_id = ?", userID).Order("created_at desc").Find(&mazes)
    
    json.NewEncoder(w).Encode(mazes)
}

func HandleDeleteMaze(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete { return }
    
    mazeID := r.URL.Query().Get("id")
    userID := middleware.GetUserID(r)

    result := db.DB.Where("id = ? AND creator_id = ?", mazeID, userID).Delete(&models.Maze{})
    
    if result.RowsAffected == 0 {
        http.Error(w, "UNAUTHORIZED_OR_NOT_FOUND", 403)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func HandleUpdateThumbnail(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", 405)
        return
    }

    var payload struct {
        ID        string `json:"id"`
        Thumbnail string `json:"thumbnail"`
    }
    json.NewDecoder(r.Body).Decode(&payload)

    result := db.DB.Model(&models.Maze{}).Where("id = ?", payload.ID).Update("thumbnail", payload.Thumbnail)
    
    if result.Error != nil {
        http.Error(w, "DB_UPDATE_FAILED", 500)
        return
    }
    w.WriteHeader(http.StatusOK)
}

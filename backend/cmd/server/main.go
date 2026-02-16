package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/maze"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
      log.Println("Warning: .env file not found, using system environment variables")
    }

	db.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/maze/get", handleGetMaze)
	mux.HandleFunc("/api/maze/generate", handleGenerateMaze)
	mux.HandleFunc("/api/maze/render", handleRenderMaze)
	mux.HandleFunc("/api/maze/solve", handleSolveMaze)

	println("GridGo API running on port 8080")
	http.ListenAndServe(":8080", enableCORS(mux))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleGenerateMaze(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // parse metadata
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    rows, _ := strconv.Atoi(r.FormValue("rows"))
    cols, _ := strconv.Atoi(r.FormValue("cols"))

    // check boundary
    if rows < 2 || rows > 300 || cols < 2 || cols > 300 {
        http.Error(w, "OUT_OF_BOUNDS: Dimensions must be between 2 and 300", http.StatusBadRequest)
        return
    }

    genType := r.FormValue("type")
    myMaze := maze.NewMaze(rows, cols)

	switch genType {
	case "image":
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image required for image-type maze", http.StatusBadRequest)
			return
		}
		defer file.Close()

		weights, err := maze.GetEdgeWeights(file, rows, cols)
		if err != nil {
			http.Error(w, "Vision processing failed", http.StatusInternalServerError)
			return
		}
		myMaze.GenerateImageMaze(weights)

	case "kruskal":
		myMaze.GenerateKruskal()

	case "recursive":
		myMaze.GenerateRecursive(0, 0)

	default:
		http.Error(w, "Invalid generation type", http.StatusBadRequest)
		return
	}

    myMaze.SetRandomStartEnd()

    gridBytes, _ := json.Marshal(myMaze.Grid)
	stats := myMaze.CalculateStats()

	dbMaze := models.Maze{
		ID:         "M-" + strconv.Itoa(rand.Intn(9000)+1000) + "-X",
		GridJSON:   string(gridBytes),
		Rows:       rows,
		Cols:       cols,
		StartRow:   myMaze.Start[0],
		StartCol:   myMaze.Start[1],
		EndRow:     myMaze.End[0],
		EndCol:     myMaze.End[1],
		DeadEnds:   stats.DeadEnds,
		Complexity: stats.Complexity,
	}

    result := db.DB.Create(&dbMaze)
    if result.Error != nil {
        http.Error(w, "Failed to save to database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(myMaze)
}

func handleRenderMaze(w http.ResponseWriter, r *http.Request) {
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

func handleSolveMaze(w http.ResponseWriter, r *http.Request) {
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

func handleGetMaze(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    mazeID := r.URL.Query().Get("id")
    if mazeID == "" {
        http.Error(w, "Maze ID required", http.StatusBadRequest)
        return
    }

    var m models.Maze
    result := db.DB.First(&m, "id = ?", mazeID)
    if result.Error != nil {
        http.Error(w, "Maze not found", http.StatusNotFound)
        return
    }

    var grid [][]maze.Cell
    json.Unmarshal([]byte(m.GridJSON), &grid)

    response := map[string]interface{}{
        "id":    m.ID,
        "rows":  m.Rows,
        "cols":  m.Cols,
        "grid":  grid,
        "start": [2]int{m.StartRow, m.StartCol},
        "end":   [2]int{m.EndRow, m.EndCol},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
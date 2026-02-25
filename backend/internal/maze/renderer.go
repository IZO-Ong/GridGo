// Package maze provides the core logic for maze manipulation.
// This file specifically handles the translation of grid data into PNG images.
package maze

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"runtime"
	"sync"
)

// prepareCanvas initializes the image buffer and background.
// It uses a +1 pixel offset to ensure that the final right and bottom borders
// are correctly rendered without being clipped.
func (m *Maze) prepareCanvas(cellSize int) *image.RGBA {
	imgWidth := m.Cols*cellSize + 1
	imgHeight := m.Rows*cellSize + 1
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill with a solid white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	return img
}

// RenderToWriter orchestrates the drawing process and streams the resulting 
// PNG directly to an io.Writer (usually an HTTP response body).
func (m *Maze) RenderToWriter(w io.Writer, cellSize int) error {
	img := m.prepareCanvas(cellSize)
	m.drawMaze(img, cellSize)
	
	return png.Encode(w, img)
}

// drawMaze implements a parallelized rendering engine. 
// It splits the grid into horizontal chunks and assigns each to a worker 
// goroutine to maximize performance on multi-core systems.
func (m *Maze) drawMaze(img *image.RGBA, cellSize int) {
	var wg sync.WaitGroup

	numCPU := runtime.NumCPU()
	rowsPerWorker := m.Rows / numCPU

	// Fallback for small mazes where parallelization is overkill
	if rowsPerWorker == 0 {
		numCPU = 1
		rowsPerWorker = m.Rows
	}

	for w := 0; w < numCPU; w++ {
		wg.Add(1)
		startY := w * rowsPerWorker
		endY := startY + rowsPerWorker

		if w == numCPU-1 {
			endY = m.Rows // Ensure all remainder rows are processed
		}

		// Perform row-based rendering in a separate Goroutine
		go func(rMin, rMax int) {
			defer wg.Done()
			for r := rMin; r < rMax; r++ {
				for c := 0; c < m.Cols; c++ {
					x := c * cellSize
					y := r * cellSize
					cell := m.Grid[r][c]

					// Highlight special points: Start (Green) and End (Red)
					if r == m.Start[0] && c == m.Start[1] {
						m.fillCell(img, x, y, cellSize, color.RGBA{144, 238, 144, 255})
					} else if r == m.End[0] && c == m.End[1] {
						m.fillCell(img, x, y, cellSize, color.RGBA{255, 99, 71, 255})
					}

					// Render each wall if it exists in the logical grid
					if cell.Walls[0] { m.paintWall(img, x, y, cellSize, 0, cell.WallWeights[0]) } // TOP
					if cell.Walls[1] { m.paintWall(img, x, y, cellSize, 1, cell.WallWeights[1]) } // RIGHT
					if cell.Walls[2] { m.paintWall(img, x, y, cellSize, 2, cell.WallWeights[2]) } // BOTTOM
					if cell.Walls[3] { m.paintWall(img, x, y, cellSize, 3, cell.WallWeights[3]) } // LEFT
				}
			}
		}(startY, endY)
	}
	wg.Wait()
}

// fillCell paints the interior pixels of a specific grid coordinate.
func (m *Maze) fillCell(img *image.RGBA, x, y, size int, col color.RGBA) {
	for i := 1; i < size; i++ {
		for j := 1; j < size; j++ {
			img.Set(x+i, y+j, col)
		}
	}
}

// paintWall handles the pixel-level drawing of a single boundary line.
// It uses the wall's weight to determine the line's color/shading.
func (m *Maze) paintWall(img *image.RGBA, x, y, cellSize, direction, weight int) {
	col := m.getWallColor(weight)

	switch direction {
	case 0: // TOP
		for i := 0; i <= cellSize; i++ { img.Set(x+i, y, col) }
	case 1: // RIGHT
		for i := 0; i <= cellSize; i++ { img.Set(x+cellSize, y+i, col) }
	case 2: // BOTTOM
		for i := 0; i <= cellSize; i++ { img.Set(x+i, y+cellSize, col) }
	case 3: // LEFT
		for i := 0; i <= cellSize; i++ { img.Set(x, y+i, col) }
	}
}

// getWallColor translates a mathematical weight into a grayscale color.
// Pure walls are black (0,0,0), while weighted walls from images vary in brightness.
func (m *Maze) getWallColor(weight int) color.RGBA {
	if weight >= 255 {
		return color.RGBA{0, 0, 0, 255}
	}
	
	// Map 0-255 weight scale to a 220-0 brightness scale for shading
	brightness := uint8(220 - (float64(weight) * (220.0 / 255.0)))
	return color.RGBA{brightness, brightness, brightness, 255}
}
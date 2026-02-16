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
func (m *Maze) prepareCanvas(cellSize int) *image.RGBA {
	// add 1 pixel to the total width/height to ensure
	// closing edges of the rightmost and bottommost cells are rendered
	imgWidth := m.Cols*cellSize + 1
	imgHeight := m.Rows*cellSize + 1
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// initialise white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	return img
}

func (m *Maze) RenderToWriter(w io.Writer, cellSize int) error {
    img := m.prepareCanvas(cellSize)
    m.drawMaze(img, cellSize)
    
    // Uses your existing logic but encodes to the HTTP stream
    return png.Encode(w, img)
}

// drawMaze iterates through the grid and paints each active wall.
func (m *Maze) drawMaze(img *image.RGBA, cellSize int) {
	var wg sync.WaitGroup

	numCPU := runtime.NumCPU()
	rowsPerWorker := m.Rows / numCPU

	// avoid spawning more workers if maze is small
	if rowsPerWorker == 0 {
		numCPU = 1
		rowsPerWorker = m.Rows
	}

	for w := range numCPU {
		wg.Add(1)
		startY := w * rowsPerWorker
		endY := startY + rowsPerWorker

		if w == numCPU-1 {
			endY = m.Rows // cover any remainder rows
		}

		// we do a goRoutine
		go func(rMin, rMax int) {
			defer wg.Done()
			for r := rMin; r < rMax; r++ {
				for c := 0; c < m.Cols; c++ {
					x := c * cellSize
					y := r * cellSize
					cell := m.Grid[r][c]

					if r == m.Start[0] && c == m.Start[1] {
						m.fillCell(img, x, y, cellSize, color.RGBA{144, 238, 144, 255}) // Light Green
					} else if r == m.End[0] && c == m.End[1] {
						m.fillCell(img, x, y, cellSize, color.RGBA{255, 99, 71, 255}) // Red
					}

					// TOP WALL
					if cell.Walls[0] {
						m.paintWall(img, x, y, cellSize, 0, cell.WallWeights[0])
					}
					// RIGHT WALL
					if cell.Walls[1] {
						m.paintWall(img, x, y, cellSize, 1, cell.WallWeights[1])
					}
					// BOTTOM WALL
					if cell.Walls[2] {
						m.paintWall(img, x, y, cellSize, 2, cell.WallWeights[2])
					}
					// LEFT WALL
					if cell.Walls[3] {
						m.paintWall(img, x, y, cellSize, 3, cell.WallWeights[3])
					}
				}
			}
		}(startY, endY)
	}
	wg.Wait()
}

// fillCell paints the interior of a maze square.
func (m *Maze) fillCell(img *image.RGBA, x, y, size int, col color.RGBA) {
	for i := 1; i < size; i++ {
		for j := 1; j < size; j++ {
			img.Set(x+i, y+j, col)
		}
	}
}

// paintWall handles the pixel-level drawing of a single boundary.
func (m *Maze) paintWall(img *image.RGBA, x, y, cellSize, direction, weight int) {
	col := m.getWallColor(weight)

	switch direction {
	case 0: // TOP
		for i := 0; i <= cellSize; i++ {
			img.Set(x+i, y, col)
		}
	case 1: // RIGHT
		for i := 0; i <= cellSize; i++ {
			img.Set(x+cellSize, y+i, col)
		}
	case 2: // BOTTOM
		for i := 0; i <= cellSize; i++ {
			img.Set(x+i, y+cellSize, col)
		}
	case 3: // LEFT
		for i := 0; i <= cellSize; i++ {
			img.Set(x, y+i, col)
		}
	}
}

// shading translates mathematical weights into RGB values.
func (m *Maze) getWallColor(weight int) color.RGBA {
    if weight >= 255 {
        return color.RGBA{0, 0, 0, 255}
    }
    
    // Linear scale: 220 down to 0
    brightness := uint8(220 - (float64(weight) * (220.0 / 255.0)))
    return color.RGBA{brightness, brightness, brightness, 255}
}
package maze

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// SaveAsImage renders the maze's structure into a PNG file
func (m *Maze) SaveAsImage(filename string, cellSize int) error {
	// add 1 pixel to the total width/height to ensure the
	// closing edges of the rightmost and bottommost cells are rendered
	imgWidth := m.Cols*cellSize + 1
	imgHeight := m.Rows*cellSize + 1

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// initialise white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// shading translates mathematical weights into RGB values
	getWallColor := func(weight int) color.RGBA {
		// high weights are rendered black
		if weight >= 1000 {
			return color.RGBA{0, 0, 0, 255}
		}

		// low weights (randomized filler) are rendered in light gray
		// modulo variance provides a slight texture to empty spaces
		intensity := uint8(230 - (weight % 30))
		return color.RGBA{intensity, intensity, intensity, 255}
	}

	// iterates through the grid and paints each active wall
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			x := c * cellSize
			y := r * cellSize

			// TOP WALL
			if m.Grid[r][c].Walls[0] {
				col := getWallColor(m.Grid[r][c].WallWeights[0])
				for i := 0; i <= cellSize; i++ {
					img.Set(x+i, y, col)
				}
			}
			// RIGHT WALL
			if m.Grid[r][c].Walls[1] {
				col := getWallColor(m.Grid[r][c].WallWeights[1])
				for i := 0; i <= cellSize; i++ {
					img.Set(x+cellSize, y+i, col)
				}
			}
			// BOTTOM WALL
			if m.Grid[r][c].Walls[2] {
				col := getWallColor(m.Grid[r][c].WallWeights[2])
				for i := 0; i <= cellSize; i++ {
					img.Set(x+i, y+cellSize, col)
				}
			}
			// LEFT WALL
			if m.Grid[r][c].Walls[3] {
				col := getWallColor(m.Grid[r][c].WallWeights[3])
				for i := 0; i <= cellSize; i++ {
					img.Set(x, y+i, col)
				}
			}
		}
	}

	// exports the in-memory buffer to a PNG
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

package maze

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func (m *Maze) SaveAsImage(filename string, cellSize int) error {
	imgWidth := m.Cols*cellSize + 1
	imgHeight := m.Rows*cellSize + 1

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	wallColor := color.Black

	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			x := c * cellSize
			y := r * cellSize

			// paint the walls (top, right, bottom, left))
			if m.Grid[r][c].Walls[0] {
				for i := 0; i <= cellSize; i++ {
					img.Set(x+i, y, wallColor)
				}
			}
			if m.Grid[r][c].Walls[1] {
				for i := 0; i <= cellSize; i++ {
					img.Set(x+cellSize, y+i, wallColor)
				}
			}
			if m.Grid[r][c].Walls[2] {
				for i := 0; i <= cellSize; i++ {
					img.Set(x+i, y+cellSize, wallColor)
				}
			}
			if m.Grid[r][c].Walls[3] {
				for i := 0; i <= cellSize; i++ {
					img.Set(x, y+i, wallColor)
				}
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

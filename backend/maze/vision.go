package maze

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

func GetEdgeWeights(path string, rows, cols int) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// sobel kernels
	gx := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	gy := [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	weights := make(map[string]int)

	for r := range rows {
		for c := range cols {
			imgX := c * width / cols
			imgY := r * height / rows

			// sobel edge magnitude
			mag := calculateSobel(img, imgX, imgY, gx, gy)

			if mag > 100 {
				weights[fmt.Sprintf("%d-%d-top", r, c)] = 100
				weights[fmt.Sprintf("%d-%d-left", r, c)] = 100
			}
		}
	}

	return weights, nil
}

func calculateSobel(img image.Image, x, y int, gx, gy [][]int) int {
	var sumX, sumY int
	bounds := img.Bounds()

	// sobel formula
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			px := x + j
			py := y + j

			if px < 0 || px >= bounds.Max.X || py < 0 || py >= bounds.Max.Y {
				continue
			}

			// luminance
			r, g, b, _ := img.At(px, py).RGBA()
			lum := int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

			sumX += lum * gx[i+1][j+1]
			sumY += lum * gy[i+1][j+1]
		}
	}

	return int(math.Sqrt(float64(sumX*sumX+sumY*sumY))) / 256
}

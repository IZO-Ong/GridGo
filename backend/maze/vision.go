package maze

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// GetEdgeWeights transforms a source image into a map of wall priorities.
// uses a Canny filter pipeline to ensure that the  outline of the
// image is preserved by assigning high weights to structural edges, which
// Kruskal's algorithm will then prioritise keeping as walls.
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

	// stage 1
	// we convert to grayscale to focus purely on luminance edges,
	// ignoring color data
	gray := image.NewGray(bounds)
	for y := range height {
		for x := 0; x < width; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	// stage 2
	// use sobel kernels to identify where intensity changes rapidly
	mags := make([][]float64, height)
	angles := make([][]float64, height)
	for i := range mags {
		mags[i] = make([]float64, width)
		angles[i] = make([]float64, width)
	}

	gx := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	gy := [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var sumX, sumY float64
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					px, py := x+j, y+i
					lum := float64(gray.GrayAt(px, py).Y)
					sumX += lum * float64(gx[i+1][j+1])
					sumY += lum * float64(gy[i+1][j+1])
				}
			}
			mags[y][x] = math.Sqrt(sumX*sumX + sumY*sumY)
			// normalise angle to 0-180
			angles[y][x] = math.Mod(math.Atan2(sumY, sumX)*180/math.Pi+180, 180)
		}
	}

	// stage 3:
	// non-maximum suppression to thin out "blurry" edges
	nmsMags := make([][]float64, height)
	for i := range nmsMags {
		nmsMags[i] = make([]float64, width)
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			angle := angles[y][x]
			mag := mags[y][x]
			var q, r float64

			// compare pixel against its neighbors along the gradient normal
			switch {
			case (angle >= 0 && angle < 22.5) || (angle >= 157.5 && angle <= 180):
				q, r = mags[y][x+1], mags[y][x-1]
			case angle >= 22.5 && angle < 67.5:
				q, r = mags[y+1][x-1], mags[y-1][x+1]
			case angle >= 67.5 && angle < 112.5:
				q, r = mags[y+1][x], mags[y-1][x]
			case angle >= 112.5 && angle < 157.5:
				q, r = mags[y-1][x-1], mags[y+1][x+1]
			}

			if mag >= q && mag >= r {
				nmsMags[y][x] = mag
			} else {
				nmsMags[y][x] = 0
			}
		}
	}

	// stage 4:
	// translate pixel magnitudes into Kruskal's weights.
	// and thresholds help maintain connectivity in the silhouette
	weights := make(map[string]int)
	highThresh := 100.0
	lowThresh := 40.0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			imgX := c * width / cols
			imgY := r * height / rows

			mag := nmsMags[imgY][imgX]

			if mag >= highThresh {
				// high weight locks priority walls into spanning tree.
				weights[fmt.Sprintf("%d-%d-top", r, c)] = 5000
				weights[fmt.Sprintf("%d-%d-left", r, c)] = 5000
			} else if mag >= lowThresh {
				// moderate weight for supporting structural details.
				weights[fmt.Sprintf("%d-%d-top", r, c)] = 1500
			}
		}
	}

	return weights, nil
}

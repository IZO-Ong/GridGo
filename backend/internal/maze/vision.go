package maze

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"runtime"
	"sync"
)

// GetEdgeWeights transforms a source image into a map of wall priorities.
// uses a Canny filter pipeline to ensure that the  outline of the
// image is preserved by assigning high weights to structural edges, which
// Kruskal's algorithm will then prioritise keeping as walls.
func GetEdgeWeights(r io.Reader, rows, cols int) (map[string]int, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	gray := convertToGrayscale(img, bounds)
	mags, angles := computeGradients(gray, width, height)
	nmsMags := applyNMS(mags, angles, width, height)
	weights := mapToWeights(nmsMags, angles, rows, cols, width, height)

	return weights, nil
}

// Convert to grayscale to focus purely on luminance edges,
// ignoring color data
func convertToGrayscale(img image.Image, bounds image.Rectangle) *image.Gray {
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	return gray
}

// Use sobel kernels to identify where intensity changes rapidly
func computeGradients(gray *image.Gray, width, height int) ([][]float64, [][]float64) {
	mags := make([][]float64, height)
	angles := make([][]float64, height)
	for i := range mags {
		mags[i] = make([]float64, width)
		angles[i] = make([]float64, width)
	}

	gx := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	gy := [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	// splits image processing using go routines
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	rowsPerWorker := (height - 2) / numCPU

	for w := 0; w < numCPU; w++ {
		wg.Add(1)
		startY := 1 + (w * rowsPerWorker)
		endY := startY + rowsPerWorker

		if w == numCPU-1 {
			endY = height - 1
		}

		go func(yMin, yMax int) {
			defer wg.Done()
			for y := yMin; y < yMax; y++ {
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
		}(startY, endY)
	}
	wg.Wait()
	return mags, angles
}

// Non-maximum suppression to thin out "blurry" edges
func applyNMS(mags, angles [][]float64, width, height int) [][]float64 {
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
	return nmsMags
}

// Translate pixel magnitudes into Kruskal's weights.
// and thresholds help maintain connectivity in the silhouette
func mapToWeights(nmsMags [][]float64, angles [][]float64, rows, cols, width, height int) map[string]int {
	weights := make(map[string]int)

	highThresh := 80.0
	lowThresh := 30.0

	for r := range rows {
		for c := range cols {
			startY, endY := r*height/rows, (r+1)*height/rows
			startX, endX := c*width/cols, (c+1)*width/cols

			maxMag := 0.0
			bestAngle := 0.0

			for y := startY; y < endY && y < height; y++ {
				for x := startX; x < endX && x < width; x++ {
					if nmsMags[y][x] > maxMag {
						maxMag = nmsMags[y][x]
						bestAngle = angles[y][x]
					}
				}
			}

			if maxMag < lowThresh {
				continue
			}

			var weight int
			if maxMag >= highThresh {
				weight = 255
			} else {
				weight = 120 + int((maxMag-lowThresh)/(highThresh-lowThresh)*135)
			}

			if (bestAngle >= 0 && bestAngle < 45) || (bestAngle >= 135 && bestAngle <= 180) {
				weights[fmt.Sprintf("%d-%d-left", r, c)] = weight
				weights[fmt.Sprintf("%d-%d-top", r, c)] = weight / 2
			} else {
				weights[fmt.Sprintf("%d-%d-top", r, c)] = weight
				weights[fmt.Sprintf("%d-%d-left", r, c)] = weight / 2
			}
		}
	}
	return weights
}
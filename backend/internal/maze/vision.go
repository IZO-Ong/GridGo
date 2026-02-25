// Package maze provides sophisticated grid manipulation.
// This file handles the Computer Vision pipeline for image-to-maze generation.
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
// It uses a Canny-inspired pipeline to identify high-contrast outlines.
func GetEdgeWeights(r io.Reader, rows, cols int) (map[string]int, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 1. Pre-processing: Focus on luminance.
	gray := convertToGrayscale(img, bounds)

	// 2. Gradient Calculation: Find where intensity changes rapidly.
	mags, angles := computeGradients(gray, width, height)

	// 3. Thinning: Ensure edges are exactly 1 pixel wide for cleaner mazes.
	nmsMags := applyNMS(mags, angles, width, height)

	// 4. Mapping: Translate mathematical gradients into Maze Weights.
	weights := mapToWeights(nmsMags, angles, rows, cols, width, height)

	return weights, nil
}

// convertToGrayscale strips color data to simplify edge detection.
func convertToGrayscale(img image.Image, bounds image.Rectangle) *image.Gray {
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

// computeGradients uses Sobel operators (3x3 kernels) to detect horizontal 
// and vertical edges. It utilizes goroutines to process rows in parallel.
func computeGradients(gray *image.Gray, width, height int) ([][]float64, [][]float64) {
	mags := make([][]float64, height)
	angles := make([][]float64, height)
	for i := range mags {
		mags[i] = make([]float64, width)
		angles[i] = make([]float64, width)
	}

	// Sobel Kernels for X and Y direction
	gx := [][]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	gy := [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	rowsPerWorker := (height - 2) / numCPU

	for w := 0; w < numCPU; w++ {
		wg.Add(1)
		startY := 1 + (w * rowsPerWorker)
		endY := startY + rowsPerWorker
		if w == numCPU-1 { endY = height - 1 }

		go func(yMin, yMax int) {
			defer wg.Done()
			for y := yMin; y < yMax; y++ {
				for x := 1; x < width-1; x++ {
					var sumX, sumY float64
					// Convolution with Sobel Kernels
					for i := -1; i <= 1; i++ {
						for j := -1; j <= 1; j++ {
							lum := float64(gray.GrayAt(x+j, y+i).Y)
							sumX += lum * float64(gx[i+1][j+1])
							sumY += lum * float64(gy[i+1][j+1])
						}
					}
					// Calculate magnitude (hypotenuse) and direction (angle)
					mags[y][x] = math.Sqrt(sumX*sumX + sumY*sumY)
					angles[y][x] = math.Mod(math.Atan2(sumY, sumX)*180/math.Pi+180, 180)
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return mags, angles
}

// applyNMS (Non-Maximum Suppression) thins out blurry gradients.
// It keeps a pixel only if it is a local maximum in the direction of the gradient.
func applyNMS(mags, angles [][]float64, width, height int) [][]float64 {
	nmsMags := make([][]float64, height)
	for i := range nmsMags { nmsMags[i] = make([]float64, width) }

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			angle := angles[y][x]
			mag := mags[y][x]
			var q, r float64

			// Determine neighboring pixels based on gradient orientation
			switch {
			case (angle >= 0 && angle < 22.5) || (angle >= 157.5 && angle <= 180):
				q, r = mags[y][x+1], mags[y][x-1] // Horizontal
			case angle >= 22.5 && angle < 67.5:
				q, r = mags[y+1][x-1], mags[y-1][x+1] // Diagonal
			case angle >= 67.5 && angle < 112.5:
				q, r = mags[y+1][x], mags[y-1][x] // Vertical
			case angle >= 112.5 && angle < 157.5:
				q, r = mags[y-1][x-1], mags[y+1][x+1] // Anti-Diagonal
			}

			if mag >= q && mag >= r {
				nmsMags[y][x] = mag
			}
		}
	}
	return nmsMags
}

// mapToWeights aggregates pixel magnitudes into discrete grid cells.
// It uses thresholds to filter noise and decide which walls are strongest.
func mapToWeights(nmsMags [][]float64, angles [][]float64, rows, cols, width, height int) map[string]int {
	weights := make(map[string]int)

	highThresh := 80.0 // Magnitudes above this are guaranteed walls
	lowThresh := 30.0  // Magnitudes below this are ignored (noise)

	for r := range rows {
		for c := range cols {
			startY, endY := r*height/rows, (r+1)*height/rows
			startX, endX := c*width/cols, (c+1)*width/cols

			maxMag := 0.0
			bestAngle := 0.0

			// Find the most significant edge within this cell's pixel area
			for y := startY; y < endY && y < height; y++ {
				for x := startX; x < endX && x < width; x++ {
					if nmsMags[y][x] > maxMag {
						maxMag = nmsMags[y][x]
						bestAngle = angles[y][x]
					}
				}
			}

			if maxMag < lowThresh { continue }

			// Interpolate weight based on magnitude
			var weight int
			if maxMag >= highThresh {
				weight = 255
			} else {
				weight = 120 + int((maxMag-lowThresh)/(highThresh-lowThresh)*135)
			}

			// Assign the weight to either Top or Left wall based on edge angle
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
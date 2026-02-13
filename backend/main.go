// Package main provides a CLI interface for the GridGo maze generation engine.
// It supports traditional randomized algorithms and computer-vision-guided
// generation based on image silhouettes.
package main

import (
	"fmt"
	"image"
	"maze-gen/maze"
	"os"
)

func main() {
	var rows, cols, choice int
	var imgPath string

	// algorithm selection
	for {
		fmt.Println("\n--- GridGo Maze Generation Menu ---")
		fmt.Println("1. Image-Based (Trace an outline from a photo)")
		fmt.Println("2. Randomized Kruskal's (Short passages, many dead ends)")
		fmt.Println("3. Recursive Backtracker (DFS - Long, winding corridors)")
		fmt.Print("Selection: ")
		if _, err := fmt.Scan(&choice); err == nil && (choice >= 1 && choice <= 3) {
			break
		}
		fmt.Println("Invalid choice, please enter 1, 2, or 3.")
	}

	// choosing dimension
	if choice == 1 {
		for {
			fmt.Print("Enter path to image file (e.g., puppy.jpg): ")
			fmt.Scan(&imgPath)

			file, err := os.Open(imgPath)
			if err != nil {
				fmt.Printf("Error: %v. Please try again.\n", err)
				continue
			}
			img, _, err := image.DecodeConfig(file)
			file.Close()
			if err != nil {
				fmt.Printf("Error decoding image config: %v. Please try again.\n", err)
				continue
			}

			fmt.Printf("Image found: %d x %d\n", img.Width, img.Height)
			fmt.Print("Use image pixel dimensions as grid size? (1=Default, 2=Custom): ")
			var sizeChoice int
			fmt.Scan(&sizeChoice)

			if sizeChoice == 1 {
				rows, cols = img.Height, img.Width
				break
			} else {
				fmt.Print("Enter custom rows: ")
				fmt.Scan(&rows)
				fmt.Print("Enter custom columns: ")
				fmt.Scan(&cols)
				if rows >= 2 && cols >= 2 {
					break
				}
				fmt.Println("Invalid custom dimensions.")
			}
		}
	} else {
		for {
			fmt.Print("Enter number of rows (min 2): ")
			if _, err := fmt.Scan(&rows); err == nil && rows >= 2 {
				break
			}
			fmt.Println("Invalid input.")
		}
		for {
			fmt.Print("Enter number of columns (min 2): ")
			if _, err := fmt.Scan(&cols); err == nil && cols >= 2 {
				break
			}
			fmt.Println("Invalid input.")
		}
	}

	myMaze := maze.NewMaze(rows, cols)

	switch choice {
	case 1:
		// get weights using canny pipeline
		weights, err := maze.GetEdgeWeights(imgPath, rows, cols)
		if err != nil {
			fmt.Printf("Critical error processing edges: %v\n", err)
			return
		}
		myMaze.GenerateImageMaze(weights)
	case 2:
		myMaze.GenerateKruskal()
	case 3:
		myMaze.GenerateRecursive(0, 0)
	}

	// post processing entrance/exit
	// clear top left wall and the bottom right wall
	myMaze.Grid[0][0].Walls[0] = false
	myMaze.Grid[rows-1][cols-1].Walls[2] = false

	// output
	fmt.Print("\nDisplay in terminal (1) or Save as PNG (2)? ")
	var displayChoice int
	fmt.Scan(&displayChoice)

	if displayChoice == 2 {
		fmt.Print("Enter output filename (e.g., result.png): ")
		var outName string
		fmt.Scan(&outName)

		err := myMaze.SaveAsImage(outName, 10)
		if err != nil {
			fmt.Printf("Error saving image: %v\n", err)
		} else {
			fmt.Printf("Successfully saved to %s\n", outName)
		}
	} else {
		myMaze.Print()
	}
}

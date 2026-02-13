package main

import (
	"fmt"
	"maze-gen/maze"
)

func main() {
	var rows, cols, choice int
	var imgPath string

	// get dimensions
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

	// select algorithm
	for {
		fmt.Println("\n--- Maze Generation Menu ---")
		fmt.Println("1. Image-Based (Trace an outline)")
		fmt.Println("2. Randomized Kruskal's (Short passages)")
		fmt.Println("3. Recursive Backtracker (DFS - Long corridors)")
		fmt.Print("Selection: ")
		if _, err := fmt.Scan(&choice); err == nil && (choice >= 1 && choice <= 3) {
			break
		}
		fmt.Println("Invalid choice, please enter 1, 2, or 3.")
	}

	myMaze := maze.NewMaze(rows, cols)

	switch choice {
	case 1:
		for {
			fmt.Print("Enter path to image file (e.g., test.jpg): ")
			fmt.Scan(&imgPath)

			// get sobel weights
			weights, err := maze.GetEdgeWeights(imgPath, rows, cols)
			if err == nil {
				myMaze.GenerateImageMaze(weights)
				break
			}
			fmt.Printf("Error loading image: %v. Please try again.\n", err)
		}
	case 2:
		myMaze.GenerateKruskal()
	case 3:
		myMaze.GenerateRecursive(0, 0)
	}

	myMaze.Grid[0][0].Walls[0] = false
	myMaze.Grid[rows-1][cols-1].Walls[2] = false

	fmt.Print("\nDisplay in terminal (1) or Save as PNG (2)? ")

	var displayChoice int

	fmt.Scan(&displayChoice)

	if displayChoice == 2 {
		fmt.Print("Enter output filename (e.g., maze.png): ")
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

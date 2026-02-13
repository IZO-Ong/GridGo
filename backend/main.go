package main

import (
	"fmt"
	"maze-gen/maze"
)

func main() {
	var rows, cols, choice int

	for {
		fmt.Print("Enter number of rows (minimum 2): ")
		_, err := fmt.Scan(&rows)

		if err == nil && rows >= 2 {
			break
		}

		fmt.Println("Invalid input, please try again.")
	}

	for {
		fmt.Print("Enter number of columns (minimum 2): ")
		_, err := fmt.Scan(&cols)

		if err == nil && cols >= 2 {
			break
		}

		fmt.Println("Invalid input, please try again.")
	}

	for {
		fmt.Println("\nChoose Generation Algorithm:")
		fmt.Println("1. Randomized Kruskal's (Short passages, many dead ends)")
		fmt.Println("2. Recursive Backtracker (Long, winding corridors)")
		fmt.Print("Selection: ")
		_, err := fmt.Scan(&choice)

		if err == nil && (choice == 1 || choice == 2) {
			break
		}
		fmt.Println("Invalid choice, please enter 1 or 2.")
	}

	myMaze := maze.NewMaze(rows, cols)

	if choice == 1 {
		myMaze.GenerateKruskal()
	} else {
		myMaze.GenerateRecursive(0, 0)
	}

	myMaze.Grid[0][0].Walls[0] = false
	myMaze.Grid[rows-1][cols-1].Walls[2] = false

	myMaze.Print()
}

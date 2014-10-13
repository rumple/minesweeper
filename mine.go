// The MIT License (MIT)
//
// Copyright (c) 2014 rumple
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Code to solve https://code.google.com/codejam/contest/5214486/dashboard
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// A cell in the grid
type Cell struct {
	neighbour int  //Neighbouring mine counts
	visited   bool //Tracks if cell is revealed to user
}

// Minesweeper grid
type Grid struct {
	size int    // size of the grid
	g    []Cell // Actual grid (row first)
}

// Check if given row and col has mine
func (g *Grid) isMine(row, col int) bool {
	if g.g[((row*g.size)+col)].neighbour == -1 {
		return true
	}
	return false
}

// Check if cell at row and vol is revealed
func (g *Grid) isVisited(row, col int) bool {
	return g.g[((row * g.size) + col)].visited
}

// Helper function to increament the neighbouring mine count for a cell
func (g *Grid) incCell(row, col int) {
	if !g.isMine(row, col) {
		g.g[((row*g.size)+col)].neighbour++
	}
}

// Get neighbouring mine count for a cell at row and col
func (g *Grid) getVal(row, col int) int {
	return g.g[((row * g.size) + col)].neighbour
}

// Mark the cell at row and col revealed to user
func (g *Grid) markVisited(row, col int) {
	g.g[((row * g.size) + col)].visited = true
}

// Show the grid
func (g *Grid) show() {
	for row := 0; row < g.size; row++ {
		for col := 0; col < g.size; col++ {
			if g.isMine(row, col) {
				fmt.Printf(" * ")
				continue
			}
			if g.isVisited(row, col) {
				fmt.Printf(" %d ", g.getVal(row, col))
			} else {
				fmt.Printf(" - ")
			}
		}
		fmt.Printf("\n")
	}
}

// Recursively reveal the cells at row and col
func (g *Grid) visit(row, col int) {
	// Base recursion case
	if g.isVisited(row, col) {
		return
	}

	// This calculates the neighbour row and col bounds
	startrow := row - 1
	if startrow < 0 {
		startrow = 0
	}
	endrow := row + 1
	if endrow == g.size {
		endrow = g.size - 1
	}
	startcol := col - 1
	if startcol < 0 {
		startcol = 0
	}
	endcol := col + 1
	if endcol == g.size {
		endcol = g.size - 1
	}

	// Mark it visitier
	g.markVisited(row, col)

	// Iterate over all neighbours and recursively visit if mine count if 0
	for nrow := startrow; nrow <= endrow; nrow++ {
		for ncol := startcol; ncol <= endcol; ncol++ {
			if nrow == row && ncol == col {
				continue
			}
			if g.isMine(row, col) {
				continue
			}
			if g.getVal(nrow, ncol) == 0 {
				g.visit(nrow, ncol)
			} else {
				g.markVisited(nrow, ncol)
			}
		}
	}
}

// Calculate neighbouring mine count for all the cells.
// This walks the grid and increaments the mine count for all the neighbouring
// cells for a mine.
func (g *Grid) getMines() {
	for row := 0; row < g.size; row++ {
		for col := 0; col < g.size; col++ {
			// If mine is found then increament the mine count for all the
			// neighbouring cells
			if g.isMine(row, col) {
				// Calculate the bounds for neighbours
				startrow := row - 1
				if startrow < 0 {
					startrow = 0
				}
				endrow := row + 1
				if endrow == g.size {
					endrow = g.size - 1
				}
				startcol := col - 1
				if startcol < 0 {
					startcol = 0
				}
				endcol := col + 1
				if endcol == g.size {
					endcol = g.size - 1
				}
				// Increament mine count for neighbours
				for nrow := startrow; nrow <= endrow; nrow++ {
					for ncol := startcol; ncol <= endcol; ncol++ {
						if nrow == row && ncol == col {
							continue
						}
						g.incCell(nrow, ncol)
					}
				}
			}
		}
	}
}

// Solve the grid.
func (g *Grid) solve() int {
	clicks := 0
	// Step 1: Click on all unvisited 0 cells.
	for row := 0; row < g.size; row++ {
		for col := 0; col < g.size; col++ {
			if g.isMine(row, col) || g.isVisited(row, col) {
				continue
			}
			if g.getVal(row, col) == 0 {
				//fmt.Println("Clicked at ", row, col)
				clicks++
				g.visit(row, col)
				//g.show()
			}

		}
	}
	// Step 2: Click on all the unvisited cells
	for row := 0; row < g.size; row++ {
		for col := 0; col < g.size; col++ {
			if g.isMine(row, col) {
				continue
			}
			if !g.isVisited(row, col) {
				g.markVisited(row, col)
				//fmt.Println("Clicked at ", row, col)
				clicks++
			}
		}
	}
	return clicks
}

// Parse the string and create the grid
func getGrid(m string) *Grid {
	var g Grid
	lines := strings.Split(m, "\n")
	g.size = len(lines) - 1
	for _, i := range m {
		if i == '.' {
			g.g = append(g.g, Cell{neighbour: 0})
		}
		if i == '*' {
			g.g = append(g.g, Cell{neighbour: -1})
		}
	}
	return &g

}

func main() {
	//file, err := os.Open("A-small-practice.in")
	file, err := os.Open("A-large-practice.in")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	ntest, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	for n := 0; n < int(ntest); n++ {
		scanner.Scan()
		size, _ := strconv.ParseInt(scanner.Text(), 10, 32)
		input := ""
		for i := 0; i < int(size); i++ {
			scanner.Scan()
			input = input + scanner.Text() + "\n"
		}
		// Create the grid
		g := getGrid(input)

		// Calculate cell mine count
		g.getMines()

		// Solve
		fmt.Printf("Case #%d: %d\n", n+1, g.solve())
	}
	return
}

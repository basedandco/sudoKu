package main

/*
#cgo CFLAGS: -I/usr/include/security
#cgo LDFLAGS: -lpam
#include <security/pam_appl.h>
#include <security/pam_modules.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
extern int pam_sm_authenticate_real(void *, int, int, char **);
extern int pam_sm_setcred_real(void *, int, int, char **);
#define pam_sm_authenticate pam_sm_authenticate_real
#define pam_sm_setcred      pam_sm_setcred_real
*/
import "C"

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	PAM_SUCCESS  = 0
	PAM_AUTH_ERR = 7
)

func generateSudoku() (board [][]int, solution [][]int) {
	board = [][]int{
		{1, 2, 3, 4},
		{3, 4, 1, 2},
		{2, 1, 4, 3},
		{4, 3, 2, 1},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(4, func(i, j int) { board[i], board[j] = board[j], board[i] })

	solution = make([][]int, 4)
	for i := range board {
		solution[i] = make([]int, 4)
		copy(solution[i], board[i])
	}

	emptyCount := r.Intn(3) + 4
	for i := 0; i < emptyCount; i++ {
		ri, ci := r.Intn(4), r.Intn(4)
		if board[ri][ci] == 0 {
			i--
			continue
		}
		board[ri][ci] = 0
	}
	return
}

func printBoard(board [][]int) {
	fmt.Println("Solve this 4x4 Sudoku (use numbers 1-4):")
	fmt.Println("┌───┬───┬───┬───┐")
	for i, row := range board {
		fmt.Print("│")
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(" _ │")
			} else {
				fmt.Printf(" %d │", cell)
			}
		}
		fmt.Println()
		if i < 3 {
			fmt.Println("├───┼───┼───┼───┤")
		}
	}
	fmt.Println("└───┴───┴───┴───┘")
}

func readUserSolution() ([][]int, error) {
	fmt.Println("Enter your solution, 4 rows of 4 numbers (space separated):")
	scanner := bufio.NewScanner(os.Stdin)
	userBoard := make([][]int, 4)
	for i := 0; i < 4; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read input")
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) != 4 {
			return nil, fmt.Errorf("row %d has %d numbers", i+1, len(fields))
		}
		row := make([]int, 4)
		for j, v := range fields {
			num, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %v", v)
			}
			if num < 1 || num > 4 {
				return nil, fmt.Errorf("number %d out of range (must be 1-4)", num)
			}
			row[j] = num
		}
		userBoard[i] = row
	}
	return userBoard, nil
}

func isValidSudoku(board [][]int) bool {
	for i := 0; i < 4; i++ {
		seen := make(map[int]bool)
		for j := 0; j < 4; j++ {
			if board[i][j] == 0 {
				continue
			}
			if seen[board[i][j]] {
				return false
			}
			seen[board[i][j]] = true
		}
	}

	for j := 0; j < 4; j++ {
		seen := make(map[int]bool)
		for i := 0; i < 4; i++ {
			if board[i][j] == 0 {
				continue
			}
			if seen[board[i][j]] {
				return false
			}
			seen[board[i][j]] = true
		}
	}

	return true
}

func checkSolution(user, solution [][]int) bool {
	if !isValidSudoku(user) {
		return false
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if user[i][j] != solution[i][j] {
				return false
			}
		}
	}
	return true
}

//export pam_sm_authenticate_real
func pam_sm_authenticate_real(pamh unsafe.Pointer, flags C.int, argc C.int, argv **C.char) C.int {
	board, solution := generateSudoku()
	printBoard(board)
	user, err := readUserSolution()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input error:", err)
		return PAM_AUTH_ERR
	}
	if checkSolution(user, solution) {
		return PAM_SUCCESS
	}
	fmt.Fprintln(os.Stderr, "Sudoku incorrect.")
	return PAM_AUTH_ERR
}

//export pam_sm_setcred_real
func pam_sm_setcred_real(pamh unsafe.Pointer, flags C.int, argc C.int, argv **C.char) C.int {
	return PAM_SUCCESS
}

func runStandalone() {
	for {
		board, solution := generateSudoku()
		printBoard(board)
		user, err := readUserSolution()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Input error:", err)
			continue
		}
		if isValidSudoku(user) {
			if checkSolution(user, solution) {
				fmt.Println("Sudoku correct! Well done!")
			} else {
				fmt.Println("Sudoku incorrect. Try again!")
			}
		} else {
			fmt.Println("Your solution violates Sudoku rules (each row/column must contain unique numbers 1-4)")
		}

		fmt.Print("Play again? (y/n): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() && strings.ToLower(scanner.Text()) != "y" {
			break
		}
	}
}

func main() {
	if os.Getenv("PAM_MODE") != "1" {
		runStandalone()
	}
}

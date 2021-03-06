package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Pair struct {
	x    int
	y    int
	name string
}

type Move struct {
	curr  Pair
	final Pair
}

type ShogiState struct {
	board  [][]string
	pieces []Pair
	parent *ShogiState
}

func (m Move) String() string {
	s := fmt.Sprintf("%d %d %d %d", m.curr.x+1, m.curr.y+1, m.final.x+1, m.final.y+1)
	return s
}

func (state ShogiState) String() string {
	var s string
	for i := 0; i < len(state.board); i++ {
		for j := 0; j < len(state.board[i]); j++ {
			if state.board[i][j] == "O" {
				s = s + fmt.Sprint(state.board[i][j]+"  ")
			} else {
				s = s + fmt.Sprint(state.board[i][j]+" ")
			}
		}
		s = s + fmt.Sprint("\n")
	}
	// var s1 string
	// for i := 0; i < len(state.pieces); i++ {
	// 	s1 = s1 + fmt.Sprintf("{x: %d, y: %d, name: %s} ", state.pieces[i].x, state.pieces[i].y, state.pieces[i].name)
	// }
	// var s2 string
	// s2 = fmt.Sprintf("\n %p\n", state.parent)

	return s //+ s1 + s2
}

func popList(list []ShogiState) []ShogiState {
	return list[1:]
}

func (state ShogiState) IsGoal(player int) bool { //returns true if a player has won
	var opponent int
	if player == 1 {
		opponent = 2
	}
	if player == 2 {
		opponent = 1
	}
	NextBoard := Succ(state, player)
	AllCheck := true //not all future moves are in check
	for i := 0; i < len(NextBoard); i++ {
		if !Check(NextBoard[i], opponent) { //do all of them involve 2 being in check
			AllCheck = false //if not then theres a move to be made
			break            //no need to keep searching
		}
	}
	if AllCheck { //all moves would still be check
		PlayerNum := strconv.Itoa(player)
		fmt.Println("player" + PlayerNum + "wins!") //other player wins
		return true
	}

	return false //valid moves can still be made
}

func Check(state ShogiState, player int) bool { //returns true if player is in check
	var opponent int
	if player == 1 {
		opponent = 2
	}
	if player == 2 {
		opponent = 1
	}
	CurrBoard := Succ(state, opponent)
	playerNum := strconv.Itoa(player)
	YourKing := "K" + playerNum
	FoundKingAll := true //found our king on all boards, no checks
	for a := 0; a < len(CurrBoard); a++ {
		FoundKingInBoard := false
		for i := 0; i < len(CurrBoard[a].board); i++ {
			for j := 0; j < len(CurrBoard[a].board[i]); j++ {
				if CurrBoard[a].board[i][j] == YourKing {
					FoundKingInBoard = true //found our king on this board
					break
				}
			}
		}
		if FoundKingInBoard == false { //didn't find it in this board
			FoundKingAll = false //there is a move where we are in check
		}
	}
	return !FoundKingAll
}

func duplicate(state ShogiState) ShogiState {
	newBoard := make([][]string, len(state.board))
	for i := 0; i < len(newBoard); i++ {
		newBoard[i] = make([]string, len(state.board[i]))
		copy(newBoard[i], state.board[i])
	}
	newPieces := make([]Pair, len(state.pieces))
	copy(newPieces, state.pieces)
	newState := ShogiState{board: newBoard, pieces: newPieces, parent: state.parent}
	return newState
}

func (state ShogiState) Equal(s ShogiState) bool {
	if len(state.board) != len(s.board) {
		return false
	}
	for i := 0; i < len(state.board); i++ {
		if len(state.board[i]) != len(s.board[i]) {
			return false
		}
		for j := 0; j < len(state.board[i]); j++ {
			if state.board[i][j] != s.board[i][j] {
				return false
			}
		}
	}

	if len(state.pieces) != len(s.pieces) {
		return false
	}
	for i := 0; i < len(state.pieces); i++ {
		for j := 0; j < len(state.pieces); j++ {
			if (state.pieces[i].x == state.pieces[j].x) != (s.pieces[i].y == s.pieces[j].y) {
				return false
			}
		}
	}
	return true
}

func Succ(state ShogiState, player int) []ShogiState {
	var final []ShogiState
	for i := 0; i < len(state.pieces); i++ {
		//Scan through all pieces and appends all possible moves of all pieces to the final slice
		if OwnsPiece(state.pieces[i].name, player) {
			switch state.pieces[i].name {
			case "P1", "P2":
				NewX := state.pieces[i].x + 0
				var NewY int
				if strings.Contains(state.pieces[i].name, "1") {
					NewY = state.pieces[i].y + 1
				}
				if strings.Contains(state.pieces[i].name, "2") {
					NewY = state.pieces[i].y - 1
				}
				NewState := MakeMove(state, player, NewX, NewY, i) //gives either a new state or if its invalid the same state
				NewState.parent = &state
				if !state.Equal(NewState) {
					final = append(final, NewState)
				}
			case "L1":
				for j := state.pieces[i].y; j < len(state.board[0]); j++ {
					NewX := state.pieces[i].x
					NewY := j
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
					if strings.Contains(state.board[NewY][NewX], "1") {
						break
					}
				}
			case "L2":
				for j := state.pieces[i].y; j >= 0; j-- {
					NewX := state.pieces[i].x
					NewY := j
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
					if strings.Contains(state.board[NewY][NewX], "2") {
						break
					}
				}
			case "N1", "N2":
				NewX := state.pieces[i].x - 1
				var NewY int
				if strings.Contains(state.pieces[i].name, "1") {
					NewY = state.pieces[i].y - 2
				}
				if strings.Contains(state.pieces[i].name, "2") {
					NewY = state.pieces[i].y + 2
				}
				NewState := MakeMove(state, player, NewX, NewY, i)
				NewState.parent = &state
				if !state.Equal(NewState) {
					final = append(final, NewState)
				}
				NewX = state.pieces[i].x + 1
				NewState = MakeMove(state, player, NewX, NewY, i)
				NewState.parent = &state
				if !state.Equal(NewState) {
					final = append(final, NewState)
				}
			case "S1", "S2":
				var moves []int
				if strings.Contains(state.pieces[i].name, "1") {
					moves = []int{-1, -1, 0, -1, 1, -1, -1, 1, 1, 1}
				}
				if strings.Contains(state.pieces[i].name, "2") {
					moves = []int{-1, 1, 0, 1, 1, 1, -1, -1, 1, -1}
				}
				for j := 0; j <= len(moves)-2; j += 2 {
					NewX := state.pieces[i].x + moves[j]
					NewY := state.pieces[i].y + moves[j+1]
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
				}
			case "G1", "G2", "P1+", "P2+", "L1+", "L2+", "N1+", "N2+", "S1+", "S2+":
				var moves []int
				if strings.Contains(state.pieces[i].name, "1") {
					moves = []int{-1, -1, 0, -1, 1, -1, 1, 0, -1, 0, 0, 1}
				}
				if strings.Contains(state.pieces[i].name, "2") {
					moves = []int{-1, 1, 0, 1, 1, 1, 0, -1, 1, 0, 0, -1}
				}
				for j := 0; j <= len(moves)-2; j += 2 {
					NewX := state.pieces[i].x + moves[j]
					NewY := state.pieces[i].y + moves[j+1]
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
				}
			case "B1", "B2":
				for j := 0; j < len(state.board); j++ {
					NewX := state.pieces[i].x + j
					NewY := state.pieces[i].y + j
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
				}
				for j := len(state.board); j >= 0; j-- {
					NewX := state.pieces[i].x - j
					NewY := state.pieces[i].y - j
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
				}
			case "R1", "R2":
				for j := state.pieces[i].x; j < len(state.board); j++ {
					NewX := j
					NewY := state.pieces[i].y
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
					var StrPlayer string
					if state.pieces[i].name == "R1" {
						StrPlayer = "1"
					} else {
						StrPlayer = "2"
					}
					if strings.Contains(state.board[NewY][NewX], StrPlayer) {
						break
					}
				}
				for j := state.pieces[i].x; j >= 0; j-- {
					NewX := j
					NewY := state.pieces[i].y
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
					var StrPlayer string
					if state.pieces[i].name == "R1" {
						StrPlayer = "1"
					} else {
						StrPlayer = "2"
					}
					if strings.Contains(state.board[NewY][NewX], StrPlayer) {
						break
					}
				}
				for g := state.pieces[i].y; g < len(state.board[0]); g++ {
					NewX := state.pieces[i].x
					NewY := g
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							fmt.Println(NewState)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
					var StrPlayer string
					if state.pieces[i].name == "R1" {
						StrPlayer = "1"
					} else {
						StrPlayer = "2"
					}
					if strings.Contains(state.board[NewY][NewX], StrPlayer) {
						break
					}
				}
				for g := state.pieces[i].y; g >= 0; g-- {
					NewX := state.pieces[i].x
					NewY := g
					if IsValid(state.board, NewX, NewY, player) {
						if state.board[NewY][NewX] != "O" {
							NewState := MakeMove(state, player, NewX, NewY, i)
							fmt.Println(NewState)
							NewState.parent = &state
							if !state.Equal(NewState) {
								final = append(final, NewState)
							}
							break
						}
					}
					var StrPlayer string
					if state.pieces[i].name == "R1" {
						StrPlayer = "1"
					} else {
						StrPlayer = "2"
					}
					if strings.Contains(state.board[NewY][NewX], StrPlayer) {
						break
					}
				}
			case "K1", "K2":
				var moves []int
				moves = []int{-1, -1, 0, -1, 1, -1, -1, 1, 1, 1, 1, 0, -1, 0, 0, 1}
				for j := 0; j <= len(moves)-2; j += 2 {
					NewX := state.pieces[i].x + moves[j]
					NewY := state.pieces[i].y + moves[j+1]
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
				}
			case "B1+", "B2+":
				for j := 0; j < len(state.board); j++ {
					NewX := state.pieces[i].x + j
					NewY := state.pieces[i].y + j
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
					if state.board[NewY][NewX] != "O" {
						break
					}
				}
				for j := len(state.board); j >= 0; j-- {
					NewX := state.pieces[i].x - j
					NewY := state.pieces[i].y - j
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
					if state.board[NewY][NewX] != "O" {
						break
					}
				}
				var moves []int
				moves = []int{0, 1, -1, 0, 1, 0, 0, -1}
				for j := 0; j <= len(moves)-2; j += 2 {
					NewX := state.pieces[i].x + moves[j]
					NewY := state.pieces[i].y + moves[j+1]
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
				}
			case "R1+", "R2+":
				for j := 0; j < len(state.board); j++ {
					NewX := j
					NewY := state.pieces[i].y
					if state.board[NewY][NewX] != "O" {
						NewState := MakeMove(state, player, NewX, NewY, i)
						NewState.parent = &state
						if !state.Equal(NewState) {
							final = append(final, NewState)
						}
						break
					}
				}
				for g := 0; g < len(state.board[0]); g++ {
					NewX := state.pieces[i].x
					NewY := g
					if state.board[NewY][NewX] != "O" {
						NewState := MakeMove(state, player, NewX, NewY, i)
						NewState.parent = &state
						if !state.Equal(NewState) {
							final = append(final, NewState)
						}
						break
					}
				}
				var moves []int
				moves = []int{1, 1, -1, 1, 1, -1, -1, -1}
				for j := 0; j <= len(moves)-2; j += 2 {
					NewX := state.pieces[i].x + moves[j]
					NewY := state.pieces[i].y + moves[j+1]
					NewState := MakeMove(state, player, NewX, NewY, i)
					NewState.parent = &state
					if !state.Equal(NewState) {
						final = append(final, NewState)
					}
				}
			}
		}
	}
	return final
}

func diff(s1, s2 ShogiState) (Move, error) {
	var m Move
	for i := 0; i < len(s1.pieces); i++ {
		if s1.pieces[i].x != s2.pieces[i].x || s1.pieces[i].y != s2.pieces[i].y {
			m.curr = s1.pieces[i]
			m.final = s2.pieces[i]
			return m, nil
		}
	}
	return m, fmt.Errorf("No diff! You're a liar! You promiced me!")
}

func getFirstMove(state ShogiState) Move {
	//Totally untested and highly dangerous
	curr := state
	final := curr
	for curr.parent != nil {
		final = curr
		curr = *(curr.parent)
	}
	m, err := diff(curr, final)
	if err != nil {
		panic("No move made, even though moves made! Someone call the Navy!")
	}

	return m
}

func Max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func Min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func auxMiniMax(state ShogiState, player, depth int, max bool) int {
	// fmt.Println(state)
	if state.IsGoal(player) {
		//Oh Mowie Wowie!
		if max {
			return 9999
		}
		return -9999
	}

	if depth == 0 {
		if player == 1 {
			return h1(state.board)
		} else {
			return h2(state.board)
		}
	} else if max {
		val := 0
		kids := Succ(state, player)
		for i := 0; i < len(kids); i++ {
			val = Max(val, auxMiniMax(kids[i], player, depth-1, false))
		}
		return val
	} else {
		val := 0
		kids := Succ(state, player)
		for i := 0; i < len(kids); i++ {
			val = Min(val, auxMiniMax(kids[i], player, depth-1, true))
		}
		return val
	}
}

func (m1 Move) MoveEqual(m2 Move) bool {
	if m1.curr.x != m2.curr.x {
		return false
	}
	if m1.curr.y != m2.curr.y {
		return false
	}
	if m1.final.x != m2.final.x {
		return false
	}
	if m1.final.y != m2.final.y {
		return false
	}
	return true
}

func (state Move) isDup(history []Move) bool {
	for i := 0; i < len(history); i++ {
		if state.MoveEqual(history[i]) {
			return true
		}
	}
	return false
}

func MiniMax(state ShogiState, player, depth int, history []Move) (Move, error) {
	//Totally untested, und highly dangerous! (waiting for Succ)
	if state.IsGoal(player) {
		return Move{}, fmt.Errorf("You won dufus! Email all your friends!")
	}

	kids := Succ(state, player)
	var vals []int
	var finalKids []ShogiState
	for i := 0; i < len(kids); i++ {
		kdiff, err := diff(state, kids[i]) //rm duplicates
		if err != nil {
			continue
		}
		if kdiff.isDup(history) {
			continue
		}

		// fmt.Println(k)
		// fmt.Println(kids[i])
		val := auxMiniMax(kids[i], player, depth-1, false)
		vals = append(vals, val)
		history = append(history, kdiff)
		finalKids = append(finalKids, kids[i])
	}

	maximum := 0
	for i := 0; i < len(vals); i++ {
		if vals[i] > vals[maximum] {
			maximum = i
		}
	}

	// fmt.Println(finalKids[maximum])
	m, err := diff(state, finalKids[maximum])
	if err != nil {
		panic("No move is max, this shouldn't be possible because maximum := 0")
	}
	return m, nil
}

func OwnsPiece(piece string, playerNum int) bool {
	strPlayerNum := strconv.Itoa(playerNum)
	if strings.Contains(piece, strPlayerNum) {
		return true
	}
	return false
}

func IsValid(board [][]string, NewX int, NewY int, player int) bool {
	if NewX >= len(board[0]) || NewX < 0 {
		return false
	}
	if NewY >= len(board) || NewY < 0 {
		return false
	}
	StrPlayer := strconv.Itoa(player)
	if strings.Contains(board[NewY][NewX], StrPlayer) {
		return false
	}
	return true
}

func CheckPromotion(Newy int, piece string) bool {
	if Newy > 5 {
		switch piece {
		case "P1":
			return true
		case "L1":
			return true
		case "N1":
			return true
		case "S1":
			return true
		case "B1":
			return true
		case "R1":
			return true
		}
	}
	if Newy < 3 {
		switch piece {
		case "P2":
			return true
		case "L2":
			return true
		case "N2":
			return true
		case "S2":
			return true
		case "B2":
			return true
		case "R2":
			return true
		}
	}
	return false
}

//validates movements and make changes, then sends back a changed state
func MakeMove(state ShogiState, player int, NewX int, NewY int, i int) ShogiState {
	newState := duplicate(state)
	if newState.pieces[i].x == NewX && newState.pieces[i].y == NewY {
		return newState
	}
	if IsValid(state.board, NewX, NewY, player) {
		if newState.board[NewY][NewX] != "O" {
			for j := 0; j < len(newState.pieces); j++ {
				if newState.pieces[j].x == NewX && newState.pieces[j].y == NewY {
					newState.pieces[j].name = ""
				}
			}
		}
		piece := newState.pieces[i].name
		newState.board[newState.pieces[i].y][newState.pieces[i].x] = "O"
		newState.board[NewY][NewX] = piece //update board
		newState.pieces[i].x = NewX        //update piece
		newState.pieces[i].y = NewY
		if !strings.Contains(piece, "+") {
			if !strings.Contains(piece, "K") {
				if CheckPromotion(NewY, newState.pieces[i].name) {
					newState.board[NewY][NewX] = piece + "+"
					newState.pieces[i].name = piece + "+"
				}
			}
		}
	}
	return newState
}

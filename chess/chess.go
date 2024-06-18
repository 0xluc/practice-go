package chess

import (
	"errors"
)

func CanKnightAttack(square1 string, square2 string) (bool, error) {
	if square1 == square2 {
		return false, errors.New("both squares are equal")
	}
	result1 := validatePosition(square1)
	if result1 == false {
		return false, errors.New("square1 not valid")
	}
	result2 := validatePosition(square2)
	if result2 == false {
		return false, errors.New("square1 not valid")
	}

	possibleMoves1 := getValidMoves(square1)

	for _, move1 := range possibleMoves1 {
		if move1 == square2 {
			return true, nil
		}
	}
	return false, nil

}

func getValidMoves(initPos string) []string {
	ch := initPos[0]
	num := initPos[1]
	validMoves := []string{}

	if validatePosition(string(ch+2) + string(num+1)) {
		validMoves = append(validMoves, string(ch+2)+string(num+1))
	}
	if validatePosition(string(ch+2) + string(num-1)) {
		validMoves = append(validMoves, string(ch+2)+string(num-1))
	}
	if validatePosition(string(ch+1) + string(num+2)) {
		validMoves = append(validMoves, string(ch+1)+string(num+2))
	}
	if validatePosition(string(ch+1) + string(num-2)) {
		validMoves = append(validMoves, string(ch+1)+string(num-2))
	}
	if validatePosition(string(ch-1) + string(num+2)) {
		validMoves = append(validMoves, string(ch-1)+string(num+2))
	}
	if validatePosition(string(ch-1) + string(num-2)) {
		validMoves = append(validMoves, string(ch-1)+string(num-2))
	}
	if validatePosition(string(ch-2) + string(num-1)) {
		validMoves = append(validMoves, string(ch-2)+string(num-1))
	}
	if validatePosition(string(ch-2) + string(num+1)) {
		validMoves = append(validMoves, string(ch-2)+string(num+1))
	}

	return validMoves

}

func validatePosition(pos string) bool {
	if len(pos) != 2 {
		return false
	}
	ch := pos[0]
	num := pos[1]
	if ch < 'a' || ch > 'h' {
		return false
	}
	if num < '1' || num > '8' {
		return false
	}
	return true
}

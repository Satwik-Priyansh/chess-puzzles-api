package rating

import (
	"math"
)

var q float64 = math.Log(10) / 400.0 //cannot be declared as const as it is a runtime eval - function call

func gRd(rd float64) float64 {
	return 1.0 / math.Sqrt(1+(3*q*q*rd*rd)/(math.Pi*math.Pi))
}
func eScore(userRating, opponentRating, opponentRD float64) float64 {
	powTerm := (-gRd(opponentRD) * (userRating - opponentRating)) / 400
	return 1.0 / (1 + math.Pow(10, powTerm))
}
func dSquared(opponentRD, e float64) float64 {
	denominator := q * q * (math.Pow(gRd(opponentRD), 2)) * e * (1 - e)
	return 1.0 / denominator
}
func updateRating(rating, rd, opponentRating, opponentRD, score float64) (newRating, newRD float64) {
	e := eScore(rating, opponentRating, opponentRD)
	d := dSquared(opponentRD, e)
	denominator_term := 1.0/(rd*rd) + 1.0/(d*d)
	return rating + (q/(denominator_term))*gRd(opponentRD)*(score-e), math.Sqrt(1.0 / denominator_term)
}
func CalculateNewRatings(userRating, userRD, puzzleRating, puzzleRD float64, solved bool) (float64, float64, float64, float64) {
	score := 0.0
	if solved {
		score = 1.0
	}
	newUserRating, newUserRD := updateRating(userRating, userRD, puzzleRating, puzzleRD, score)
	newPuzzleRating, newPuzzleRD := updateRating(puzzleRating, puzzleRD, userRating, userRD, 1-score)
	return newUserRating, newUserRD, newPuzzleRating, newPuzzleRD
}

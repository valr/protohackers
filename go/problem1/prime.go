package problem1

import "math/big"

const (
	numberOfTest = 20
)

func IsPrime(number float64) (isPrime bool) {
	bigFloat := big.NewFloat(number)
	bigInt, accuracy := bigFloat.Int(nil)
	if accuracy == big.Exact {
		isPrime = bigInt.ProbablyPrime(numberOfTest)
	}
	return
}

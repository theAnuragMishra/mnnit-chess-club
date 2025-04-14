package utils

import (
	"math"
)

const (
	tau     = 0.5
	epsilon = 0.000001
	scale   = 173.7178
)

type Player struct {
	Rating     float64
	RD         float64
	Volatility float64
}

func UpdateMatch(playerA, playerB Player, resultA float64) (Player, Player) {
	// Convert to Glicko-2 scale
	muA := (playerA.Rating - 1500) / scale
	phiA := playerA.RD / scale

	muB := (playerB.Rating - 1500) / scale
	phiB := playerB.RD / scale

	g := func(phi float64) float64 {
		return 1 / math.Sqrt(1+3*phi*phi/math.Pi/math.Pi)
	}
	E := func(mu1, mu2, phi2 float64) float64 {
		return 1 / (1 + math.Exp(-g(phi2)*(mu1-mu2)))
	}

	gB := g(phiB)
	EAB := E(muA, muB, phiB)
	vA := 1 / (gB * gB * EAB * (1 - EAB))
	deltaA := vA * gB * (resultA - EAB)

	sigmaPrimeA := updateVolatility(deltaA, phiA, vA, playerA.Volatility)
	phiStarA := math.Sqrt(phiA*phiA + sigmaPrimeA*sigmaPrimeA)
	phiPrimeA := 1 / math.Sqrt(1/(phiStarA*phiStarA)+1/vA)
	muPrimeA := muA + phiPrimeA*phiPrimeA*gB*(resultA-EAB)

	// Do same for player B with reversed result
	resultB := 1 - resultA
	gA := g(phiA)
	EBA := E(muB, muA, phiA)
	vB := 1 / (gA * gA * EBA * (1 - EBA))
	deltaB := vB * gA * (resultB - EBA)

	sigmaPrimeB := updateVolatility(deltaB, phiB, vB, playerB.Volatility)
	phiStarB := math.Sqrt(phiB*phiB + sigmaPrimeB*sigmaPrimeB)
	phiPrimeB := 1 / math.Sqrt(1/(phiStarB*phiStarB)+1/vB)
	muPrimeB := muB + phiPrimeB*phiPrimeB*gA*(resultB-EBA)

	// Convert back to Glicko scale
	playerA.Rating = muPrimeA*scale + 1500
	playerA.RD = phiPrimeA * scale
	playerA.Volatility = sigmaPrimeA

	playerB.Rating = muPrimeB*scale + 1500
	playerB.RD = phiPrimeB * scale
	playerB.Volatility = sigmaPrimeB

	return playerA, playerB
}

func updateVolatility(delta, phi, v, sigma float64) float64 {
	a := math.Log(sigma * sigma)
	A := a
	var B float64
	if delta*delta > phi*phi+v {
		B = math.Log(delta*delta - phi*phi - v)
	} else {
		k := 1.0
		for f(a-k*tau, delta, phi, v, a) < 0 {
			k++
		}
		B = a - k*tau
	}

	fA := f(A, delta, phi, v, a)
	fB := f(B, delta, phi, v, a)

	for math.Abs(B-A) > epsilon {
		C := A + (A-B)*fA/(fB-fA)
		fC := f(C, delta, phi, v, a)
		if fC*fB < 0 {
			A = B
			fA = fB
		} else {
			fA /= 2
		}
		B = C
		fB = fC
	}
	return math.Exp(A / 2)
}

func f(x, delta, phi, v, a float64) float64 {
	expX := math.Exp(x)
	num := expX * (delta*delta - phi*phi - v - expX)
	den := 2 * (phi*phi + v + expX) * (phi*phi + v + expX)
	return num/den - (x-a)/(tau*tau)
}

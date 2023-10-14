package main

import (
	"math"
	"math/big"
)

type params struct {
	network   int
	faulty    int
	honest    int
	endorsers int
	quorum    float64
}

// Probability of bad nodes pushing through a bad block
func GetFailureProbability(p params) float64 {
	// Safety breaks if there are enough bad nodes to form a
	// quorum, i.e. the number of good nodes is (1-q)*E
	failurePoint := int(math.Round((1 - p.quorum) * float64(p.endorsers)))
	return hypergeoCDF(p, failurePoint)
}

// Probability of having liveness failure in the endorser set
func GetLivenessProbability(p params) float64 {
	// Liveness breaks if there aren't enough good nodes to form a
	// quorum, i.e. the number of good nodes is q*E - 1
	failurePoint := int(math.Round(p.quorum*float64(p.endorsers)) - 1)
	return hypergeoCDF(p, failurePoint)
}

// Find the minimum endorser quorum for a given endorser set
// that satisfies the given failure probability
func FindQ(p params, start, epsilon float64) (int, float64) {
	for q := start; q <= 1; q += 0.01 {
		p.quorum = q
		prob := GetFailureProbability(p)
		if prob <= epsilon {
			return int(q * 100), prob
		}
	}
	return 101, 0.0
}

// Find the minimum endorser set size for a given quorum percentage
// that satisfies the given failure probability
func FindE(p params, start, step int, epsilon float64) (int, float64) {
	for e := start; e <= p.network; e += step {
		p.endorsers = e
		prob := GetFailureProbability(p)
		if prob <= epsilon {
			return e, prob
		}
	}
	return p.network, 0.0
}

// Compute the Hypergeometric Cumulative Distribution Function
func hypergeoCDF(p params, x int) float64 {
	sum := big.NewInt(0)

	// N choose totalEndorsers
	subsamples := big.NewInt(1)
	subsamples.Binomial(int64(p.network), int64(p.endorsers))
	// Faulty choose totalEndorsers - i
	faultyComb := big.NewInt(1)
	// Honest choose i
	honestComb := big.NewInt(1)

	// Probability that bad nodes can convince half of good nodes to achieve
	// at least the required minimum number of signatures in order to approve a block.
	for i := 0; i <= x; i++ {
		honestComb.Binomial(int64(p.honest), int64(i))
		faultyComb.Binomial(int64(p.faulty), int64(p.endorsers-i))
		honestComb.Mul(honestComb, faultyComb)
		sum.Add(sum, honestComb)
	}

	probability := big.NewRat(1, 1)
	probability.SetFrac(sum, subsamples)
	f, _ := probability.Float64()
	return f
}

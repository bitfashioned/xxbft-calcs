package main

import (
	"fmt"
	"math"
)

type result struct {
	T int
	F float64
	L float64
}

func main() {
	// Ask user to input the parameters
	var p params
	fmt.Println("Enter the network size")
	fmt.Scanf("%d", &p.network)

	fmt.Println("Enter the percentage of honest nodes")
	var h int
	fmt.Scanf("%d", &h)
	// Compute the number of honest nodes
	p.honest = int(math.Ceil(float64(p.network) * float64(h) / 100))
	// Compute the number of faulty nodes
	p.faulty = p.network - p.honest

	fmt.Printf("Nodes: %d\n", p.network)
	fmt.Printf("GOOD: %d\n", p.honest)
	fmt.Printf("BAD: %d\n\n", p.faulty)

	// Choose the caclulation to perform
	fmt.Println("Select the calculation to perform")
	fmt.Println("  1 -> Failure probabilities for multiple quorums for a given endorser set size")
	fmt.Println("  2 -> Needed quorum to achieve a given failure probability for a given endorser set size")
	fmt.Println("  3 -> Needed endorser set size to achieve a given failure probability for a given quorum")
	var choice int
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		// Ask user to input the endorser set size
		fmt.Println("Enter the endorser set size")
		fmt.Scanf("%d", &p.endorsers)

		// Ask user to input the lowest quorum percentage
		fmt.Println("Enter the lowest quorum percentage")
		var low int
		fmt.Scanf("%d", &low)

		// Ask user to input the highest quorum percentage
		fmt.Println("Enter the highest quorum percentage")
		var high int
		fmt.Scanf("%d", &high)

		// Compute the failure probability for each quorum percentage
		results := make([]result, high-low+1)
		for q := low; q <= high; q++ {
			p.quorum = float64(q) / 100
			ps := GetFailureProbability(p)
			pl := GetLivenessProbability(p)
			results[q-low] = result{T: q, F: ps, L: pl}
		}
		fmt.Println()
		fmt.Println()
		for _, val := range results {
			fmt.Printf("Quorum percentage %d%%: p(safety) = %e, p(liveness)= %e\n", val.T, val.F, val.L)
		}
	case 2:
		// Ask user to input the endorser set size
		fmt.Println("Enter the endorser set size")
		fmt.Scanf("%d", &p.endorsers)

		// Ask user to input the failure probability
		fmt.Println("Enter the failure probability")
		var epsilon float64
		fmt.Scanf("%e", &epsilon)

		// Ask user to input the lowest quorum percentage
		fmt.Println("Enter the lowest quorum percentage")
		var start int
		fmt.Scanf("%d", &start)

		// Compute the needed quorum percentage
		q, prob := FindQ(p, float64(start)/100, epsilon)
		// Compute the liveness probability
		p.quorum = float64(q) / 100
		pl := GetLivenessProbability(p)
		fmt.Printf("\n\nQuorum percentage %d%%: p(safety) = %e, p(liveness)= %e\n", q, prob, pl)
	case 3:
		// Ask user to input the failure probability
		fmt.Println("Enter the failure probability")
		var epsilon float64
		fmt.Scanf("%e", &epsilon)

		// Ask user to input the quorum percentage
		fmt.Println("Enter the quorum percentage")
		var q int
		fmt.Scanf("%d", &q)
		p.quorum = float64(q) / 100

		// Ask user to input the lowest endorser set size
		fmt.Println("Enter the lowest endorser set size")
		var start int
		fmt.Scanf("%d", &start)

		// Ask user to input the endorser set size step increment
		fmt.Println("Enter the endorser set size step increment")
		var step int
		fmt.Scanf("%d", &step)

		// Compute the needed endorser set size
		e, prob := FindE(p, start, step, epsilon)
		// Compute the liveness probability
		p.endorsers = e
		pl := GetLivenessProbability(p)
		fmt.Printf("\n\nEndorser set size %d: p(safety) = %e, p(liveness)= %e\n", e, prob, pl)
	default:
		fmt.Println("Invalid choice")
	}
}

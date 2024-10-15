package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type result struct {
	T int
	F float64
	L float64
}

// Helper function to format numbers with 'k' or 'm'
func formatNumber(value float64) string {
	value = math.Round(value*10000) / 10000 // Round to 4 decimal places
	if value >= 1000000000000 {
		if value/1000000000000 < 10 {
			return fmt.Sprintf("%4.0fT", value/1000000000000) // Pad with zeros
		} else {
			return "> 10T"
		}
	} else if value >= 1000000000 {
		return fmt.Sprintf("%4.0fB", value/1000000000) // Pad with zeros
	} else if value >= 1000000 {
		return fmt.Sprintf("%4.0fM", value/1000000) // Pad with zeros
	} else if value >= 1000 {
		return fmt.Sprintf("%4.0fK", value/1000) // Pad with zeros
	} else if value < 1 {
		return "< 01"
	}
	return fmt.Sprintf("%4.0f", value) // Pad with zeros
}

// Helper function to read input with a default value using Scanf and return an integer or float64
func readInput(prompt string, defaultValue string) interface{} {
	fmt.Printf("%s (default: %s): ", prompt, defaultValue)
	var input string
	_, err := fmt.Scanf("%s\n", &input)
	if err != nil || input == "" {
		input = defaultValue
	}
	// Check if input contains "e" (scientific notation)
	if strings.Contains(input, "e") {
		value, _ := strconv.ParseFloat(input, 64)
		return value
	}
	value, _ := strconv.Atoi(input)
	return value
}

func readInputFloat(prompt string, defaultValue string) float64 {
	return readInput(prompt, defaultValue).(float64)
}

func readInputInt(prompt string, defaultValue string) int {
	return readInput(prompt, defaultValue).(int)
}

func main() {
	var p params

	// Default values
	defaultNetwork := "1000"
	defaultHonest := "67"
	defaultBlockTime := "1"
	defaultChoice := "1"
	defaultEndorsers := "150"
	defaultQuorum := "80"
	defaultLowQuorum := "50"
	defaultHighQuorum := "80"
	defaultEpsilon := "1e-14"
	defaultLowEndorsers := "100"
	defaultStepSizeEndorsers := "10"
	defaultBits := "1"

	// Ask user to input the parameters with default values
	p.network = readInputInt("Enter the network size", defaultNetwork)
	h := readInputInt("Enter the percentage of honest nodes", defaultHonest)
	p.honest = int(math.Ceil(float64(p.network) * float64(h) / 100))
	// Compute the number of faulty nodes
	p.faulty = p.network - p.honest

	// global var with number of seconds in a year
	secondsInYear := float64(31536000)
	// global var with number of seconds in a day
	secondsInDay := float64(86400)

	// Enter estimated block time
	t := readInputInt("Enter the estimated block time in seconds", defaultBlockTime)

	// Choose the caclulation to perform
	fmt.Println("Possible types of calculations:")
	fmt.Println("  1 -> Failure probabilities for multiple quorums for a given endorser set size")
	fmt.Println("  2 -> Needed quorum to achieve a given failure probability for a given endorser set size")
	fmt.Println("  3 -> Needed endorser set size to achieve a given failure probability for a given quorum")
	fmt.Println("  4 -> Failure probabilities for multiple quorums for a given endorser set size with biasable randomness")
	fmt.Println()
	choice := readInputInt("Select the calculation to perform from the above list", defaultChoice)
	fmt.Println()

	fmt.Println("Chosen parameters:")
	fmt.Printf("=> Total Nodes: %d\n", p.network)
	fmt.Printf("=> # Good Nodes: %d\n", p.honest)
	fmt.Printf("=> # Bad Nodes: %d\n", p.faulty)
	fmt.Printf("=> Block time: %d seconds\n", t)
	fmt.Printf("=> Calculation Type: %d\n", choice)
	fmt.Println()

	switch choice {
	case 1:
		// Ask user to input the endorser set size
		p.endorsers = readInputInt("Enter the endorser set size", defaultEndorsers)
		// Ask user to input the lowest quorum percentage
		low := readInputInt("Enter the lowest quorum percentage", defaultLowQuorum)
		// Ask user to input the highest quorum percentage
		high := readInputInt("Enter the highest quorum percentage", defaultHighQuorum)

		// Compute the failure probability for each quorum percentage
		results := make([]result, high-low+1)
		for q := low; q <= high; q++ {
			p.quorum = float64(q) / 100
			ps := GetFailureProbability(p)
			pl := GetLivenessProbability(p)
			results[q-low] = result{T: q, F: ps, L: pl}
		}

		// Print the calculation specific parameters
		fmt.Println()
		fmt.Println("Calculation specific parameters:")
		fmt.Printf("=> Quorum percentage: %d%% - %d%%\n", low, high)
		fmt.Printf("=> Endorser set size: %d\n", p.endorsers)
		fmt.Println()

		for _, val := range results {
			// Calculate the values
			safetyYears := float64(t) / (secondsInYear * val.F)
			livenessDays := float64(t) / (secondsInDay * val.L)

			// Use the helper function to format the values
			safetyYearsStr := formatNumber(safetyYears)
			livenessDaysStr := formatNumber(livenessDays)

			// Print the formatted values with percentages
			fmt.Printf("Quorum percentage %d%%: p(safety) = %e (%3.0f%%) (%s years),\t p(liveness)= %e (%3.0f%%) (%s days),\t finality = %.2f seconds\n", val.T, val.F, val.F*100, safetyYearsStr, val.L, val.L*100, livenessDaysStr, float64(t))
		}
	case 2:
		// Ask user to input the endorser set size
		p.endorsers = readInputInt("Enter the endorser set size", defaultEndorsers)
		// Ask user to input the failure probability
		epsilon := readInputFloat("Enter the failure probability", defaultEpsilon)
		// Ask user to input the lowest quorum percentage
		start := readInputInt("Enter the lowest quorum percentage", defaultLowQuorum)

		// Print the calculation specific parameters
		fmt.Println()
		fmt.Println("Calculation specific parameters:")
		fmt.Printf("=> Endorser set size: %d\n", p.endorsers)
		fmt.Printf("=> Failure probability: %e\n", epsilon)
		fmt.Printf("=> Low Quorum percentage: %d%%\n", start)
		fmt.Println()

		// Compute the needed quorum percentage
		q, prob := FindQ(p, float64(start)/100, epsilon)
		// Compute the liveness probability
		p.quorum = float64(q) / 100
		pl := GetLivenessProbability(p)
		safetyYears := float64(t) / (secondsInYear * prob)
		livenessDays := float64(t) / (secondsInDay * pl)
		safetyYearsStr := formatNumber(safetyYears)
		livenessDaysStr := formatNumber(livenessDays)
		fmt.Printf("\n\nQuorum percentage %d%%: p(safety) = %e (%3.0f%%) (%s years), p(liveness)= %e (%3.0f%%) (%s days)\n", q, prob, prob*100, safetyYearsStr, pl, pl*100, livenessDaysStr)
	case 3:
		// Ask user to input the failure probability
		epsilon := readInputFloat("Enter the failure probability", defaultEpsilon)
		// Ask user to input the quorum percentage
		q := readInputInt("Enter the quorum percentage", defaultQuorum)
		p.quorum = float64(q) / 100

		// Ask user to input the lowest endorser set size
		start := readInputInt("Enter the lowest endorser set size", defaultLowEndorsers)

		// Ask user to input the endorser set size step increment
		step := readInputInt("Enter the endorser set size step increment", defaultStepSizeEndorsers)

		// Print the calculation specific parameters
		fmt.Println()
		fmt.Println("Calculation specific parameters:")
		fmt.Printf("=> Failure probability: %e\n", epsilon)
		fmt.Printf("=> Quorum percentage: %d%%\n", q)
		fmt.Printf("=> Lowest Endorser set size: %d\n", start)
		fmt.Printf("=> Endorser set size step increment: %d\n\n", step)
		fmt.Println()

		// Compute the needed endorser set size
		e, prob := FindE(p, start, step, float64(epsilon))
		// Compute the liveness probability
		p.endorsers = e
		pl := GetLivenessProbability(p)
		safetyYears := float64(t) / (secondsInYear * prob)
		livenessDays := float64(t) / (secondsInDay * pl)
		safetyYearsStr := formatNumber(safetyYears)
		livenessDaysStr := formatNumber(livenessDays)
		fmt.Printf("\n\nEndorser set size %d: p(safety) = %e (%3.0f%%) (%s years), p(liveness)= %e (%3.0f%%) (%s days)\n", e, prob, prob*100, safetyYearsStr, pl, pl*100, livenessDaysStr)
	case 4:
		// Ask user to input the endorser set size
		p.endorsers = readInputInt("Enter the endorser set size", defaultEndorsers)

		// Ask user to input the lowest quorum percentage
		low := readInputInt("Enter the lowest quorum percentage", defaultLowQuorum)

		// Ask user to input the highest quorum percentage
		high := readInputInt("Enter the highest quorum percentage", defaultHighQuorum)

		// Ask user to input he number of biasable bits of randomness
		bits := readInputInt("Enter the number of biasable bits of randomness", defaultBits)

		// Print the calculation specific parameters
		fmt.Println()
		fmt.Println("Calculation specific parameters:")
		fmt.Printf("=> Endorser set size: %d\n", p.endorsers)
		fmt.Printf("=> Quorum percentage: %d%% - %d%%\n", low, high)
		fmt.Printf("=> Number of biasable bits of randomness: %d\n\n", bits)
		fmt.Println()

		// Compute the failure probability for each quorum percentage
		results := make([]result, high-low+1)
		for q := low; q <= high; q++ {
			p.quorum = float64(q) / 100
			ps := GetFailureProbability(p)
			ps = 1 - math.Pow(1-ps, math.Pow(2, float64(bits)))
			pl := GetLivenessProbability(p)
			pl = 1 - math.Pow(1-pl, math.Pow(2, float64(bits)))
			results[q-low] = result{T: q, F: ps, L: pl}
		}

		for _, val := range results {
			safetyYears := float64(t) / (secondsInYear * val.F)
			livenessDays := float64(t) / (secondsInDay * val.L)
			safetyYearsStr := formatNumber(safetyYears)
			livenessDaysStr := formatNumber(livenessDays)
			fmt.Printf("Quorum percentage %d%%: p(safety) = %e (%3.0f%%) (%s years), p(liveness)= %e (%3.0f%%) (%s days)\n", val.T, val.F, val.F*100, safetyYearsStr, val.L, val.L*100, livenessDaysStr)
		}
	default:
		fmt.Println("Invalid choice")
	}
}

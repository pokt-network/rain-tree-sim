package main

import (
	"fmt"
	"log"
	"math"
)

type Results struct {
	NumberOfNodes              uint64
	Levels                     uint
	AverageRedundancy          float64
	NonDeadCoveragePercentage  float64
	DeadCount                  uint64
	ConsecutiveLevelZeroMatrix map[int]int
	Communications             uint64
	IndividualNodeData         AddressBook
}

type ResultsCSV [][]string

func NewResultsCSV(data []Results) (rcsv ResultsCSV) {
	rcsv = make([][]string, 0)
	dataLen := len(data)
	if dataLen < 1 {
		return
	}
	titles := []string{"Nodes", "Levels", "Comms", "Redundancy", "Coverage", "Missed", "LongestMiss"}
	rcsv = append(rcsv, titles)
	for _, res := range data {
		missed, longestMiss := 0, 0
		for miss, occur := range res.ConsecutiveLevelZeroMatrix {
			missed += miss * occur
			if longestMiss < miss {
				longestMiss = miss
			}
		}
		rcsv = append(rcsv, []string{
			fmt.Sprintf("%d", res.NumberOfNodes), fmt.Sprintf("%d", res.Levels),
			fmt.Sprintf("%d", res.Communications), fmt.Sprintf("%f", res.AverageRedundancy),
			fmt.Sprintf("%d", int(res.NonDeadCoveragePercentage)), fmt.Sprintf("%d", missed), fmt.Sprintf("%d", longestMiss)})
	}
	return
}

func GatherData(c *Config, globalAddressBook AddressBook) (r Results) {
	log.Println("Gathering network level data")
	r = Results{
		ConsecutiveLevelZeroMatrix: make(map[int]int),
	}
	fullListSize := float64(len(globalAddressBook))
	r.Levels = uint(math.Ceil(math.Round(logBase3(fullListSize)*100) / 100))
	r.NumberOfNodes = c.NumberOfNodes
	// show individual node data?
	if c.ShowIndividualNodeResults {
		r.IndividualNodeData = globalAddressBook
	}
	totalRedundancy, totalNonDeadMiss, currentConsecutiveZeroCount := float64(0), float64(0), 0
	// calculate network wide stats using the accumulation of the individual nodes
	for _, node := range globalAddressBook {
		// if no message received, increment current consecutive zero count
		if node.MessagesReceived == 0 {
			currentConsecutiveZeroCount++
		}
		// if node is dead, we don't track data any further
		if node.IsDead {
			r.DeadCount++
			continue
		}
		// calculate non-dead misses
		if node.MessagesReceived == 0 {
			totalNonDeadMiss += 1
		} else {
			// track total redundancy; 1 redundancy is a minimum at this point
			totalRedundancy += float64(node.MessagesReceived)
			// if we've been tracking consecutive zeroes
			if currentConsecutiveZeroCount != 0 {
				// break the consecutive zeroes & track
				r.ConsecutiveLevelZeroMatrix[currentConsecutiveZeroCount] += 1
				currentConsecutiveZeroCount = 0
			}
		}
	}
	if currentConsecutiveZeroCount != 0 {
		// break the consecutive zeroes & track
		r.ConsecutiveLevelZeroMatrix[currentConsecutiveZeroCount] += 1
		currentConsecutiveZeroCount = 0
	}
	// calculate average redundancy
	r.AverageRedundancy = totalRedundancy / float64(c.NumberOfNodes-r.DeadCount-uint64(totalNonDeadMiss))
	r.NonDeadCoveragePercentage = math.Round((1.0 - totalNonDeadMiss/float64(c.NumberOfNodes-r.DeadCount)) * 100)
	r.Communications = uint64(totalRedundancy - 1)
	return
}

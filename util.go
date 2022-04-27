package main

import (
	cryptRand "crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"

	"gonum.org/v1/gonum/stat/distuv"
)

func getPartialViewershipCurve(c *Config) []int {
	log.Println("Calculating the 'curve' for the partial viewership distribution")
	dist := distuv.UnitNormal
	dist.Mu = float64(c.TargetPartialViewershipPercentage)
	dist.Sigma = float64(c.PartialViewershipStdDev)
	curveArr := make([]int, c.NumberOfNodes)
	for i := range curveArr {
		if c.InvertCurve {
			curveArr[i] = int(math.Round(dist.Rand())+45) % 100
		} else {
			curveArr[i] = int(math.Round(dist.Rand()))
		}
		// set max & min for practicality
		// max = 100 & min = 10
		if curveArr[i] > 100 {
			curveArr[i] = 100
		} else if curveArr[i] <= 10 {
			curveArr[i] = 10
		}
	}
	return curveArr
}

func printResults(results []Results, c *Config) {
	log.Println("Printing results")
	resBytes, err := json.MarshalIndent(&results, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(c.ResultFileOutputName+".json", resBytes, 0777)
	if err != nil {
		log.Println("Error trying to print result resBytes: ", err.Error())
		log.Println(string(resBytes))
	}
	// write csv too
	f, err := os.Create(c.ResultFileOutputName + ".csv")
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range NewResultsCSV(results) {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}

func getIsDead(index uint64, c *Config) bool {
	randSeed()
	if c.FixedDeadNodes {
		for _, in := range c.FixedDeadNodesIndexArray {
			if index == in {
				return true
			}
		}
		return false
	}
	i := rand.Intn(99) + 1
	if uint8(i) <= c.DeadNodePercentage {
		return true
	}
	return false
}

func logBase3(x float64) float64 {
	return math.Log(x) / math.Log(3.0)
}

func randSeed() {
	seed, err := cryptRand.Int(cryptRand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	rand.Seed(seed.Int64())
}

// simulated address

type Address [20]byte

func (a *Address) String() string {
	return hex.EncodeToString(a[:])
}

func newAddress() string {
	var addr Address
	_, err := rand.Read(addr[:])
	if err != nil {
		_ = err
	}
	return addr.String()
}

// simulated hash

type Hash [32]byte

func newHash() string {
	randSeed()
	var hash Hash
	_, err := rand.Read(hash[:])
	if err != nil {
		_ = err
	}
	return hash.String()
}

func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}

package main

import (
	cryptRand "crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"

	"gonum.org/v1/gonum/stat/distuv"
)

const (
	// set max & min for practicality
	partialViewershipCurveMax = 100
	partialViewershipCurveMin = 10
)

func generatePartialViewershipCurve(c *Config) []int {
	log.Println("Calculating the 'curve' for the partial viewership distribution")
	dist := distuv.UnitNormal
	dist.Mu = float64(c.PartialViewershipMedian)
	dist.Sigma = float64(c.PartialViewershipStdDev)
	curveArr := make([]int, c.NumberOfNodes)
	for i := range curveArr {
		if c.InvertCurve {
			// Discuss: what kind of inversion is this?
			curveArr[i] = int(math.Round(dist.Rand())+45) % 100
		} else {
			curveArr[i] = int(math.Round(dist.Rand()))
		}
		// Cap at max
		if curveArr[i] > partialViewershipCurveMax {
			curveArr[i] = partialViewershipCurveMax
		}
		// Cap at min
		if curveArr[i] <= partialViewershipCurveMin {
			curveArr[i] = partialViewershipCurveMin
		}
	}
	return curveArr
}

func getIsDead(c *Config, index uint64) bool {
	// TODO: switch to `slices.Contains` in Go 1.18
	if c.FixedDeadNodes {
		for _, in := range c.FixedDeadNodesIndexArray {
			if index == in {
				return true
			}
		}
		return false
	}

	// If no configs were provided, randomize the liveness of the node based on DeadNodePercentage
	setRandSeed()
	i := rand.Intn(99) + 1 // TODO: Why not just do `return rand.Intn(100) <= c.DeadNodePercentage`?
	return uint8(i) <= c.DeadNodePercentage
}

func logBase3(x float64) float64 {
	return math.Log(x) / math.Log(3.0)
}

func setRandSeed() {
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

// simulated hash

type Hash [32]byte

func newHash() string {
	setRandSeed()
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

func PrintPartialViewershipCurveToFile(curve []int) {
	bz, _ := json.Marshal(curve)
	log.Println(string(bz))
}

func DumpPartialViewershipCurveToFile(curve []int) {
	f, err := os.Create(fmt.Sprintf("curves/partial_viewership_curve.%s.csv", randomFilePrefix()))
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	st := strings.Fields(strings.Trim(fmt.Sprint(curve), "[]"))
	w.Write(st)
}

func randomFilePrefix() string {
	randBytes := make([]byte, 4)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// TODO: Use struct tags so that json contains snake_case keys rather than CamelCase
// TODO: Rename these variables (e.g. bool vars appropriately)
type Config struct {
	NumNodesFirstSimulation uint64
	NumNodesLastSimulation  uint64

	DeadNodePercentage       uint8
	FixedDeadNodes           bool
	FixedDeadNodesIndexArray []uint64

	InvertCurve               bool
	FixedViewershipPercentage bool
	FixedViewershipCurveArray []int

	RandomizePartialAddressBooks bool
	PartialViewershipMedian      uint8
	PartialViewershipStdDev      int8

	RedundancyLayerRightOn bool
	RedundancyLayerLeftOn  bool

	OriginatorIndex int64
	MaxHotlist      uint

	ShowIndividualNodeSimResult           bool
	ShowIndividualNodePartialAddressBooks bool
	ResultFileOutputName                  string

	// NOT part of the config.json - used for implementation
	NumberOfNodes uint64
}

func LoadConfigFile() (c *Config) {
	config := Config{}
	jsonFile, err := os.Open("config.json")
	if err != nil {
		if err != nil {
			log.Println("Error opening the config file; ", err.Error())
			os.Exit(1)
		}
	}
	defer jsonFile.Close()

	bz, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("Error reading the config file; ", err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(bz, &config)
	if err != nil {
		log.Println("Error unmarshalling the config file; ", err.Error())
		os.Exit(1)
	}
	if err = config.hydrate(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	if ok, err := config.IsValid(); !ok {
		log.Println(err.Error())
		os.Exit(1)
	}
	return &config
}

func (c *Config) hydrate() error {
	if c.ResultFileOutputName == "" {
		c.ResultFileOutputName = "results.json"
	}
	return nil
}

func (c *Config) IsValid() (bool, error) {
	if c.NumNodesLastSimulation != 0 && c.NumNodesFirstSimulation > c.NumNodesLastSimulation {
		return false, errors.New("number of nodes must be greater than or equal to EndingNumberOfNodes; set EndingNumberOfNodes = to numberOfNodes to disable")
	}
	if c.DeadNodePercentage > 100 {
		return false, errors.New("dead node percentage can't be greater than 100%")
	}
	if c.FixedViewershipPercentage {
		if len(c.FixedViewershipCurveArray) != int(c.NumNodesFirstSimulation) {
			return false, errors.New("target viewership array length must equal number of nodes")
		}
	}
	return true, nil
}

func (c *Config) String() string {
	bz, _ := json.MarshalIndent(c, "", "    ")
	return string(bz)
}

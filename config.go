package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	NumberOfNodes                         uint64
	EndingNumberOfNodes                   uint64
	DeadNodePercentage                    uint8
	FixedDeadNodes                        bool
	FixedDeadNodesIndexArray              []uint64
	ViewershipPercentageFixed             bool
	ViewershipCurveArray                  []int
	TargetPartialViewershipPercentage     uint8
	PartialViewershipStdDev               int8
	InvertCurve                           bool
	RedundancyLayerRightOn                bool
	RedundancyLayerLeftOn                 bool
	MaxHotlist                            uint
	ShowIndividualNodeResults             bool
	ShowIndividualNodePartialAddressBooks bool
	ResultFileOutputName                  string
	OriginatorIndex                       int64
	RandomizePartialAddressBooks          bool
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
	if ok, err := config.IsValid(); !ok {
		log.Println(err.Error())
		os.Exit(1)
	}
	return &config
}

func (c *Config) IsValid() (bool, error) {
	if c.ResultFileOutputName == "" {
		c.ResultFileOutputName = "results.js¬¬on"
	}
	if c.EndingNumberOfNodes != 0 && c.NumberOfNodes > c.EndingNumberOfNodes {
		return false, errors.New("number of nodes must be greater than or equal to EndingNumberOfNodes; set EndingNumberOfNodes = to numberOfNodes to disable")
	}
	if c.DeadNodePercentage > 100 {
		return false, errors.New("dead node percentage can't be greater than 100")
	}
	if c.ViewershipPercentageFixed {
		if len(c.ViewershipCurveArray) != int(c.NumberOfNodes) {
			return false, errors.New("target viewership array length must equal number of nodes")
		}
	}
	return true, nil
}

func (c *Config) String() string {
	bz, _ := json.MarshalIndent(c, "", "    ")
	return string(bz)
}

package main

import (
	"log"
	"math/rand"
	"sort"
)

func NewRainTreeNetwork(c *Config) {
	log.Printf("Starting a rain tree network with the configuration: %s\n", c.String())
	results := make([]Results, 0)
	// setup a new rain tree network
	for i := c.NumberOfNodes; i <= c.EndingNumberOfNodes; i++ {
		conf := *c
		conf.NumberOfNodes = i
		globalAddressBook := PopulateAddressBook(&conf)
		SetupOriginator(&conf, globalAddressBook)
		log.Println("Running through global address book & executing message queue lexicographically")
		for {
			actionDone := RunQueue(&conf, globalAddressBook)
			if !actionDone {
				log.Println("No action done; All queues are exhausted")
				break
			}
		}
		results = append(results, GatherData(&conf, globalAddressBook))
	}
	printResults(results, c)
}

func SetupOriginator(c *Config, globalAddressBook AddressBook) {
	randSeed()
	maxIndex := int(c.NumberOfNodes - 1)
	origIndex := rand.Intn(maxIndex)
	if c.OriginatorIndex != -1 {
		origIndex = int(c.OriginatorIndex)
	}
	origNode := globalAddressBook[origIndex]
	log.Printf("Originator is %s at index %d out of %d\n", origNode.Address, origIndex, maxIndex)
	SendMessage(origNode.Address, NewMessage(origNode.PartialAddressBook), c, globalAddressBook, false)
}

func RunQueue(c *Config, globalAddressBook AddressBook) (actionDone bool) {
	// for every node in the global address book
	for i, n := range globalAddressBook {
		if n.Message == nil {
			continue
		}
		m := globalAddressBook[i].Message.Copy()
		globalAddressBook[i].Message = nil
		if m.Hash == "" {
			continue
		}
		// track if action is done
		if m.Level == -2 {
			continue
		}
		actionDone = true
		// handle each m
		PropagateMessage(n, m, globalAddressBook, c)
	}
	return
}

func PropagateMessage(node Node, message Message, globalAddressBook AddressBook, c *Config) {
	var isOriginator bool
	partialAddressBook, nodePosition := node.PartialAddressBook, node.PartialAddressBookPosition
	// get full list size (relative to sender)
	partialAddressBookSize := len(partialAddressBook)
	if message.Level == message.NetworkLevels {
		isOriginator = true
	}
	networkLevels := CalculateLevels(partialAddressBook)
	// calibrate levels (if message contains incorrect information)
	// levelWithDec, levelWithoutDec := CalibrateLevels(networkLevels, message)

	levelWithoutDec := int64(message.Level)
	levelWithDec := int64(message.Level - 1)

	message.Level = int(levelWithDec)
	message.NetworkLevels = int(networkLevels)
	// redundancy layer logic
	if message.Level <= -1 {
		// use partial address book size as `target list length` (definitional)
		targetAIndex, targetBIndex := GetTargetIndices(nodePosition, partialAddressBookSize, partialAddressBookSize)
		// get the addresses from the indices
		targetAAddress := partialAddressBook[targetAIndex].Address
		targetBAddress := partialAddressBook[targetBIndex].Address
		if c.RedundancyLayerLeftOn {
			SendMessage(targetAAddress, message, c, globalAddressBook, false)
		}
		if c.RedundancyLayerRightOn {
			SendMessage(targetBAddress, message, c, globalAddressBook, false)
		}
		SendMessage(node.Address, message, c, globalAddressBook, true)
		return
	} // if not level 0
	// calculate the target list length
	targetListLength := GetTargetListLength(partialAddressBookSize, int64(networkLevels), levelWithoutDec)
	// get the targets
	targetAIndex, targetBIndex := GetTargetIndices(nodePosition, targetListLength, partialAddressBookSize)
	log.Printf("Position: %d, TargetA: %d, TargetB: %d, TargetListLen: %d, NetworkLevels: %d\n", nodePosition, targetAIndex, targetBIndex, targetListLength, networkLevels)
	// originator logic (acks & hotlist)
	if isOriginator {
		targetAAddress, targetBAddress := GetRecursiveTargets(c, targetAIndex, targetBIndex, partialAddressBookSize, partialAddressBook)
		SendMessage(targetAAddress, message, c, globalAddressBook, false)
		SendMessage(targetBAddress, message, c, globalAddressBook, false)
	} else { // every other non-originator cases
		targetAAddress := partialAddressBook[targetAIndex].Address
		targetBAddress := partialAddressBook[targetBIndex].Address
		SendMessage(targetAAddress, message, c, globalAddressBook, false)
		SendMessage(targetBAddress, message, c, globalAddressBook, false)
	}
	// self queue to move down with everyone else
	SendMessage(node.Address, message, c, globalAddressBook, true)
	return
}

func SendMessage(nodeID string, message Message, c *Config, globalAddressBook AddressBook, isSelfMessage bool) {
	// no identity; no action
	if nodeID == "" {
		return
	}
	// search for the node index in the global list
	index := sort.Search(int(c.NumberOfNodes), func(i int) bool {
		return globalAddressBook[i].Address >= nodeID
	})
	globalAddressBook[index] = globalAddressBook[index].WithMessage(message, isSelfMessage)
}

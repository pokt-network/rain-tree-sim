package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

// TODO: A better design would be to create a Network struct and have `globalAddressBook` and `Config`
// attributes, so we don't have don't have to pass them around to all the functions here.

func NewRainTreeNetwork(c *Config) {
	results := make(SimResults, 0)
	for numNodesToSimulate := c.NumNodesFirstSimulation; numNodesToSimulate <= c.NumNodesLastSimulation; numNodesToSimulate++ {
		conf := *c // copy config
		conf.NumberOfNodes = numNodesToSimulate

		globalAddressBook := GenerateGlobalAddressBook(&conf)
		InitializeAllNodesInGlobalAddressBook(&conf, globalAddressBook)
		origNode := getOriginatorNode(&conf, globalAddressBook)

		// Create initial message from the originator node to distribute through the network
		msg := NewMessage(origNode.PartialAddressBook)
		sendMessage(&conf, origNode.Address, msg, globalAddressBook, false)
		log.Println("Running through global address book & executing message queue lexicographically")
		for {
			if actionDone := runQueue(&conf, globalAddressBook); !actionDone {
				log.Println("No actions were performed at all because all queues were empty.")
				break
			}
		}
		results = append(results, GatherData(&conf, globalAddressBook))
	}
	results.Print(c)
}

func getOriginatorNode(c *Config, globalAddressBook AddressBook) *Node {
	maxIndex := int(c.NumberOfNodes - 1)

	// Either randomly choose the originator index or use the one specified in the configs
	var origIndex int
	if c.OriginatorIndex == -1 {
		setRandSeed()
		origIndex = rand.Intn(maxIndex)
	} else {
		origIndex = int(c.OriginatorIndex)
	}
	origNode := globalAddressBook[origIndex]
	log.Printf("Originator node is %s at index %d out of %d nodes\n", origNode.Address, origIndex, maxIndex)
	return origNode
}

func runQueue(c *Config, globalAddressBook AddressBook) (actionDone bool) {
	for i, node := range globalAddressBook {
		if node.Message == nil {
			continue
		}
		// TODO: Because we pass by value below, we don't need to do a copy here
		m := globalAddressBook[i].Message.Copy()
		globalAddressBook[i].Message = nil
		fmt.Println(i, m.Level, m.Hash)
		// TODO: What does an empty hash represent?
		if m.Hash == "" {
			continue
		}
		// TODO: What does this magic number represent?
		if m.Level == -2 {
			continue
		}
		// Track if action is done
		// TODO: Don't fully understand this part
		actionDone = true
		// propagate message to each node in the global address book
		propagateMessage(c, node, m, globalAddressBook)
	}
	return
}

func propagateMessage(c *Config, node *Node, message Message, globalAddressBook AddressBook) {
	isOriginator := message.Level == message.TotalNumNetworkLevels
	partialAddressBook, nodePosition := node.PartialAddressBook, node.PartialAddressBookPosition
	// Get full list size (relative to sender)
	partialAddressBookSize := len(partialAddressBook)
	numNetworkLevels := CalculateNumLevelsInAddrBook(partialAddressBook)
	// Calibrate levels (if message contains incorrect information)
	// levelWithDec, levelWithoutDec := CalibrateLevels(numNetworkLevels, message)
	levelWithoutDec := int64(message.Level)
	levelWithDec := int64(message.Level - 1)
	message.Level = int(levelWithDec)
	message.TotalNumNetworkLevels = int(numNetworkLevels)
	// TODO(olshansky): Haven't explored redundancy yet
	// Redundancy layer logic
	if message.Level <= -1 {
		// use partial address book size as `target list length` (definitional)
		targetAIndex, targetBIndex := GetTargetIndices(nodePosition, partialAddressBookSize, partialAddressBookSize)
		// get the addresses from the indices
		targetAAddress := partialAddressBook[targetAIndex].Address
		targetBAddress := partialAddressBook[targetBIndex].Address
		if c.RedundancyLayerLeftOn {
			sendMessage(c, targetAAddress, message, globalAddressBook, false)
		}
		if c.RedundancyLayerRightOn {
			sendMessage(c, targetBAddress, message, globalAddressBook, false)
		}
		sendMessage(c, node.Address, message, globalAddressBook, true)
		return
	}
	// calculate the target list length
	// TODO: Does this not result in an endless loop if you look at the last line of this function?
	targetListLength := GetTargetListLength(partialAddressBookSize, int64(numNetworkLevels), levelWithoutDec)
	// get the targets
	targetAIndex, targetBIndex := GetTargetIndices(nodePosition, targetListLength, partialAddressBookSize)
	log.Printf("Position: %d, TargetA: %d, TargetB: %d, TargetListLen: %d, TotalNumNetworkLevels: %d\n", nodePosition, targetAIndex, targetBIndex, targetListLength, numNetworkLevels)
	// originator logic (acks & hotlist)
	if isOriginator {
		targetAAddress, targetBAddress := GetRecursiveTargets(c, targetAIndex, targetBIndex, partialAddressBookSize, partialAddressBook)
		sendMessage(c, targetAAddress, message, globalAddressBook, false)
		sendMessage(c, targetBAddress, message, globalAddressBook, false)
	} else { // every other non-originator cases
		targetAAddress := partialAddressBook[targetAIndex].Address
		targetBAddress := partialAddressBook[targetBIndex].Address
		sendMessage(c, targetAAddress, message, globalAddressBook, false)
		sendMessage(c, targetBAddress, message, globalAddressBook, false)

	}
	// self queue to move down with everyone else
	sendMessage(c, node.Address, message, globalAddressBook, true)
}

func sendMessage(c *Config, nodeAddr string, message Message, globalAddressBook AddressBook, isSelfMessage bool) {
	// no identity => no action
	if nodeAddr == "" {
		return
	}
	// search for the node index in the global list
	index := sort.Search(int(c.NumberOfNodes), func(i int) bool {
		return globalAddressBook[i].Address >= nodeAddr
	})
	nodeWithMessage := globalAddressBook[index].CopyNodeWithMessage(message, isSelfMessage)
	globalAddressBook[index] = &nodeWithMessage
}

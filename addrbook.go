package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"sort"
)

type AddressBook []*Node

type ExportableAddressBook []NodeExportable

type ExportedAddressBook []NodeExported

// TODO: Should this return a pointer?
func GenerateGlobalAddressBook(c *Config) (globalAddrBook AddressBook) {
	// 1. Create a list of random addresses
	globalAddrBook = make(AddressBook, c.NumberOfNodes)
	for i := uint64(0); i < c.NumberOfNodes; i++ {
		globalAddrBook[i] = CreateNode()
	}
	// 2. Sort addresses lexicographically
	sort.Slice(globalAddrBook, func(i int, j int) bool {
		return globalAddrBook[i].Address < globalAddrBook[j].Address
	})
	// 3. Make certain peers / addresses unreachable based on provided configs
	for i := uint64(0); i < c.NumberOfNodes; i++ {
		globalAddrBook[i].IsDead = getIsDead(c, i)
	}
	return
}

func InitializeAllNodesInGlobalAddressBook(c *Config, globalAddrBook AddressBook) {
	partialViewershipPercentages := generatePartialViewershipPercentages(c)
	for i := uint64(0); i < c.NumberOfNodes; i++ {
		globalAddrBookCopy := make(AddressBook, len(globalAddrBook))
		copy(globalAddrBookCopy, globalAddrBook)
		partialAddressBook := getPartialExportableAddressBook(c, globalAddrBook[i], globalAddrBookCopy, float64(partialViewershipPercentages[i]))
		globalAddrBook[i].InitNode(i, partialAddressBook, uint8(partialViewershipPercentages[i]), c.ShowIndividualNodePartialAddressBooks)
	}
}

func GetNodePosition(nodeAddr string, a ExportableAddressBook) int {
	return sort.Search(len(a), func(i int) bool {
		return a[i].Address >= nodeAddr
	})
}

// TODO: Still no clue where Base3 comes from since we're not actually building a ternary tree
func CalculateNumLevelsInAddrBook(a ExportableAddressBook) uint {
	fullListSize := float64(len(a))
	return uint(math.Ceil(math.Round(logBase3(fullListSize)*100) / 100))
}

func generatePartialViewershipPercentages(c *Config) []int {
	// Generate a random partial viewership curve or use the fix one
	var partialViewershipPercentages []int
	if !c.FixedViewershipPercentage {
		partialViewershipPercentages = generatePartialViewershipCurve(c)
	} else {
		partialViewershipPercentages = c.FixedViewershipCurveArray
	}

	// Potentially shuffle the partial viewership curve
	if c.RandomizePartialAddressBooks {
		setRandSeed()
		rand.Shuffle(len(partialViewershipPercentages), func(i, j int) {
			partialViewershipPercentages[i], partialViewershipPercentages[j] =
				partialViewershipPercentages[j], partialViewershipPercentages[i]
		})
	}

	DumpPartialViewershipCurveToFile(partialViewershipPercentages)
	PrintPartialViewershipCurveToFile(partialViewershipPercentages)

	return partialViewershipPercentages
}

// Truncate the node's address book based on its partial visibility percentage of the network
func getPartialExportableAddressBook(c *Config, selfNode *Node, globalAddrBook AddressBook, partialViewershipPercentage float64) ExportableAddressBook {
	// shuffle the original address book; this affect the original pointer - address book is a slice and passed by reference
	if c.RandomizePartialAddressBooks {
		setRandSeed()
		rand.Shuffle(len(globalAddrBook), func(i, j int) { globalAddrBook[i], globalAddrBook[j] = globalAddrBook[j], globalAddrBook[i] })
	}

	// truncate the address book based on partialViewershipPercentage
	maxIndex := int(math.Trunc(float64(partialViewershipPercentage) / 100 * float64(len(globalAddrBook))))
	newAddrBook := globalAddrBook[:maxIndex]
	newAddrBook = ensureSelfInAddrBook(selfNode, newAddrBook)

	// sort the new address book lexicographically
	sort.Slice(newAddrBook, func(i int, j int) bool {
		return newAddrBook[i].Address < newAddrBook[j].Address
	})

	// Prepare an exportable address book which has a partial view of the network
	return newAddrBook.GetExportableAddrBook()
}

func (a AddressBook) GetExportableAddrBook() ExportableAddressBook {
	exportableAddressBook := make(ExportableAddressBook, 0)
	for _, n := range a {
		exportableAddressBook = append(exportableAddressBook, NodeExportable{
			Address: n.Address,
			IsDead:  n.IsDead,
		})
	}
	return exportableAddressBook
}

// Ensures that self node was not removed from its own address book
func ensureSelfInAddrBook(selfNode *Node, partialAddressBook AddressBook) AddressBook {
	addrBookLen := len(partialAddressBook)
	index := sort.Search(addrBookLen, func(i int) bool {
		return partialAddressBook[i].Address >= selfNode.Address
	})
	if index < addrBookLen && partialAddressBook[index].Address == selfNode.Address {
		return partialAddressBook
	}
	partialAddressBook = append(partialAddressBook, selfNode)
	return partialAddressBook
}

func (a *AddressBook) MarshalJSON() ([]byte, error) {
	nodesExported := make(ExportedAddressBook, 0)
	for _, node := range *a {
		ne := NodeExported{
			GlobalPosition:              node.GlobalPosition,
			Address:                     node.Address,
			Redundancy:                  uint64(node.MessagesReceived),
			PartialViewershipPercentage: node.PartialViewershipPercentage,
			PartialAddressBook:          nil,
			IsDead:                      node.IsDead,
		}
		if node.exportPartialAddressBook {
			ne.PartialAddressBook = node.PartialAddressBook
		}
		nodesExported = append(nodesExported, ne)
	}
	return json.Marshal(nodesExported)
}

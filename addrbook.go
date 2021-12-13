package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"
)

type AddressBook []Node

type ExportableAddressBook []NodeExportable

func PopulateAddressBook(c *Config) (addrBook AddressBook) {
	log.Println("Populating the global address book by creating nodes")
	partialViewershipPercentages := getPartialViewershipCurve(c)
	if c.ViewershipPercentageFixed {
		partialViewershipPercentages = c.ViewershipCurveArray
	}
	if c.RandomizePartialAddressBooks {
		randSeed()
		rand.Shuffle(len(partialViewershipPercentages), func(i, j int) {
			partialViewershipPercentages[i], partialViewershipPercentages[j] =
				partialViewershipPercentages[j], partialViewershipPercentages[i]
		})
	}
	bz, _ := json.Marshal(partialViewershipPercentages)
	log.Println(string(bz))
	addrBook = make([]Node, c.NumberOfNodes)
	// populate global address book
	for i := uint64(0); i < c.NumberOfNodes; i++ {
		addrBook[i].CreateNode(newAddress())
	}
	sort.Slice(addrBook, func(i int, j int) bool {
		return addrBook[i].Address < addrBook[j].Address
	})
	for i := uint64(0); i < c.NumberOfNodes; i++ {
		addrBook[i].IsDead = getIsDead(i, c)
	}
	for i := uint64(0); i < c.NumberOfNodes; i++ {
		addrBookCopy := make(AddressBook, len(addrBook))
		copy(addrBookCopy, addrBook)
		partialAddressBook := GetPartialAddressBook(c, addrBook[i], addrBookCopy, float64(partialViewershipPercentages[i]))
		addrBook[i].InitNode(i, partialAddressBook, uint8(partialViewershipPercentages[i]), c.ShowIndividualNodePartialAddressBooks)
	}
	return
}

func GetPartialAddressBook(c *Config, selfNode Node, a AddressBook, partialViewershipPercentage float64) (partial ExportableAddressBook) {
	maxIndex := int(math.Trunc(float64(partialViewershipPercentage) / 100 * float64(len(a))))
	randSeed()
	// shuffle the original address book (this affect the original pointer)
	if c.RandomizePartialAddressBooks {
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	}
	// truncate the address book from 0 to maxI dex
	newAddrBook := a[:maxIndex]
	// sort the new address book lexicographically
	sort.Slice(newAddrBook, func(i int, j int) bool {
		return newAddrBook[i].Address < newAddrBook[j].Address
	})
	// make sure self node is within the address book
	newAddrBook = EnsureSelf(selfNode, newAddrBook)
	partial = make([]NodeExportable, 0)
	for _, n := range newAddrBook {
		partial = append(partial, NodeExportable{
			Address: n.Address,
			IsDead:  n.IsDead,
		})
	}
	return
}

func GetNodePosition(nodeID string, a ExportableAddressBook) (index int) {
	index = sort.Search(len(a), func(i int) bool {
		return a[i].Address >= nodeID
	})
	return
}

func CalculateLevels(a ExportableAddressBook) uint {
	fullListSize := float64(len(a))
	return uint(math.Ceil(math.Round(logBase3(fullListSize)*100) / 100))
}

func EnsureSelf(selfNode Node, partialAddressBook AddressBook) AddressBook {
	addrBookLen := len(partialAddressBook)
	index := sort.Search(addrBookLen, func(i int) bool {
		return partialAddressBook[i].Address >= selfNode.Address
	})
	if index < addrBookLen && partialAddressBook[index].Address == selfNode.Address {
		return partialAddressBook
	}
	partialAddressBook = append(partialAddressBook, selfNode)
	// sort the new address book lexicographically
	sort.Slice(partialAddressBook, func(i int, j int) bool {
		return strings.Compare(partialAddressBook[i].Address, partialAddressBook[j].Address) == -1
	})
	return partialAddressBook
}

func (a *AddressBook) MarshalJSON() ([]byte, error) {
	nodesExported := make([]NodeExported, 0)
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

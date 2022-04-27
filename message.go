package main

import (
	"log"
	"math"
)

const (
	TargetAddrBookCoverageAtlevel = float64(2) / float64(3)
	FirstTargetRelativeIndex      = float64(1) / float64(3)
	SecondTargetRelativeIndex     = float64(2) / float64(3)
)

type Message struct {
	Hash                  string
	Level                 int // 0 TODO: is this the current level?
	TotalNumNetworkLevels int // 7 TODO: is this the max number of levels?
}

func (m *Message) Copy() Message {
	return Message{
		Hash:                  m.Hash,
		Level:                 m.Level,
		TotalNumNetworkLevels: m.TotalNumNetworkLevels,
	}
}

func NewMessage(a ExportableAddressBook) Message {
	networkLevel := CalculateNumLevelsInAddrBook(a)
	return Message{
		Hash:                  newHash(),
		Level:                 int(networkLevel),
		TotalNumNetworkLevels: int(networkLevel),
	}
}

// TODO: Why are we just not subtracting one until we're at 0?
func CalibrateLevels(networkLevel uint, m Message) (levelWithDecrement, levelWithoutDecrement int64) {
	levelWithDecrement = int64(uint(float64(networkLevel) / float64(m.TotalNumNetworkLevels) * (float64(m.Level) - 1)))
	levelWithoutDecrement = int64(uint(float64(networkLevel) / float64(m.TotalNumNetworkLevels) * (float64(m.Level))))
	return
}

func GetTargetListLength(fullListSize int, TotalNumNetworkLevels, messageLevel int64) int {
	m := math.Pow(TargetAddrBookCoverageAtlevel, float64(TotalNumNetworkLevels-messageLevel))
	return int(math.Ceil(float64(fullListSize) * m))
}

func GetTargetIndices(nodePosition, targetListLength, partialAddressBookSize int) (targetA, targetB int) {
	targetA = (nodePosition + int(math.Round(float64(targetListLength)*FirstTargetRelativeIndex))) % partialAddressBookSize
	targetB = (nodePosition + int(math.Round(float64(targetListLength)*SecondTargetRelativeIndex))) % partialAddressBookSize
	return
}

func GetRecursiveTargets(c *Config, targetAIndex, targetBIndex, partialAddressBookSize int, partialAddressBook ExportableAddressBook) (targetA, targetB string) {
	targetA = RecursiveTargetFinder(c.MaxHotlist, -1, targetAIndex, partialAddressBookSize, partialAddressBook)
	targetB = RecursiveTargetFinder(c.MaxHotlist, -1, targetBIndex, partialAddressBookSize, partialAddressBook)
	return
}

func RecursiveTargetFinder(max uint, count, targetIndex, addrBookLen int, a ExportableAddressBook) (target string) {
	if count >= int(max) {
		log.Println("WARNING: First level `hotlist` exhausted with consecutive bad nodes")
		return ""
	}
	if !a[targetIndex].IsDead {
		return a[targetIndex].Address
	}
	newTargetA := (targetIndex + 1) % addrBookLen
	count++
	return RecursiveTargetFinder(max, count, newTargetA, addrBookLen, a)
}

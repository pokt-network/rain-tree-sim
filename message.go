package main

import (
	"log"
	"math"
)

type Message struct {
	Hash          string
	Level         int // 0
	NetworkLevels int // 7
}

func (m *Message) Copy() Message {
	return Message{
		Hash:          m.Hash,
		Level:         m.Level,
		NetworkLevels: m.NetworkLevels,
	}
}

func NewMessage(a ExportableAddressBook) Message {
	networkLevel := CalculateLevels(a)
	return Message{
		Hash:          newHash(),
		Level:         int(networkLevel),
		NetworkLevels: int(networkLevel),
	}
}

func CalibrateLevels(networkLevel uint, m Message) (levelWithDecrement, levelWithoutDecrement int64) {
	levelWithDecrement = int64(uint(float64(networkLevel) / float64(m.NetworkLevels) * (float64(m.Level) - 1)))
	levelWithoutDecrement = int64(uint(float64(networkLevel) / float64(m.NetworkLevels) * (float64(m.Level))))
	return
}

func GetTargetListLength(fullListSize int, networkLevels, messageLevel int64) int {
	m := math.Pow(float64(2)/float64(3), float64(networkLevels-messageLevel))
	return int(math.Ceil(float64(fullListSize) * m))
}

func GetTargetIndices(nodePosition, targetListLength, partialAddressBookSize int) (targetA, targetB int) {
	targetA = (nodePosition + int(math.Round(float64(targetListLength)/3))) % partialAddressBookSize
	targetB = (nodePosition + int(math.Round(float64(targetListLength)/1.5))) % partialAddressBookSize
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

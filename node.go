package main

import "math/rand"

type Node struct {
	GlobalPosition              uint64
	Address                     string
	PartialViewershipPercentage uint8
	PartialAddressBookPosition  int
	PartialAddressBook          ExportableAddressBook
	IsDead                      bool
	MessagesReceived            int
	Message                     Message
	exportPartialAddressBook    bool
}

type NodeExportable struct {
	Address string
	IsDead  bool
}

type NodeExported struct {
	GlobalPosition              uint64
	Address                     string
	Redundancy                  uint64
	PartialViewershipPercentage uint8
	PartialAddressBook          ExportableAddressBook
	IsDead                      bool
}

func CreateNode() *Node {
	return &Node{
		Address:                    newRandomAddress(),
		PartialAddressBookPosition: -1,
	}
}

func (n *Node) InitNode(globalPosition uint64, partialAddrBook ExportableAddressBook, partialVieweshipPer uint8, exportPartialAddressBook bool) {
	n.GlobalPosition = globalPosition
	n.PartialAddressBook = partialAddrBook
	n.PartialAddressBookPosition = GetNodePosition(n.Address, partialAddrBook)
	n.PartialViewershipPercentage = partialVieweshipPer
	n.exportPartialAddressBook = exportPartialAddressBook
}

// TODO: Not a fan of this function pattern, so rethink it once things work
func (n Node) CopyNodeWithMessage(m Message, isSelfMessage bool) (node Node) {
	// If dead => don't queue
	if n.IsDead {
		return n
	}
	// no dup checks on self queues
	if !isSelfMessage {
		// track
		n.MessagesReceived = n.MessagesReceived + 1
		// if dup => track but don't queue
		if n.MessagesReceived > 1 {
			return n
		}
	}
	// queue the message on this node
	n.Message = m
	return n
}

func newRandomAddress() string {
	var addr Address
	_, err := rand.Read(addr[:])
	if err != nil {
		_ = err
	}
	return addr.String()
}

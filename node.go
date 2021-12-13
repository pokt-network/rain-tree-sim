package main

type Node struct {
	GlobalPosition              uint64
	Address                     string
	PartialViewershipPercentage uint8
	PartialAddressBookPosition int
	PartialAddressBook         ExportableAddressBook
	IsDead                     bool
	MessagesReceived            int
	Message                     Message
	exportPartialAddressBook    bool
}

type NodeExportable struct {
	Address string
	IsDead  bool
}

func (n *Node) CreateNode(address string) {
	n.Address = address
	n.PartialAddressBookPosition = -1
}

func (n *Node) InitNode(globalPosition uint64, partialAddrBook ExportableAddressBook, pap uint8, exportPartialAddressBook bool) {
	n.GlobalPosition = globalPosition
	n.PartialAddressBook = partialAddrBook
	n.PartialAddressBookPosition = GetNodePosition(n.Address, partialAddrBook)
	n.PartialViewershipPercentage = pap
	n.exportPartialAddressBook = exportPartialAddressBook
}

func (n Node) WithMessage(m Message, isSelfMessage bool) (node Node) {
	// if dead... don't queue
	if n.IsDead {
		return n
	}
	// no dup checks on self queues
	if !isSelfMessage {
		// track
		n.MessagesReceived = n.MessagesReceived + 1
		// if dup... track -> but don't queue
		if n.MessagesReceived > 1 {
			return n
		}
	}
	// queue the message
	n.Message = m
	// sort by highest level message; (might not be needed)
	return n
}

type NodeExported struct {
	GlobalPosition              uint64
	Address                     string
	Redundancy                  uint64
	PartialViewershipPercentage uint8
	PartialAddressBook          ExportableAddressBook
	IsDead                      bool
}

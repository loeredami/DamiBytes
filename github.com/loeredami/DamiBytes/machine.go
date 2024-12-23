package main

type Machine struct {
	memSize uint64
	procC   uint64
	threadC uint64
	regC    uint64

	regs []uint64

	memory []uint8

	address_map map[uint64]uint64

	stack8    []uint8
	stack16   []uint16
	stack32   []uint32
	stack64   []uint64
	processes []*MachineProcess
}

type MachineProcess struct {
	programP uint64
	machine  *Machine
	pID      uint64
}

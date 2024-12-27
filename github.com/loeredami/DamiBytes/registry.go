package main

func (machine *Machine) getReg_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	registryIdx := uint64(0)

	if machine.bit64 {
		registryIdx = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		machine.stack64 = append(machine.stack64, machine.regs[registryIdx])
	} else {
		registryIdx = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		machine.stack32 = append(machine.stack32, uint32(machine.regs[registryIdx]))
	}

	machine.stackOpen = true
}

func (machine *Machine) setReg_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	registryIdx := uint64(0)
	value := uint64(0)

	if machine.bit64 {
		registryIdx = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]

		value = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		registryIdx = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]

		value = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	machine.regs[registryIdx] = value
}
package main

func (machine *Machine) getReg_inst(proc *MachineProcess) {
	registryIdx := uint64(0)

	if proc.bit64 {
		registryIdx = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		proc.stack64 = append(proc.stack64, machine.regs[registryIdx])
	} else {
		registryIdx = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		proc.stack32 = append(proc.stack32, uint32(machine.regs[registryIdx]))
	}
}

func (machine *Machine) setReg_inst(proc *MachineProcess) {
	registryIdx := uint64(0)
	value := uint64(0)

	if proc.bit64 {
		registryIdx = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]

		value = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		registryIdx = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]

		value = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	machine.regs[registryIdx] = value
}
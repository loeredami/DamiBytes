package main

func (machine *Machine) comp_inst(proc *MachineProcess) {
	var val1 uint64
	var val2 uint64

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.int_pop()

		val2 = proc.stack64[len(proc.stack64)-1]
		proc.int_pop()
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.int_pop()

		val2 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.int_pop()
	}

	res := makeComparisonResult(val1, val2)

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, uint64(res))
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}

func (machine *Machine) bitAnd_inst(proc *MachineProcess) {
	val1 := uint64(0)
	val2 := uint64(0)

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		val2 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		val2 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	res := val1 & val2

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, res)
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}
func (machine *Machine) bitOr_inst(proc *MachineProcess) {
	val1 := uint64(0)
	val2 := uint64(0)

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		val2 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		val2 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	res := val1 | val2

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, res)
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}

func (machine *Machine) bitLShift_inst(proc *MachineProcess) {
	val1 := uint64(0)
	val2 := uint64(0)

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		val2 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		val2 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	res := val1 << val2

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, res)
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}

func (machine *Machine) bitRShift_inst(proc *MachineProcess) {
	val1 := uint64(0)
	val2 := uint64(0)

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		val2 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		val2 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	res := val1 >> val2

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, res)
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}

func (machine *Machine) bitNot_inst(proc *MachineProcess) {
	val1 := uint64(0)

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	res := ^val1

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, res)
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}

func (machine *Machine) bitXor_inst(proc *MachineProcess) {
	val1 := uint64(0)
	val2 := uint64(0)

	if proc.bit64 {
		val1 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		val2 = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		val1 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		val2 = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	res := val1 ^ val2

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, res)
	} else {
		proc.stack32 = append(proc.stack32, uint32(res))
	}
}
func (machine *Machine) push_inst(payload uint64, proc *MachineProcess) {
	if proc.bit64 {
		proc.stack64 = append(proc.stack64, payload)
	} else {
		proc.stack32 = append(proc.stack32, uint32(payload))
	}
}

func (machine *Machine) pop_inst(payload uint64, proc *MachineProcess) {
	registryIdx := payload
	value := uint64(0)

	if proc.bit64 {
		value = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		value = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}
	machine.regs[registryIdx] = value
}

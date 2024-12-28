package main

func (machine *Machine) comp_inst() {
	var val1 uint64
	var val2 uint64

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.int_pop()

		val2 = machine.stack64[len(machine.stack64)-1]
		machine.int_pop()
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.int_pop()

		val2 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.int_pop()
	}

	res := makeComparisonResult(val1, val2)

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, uint64(res))
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}
}

func (machine *Machine) bitAnd_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	val1 := uint64(0)
	val2 := uint64(0)

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		val2 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		val2 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	res := val1 & val2

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, res)
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}

	machine.stackOpen = true
}
func (machine *Machine) bitOr_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	val1 := uint64(0)
	val2 := uint64(0)

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		val2 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		val2 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	res := val1 | val2

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, res)
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}

	machine.stackOpen = true
}

func (machine *Machine) bitLShift_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	val1 := uint64(0)
	val2 := uint64(0)

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		val2 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		val2 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	res := val1 << val2

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, res)
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}

	machine.stackOpen = true
}

func (machine *Machine) bitRShift_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	val1 := uint64(0)
	val2 := uint64(0)

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		val2 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		val2 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	res := val1 >> val2

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, res)
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}

	machine.stackOpen = true
}

func (machine *Machine) bitNot_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	val1 := uint64(0)

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	res := ^val1

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, res)
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}

	machine.stackOpen = true
}

func (machine *Machine) bitXor_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	val1 := uint64(0)
	val2 := uint64(0)

	if machine.bit64 {
		val1 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		val2 = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		val1 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		val2 = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	res := val1 ^ val2

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, res)
	} else {
		machine.stack32 = append(machine.stack32, uint32(res))
	}

	machine.stackOpen = true
}
func (machine *Machine) push_inst(payload uint64) {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, payload)
	} else {
		machine.stack32 = append(machine.stack32, uint32(payload))
	}

	machine.stackOpen = true
}

func (machine *Machine) pop_inst(payload uint64) {
	registryIdx := payload

	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	value := uint64(0)

	if machine.bit64 {
		value = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		value = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	machine.regs[registryIdx] = value
}

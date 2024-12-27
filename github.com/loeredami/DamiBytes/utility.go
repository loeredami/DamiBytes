package main

func (machine *Machine) int_pop() {
	if machine.bit64 {
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}
}

func (machine *Machine) decodePayload(instruction []byte, startBit uint64, bitLength uint64) uint64 {
	var result uint64 = 0
	for i := range bitLength {
		result = result | (((((uint64)(instruction[startBit/8])) >> (startBit % 8)) & 1) << i)
		startBit++
	}
	return result
}

func (machine *Machine) handle_instruction(address uint64, proc *MachineProcess) {
	if machine.bit64 {
		inst1 := machine.memory[address]
		inst2 := machine.memory[address+1]

		payload1 := machine.memory[address+2]
		payload2 := machine.memory[address+3]
		payload3 := machine.memory[address+4]
		payload4 := machine.memory[address+5]
		payload5 := machine.memory[address+6]
		payload6 := machine.memory[address+7]

		var instruction uint16 = (uint16(inst1) << 8) + uint16(inst2)

		var payload uint64 = (uint64(payload1) << (8 * 5)) +
			(uint64(payload2) << (8 * 4)) +
			(uint64(payload3) << (8 * 3)) +
			(uint64(payload4) << (8 * 2)) +
			(uint64(payload5) << 8) +
			uint64(payload6)

		machine.run_instruction(instruction, payload, proc)
	} else {
		inst := machine.memory[address]

		payload1 := machine.memory[address+1]
		payload2 := machine.memory[address+2]
		payload3 := machine.memory[address+3]

		var instruction uint16 = uint16(inst)

		var payload uint64 = (uint64(payload3) << (8 * 2)) +
			(uint64(payload2) << 8) +
			uint64(payload1)

		machine.run_instruction(instruction, payload, proc)
	}
}

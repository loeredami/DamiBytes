package main

const (
	PROCESS_ACTIVE   uint8 = 0x00000001
	PROCESS_SLEEPING uint8 = 0x00000010
)

type Machine struct {
	memSize uint64
	procC   uint16
	threadC uint16
	regC    uint64

	lastPID uint64

	regs []uint64

	memory []uint8

	external_map map[uint64]uint64

	stack8    []uint8
	stack16   []uint16
	stack32   []uint32
	stack64   []uint64
	processes []*MachineProcess
	bit64     bool
}

type MachineProcess struct {
	programP uint64
	machine  *Machine
	pID      uint64
	state    uint8
}

func makeMachine(memSize uint64, processorCount uint16, threadC uint16, regC uint64) *Machine {
	registers := []uint64{}
	for i := uint64(0); i < regC; i++ {
		registers = append(registers, 0)
	}

	memory := []uint8{}
	for i := uint64(0); i < memSize; i++ {
		memory = append(memory, 0)
	}

	external_map := map[uint64]uint64{}

	stack8 := []uint8{}
	stack16 := []uint16{}
	stack32 := []uint32{}
	stack64 := []uint64{}
	processes := []*MachineProcess{}

	return &Machine{
		memSize, processorCount, threadC, regC, 0, registers, memory, external_map,
		stack8, stack16, stack32, stack64, processes, false,
	}
}

func (machine *Machine) makeProcess(programP uint64) {
	machine.lastPID += 1
	PID := machine.lastPID
	state := PROCESS_ACTIVE
	machine.processes = append(machine.processes, &MachineProcess{
		programP: programP,
		machine:  machine,
		pID:      PID,
		state:    state,
	})
}

func (machine *Machine) should_shutdown() bool {
	return false
}

func (machine *Machine) run_instruction(instruction uint16, payload uint64) {

}

func (machine *Machine) handle_instruction(address uint64) {
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

		machine.run_instruction(instruction, payload)
	} else {
		inst := machine.memory[address]

		payload1 := machine.memory[address+1]
		payload2 := machine.memory[address+2]
		payload3 := machine.memory[address+3]

		var instruction uint16 = uint16(inst)

		var payload uint64 = (uint64(payload3) << (8 * 2)) +
			(uint64(payload2) << 8) +
			uint64(payload1)

		machine.run_instruction(instruction, payload)
	}
}

func (machine *Machine) tick() {
	flagged_for_removal := []int{}

	for_processing := []*MachineProcess{}

	for i, proc := range machine.processes {
		if (proc.state & PROCESS_ACTIVE) == 0 {
			flagged_for_removal = append(flagged_for_removal, i)
			continue
		}
		if (proc.state & PROCESS_SLEEPING) != 0 {
			continue
		}
		for_processing = append(for_processing, proc)
	}

	for _, val := range flagged_for_removal {
		machine.processes = append(machine.processes[:val], machine.processes[val+1:]...)
	}

	// Not complete
}

func (machine *Machine) run() {
	machine.makeProcess(0)
	for {
		if machine.should_shutdown() {
			break
		}
		machine.tick()
	}
}
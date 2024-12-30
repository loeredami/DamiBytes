package main

import (
	"sync"
)

const (
	PROCESS_ACTIVE   byte = 0x00000001
	PROCESS_SLEEPING byte = 0x00000010
)

type Machine struct {
	memSize uint64
	procC   uint16
	regC    uint64

	lastPID uint64

	regs []uint64

	memory []byte

	processes []*MachineProcess
}

type MachineProcess struct {
	programP *byte
	machine  *Machine
	pID      uint64
	state    byte

	stack32   []uint32
	stack64   []uint64
	bit64     bool
}

func makeMachine(memSize uint64, processorCount uint16, regC uint64) *Machine {
	registers := make([]uint64, regC)
	memory := make([]byte, memSize)

	processes := []*MachineProcess{}

	return &Machine{
		memSize, processorCount, regC, 0, registers, memory,
		processes,
	}
}

func (machine *Machine) streamInProgram(program []byte, start uint64) {
	for i := range len(program) {
		machine.memory[int(start)+i] = program[i]
	}
}

func (machine *Machine) makeProcess(programP *byte) {
	machine.lastPID += 1
	PID := machine.lastPID
	state := PROCESS_ACTIVE
	machine.processes = append(machine.processes, &MachineProcess{
		programP: programP,
		machine:  machine,
		pID:      PID,
		state:    state,
		stack32:  []uint32{},
		stack64:  []uint64{},
		bit64:    false,
	})
}

func (machine *Machine) should_shutdown() bool {
	return len(machine.processes) < 1 
}

func (machine *Machine) run_instruction(instruction uint8, proc *MachineProcess) {
	switch instruction {
	case 0x0000: break 
	case 0x0001: machine.add_inst(proc)
	case 0x0002: machine.sub_inst(proc)
	case 0x0003: machine.mul_inst(proc)
	case 0x0004: machine.div_inst(proc)
	case 0x0005: machine.mod_inst(proc)
	case 0x0006: machine.addF_inst(proc)
	case 0x0007: machine.subF_inst(proc)
	case 0x0008: machine.mulF_inst(proc)
	case 0x0009: machine.divF_inst(proc)
	case 0x000A: machine.store_inst(proc)
	case 0x000B: machine.load_inst(proc)
	case 0x000C: machine.syscallHandle(proc)
	case 0x000D: machine.jump_inst(proc)
	case 0x000E: machine.comp_inst(proc)
	case 0x000F: machine.if_inst(proc)
	case 0x0010: machine.inc_instruction(proc)
	case 0x0011: machine.dec_instruction(proc)
	case 0x0012: machine.bitAnd_inst(proc)
	case 0x0013: machine.bitOr_inst(proc)
	case 0x0014: machine.bitNot_inst(proc)
	case 0x0015: machine.bitXor_inst(proc)
	case 0x0016: machine.interruptOn_inst(proc)
	case 0x0017: machine.interruptOf_inst(proc)
	case 0x0018: // in
	case 0x0019: // out
	case 0x001A: machine.bitLShift_inst(proc)
	case 0x001B: machine.bitRShift_inst(proc)
	case 0x001C: machine.ext_inst(proc)
	case 0x001D: machine.push_inst(proc)
	case 0x001E: machine.pop_inst(proc)
	case 0x001F: machine.getReg_inst(proc)
	case 0x0020: machine.setReg_inst(proc)
	case 0x0021: machine.ptrHere_inst(proc)
	case 0x0022: machine.free_inst(proc)
	case 0x0023: machine.go_inst(proc)
	case 0x0024: machine.pID_inst(proc)
	case 0x0025: machine.memIncr_inst(proc)
	case 0x0026: machine.memDec_inst(proc)
	case 0x0027: machine.bits_inst(proc)
	case 0x0028: machine.machineData_inst(proc)
	case 0x0029: machine.func_inst(proc)
	case 0x002A: machine.exit_inst(proc)
	default: break
	}
}

func (machine *Machine) doWork(processes []*MachineProcess, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, proc := range processes {
		machine.handle_instruction(proc.programP, proc)
		}
	ch <- 0
}

func (machine *Machine) tick() {
	flagged_for_removal := []int{}


	for i, proc := range machine.processes {
		if (proc.state & PROCESS_ACTIVE) == 0 {
			flagged_for_removal = append(flagged_for_removal, i)
			continue
		}
		if (proc.state & PROCESS_SLEEPING) != 0 {
			continue
		}
	}

	for _, val := range flagged_for_removal {
		machine.processes = append(machine.processes[:val], machine.processes[val+1:]...)
	}

	procSplit := [][]*MachineProcess{}

	for range machine.procC {
		procSplit = append(procSplit, []*MachineProcess{})
	}

	for i, proc := range machine.processes {
		procSplit[i%int(machine.procC)] = append(procSplit[i%int(machine.procC)], proc)
	}

	ch := make(chan int)
	var wg sync.WaitGroup

	for _, processes := range procSplit {
		wg.Add(1)
		go machine.doWork(processes, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := []int{}

	for result := range ch {
		results = append(results, result)
	}

	for i, result := range results {
		if result != 0 {
			panic([]int{i, result})
		}
	}
}

func (machine *Machine) run() {
	machine.makeProcess(&machine.memory[0])
	for {
		if machine.should_shutdown() {
			break
		}
		machine.tick()
	}
}
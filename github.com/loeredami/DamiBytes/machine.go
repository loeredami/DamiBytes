package main

import (
	"encoding/binary"
	"fmt"
	"sync"
	"syscall"
	"unsafe"
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

	external_map map[uint64]uint64

	stackOpen bool
	stack32   []uint32
	stack64   []uint64
	processes []*MachineProcess
	bit64     bool
}

type MachineProcess struct {
	programP uint64
	machine  *Machine
	pID      uint64
	state    byte
}

func makeMachine(memSize uint64, processorCount uint16, regC uint64) *Machine {
	registers := []uint64{}
	for i := uint64(0); i < regC; i++ {
		registers = append(registers, 0)
	}

	memory := []byte{}
	for i := uint64(0); i < memSize; i++ {
		memory = append(memory, 0)
	}

	external_map := map[uint64]uint64{}

	stack32 := []uint32{}
	stack64 := []uint64{}
	processes := []*MachineProcess{}

	return &Machine{
		memSize, processorCount, regC, 0, registers, memory, external_map,
		true, stack32, stack64, processes, false,
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

func(machine *Machine) int_pop() {
	if machine.bit64 {
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}
}

func (machine *Machine) add_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1+num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1+num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) sub_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1-num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1-num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) mul_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1*num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1*num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) div_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1/num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1/num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) mod_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1%num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1%num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) addF_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1+num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1+num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) subF_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1-num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1-num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) mulF_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1*num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1*num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) divF_inst() { 
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1/num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1/num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) store_inst(payload uint64) {
	instructionSize := 16
	if !machine.bit64 {
		instructionSize = 8
	}
	valueSizeBits := 4
	if !machine.bit64 {
		valueSizeBits = 3
	}
	addressBits := 44
	if !machine.bit64 {
		addressBits = 21
	}

	payloadAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(payloadAsBytes, payload)

	valueSize := machine.decodePayload(payloadAsBytes, uint64(instructionSize), uint64(valueSizeBits))
	address := machine.decodePayload(payloadAsBytes, uint64(instructionSize) + uint64(valueSizeBits), uint64(addressBits))

	value := machine.stack64[len(machine.stack64)-1]
	if !machine.bit64 {
		value = uint64(machine.stack32[len(machine.stack32)-1])
	}
	machine.int_pop()

	machine.store_value(address, value, valueSize)
}

func (machine *Machine) store_value(address uint64, value uint64, valueSize uint64) {
	for i := range valueSize {
		machine.memory[address + i] = (byte(value) >> (i*8)) & 0xFF
	}
}

func (machine *Machine) decodePayload(instruction []byte, startBit uint64, bitLength uint64) uint64 {
	var result uint64 = 0
	for i := range bitLength {
		result = result | (((((uint64)(instruction[startBit / 8])) >> (startBit % 8)) & 1) << i)
		startBit++
	}
	return result
}

func (machine *Machine) load_inst(payload uint64) {
	instructionSize := uint64(16)
    if !machine.bit64 {
        instructionSize = 8
    }

    valueSizeBits := uint64(4)
    addressBits := uint64(44)
    if !machine.bit64 {
        valueSizeBits = 3
        addressBits = 21
    }

	payloadBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(payloadBytes, payload)

	valueSize := machine.decodePayload(payloadBytes, instructionSize, valueSizeBits)
    address := machine.decodePayload(payloadBytes, instructionSize+valueSizeBits, addressBits)

	value := machine.loadValueFromMemory(address, valueSize)

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, uint64(value))
	} else {
		machine.stack32 = append(machine.stack32, uint32(value))
	}
}


func (machine *Machine) loadValueFromMemory(address uint64, valueSize uint64) uint64 {
    if address+valueSize > uint64(len(machine.memory)) {
        panic("memory out of bounds")
    }

    value := uint64(0)
    for i := uint64(0); i < valueSize; i++ {
        value |= uint64(machine.memory[address+i]) << (i * 8)
    }
    return value
}

func (machine *Machine) run_instruction(instruction uint16, payload uint64) {
	switch instruction {
	case 0x0000: break 
	case 0x0001: machine.add_inst()
	case 0x0002: machine.sub_inst()
	case 0x0003: machine.mul_inst()
	case 0x0004: machine.div_inst()
	case 0x0005: machine.mod_inst()
	case 0x0006: machine.addF_inst()
	case 0x0007: machine.subF_inst()
	case 0x0008: machine.mulF_inst()
	case 0x0009: machine.divF_inst()
	case 0x000A: machine.store_inst(payload)
	case 0x000B: machine.load_inst(payload)
	case 0x000C: machine.syscallHandle()
	default: break
	}
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

func (machine *Machine) syscallHandle() {
	if machine.bit64 {
		if len(machine.stack64) < 1 {
			panic("syscallHandler: no syscall number on stack")
		}

		syscallNumber := machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]

		if len(machine.stack64) < 1 {
			panic("syscallHandler: missing argument count")
		}

		argCount := machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]

		if len(machine.stack64) < int(argCount) {
			panic("syscallHandler: not enough arguments on stack")
		}

		args := make([]uintptr, argCount)
		
		for i := 0; i < int(argCount); i++ {
			args[i] = uintptr(machine.stack64[len(machine.stack64)-1])
			machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		}

		ret, _, err := syscall.SyscallN(uintptr(syscallNumber), args...)

		if err != 0 {
			panic(fmt.Sprintf("syscall failed: %v", err))
		}

		machine.stack64 = append(machine.stack64, uint64(ret))
	} else {
		if len(machine.stack32) < 1 {
			panic("syscallHandler: no syscall number on stack")
		}

		syscallNumber := machine.stack32[len(machine.stack32)-1]
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]

		if len(machine.stack32) < 1 {
			panic("syscallHandler: missing argument count")
		}

		argCount := machine.stack32[len(machine.stack32)-1]
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]

		if len(machine.stack32) < int(argCount) {
			panic("syscallHandler: not enough arguments on stack")
		}

		args := make([]uintptr, argCount)
		
		for i := 0; i < int(argCount); i++ {
			args[i] = uintptr(machine.stack32[len(machine.stack32)-1])
			machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		}

		ret, _, err := syscall.SyscallN(uintptr(syscallNumber), args...)

		if err != 0 {
			panic(fmt.Sprintf("syscall failed: %v", err))
		}

		machine.stack32 = append(machine.stack32, uint32(ret))
	}
}

func (machine *Machine) doWork(processes []*MachineProcess, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, proc := range processes {
		machine.handle_instruction(proc.programP)
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
		proc.programP++
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
	machine.makeProcess(0)
	for {
		if machine.should_shutdown() {
			break
		}
		machine.tick()
	}
}
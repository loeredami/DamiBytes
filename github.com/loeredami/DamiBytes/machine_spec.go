package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/pbnjay/memory"
	"golang.org/x/sys/windows"
)

func (machine *Machine) ext_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false

	dllNamePtr := uint64(0)
	funcNamePtr := uint64(0)

	if machine.bit64 {
		dllNamePtr = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]

		funcNamePtr = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		dllNamePtr = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]

		funcNamePtr = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	dll, err := windows.LoadDLL(getStringFromByte((*byte)(unsafe.Pointer(uintptr(dllNamePtr)))))

	if err != nil {
		panic(err)
	}

	proc, err := dll.FindProc(getStringFromByte((*byte)(unsafe.Pointer(uintptr(funcNamePtr)))))
	
	if err != nil {
		panic(err)
	}

	syscallNumber := proc.Addr()

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, uint64(syscallNumber))
	} else {
		machine.stack32 = append(machine.stack32, uint32(syscallNumber))
	}

	machine.stackOpen = true
}

func (machine *Machine) exit_inst(proc *MachineProcess) {
	proc.state = proc.state & (^PROCESS_ACTIVE)
}

func (machine *Machine) jump_inst(procC *MachineProcess) {
	address := uint64(0)

	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		address = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		address = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true


	procC.programP = (*byte)(unsafe.Pointer((uintptr)(address)))
}

func (machine *Machine) if_inst(payload uint64, proc *MachineProcess) {
	instructionSize := uint64(16)
	if !machine.bit64 {
		instructionSize = 8
	}

	compSizeBits := uint64(4)
	addressBits := uint64(44)
	if !machine.bit64 {
		compSizeBits = 3
		addressBits = 21
	}

	payloadBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(payloadBytes, payload)

	compType := machine.decodePayload(payloadBytes, instructionSize, compSizeBits)
	memoryOffset := machine.decodePayload(payloadBytes, instructionSize+compSizeBits, addressBits)

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

	if (compType & uint64(comparisonResults.isEqual)) != 0 {
		if (value & uint64(comparisonResults.isEqual)) != 0 {
			proc.programP = (*byte)(unsafe.Pointer(uintptr(memoryOffset) - 1))
		}
	} else if (compType & uint64(comparisonResults.isGreater)) != 0 {
		if (value & uint64(comparisonResults.isGreater)) != 0 {
			proc.programP = (*byte)(unsafe.Pointer(uintptr(memoryOffset) - 1))
		}
	}
}

func (machine *Machine) interruptOn_inst(proc *MachineProcess) {
	proc.state |= PROCESS_SLEEPING
}

func (machine *Machine) interruptOf_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false

	PID := uint64(0)

	if machine.bit64 {
		PID = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		PID = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	for _, proc := range machine.processes {
		if proc.pID == PID {
			proc.state = proc.state & (^PROCESS_SLEEPING)
			break
		}
	}
}

func (machine *Machine) go_inst() {
	// Starting a new process, usually always end up being run on another thread.
	address := uint64(0)
	
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		address = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		address = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	machine.makeProcess((*byte)(unsafe.Pointer(uintptr(address)-1)))
}

func (machine *Machine) pID_inst(proc *MachineProcess) {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, proc.pID)
	} else {
		machine.stack32 = append(machine.stack32, uint32(proc.pID))
	}

	machine.stackOpen = true
}

func (machine *Machine) bits_inst(payload uint64) {
	if payload == 32 {
		machine.bit64 = false
	} else if payload == 64 {
		machine.bit64 = true
	} else {
		machine.bit64 = !machine.bit64
	}
}

func (machine *Machine) machineData_inst() {
	var key uint64

	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		key = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		key = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}
	var messagePtr *byte

	switch key {
	case 0x0000:
		if machine.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack64[len(machine.stack64)-1])))
			machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack32[len(machine.stack32)-1])))
			machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		}

		streamStringToBytePointer(messagePtr, runtime.GOOS + string(rune(0x00)))
	case 0x0001:
		if machine.bit64 {
			machine.stack64 = append(machine.stack64, uint64(runtime.NumCPU()))
		} else {
			machine.stack32 = append(machine.stack32, uint32(runtime.NumCPU()))
		}
	case 0x0002:
		if machine.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack64[len(machine.stack64)-1])))
			machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack32[len(machine.stack32)-1])))
			machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		}
		streamStringToBytePointer(messagePtr, runtime.GOARCH + string(rune(0x00)))
	case 0x0003:
		if machine.bit64 {
			machine.stack64 = append(machine.stack64, memory.TotalMemory())
		} else {
			machine.stack32 = append(machine.stack32, uint32(memory.TotalMemory()))
		}
	case 0x0004:
		if machine.bit64 {
			machine.stack64 = append(machine.stack64, memory.FreeMemory())
		} else {
			machine.stack32 = append(machine.stack32, uint32(memory.FreeMemory()))
		}
	case 0x0005:
		if machine.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack64[len(machine.stack64)-1])))
			machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack32[len(machine.stack32)-1])))
			machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		}
		streamStringToBytePointer(messagePtr, string(runtime.CPUProfile()) + string(rune(0x00)))
	case 0x0006:
		if machine.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack64[len(machine.stack64)-1])))
			machine.stack64 = machine.stack64[:len(machine.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(machine.stack32[len(machine.stack32)-1])))
			machine.stack32 = machine.stack32[:len(machine.stack32)-1]
		}
		path, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		streamStringToBytePointer(messagePtr, path + string(rune(0x00)))
	case 0x0007:
		if machine.bit64 {
			machine.stack64 = append(machine.stack64, uint64(syscall.Stdin))
		} else {
			machine.stack32 = append(machine.stack32, uint32(syscall.Stdin))
		}
	case 0x0008:
		if machine.bit64 {
			machine.stack64 = append(machine.stack64, uint64(syscall.Stdout))
		} else {
			machine.stack32 = append(machine.stack32, uint32(syscall.Stdout))
		}
	}

	machine.stackOpen = true
}


func (machine *Machine) syscallHandle() {
	if machine.bit64 {

		for {
			if machine.stackOpen {
				break
			}
		}

		machine.stackOpen = false

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

		machine.stackOpen = true
	} else {

		for {
			if machine.stackOpen {
				break
			}
		}

		machine.stackOpen = false

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
		machine.stackOpen = true
	}
}


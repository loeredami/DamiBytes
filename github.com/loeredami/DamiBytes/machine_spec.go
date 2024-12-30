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

func (machine *Machine) ext_inst(proc *MachineProcess) {
	dllNamePtr := uint64(0)

	if proc.bit64 {
		dllNamePtr = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		dllNamePtr = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	dll, err := windows.LoadDLL(getStringFromByte((*byte)(unsafe.Pointer(uintptr(dllNamePtr)))))

	if err != nil {
		panic(err)
	}

	dllAddr := uintptr(unsafe.Pointer(dll))

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, uint64(dllAddr))
	} else {
		proc.stack32 = append(proc.stack32, uint32(dllAddr))
	}
}

func (machine *Machine) func_inst(proc *MachineProcess) {
	dllPtr := uint64(0)
	funcNamePtr := uint64(0)

	if proc.bit64 {
		dllPtr = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		funcNamePtr = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		dllPtr = uint64(proc.stack32[len(proc.stack64)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		funcNamePtr = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	dll := *(*windows.DLL)(unsafe.Pointer(uintptr(dllPtr)))

	procW, err := dll.FindProc(getStringFromByte((*byte)(unsafe.Pointer(uintptr(funcNamePtr)))))

	if err != nil {
		panic(err)
	}

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, uint64(procW.Addr()))
	} else {
		proc.stack32 = append(proc.stack32, uint32(procW.Addr()))
	}
}

func (machine *Machine) exit_inst(proc *MachineProcess) {
	proc.state = proc.state & (^PROCESS_ACTIVE)
}

func (machine *Machine) jump_inst(procC *MachineProcess) {
	address := uint64(0)
	if procC.bit64 {
		address = procC.stack64[len(procC.stack64)-1]
		procC.stack64 = procC.stack64[:len(procC.stack64)-1]
	} else {
		address = uint64(procC.stack32[len(procC.stack32)-1])
		procC.stack32 = procC.stack32[:len(procC.stack32)-1]
	}
	procC.programP = (*byte)(unsafe.Pointer((uintptr)(address)))
}

func (machine *Machine) if_inst(proc *MachineProcess) {
	var payload, payload2 uint64
	if proc.bit64 {
		payload = binary.BigEndian.Uint64(GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		8,
		))
		payload2 = binary.BigEndian.Uint64(GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8)),
		0,
		8,
		))
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+16))
	} else {
		payload = binary.BigEndian.Uint64(GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		4,
		))
		payload2 = binary.BigEndian.Uint64(GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+4)),
		0,
		4,
		))
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8))
	}

	compType := payload

	memoryOffset := payload2

	value := uint64(0)

	if proc.bit64 {
		value = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		value = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}


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

func (machine *Machine) interruptOf_inst(proc *MachineProcess) {
	PID := uint64(0)

	if proc.bit64 {
		PID = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		PID = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}
	for _, proc := range machine.processes {
		if proc.pID == PID {
			proc.state = proc.state & (^PROCESS_SLEEPING)
			break
		}
	}
}

func (machine *Machine) go_inst(proc *MachineProcess) {
	// Starting a new process, usually always end up being run on another thread.
	address := uint64(0)

	if proc.bit64 {
		address = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		address = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}
	machine.makeProcess((*byte)(unsafe.Pointer(uintptr(address)-1)))
}

func (machine *Machine) pID_inst(proc *MachineProcess) {
	if proc.bit64 {
		proc.stack64 = append(proc.stack64, proc.pID)
	} else {
		proc.stack32 = append(proc.stack32, uint32(proc.pID))
	}
}

func (machine *Machine) bits_inst(proc *MachineProcess) {
	var payload uint64
	if proc.bit64 {
		payload = binary.BigEndian.Uint64(GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		8,
		))
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8))

	} else {
		payload = uint64(binary.BigEndian.Uint32(GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		4,
		)))

		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+4))
	}
	if payload == 32 {
		proc.bit64 = false
	} else if payload == 64 {
		proc.bit64 = true
	} else {
		proc.bit64 = !proc.bit64
	}
}

func (machine *Machine) machineData_inst(proc *MachineProcess) {
	var key uint64
	if proc.bit64 {
		key = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		key = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}
	var messagePtr *byte

	switch key {
	case 0x0000:
		if proc.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack64[len(proc.stack64)-1])))
			proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack32[len(proc.stack32)-1])))
			proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		}

		streamStringToBytePointer(messagePtr, runtime.GOOS + string(rune(0x00)))
	case 0x0001:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, uint64(runtime.NumCPU()))
		} else {
			proc.stack32 = append(proc.stack32, uint32(runtime.NumCPU()))
		}
	case 0x0002:
		if proc.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack64[len(proc.stack64)-1])))
			proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack32[len(proc.stack32)-1])))
			proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		}
		streamStringToBytePointer(messagePtr, runtime.GOARCH + string(rune(0x00)))
	case 0x0003:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, memory.TotalMemory())
		} else {
			proc.stack32 = append(proc.stack32, uint32(memory.TotalMemory()))
		}
	case 0x0004:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, memory.FreeMemory())
		} else {
			proc.stack32 = append(proc.stack32, uint32(memory.FreeMemory()))
		}
	case 0x0005:
		if proc.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack64[len(proc.stack64)-1])))
			proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack32[len(proc.stack32)-1])))
			proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		}
		streamStringToBytePointer(messagePtr, string(runtime.CPUProfile()) + string(rune(0x00)))
	case 0x0006:
		if proc.bit64 {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack64[len(proc.stack64)-1])))
			proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		} else {
			messagePtr = (*byte)(unsafe.Pointer(uintptr(proc.stack32[len(proc.stack32)-1])))
			proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		}
		path, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		streamStringToBytePointer(messagePtr, path + string(rune(0x00)))
	case 0x0007:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, uint64(syscall.Stdin))
		} else {
			proc.stack32 = append(proc.stack32, uint32(syscall.Stdin))
		}
	case 0x0008:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, uint64(syscall.Stdout))
		} else {
			proc.stack32 = append(proc.stack32, uint32(syscall.Stdout))
		}
	case 0x0009:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, uint64(len(machine.memory)))
		} else {
			proc.stack32 = append(proc.stack32, uint32(len(machine.memory)))
		}
	case 0x000A:
		if proc.bit64 {
			proc.stack64 = append(proc.stack64, uint64(len(proc.stack64)))
		} else {
			proc.stack32 = append(proc.stack32, uint32(len(proc.stack32)))
		}
	}
}


func (machine *Machine) syscallHandle(proc *MachineProcess) {
	if proc.bit64 {
		if len(proc.stack64) < 1 {
			panic("syscallHandler: no syscall number on stack")
		}

		syscallNumber := proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]

		if len(proc.stack64) < 1 {
			panic("syscallHandler: missing argument count")
		}

		argCount := proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]

		if len(proc.stack64) < int(argCount) {
			panic("syscallHandler: not enough arguments on stack")
		}

		args := make([]uintptr, argCount)
		
		for i := 0; i < int(argCount); i++ {
			args[i] = uintptr(proc.stack64[len(proc.stack64)-1])
			proc.stack64 = proc.stack64[:len(proc.stack64)-1]
		}


		ret, _, err := syscall.SyscallN(uintptr(syscallNumber), args...)

		if err != 0 {
			panic(fmt.Sprintf("syscall failed: %v", err))
		}

		proc.stack64 = append(proc.stack64, uint64(ret))
	} else {
		if len(proc.stack32) < 1 {
			panic("syscallHandler: no syscall number on stack")
		}

		syscallNumber := proc.stack32[len(proc.stack32)-1]
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]

		if len(proc.stack32) < 1 {
			panic("syscallHandler: missing argument count")
		}

		argCount := proc.stack32[len(proc.stack32)-1]
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]

		if len(proc.stack32) < int(argCount) {
			panic("syscallHandler: not enough arguments on stack")
		}

		args := make([]uintptr, argCount)
		
		for i := 0; i < int(argCount); i++ {
			args[i] = uintptr(proc.stack32[len(proc.stack32)-1])
			proc.stack32 = proc.stack32[:len(proc.stack32)-1]
		}

		
		ret, _, err := syscall.SyscallN(uintptr(syscallNumber), args...)

		if err != 0 {
			panic(fmt.Sprintf("syscall failed: %v", err))
		}

		proc.stack32 = append(proc.stack32, uint32(ret))
	}
}


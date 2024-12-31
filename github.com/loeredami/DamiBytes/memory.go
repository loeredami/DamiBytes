package main

import (
	"encoding/binary"
	"unsafe"
)

func streamStringToBytePointer(startPtr *byte, str string) {
	for i := 0; i < len(str); i++ {
		*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(startPtr)) + uintptr(i))) = str[i]
	}
}

func getStringFromByte(startPtr *byte) string {
	var str string = ""
	var i int = 0
	for {
		if *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(startPtr))+uintptr(i))) == 0x0 {
			break
		}
		str += string(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(startPtr))+uintptr(i))))
		i += 1
	}
	return str
}
/*
func (machine *Machine) here_inst(proc *MachineProcess) {
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, uint64(uintptr(unsafe.Pointer(proc.programP))))
	} else {
		machine.stack32 = append(machine.stack32, uint32(uintptr(unsafe.Pointer(proc.programP))))
	}

	machine.stackOpen = true
}
*/

func (machine *Machine) memIncr_inst(proc *MachineProcess) {
	var size uint64 = 0
	if proc.bit64 {
		size = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		size = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	machine.memory = append(machine.memory, make([]byte, size)...)
}

func (machine *Machine) memDec_inst(proc *MachineProcess) {
	var size uint64 = 0
	if proc.bit64 {
		size = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		size = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}
	machine.memory = machine.memory[:len(machine.memory)-int(size)]
}

func (machine *Machine) free_inst(proc *MachineProcess) {
	address := uint64(0)
	size := uint64(0)

	if proc.bit64 {
		address = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]

		size = proc.stack64[len(proc.stack64)-1]
		proc.stack64 = proc.stack64[:len(proc.stack64)-1]
	} else {
		address = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]

		size = uint64(proc.stack32[len(proc.stack32)-1])
		proc.stack32 = proc.stack32[:len(proc.stack32)-1]
	}

	for i := address; i < address+size; i++ {
		machine.memory[i] = 0x0000 // 0 but pretty hex formatting
	}
}

func (machine *Machine) ptrHere_inst(proc *MachineProcess) {
	var payload uint64
	if proc.bit64 {
		bytes := GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		8,
		)
		payload = binary.BigEndian.Uint64(bytes)

		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8))
	} else {
		bytes := GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))-4)),
		0,
		4,
		)

		payload = uint64(binary.BigEndian.Uint32(bytes))

		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+4))
	}

	value := uintptr(unsafe.Pointer(&machine.memory[payload]))

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, uint64(value))
	} else {
		proc.stack32 = append(proc.stack32, uint32(value))
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

func (machine *Machine) store_value(address uint64, value uint64, valueSize uint64) {
	for i := range valueSize {
		machine.memory[address+i] = (byte(value >> (i*8))) & 0xFF
	}
}


func (machine *Machine) store_inst(proc *MachineProcess) {
	var payload, payload2 uint64
	var bytes, bytes2 []byte
	if proc.bit64 {
		bytes = GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		8,
		)

		bytes2 =GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8)),
		0,
		8,
		)

		payload = binary.BigEndian.Uint64(bytes)
		payload2 = binary.BigEndian.Uint64(bytes2)
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+16))
	} else {
		bytes = GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		4,
		)

		bytes2 =GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8)),
		0,
		4,
		)

		payload = uint64(binary.BigEndian.Uint32(bytes))
		payload2 = uint64(binary.BigEndian.Uint32(bytes2))
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+16))
	}


	valueSize := payload
	address := payload2
	
	var value uint64

	if !proc.bit64 {
		value = uint64(proc.stack32[len(proc.stack32)-1])
	} else {
		value = proc.stack64[len(proc.stack64)-1]
	}
	proc.int_pop()

	machine.store_value(address, value, valueSize)
}

func (machine *Machine) load_inst(proc *MachineProcess) {
	var payload, payload2 uint64
	var bytes, bytes2 []byte
	if proc.bit64 {
		bytes = GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		8,
		)

		bytes2 =GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8)),
		0,
		8,
		)

		payload = binary.BigEndian.Uint64(bytes)
		payload2 = binary.BigEndian.Uint64(bytes2)
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+16))
	} else {
		bytes = GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP)))),
		0,
		4,
		)

		bytes2 =GetBytesFromPointer(
		(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+8)),
		0,
		4,
		)

		payload = uint64(binary.BigEndian.Uint32(bytes))
		payload2 = uint64(binary.BigEndian.Uint32(bytes2))
		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+16))
	}

	valueSize := payload

	address := payload2

	value := machine.loadValueFromMemory(address, valueSize)

	if proc.bit64 {
		proc.stack64 = append(proc.stack64, uint64(value))
	} else {
		proc.stack32 = append(proc.stack32, uint32(value))
	}
}
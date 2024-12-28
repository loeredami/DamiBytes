package main

import (
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

func (machine *Machine) memIncr_inst() {
	var size uint64 = 0
	
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		size = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		size = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	machine.memory = append(machine.memory, make([]byte, size)...)
}

func (machine *Machine) memDec_inst() {
	var size uint64 = 0
	
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		size = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		size = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	machine.memory = machine.memory[:len(machine.memory)-int(size)]
}

func (machine *Machine) free_inst() {
	address := uint64(0)
	size := uint64(0)

	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		address = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]

		size = machine.stack64[len(machine.stack64)-1]
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		address = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]

		size = uint64(machine.stack32[len(machine.stack32)-1])
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}

	machine.stackOpen = true

	for i := address; i < address+size; i++ {
		machine.memory[i] = 0x0000 // 0 but pretty hex formatting
	}
}

func (machine *Machine) ptrHere_inst(payload uint64) {
	value := uintptr(unsafe.Pointer(&machine.memory[payload]))

	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, uint64(value))
	} else {
		machine.stack32 = append(machine.stack32, uint32(value))
	}
	
	machine.stackOpen = true
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


func (machine *Machine) store_inst(payload uint64) {
	addressBits := 44
	if !machine.bit64 {
			addressBits = 21
	}

	// Extract valueSize
	valueSize := (payload >> (addressBits)) & 0x0F 

	// Extract address
	address := payload & ((1 << addressBits) - 1) 
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	
	var value uint64

	if !machine.bit64 {
		value = uint64(machine.stack32[len(machine.stack32)-1])
	} else {
		value = machine.stack64[len(machine.stack64)-1]
	}
	machine.int_pop()
	machine.stackOpen = true

	machine.store_value(address, value, valueSize)
}

func (machine *Machine) load_inst(payload uint64) {
	addressBits := 44
	if !machine.bit64 {
			addressBits = 21
	}

	valueSize := (payload >> (addressBits)) & 0x0F 

	address := payload & ((1 << addressBits) - 1) 
	for {
		if machine.stackOpen {
			break
		}
	}

	machine.stackOpen = false

	value := machine.loadValueFromMemory(address, valueSize)

	if machine.bit64 {
		machine.stack64 = append(machine.stack64, uint64(value))
	} else {
		machine.stack32 = append(machine.stack32, uint32(value))
	}

	machine.stackOpen =true
}
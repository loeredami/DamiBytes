package main

import (
	"unsafe"
)

func (machine *MachineProcess) int_pop() {
	if machine.bit64 {
		machine.stack64 = machine.stack64[:len(machine.stack64)-1]
	} else {
		machine.stack32 = machine.stack32[:len(machine.stack32)-1]
	}
}

func PointerToNextByte(ptr *byte) *byte {
        // Calculate the address of the next byte
        nextBytePtr := unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1)

        // Cast the pointer back to *byte
        return (*byte)(nextBytePtr)
}

func GetBytesFromPointer(ptr *byte, offset int, size int) []byte {
        if size <= 0 || offset < 0 {
                return nil
        }

        // Ensure the pointer is within a valid memory range
        // This check is crucial for safety, but might require more robust logic
        // depending on your specific use case.
        if ptr == nil {
                return nil
        }

        // Calculate the adjusted pointer with the offset
        offsetPtr := unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(offset))

        // Create a slice pointing to the memory region
        slice := (*[1 << 30]byte)(offsetPtr)[:size:size] 
        return slice
}

func (machine *Machine) handle_instruction(address *byte, proc *MachineProcess) {
        inst := uint8(GetBytesFromPointer(address, 0, 1)[0])

        proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+uintptr(1)))

        machine.run_instruction(inst, proc)
}

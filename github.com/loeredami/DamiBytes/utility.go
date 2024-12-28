package main

import (
	"encoding/binary"
	"slices"
	"unsafe"
)

func (machine *Machine) int_pop() {
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

func (machine *Machine) decodePayload(instruction []byte, startBit uint64, bitLength uint64) uint64 {
        // Calculate starting and ending byte indices
        startByte := startBit / 8
        endByte := (startBit + bitLength - 1) / 8

        // Handle out-of-bounds access
        if endByte >= uint64(len(instruction)) {
                panic("decodePayload: Out of bounds access")
        }

        // Calculate bit offsets within the starting and ending bytes
        startBitOffset := startBit % 8
        endBitOffset := (startBit + bitLength - 1) % 8

        // Extract the field value
        var result uint64
        for i := startByte; i <= endByte; i++ {
                if i == startByte {
                        // Extract bits from the starting byte
                        result |= uint64(instruction[i]) >> startBitOffset
                } else if i == endByte {
                        // Extract bits from the ending byte
                        result |= uint64(instruction[i]) << (8 - endBitOffset)
                } else {
                        // Extract all bits from intermediate bytes
                        result |= uint64(instruction[i]) << ((i - startByte) * 8)
                }
        }

        // Create a mask to extract the desired bits
        mask := uint64(0xFF) >> (8 - bitLength) 
        return result & mask
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
	if machine.bit64 {
		full_inst := GetBytesFromPointer(address, 0, 8)
		slices.Reverse(full_inst)
		
		instruction := binary.BigEndian.Uint64(full_inst) & 0xFFFF000000000000

		payload := binary.BigEndian.Uint64(full_inst) & 0x0000FFFFFFFFFFFF

		add := 8

		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+uintptr(add)))
	

		machine.run_instruction(uint16(instruction>>(32+16)), payload, proc)
	} else {
		full_inst := GetBytesFromPointer(address, 0, 4)
		slices.Reverse(full_inst)

		instruction := binary.BigEndian.Uint32(full_inst) & 0xFF000000

		payload := binary.BigEndian.Uint32(full_inst) &  0x00FFFFFF

		add := 4

		proc.programP = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(proc.programP))+uintptr(add)))
	
		machine.run_instruction(uint16(instruction>>(16+8)), uint64(payload), proc)
	}
}

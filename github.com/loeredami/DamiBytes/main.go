package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	// Example: Virtual Machine setup
	machine := makeMachine(1024, 4, 64)

	machine.bit64 = true
	
	// Example syscall: Write to standard output
	message := []byte("Hello, World!\n")
	messageAddr := uintptr(unsafe.Pointer(&message[0]))
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		panic(err)
	}
	proc, err := dll.FindProc("WriteFile")
	if err != nil {
		panic(err)
	}

	syscallNumber := proc.Addr()


	// Push syscall number, argument count, and arguments
	machine.stack64 = append(machine.stack64,
		0,                      // Overlapped (NULL)
		uint64(len(message)),   // Length of message
		uint64(messageAddr),    // Pointer to buffer
		uint64(syscall.Stdout), // Handle (stdout)
		4,                      // Number of arguments
		uint64(syscallNumber),  // Syscall number (WriteFile)
	)

	// Execute syscall
	machine.syscallHandle()

	// Print result
	fmt.Println("Bytes written:", machine.stack64[len(machine.stack64)-1])
}
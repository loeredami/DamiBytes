package main

import (
	"os"
)

func main() {
	// Temporary Hello World Example Code
	dat, err := os.ReadFile("./sdl.da")

	if err != nil {
		panic(err)
	}

	lexer := makeLexer(string(dat))

	tokens := lexer.lex()

	parser := makeParser(tokens)

	parsed := parser.parse()

	err = os.WriteFile("./sdl.dabin", parsed, os.FileMode(int(0777)))

	if err != nil {
		panic(err)
	}

	dat, err = os.ReadFile("./sdl.dabin")

	
	if err != nil {
		panic(err)
	}

	machine := makeMachine(uint64(len(dat)+128), 1, 2)

	machine.streamInProgram(dat, 0)

	machine.run()

	/*
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

		machine.stack64 = append(machine.stack64,
			0,                      // Overlapped (NULL)
			uint64(len(message)),   // Length of message
			uint64(messageAddr),    // Pointer to buffer
			uint64(syscall.Stdout), // Handle (stdout)
			4,                      // Number of arguments
			uint64(syscallNumber),  // Syscall number (WriteFile)
		)
		machine.syscallHandle()

		fmt.Println("Bytes written:", machine.stack64[len(machine.stack64)-1])
	*/
}
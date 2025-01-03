package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type MachineInfo struct {
	ProC int `json:"processor_count"`
	RegC int `json:"registry_count"`
	MemSize int `json:"extra_memory_size"`
}

var machine_info MachineInfo


func checkAndLoadInConfigs() {
	data, err := os.ReadFile("config.json")

	if err != nil {
		data = []byte(
`{
	"processor_count": 2,
	"registry_count": 4,
	"extra_memory_size": 1024
}`)
	}

	err = json.Unmarshal(data, &machine_info)
	
	if err != nil {
		fmt.Println("When loading config", err)
		machine_info = MachineInfo{2, 4, 1024}
	}
}

func buildMachineUsingConfigs(program []byte) *Machine {
	machine := makeMachine(uint64(len(program))+uint64(machine_info.MemSize), uint16(machine_info.ProC), uint64(machine_info.RegC))
	
	machine.streamInProgram(program, 0)

	return machine
}

func loadProgramWithNewMachine(fileName string) *Machine {
	data, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return buildMachineUsingConfigs(data)
}

func compileProgram(fileName, outputPath string) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	lexer := makeLexer(string(data))

	tokens := lexer.lex()

	parser := makeParser(tokens)

	parsed := parser.parse()

	err = os.WriteFile(outputPath, parsed, os.FileMode(uint32(0777)))

	if err != nil {
		panic(err)
	}

	fmt.Println("Compiled " + fileName+ " and exported to " + outputPath)
}

func main() {
	checkAndLoadInConfigs()
	args := os.Args[1:]


	if len(args) > 0 {
		pathSplit := strings.Split(args[0], ".")
		fileEnding := pathSplit[len(pathSplit)-1]  
		if fileEnding == "dabin" {
			loadProgramWithNewMachine(args[0]).run()
		} else if fileEnding == "da" {
			output := strings.Split(args[0], ".")[0]
			output += ".dabin"
			if len(args) > 1 {
				output = args[1]
			}
			compileProgram(args[0], output)
		} else {
			fmt.Println("No valid file ending. Taking [da | dabin].")
		}
	} else {
		if _, err := os.Stat("bios.dabin"); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Could not find a bios.dabin file to execute.")
		} else {
			loadProgramWithNewMachine("bios.dabin").run()
		}
	}
}
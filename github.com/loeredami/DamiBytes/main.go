package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type MachineInfo struct {
	procC uint64 `toml:processor_count`
	registryCount uint64 `toml:registry_count`
	extraMemSize uint64 `toml:extra_memory_size`
}

var machine_info MachineInfo

type MachineConfig struct {
	machineInfo MachineInfo `toml:MachineInfo`
}

func checkAndLoadInConfigs() {
	var conf MachineConfig
	data, err := os.ReadFile("config.toml")

	if err != nil {
		data = []byte(
`[MachineInfo]
processor_count = 2
registry_count = 4
extra_memory_size = 1024
`)
	}

	if _, err := toml.Decode(string(data), &conf); err == nil {
		conf = MachineConfig{machineInfo: MachineInfo{procC: uint64(1), registryCount:  uint64(4), extraMemSize: uint64(1024)}}
		machine_info = conf.machineInfo
	} else {
		fmt.Println("When loading config", err)
		machine_info = MachineInfo{2, 4, 1024}
	}

}

func buildMachineUsingConfigs(program []byte) *Machine {
	machine := makeMachine(uint64(len(program))+machine_info.extraMemSize, uint16(machine_info.procC), machine_info.registryCount)
	
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
		if strings.Split(args[0], ".")[1] == "dabin" {
			loadProgramWithNewMachine(args[0]).run()
		} else if strings.Split(args[0], ".")[1] == "da" {
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
			loadProgramWithNewMachine(args[0]).run()
		}
	}
}
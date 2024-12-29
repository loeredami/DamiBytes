# About
The DaBytes or Dami Bytes Virtual Machine, is not a virtual machine in the traditional sense that emulates all the typical hardware you would actually see on a machine.

Rather it emulates a emulator that emulates a virtual machine, which may sound like an absolutely stupid idea, until you realize this is mostly just a joke of a definition.

This really is just a byte code executer, with it's own low level language, which has the ability to load in SDL libraries.

This makes the compatibility of each byte code file dependent of which DLL files the programmer decides to use.

DaBytes, is a virtual machine, that was designed so new Programming languages can compile into DaBytes byte code first, instead of putting all of the work into compiling into actual Machine Code. 

This can save a lot of time, when learning how to compile into machine code in a simplified level, instead of learning all of these different instruction sets for each architecture.

# Why does this exist?
DaBytes was originally created for a higher level programming language I will be working on next. In fact I encourage for others to do similar things with this virtual virtual machine. Virtual Squared Machine?

Hey listen, I don't have this terminology really figured out, I really just code, ok?

# Technical Stuff, because someone cares
DaBytes has a 64 Bit and 32 Bit mode, it first boots into 32 Bit mode, and you 100% of the time have to manually switch to 64 Nit mode, because the 32 Bit mode essentially can not do anything, as on Windows 11 there is only 64 Bit, and pointers usually always will be in 64 Bit. (Windows 11 btw. I wish support stayed for 10, and I had an proper SSD for dual booting Linux. Please give me money I have my Ko-Fi linked on the side.)

If you somehow get this to run on a 32 Bit machine, you may have to rewrite everything to use 32 Bit, or figure out a way to automatically detect the proper Bit mode to use.

There is a 64Bit Stack, and a 32Bit Stack, in reality they are lists, however they are treated as stacks by the Virtual Squared Machine. In 64 Bit Mode only the 64Bit Stack will be used, and in 32Bit Mode only the 32Bit Stack will be used.

It will try to load in a config.toml, and otherwise try to load it's default settings:
```toml
[MachineInfo]
processor_count = 2
registry_count = 4
extra_memory_size = 1024
```

Both modes have access to the same memory, however how many bytes to read need to be specific each time.

Like any old regular virtual machine, both the program and the "variables" are stored in the same bytes array.

## Byte array?
Yes, I am using a byte array to store all of the instructions, it made it easier to store string data into the programs as well, and link certain memory addresses for the kernel32 dll on windows.

Hold on a f*cking second... if it's kernel **32**, why the f*ck do I have a problem loading in 32Bit f*cking--

*The previous writer was fired, we excuse any inconveniences the older writer may have caused.*
##### -- *We also did not have the budget to rewrite the existing Documentation.* --

# Usage
Without any further arguments the VSquared Machine, will try to load *desperately* `bios.dabin`.

Otherwise the program can only compile one file at a time, because why the hell would you need multiple files? Organization? How about you organize your room, you insignificant speck of dust.

#### -- *Loere Dami Studios is not liable for hurt feelings.* -- 

## Compiling DaAssembly into DaBytes byte Code (or DaBytes Squared Code)
If output path is not given it will replace the `da` with `dabin`. (Not da trashbin.)

#### -- Note: The specific file MUST end with .da 

`./DamiBytes.exe Hello.da <Optional output path>`

## Running dabin files other than bios.

`./DamiBytes.exe Hello.dabin`

# DaAssembly
DaAssembly can be imagined like assembly, but simplified, directly streamed into memory.

Since any data can be interpreted as an instruction you need to be careful, and check every section that is meant for data will be jumped over.

Every part of the programming, can be jumped to or linked to whether it is an instruction or not, using labels.

```sh
bits 64

# Pushes the pointer to main, it will to the stack
ptr main

# Gets the last thing from the stack, treat it as an address, and jump to it, it will automatically be popped
jump

main:
    exit
```

Labels are entirely handled by the compiler, they basically do not exist in the program itself, instead it uses witchcraft to jump to the given memory address.

In other words, the compiler knows how the machine will find the memory address before hand, and insert the appropriate value.

Usually everything will be treated as data, that will be turned into values.

This includes instructions, instructions however will collect their required arguments, to bake it's payload into it's own value.

Strings will be treated as bytes array and therefore will take less space per character, rather numbers taking the amount of bits depending on the bits mode stated above it, or whatever the compiler was set to at the moment.

```sh
bits 64

# Jumping so "Hello, f" will not be treated as an instruction
ptr main
jump

# This label will point to the H of Hello
message: "Hello, from DaBytes!" 0

main:
    exit
```

# More Documentation Soon
As there is not proper documentation yet, a release is also just something for the future, however I am considering a Snapshot release for the new Year 2025.

I am leaving the language mostly undocumented for now, there might be a wiki soon, for all of the instruction codes, instruction names and what not.

The instruction_codes.txt may be a good enough reference for now, however most of them are not documented correctly, so it may not be the correct instruction name, or not showing the bit format required for the payload, as well as not showing the correct arguments for the compiler.

If you want to brute force a program, or some working byte code, a good starting point would be the ext instruction for loading in a dll. Or the func instruction for loading a function from the dll.

# Build
build it in the *folder* `github.com/loeredami/DamiBytes` with `go build .`, you may require to check the `go.mod` for installing some little packages nobody can ever be bothered writing themselves. Except of course those who wrote the package.


package main

import (
	"encoding/binary"
	"slices"
)

type Parser struct {
	tokens        []Token
	idx           uint64
	done          bool
	curToken      Token
	labelMap      map[string]uint64 // Keeps where certain labels are kept.
	bit64         bool
	filled_memory []byte
}

func IsValidInstruction(category string) bool {
    switch category {
	case "add",
		"sub",
		"mul", 
		"div", 
		"mod",
		"addf",
		"subf", 
		"mulf", 
		"divf", 
		"store", 
		"load", 
		"syscall",
		"jump", 
		"comp",
		"if", 
		"incr",
		"decr",
		"and", 
		"or",
		"not",
		"xor",
		"int",
		"continue",
		//case "in": parser.simpleInstNoPayload(0x0018)
		//case "out": parser.simpleInstNoPayload(0x0018)
		"lshift",
		"rshift", 
		"ext",
		"push", 
		"pop",
		"getR", 
		"setR",
		"ptr",
		"free",
		"go",
		"pid",
		"mincr",
		"mdecr", 
		"bits", 
		"data",
		"func",
		"exit": 
	return true
    }
    return false
}

func makeParser(tokens []Token) *Parser {
	parser := &Parser{tokens, 0, false, &TokenWithNumber{0}, map[string]uint64{}, false, []byte{}}

	parser.updateToken()

	return parser
}

func (parser *Parser) updateToken() {
	if parser.idx >= uint64(len(parser.tokens)) {
		parser.done = true
		return
	}

	parser.curToken = parser.tokens[parser.idx]
}

func (parser *Parser) advance() {
	parser.idx += 1
	parser.updateToken()
}

func (parser *Parser) preProcessLabelMap() {
	real_idx := 0
	was_just_changed := false
	for idx, token := range parser.tokens {
		if was_just_changed {
			was_just_changed = false
			continue
		}
		if token.getType() == tokenTypes.Identifier {

			if !IsValidInstruction(token.(*TokenWithString).val) {
				if parser.bit64 {
					real_idx += 8
				} else {
					real_idx += 4
				}

				continue
			}

			if token.(*TokenWithString).val == "bits" {
				if parser.bit64 {
					real_idx += 8
				} else {
					real_idx += 4
				}
				if parser.tokens[idx+1].(*TokenWithNumber).val == 64 {
					parser.bit64 = true
				} else {
					parser.bit64 = false
				}
				was_just_changed = true
			}

			real_idx += 1
			continue
		}
		if token.getType() == tokenTypes.Label {
			if parser.bit64 {
				parser.labelMap[token.(*TokenWithString).val] = uint64(real_idx)
			} else {
				parser.labelMap[token.(*TokenWithString).val] = uint64(real_idx)
			}
			continue
		}
		if token.getType() == tokenTypes.String {
			real_idx += len(token.(*TokenWithString).val)
			continue
		}
		if parser.bit64 {
			real_idx += 8
		} else {
			real_idx += 4
		}
	}
	parser.bit64 = false
}

func (parser *Parser) handleRawString() {
	val := parser.curToken.(*TokenWithString).val
	bytes := []byte(val)
	parser.filled_memory = append(parser.filled_memory, bytes...)
	parser.advance()
}

func (parser *Parser) handleRawNumber() {
	val := parser.curToken.(*TokenWithNumber).val
	var bytes []byte

	if parser.bit64 {
		bytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, val)
	} else {
		bytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, uint32(val))
	}

	slices.Reverse(bytes)

	parser.filled_memory = append(parser.filled_memory, bytes...)
	parser.advance()
}

func (parser *Parser) quicklyPutInst(inst uint8) {
	parser.filled_memory = append(parser.filled_memory, inst)
}


func (parser *Parser) simpleInstNoPayload(inst uint8) {
	parser.quicklyPutInst(inst)
	parser.advance()
}


func (parser *Parser) handleInstruction() {
	switch parser.curToken.(*TokenWithString).val {
	case "add": parser.simpleInstNoPayload(0x01)
	case "sub": parser.simpleInstNoPayload(0x02)
	case "mul": parser.simpleInstNoPayload(0x03)
	case "div": parser.simpleInstNoPayload(0x04)
	case "mod": parser.simpleInstNoPayload(0x05)
	case "addf": parser.simpleInstNoPayload(0x06)
	case "subf": parser.simpleInstNoPayload(0x07)
	case "mulf": parser.simpleInstNoPayload(0x08)
	case "divf": parser.simpleInstNoPayload(0x09)
	case "store": parser.simpleInstNoPayload(0x0A)
	case "load": parser.simpleInstNoPayload(0x0B)
	case "syscall": parser.simpleInstNoPayload(0x0C)
	case "jump": parser.simpleInstNoPayload(0x0D)
	case "comp": parser.simpleInstNoPayload(0x0E)
	case "if": parser.simpleInstNoPayload(0x0F)
	case "incr": parser.simpleInstNoPayload(0x10)
	case "decr": parser.simpleInstNoPayload(0x11)
	case "and": parser.simpleInstNoPayload(0x12)
	case "or": parser.simpleInstNoPayload(0x13)
	case "not": parser.simpleInstNoPayload(0x14)
	case "xor": parser.simpleInstNoPayload(0x15)
	case "int": parser.simpleInstNoPayload(0x16)
	case "continue": parser.simpleInstNoPayload(0x17)
	//case "in": parser.simpleInstNoPayload(0x0018)
	//case "out": parser.simpleInstNoPayload(0x0018)
	case "lshift": parser.simpleInstNoPayload(0x1A)
	case "rshift": parser.simpleInstNoPayload(0x1B)
	case "ext": parser.simpleInstNoPayload(0x1C)
	case "push": parser.simpleInstNoPayload(0x1D)
	case "pop": parser.simpleInstNoPayload(0x01E)
	case "getR": parser.simpleInstNoPayload(0x1F)
	case "setR": parser.simpleInstNoPayload(0x20)
	case "ptr": parser.simpleInstNoPayload(0x21)
	case "free": parser.simpleInstNoPayload(0x22)
	case "go": parser.simpleInstNoPayload(0x23)
	case "pid": parser.simpleInstNoPayload(0x24)
	case "mincr": parser.simpleInstNoPayload(0x25)
	case "mdecr": parser.simpleInstNoPayload(0x26)
	case "bits": 
		parser.simpleInstNoPayload(0x27)
		parser.handleRawNumber()
		if parser.tokens[parser.idx-1].(*TokenWithNumber).val == 64{
			parser.bit64 = true
		} else {
			parser.bit64 = false
		}


	case "data": parser.simpleInstNoPayload(0x0028)
	case "func": parser.simpleInstNoPayload(0x0029)
	case "exit": parser.simpleInstNoPayload(0x002A)

	default:
		val := parser.labelMap[parser.curToken.(*TokenWithString).val]
		var bytes []byte

		if parser.bit64 {
			bytes = make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, val)
		} else {
			bytes = make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, uint32(val))
		}

		slices.Reverse(bytes)

		parser.filled_memory = append(parser.filled_memory, bytes...)
		parser.advance()
	}
}

func (parser *Parser) parse() []byte {
	parser.preProcessLabelMap()
	for {
		if parser.done {
			break
		}
		switch parser.curToken.getType() {
		case tokenTypes.Number: parser.handleRawNumber()
		case tokenTypes.Identifier: parser.handleInstruction()
		case tokenTypes.Label: parser.advance()
		case tokenTypes.String: parser.handleRawString()
		}
	}

	return parser.filled_memory
}
package main

import "encoding/binary"

type Parser struct {
	tokens        []Token
	idx           uint64
	done          bool
	curToken      Token
	labelMap      map[string]uint64 // Keeps where certain labels are kept.
	bit64         bool
	filled_memory []byte
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
				continue
			}
		}
		if token.getType() == tokenTypes.Label {
			if parser.bit64 {
				parser.labelMap[token.(*TokenWithString).val] = uint64(real_idx) - 8

			} else {
				parser.labelMap[token.(*TokenWithString).val] = uint64(real_idx) - 4
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

	parser.filled_memory = append(parser.filled_memory, bytes...)
	parser.advance()
}

func (parser *Parser) quicklyPutInst(inst, payload uint64) {
	var in_bytes []byte
	if parser.bit64 {
		in_bytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(in_bytes, (inst<<(8*6))+payload)
	} else {
		in_bytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(in_bytes, uint32((inst<<(8*3))+payload))	
	}

	parser.filled_memory = append(parser.filled_memory, in_bytes...)
}


func (parser *Parser) simpleInstNoPayload(inst uint64) {
	parser.quicklyPutInst(inst, 0x0)
	parser.advance()
}

func (parser *Parser) storeInst() {
        parser.advance()
        size := parser.curToken.(*TokenWithNumber).val 

        // Ensure size is within the valid range
        if size > 15 { // Assuming size is 4 bits
                panic("Invalid size value")
        }

        parser.advance()
        memAddr := parser.labelMap[parser.curToken.(*TokenWithString).val]

        var payload uint64

        // Combine size and memAddr with proper bit shifting and masking
        if parser.bit64 {
                sizeBits := uint64(size) & 0x0F 
                payload = (sizeBits << 44) | memAddr 
        } else {
                sizeBits := uint64(size) & 0x07 
                payload = (sizeBits << 21) | memAddr 
        }

        parser.quicklyPutInst(0x000A, payload) 
        parser.advance()
}

func (parser *Parser) loadInst() {
        parser.advance()
        size := parser.curToken.(*TokenWithNumber).val 

        // Ensure size is within the valid range
        if size > 15 { // Assuming size is 4 bits
                panic("Invalid size value")
        }

        parser.advance()
        memAddr := parser.labelMap[parser.curToken.(*TokenWithString).val]

        var payload uint64

        if parser.bit64 {
                // Use bitwise AND to ensure size fits within 4 bits
                sizeBits := uint64(size) & 0x0F 
                payload = (sizeBits << 44) + memAddr 
        } else {
                // Use bitwise AND to ensure size fits within 3 bits
                sizeBits := uint64(size) & 0x07 
                payload = (sizeBits << 21) + memAddr 
        }

        parser.quicklyPutInst(0x000B, payload) 
        parser.advance()
}

func (parser *Parser) nextIsPayload(inst uint64) {
	parser.advance()
	var payload uint64
	if parser.curToken.getType() == tokenTypes.Identifier {
		payload = parser.labelMap[parser.curToken.(*TokenWithString).val]
	} else {
		payload = parser.curToken.(*TokenWithNumber).val
	}
	parser.advance()

	parser.quicklyPutInst(inst, payload)
}

func (parser *Parser) handleInstruction() {
	switch parser.curToken.(*TokenWithString).val {
	case "add": parser.simpleInstNoPayload(0x0001)
	case "sub": parser.simpleInstNoPayload(0x0002)
	case "mul": parser.simpleInstNoPayload(0x0003)
	case "div": parser.simpleInstNoPayload(0x0004)
	case "mod": parser.simpleInstNoPayload(0x0005)
	case "addf": parser.simpleInstNoPayload(0x0006)
	case "subf": parser.simpleInstNoPayload(0x0007)
	case "mulf": parser.simpleInstNoPayload(0x0008)
	case "divf": parser.simpleInstNoPayload(0x0009)
	case "store": parser.storeInst()
	case "load": parser.loadInst()
	case "syscall": parser.simpleInstNoPayload(0x000C)
	case "jump": parser.simpleInstNoPayload(0x000D)
	case "comp": parser.simpleInstNoPayload(0x000E)
	case "if": parser.nextIsPayload(0x000F)
	case "incr": parser.simpleInstNoPayload(0x0010)
	case "decr": parser.simpleInstNoPayload(0x0011)
	case "and": parser.simpleInstNoPayload(0x0012)
	case "or": parser.simpleInstNoPayload(0x0013)
	case "not": parser.simpleInstNoPayload(0x0014)
	case "xor": parser.simpleInstNoPayload(0x0015)
	case "int": parser.simpleInstNoPayload(0x0016)
	case "continue": parser.simpleInstNoPayload(0x0017)
	//case "in": parser.simpleInstNoPayload(0x0018)
	//case "out": parser.simpleInstNoPayload(0x0018)
	case "lshift": parser.simpleInstNoPayload(0x001A)
	case "rshift": parser.simpleInstNoPayload(0x001B)
	case "ext": parser.simpleInstNoPayload(0x001C)
	case "push": parser.nextIsPayload(0x001D)
	case "pop": parser.simpleInstNoPayload(0x001E)
	case "getR": parser.simpleInstNoPayload(0x001F)
	case "setR": parser.simpleInstNoPayload(0x0020)
	case "ptr": parser.nextIsPayload(0x0021)
	case "free": parser.simpleInstNoPayload(0x0022)
	case "go": parser.simpleInstNoPayload(0x0023)
	case "pid": parser.simpleInstNoPayload(0x0024)
	case "mincr": parser.nextIsPayload(0x0025)
	case "mdecr": parser.nextIsPayload(0x0026)
	case "bits": 
		parser.nextIsPayload(0x0027)
		if parser.tokens[parser.idx-1].(*TokenWithNumber).val == 64{
			parser.bit64 = true
		} else {
			parser.bit64 = false
		}

	case "data": parser.simpleInstNoPayload(0x0028)
	case "exit": parser.simpleInstNoPayload(0x002A)

	default:
		parser.quicklyPutInst(0, 0)
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
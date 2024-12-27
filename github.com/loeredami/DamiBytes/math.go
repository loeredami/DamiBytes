package main

import "unsafe"

func (machine *Machine) inc_instruction() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false

	if machine.bit64 {
		machine.stack64[len(machine.stack64)-1] += 1
	} else {
		machine.stack32[len(machine.stack32)-1] += 1
	}

	machine.stackOpen = true
}

func (machine *Machine) dec_instruction() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false

	if machine.bit64 {
		machine.stack64[len(machine.stack64)-1] -= 1
	} else {
		machine.stack32[len(machine.stack32)-1] -= 1
	}

	machine.stackOpen = true
}

func (machine *Machine) add_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1+num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1+num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) sub_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1-num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1-num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) mul_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1*num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1*num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) div_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1/num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1/num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) mod_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := machine.stack64[len(machine.stack64)-1]
		num2 := machine.stack64[len(machine.stack64)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack64 = append(machine.stack64, num1%num2)
	} else {
		num1 := machine.stack32[len(machine.stack32)-1]
		num2 := machine.stack32[len(machine.stack32)-2]
		machine.int_pop()
		machine.int_pop()

		machine.stack32 = append(machine.stack32, num1%num2)
	}
	machine.stackOpen = true
}

func (machine *Machine) addF_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 + num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 + num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) subF_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 - num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 - num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) mulF_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 * num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 * num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}

func (machine *Machine) divF_inst() {
	for {
		if machine.stackOpen {
			break
		}
	}
	machine.stackOpen = false
	if machine.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&machine.stack64[len(machine.stack64)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 / num2

		machine.stack64 = append(machine.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&machine.stack32[len(machine.stack32)-2]))
		machine.int_pop()
		machine.int_pop()

		res := num1 / num2

		machine.stack32 = append(machine.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
	machine.stackOpen = true
}
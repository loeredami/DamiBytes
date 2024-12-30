package main

import "unsafe"

func (machine *Machine) inc_instruction(proc *MachineProcess) {
	if proc.bit64 {
		proc.stack64[len(proc.stack64)-1] += 1
	} else {
		proc.stack32[len(proc.stack32)-1] += 1
	}
}

func (machine *Machine) dec_instruction(proc *MachineProcess) {
	if proc.bit64 {
		proc.stack64[len(proc.stack64)-1] -= 1
	} else {
		proc.stack32[len(proc.stack32)-1] -= 1
	}
}

func (machine *Machine) add_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := proc.stack64[len(proc.stack64)-1]
		num2 := proc.stack64[len(proc.stack64)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack64 = append(proc.stack64, num1+num2)
	} else {
		num1 := proc.stack32[len(proc.stack32)-1]
		num2 := proc.stack32[len(proc.stack32)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack32 = append(proc.stack32, num1+num2)
	}
}

func (machine *Machine) sub_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := proc.stack64[len(proc.stack64)-1]
		num2 := proc.stack64[len(proc.stack64)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack64 = append(proc.stack64, num1-num2)
	} else {
		num1 := proc.stack32[len(proc.stack32)-1]
		num2 := proc.stack32[len(proc.stack32)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack32 = append(proc.stack32, num1-num2)
	}
}

func (machine *Machine) mul_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := proc.stack64[len(proc.stack64)-1]
		num2 := proc.stack64[len(proc.stack64)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack64 = append(proc.stack64, num1*num2)
	} else {
		num1 := proc.stack32[len(proc.stack32)-1]
		num2 := proc.stack32[len(proc.stack32)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack32 = append(proc.stack32, num1*num2)
	}
}

func (machine *Machine) div_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := proc.stack64[len(proc.stack64)-1]
		num2 := proc.stack64[len(proc.stack64)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack64 = append(proc.stack64, num1/num2)
	} else {
		num1 := proc.stack32[len(proc.stack32)-1]
		num2 := proc.stack32[len(proc.stack32)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack32 = append(proc.stack32, num1/num2)
	}
}

func (machine *Machine) mod_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := proc.stack64[len(proc.stack64)-1]
		num2 := proc.stack64[len(proc.stack64)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack64 = append(proc.stack64, num1%num2)
	} else {
		num1 := proc.stack32[len(proc.stack32)-1]
		num2 := proc.stack32[len(proc.stack32)-2]
		proc.int_pop()
		proc.int_pop()

		proc.stack32 = append(proc.stack32, num1%num2)
	}
}

func (machine *Machine) addF_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 + num2

		proc.stack64 = append(proc.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 + num2

		proc.stack32 = append(proc.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
}

func (machine *Machine) subF_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 - num2

		proc.stack64 = append(proc.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 - num2

		proc.stack32 = append(proc.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
}

func (machine *Machine) mulF_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 * num2

		proc.stack64 = append(proc.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 * num2

		proc.stack32 = append(proc.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
}

func (machine *Machine) divF_inst(proc *MachineProcess) {
	if proc.bit64 {
		num1 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-1]))
		num2 := *(*float64)(unsafe.Pointer(&proc.stack64[len(proc.stack64)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 / num2

		proc.stack64 = append(proc.stack64, *(*uint64)(unsafe.Pointer(&(res))))
	} else {
		num1 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-1]))
		num2 := *(*float32)(unsafe.Pointer(&proc.stack32[len(proc.stack32)-2]))
		proc.int_pop()
		proc.int_pop()

		res := num1 / num2

		proc.stack32 = append(proc.stack32, *(*uint32)(unsafe.Pointer(&(res))))
	}
}
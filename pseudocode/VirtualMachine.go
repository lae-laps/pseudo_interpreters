// Virtual Machine which will run the bytecodes generated by de interpreter

package main

import (
	"fmt"
	"math"
)

type VirtualMachine struct {
	ARCH 		int
	HEAP_SIZE 	int
	STACK_SIZE 	int

	heap	 	[]byte
	stack		[]byte

	interrupt			int
	clockCycles			int
	stackPointer 		int
	instructionPointer 	int
	eax 				byte
	ebx					byte
	ecx					byte
	edx					byte

	bytecodes			[][]byte
	flagIdentifier		byte
	codeFlags			map[byte]int
}

func (self *VirtualMachine) initialize() {
	self.ARCH = 32
	self.HEAP_SIZE = int(math.Pow(2, 8))
	self.STACK_SIZE = int(math.Pow(2, 8))

	self.freeHeap()
	self.freeStack()

	self.interrupt			= 0
	self.clockCycles		= 0
	self.stackPointer 		= 0
	self.instructionPointer = 0
	self.eax 				= 0
	self.ebx 				= 0
	self.ecx 				= 0
	self.edx 				= 0

	self.flagIdentifier		= 0xFF					// Identifier for flags

	self.parseCodeFlags()
}

func (self *VirtualMachine) setBytecodes(bytecodes [][]byte) {
	// TODO: check if they are valid

	self.bytecodes = bytecodes
}

func (self *VirtualMachine) parseCodeFlags() {

	self.codeFlags = make(map[byte]int)

	for i, x := range self.bytecodes {
		if x[0] == self.flagIdentifier {
			if self.bytecodes[i + 1][0] == self.flagIdentifier {
				throwError("doubled flag tag but no identifier provided")
				abort(3)
			}
			self.codeFlags[x[1]] = i
		}
	}
}

func (self *VirtualMachine) setStackPointer(value int) {
	if value >= self.HEAP_SIZE {
		throwError("segmentation fault : stack overflow")
		abort(3)
	} else if value < 0 {
		throwError("stack pointer value cannot be lower than 0")
		abort(3)
	}
	self.stackPointer = value
}

func (self *VirtualMachine) setInstructionPointer(value int) {
	if (value >= len(self.bytecodes)) || (value < 0) {
		throwError("segmentation fault : invalid instruction pointer value")
		abort(3)
	}
	self.instructionPointer = value
}

func (self *VirtualMachine) push(value byte) {
	self.setStackPointer(self.stackPointer + 1)
	self.stack[self.stackPointer] = value
}

func (self *VirtualMachine) pop() byte {
	self.setStackPointer(self.stackPointer - 1)
	return self.stack[self.stackPointer + 1]
}

func (self *VirtualMachine) setReg(reg int, value byte) {
	/* Sets value for a general purpose register
		0 -> EAX
		1 -> EBX
		2 -> ECX
		3 -> EDX
	*/

	if int(value) > int(math.Pow(2, float64(self.ARCH))) {
		throwError("value for register exceeds architecture size")
		abort(3)
	}

	switch reg {
	case 0:
		self.eax = value
		break
	case 1:
		self.ebx = value
		break
	case 2:
		self.ecx = value
		break
	case 3:
		self.edx = value
		break
	default:
		throwError("invalid register : expected a number between 0 and 3")
		abort(3)
	}
}

func (self *VirtualMachine) freeHeap() {
	self.heap = make([]byte, self.HEAP_SIZE)
}

func (self *VirtualMachine) freeStack() {
	self.stack = make([]byte, self.STACK_SIZE)
}

func (self *VirtualMachine) run() {
	
	if len(self.bytecodes) == 0 {
		throwError("no bytecodes to run")
		abort(3)
	}

	self.setStackPointer(0)
	self.setInstructionPointer(0)

	self.parseCodeFlags()

	for self.interrupt == 0 {
		if self.clockCycles == 0 {self.nextCycle()}
		self.clockCycles += 1
		self.setInstructionPointer(self.instructionPointer + 1)
		self.nextCycle()
	}

	debug(fmt.Sprintf("program finished with exit code : %d", self.interrupt))
}

func (self *VirtualMachine) nextCycle() {

	frame := self.bytecodes[self.instructionPointer]

	if len(frame) == 0 {
		return
	}

	if frame[0] == self.flagIdentifier {
		self.setInstructionPointer(self.instructionPointer + 1)
		return
	}

	switch frame[0] {
		// TODO: complete this
	}

}
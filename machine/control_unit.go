package machine

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
)

type ControlUnit struct {
	intCtrl            InterruptionController
	logs               string
	programMemory      []Operation
	instructionPointer int
	tempPointer        int
	instructionReg     Operation
	instructionCounter int
	dataPath           DataPath
	curTick            int
	halted             bool
	ei                 bool
	handlingInterrupt  bool
}

func NewControlUnit(program []Operation, dataPath DataPath) *ControlUnit {
	return &ControlUnit{
		intCtrl:            *NewInterruptionController(),
		logs:               "",
		programMemory:      program,
		instructionPointer: 0,
		instructionCounter: 0,
		dataPath:           dataPath,
		curTick:            0,
		halted:             false,
		ei:                 false,
		handlingInterrupt:  false,
	}
}

type InterruptionController struct {
	interrupt bool
	isrAddr   int
}

func NewInterruptionController() *InterruptionController {
	return &InterruptionController{
		interrupt: false,
		isrAddr:   0,
	}
}

func (cu *ControlUnit) printState() {
	out := fmt.Sprintf("TICK: %3d | IC: %3d | CMD: %4s | ARG: %8d | AC: %8d | DR: %8d | AR: %3d | INT: %t \n",
		cu.curTick, cu.instructionCounter, cu.instructionReg.Cmd, cu.instructionReg.Arg, cu.dataPath.accReg, cu.dataPath.dataReg, cu.dataPath.addressReg, cu.intCtrl.interrupt)
	cu.logs = fmt.Sprintf(cu.logs + out)
	print(out)
}

func sortedKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (cu *ControlUnit) tick() {
	cu.printState()
	cu.curTick += 1
}

func (cu *ControlUnit) incrementIC() {
	cu.instructionCounter++
}

func (cu *ControlUnit) setEI(val bool) {
	cu.ei = val
}

func (cu *ControlUnit) checkExit() {
	if cu.halted {
		os.Exit(0)
	}
}

func (cu *ControlUnit) latchInstructionPointer(val int) {
	cu.instructionPointer = val
}

func (cu *ControlUnit) latchTempPointer() {
	cu.tempPointer = cu.instructionPointer
}

func (cu *ControlUnit) incrementInstructionPointer() {
	cu.instructionPointer++
}

func (cu *ControlUnit) checkInterrupt() {
	sortedBuf := sortedKeys(cu.dataPath.portCtrl.iBuf)
	if len(sortedBuf) > 0 {
		tick := sortedBuf[0]
		char := cu.dataPath.portCtrl.iBuf[tick]
		for addr, port := range cu.dataPath.portCtrl.isrMap {
			if tick <= cu.curTick && port == cu.dataPath.portCtrl.iPort && cu.ei {
				cu.dataPath.portCtrl.interruptionRequest(&cu.intCtrl, addr)
				cu.dataPath.portCtrl.bus = char
				break
			}
		}
	}
}

func (cu *ControlUnit) handleInterrupt() {
	//logrus.Info(cu.dataPath.intCtrl.interrupt, " ", cu.ei, " ", cu.handlingInterrupt)
	if !cu.intCtrl.interrupt || !cu.ei || cu.handlingInterrupt {
		return
	}
	//logrus.Info("here")
	cu.latchTempPointer()
	cu.latchInstructionPointer(cu.intCtrl.isrAddr)
	cu.tick()
	cu.handlingInterrupt = true
}

func (cu *ControlUnit) exitInterrupt() {
	cu.handlingInterrupt = false
	cu.latchInstructionPointer(cu.tempPointer)
	cu.intCtrl.unsetInterruption()

	sortedBuf := sortedKeys(cu.dataPath.portCtrl.iBuf)
	tick := sortedBuf[0]
	delete(cu.dataPath.portCtrl.iBuf, tick)
	cu.tick()
}

func (cu *ControlUnit) instructionFetch() {
	cu.instructionReg = cu.programMemory[cu.instructionPointer]
	cu.tick()
}

func (cu *ControlUnit) operandFetch() {
	if cu.instructionReg.Cmd.EnumIndex() <= 10 {
		arg := cu.instructionReg.Arg
		if cu.instructionReg.AddrMode == DIRECT { // аргумент - константа
			cu.dataPath.latchDataReg(DRir, &cu.instructionReg.Arg)
			cu.tick()
		} else if cu.instructionReg.AddrMode == DEFAULT { // аргумент - адрес операнда
			cu.dataPath.latchAddressReg(arg)
			cu.tick()
			cu.dataPath.latchDataReg(DRmem, nil)

		} else { // аргумент - адрес адреса
			cu.dataPath.latchAddressReg(arg)
			cu.tick()
			cu.dataPath.latchDataReg(DRmem, nil)
			cu.tick()
			cu.dataPath.latchAddressReg(cu.dataPath.dataReg)
			cu.tick()
			cu.dataPath.latchDataReg(DRmem, nil)
			cu.tick()
		}
	}
}

func (cu *ControlUnit) decodeExecuteCFInstruction(operation Operation) bool {
	if operation.Cmd == HLT {
		cu.halted = true
		return true
	}
	if operation.Cmd == JMP {
		cu.instructionPointer = cu.programMemory[operation.Arg].Idx
		cu.tick()
		return true
	}
	if operation.Cmd == JZ {
		if cu.dataPath.zeroFLag {
			cu.instructionPointer = cu.programMemory[operation.Arg].Idx
		} else {
			cu.incrementInstructionPointer()
		}
		cu.tick()
		return true
	}
	if operation.Cmd == JE {
		if cu.dataPath.evenFlag {
			cu.instructionPointer = cu.programMemory[operation.Arg].Idx
		} else {
			cu.incrementInstructionPointer()
		}
		cu.tick()
		return true
	}
	if operation.Cmd == JN {
		if cu.dataPath.negFlag {
			cu.instructionPointer = cu.programMemory[operation.Arg].Idx
		} else {
			cu.incrementInstructionPointer()
		}
		cu.tick()
		return true
	}
	if operation.Cmd == CMP {
		cu.dataPath.latchDataReg(DRir, &operation.Arg)
		cu.tick()
		cu.dataPath.sub()
		cu.dataPath.setFlags()
		cu.incrementInstructionPointer()
		cu.tick()
		return true
	}
	return false
}

func (cu *ControlUnit) decodeExecuteInstruction() {
	cu.instructionFetch() // 1 tick
	if cu.decodeExecuteCFInstruction(cu.instructionReg) {
		cu.incrementIC()
		return
	}
	cu.operandFetch() // 0 or 2 ticks
	opcode := cu.instructionReg.Cmd
	if opcode.EnumIndex() == LD {
		cu.dataPath.latchAcc(cu.dataPath.dataReg)
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == ST {
		cu.dataPath.latchAddressReg(cu.dataPath.dataReg)
		cu.tick()
		cu.dataPath.latchDataReg(DRacc, nil)
		cu.tick()
		cu.dataPath.saveToMemory()
		cu.tick()
	}
	if opcode.EnumIndex() == IN {
		if !cu.handlingInterrupt {
			logrus.Fatal("Incorrect IN usage, interrupts disabled")
		}
		cu.dataPath.in()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == OUT {
		if cu.dataPath.dataReg != cu.dataPath.portCtrl.oPort {
			logrus.Fatal("Incorrect OUT usage, wrong port")
		}
		cu.dataPath.out()
		cu.tick()
	}
	if opcode.EnumIndex() == IRET {
		if !cu.handlingInterrupt {
			logrus.Fatal("Cannot exit interrupt")
		}
		cu.incrementIC()
		cu.exitInterrupt()
		return
	}
	if opcode.EnumIndex() == INC {
		cu.dataPath.inc()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == DEC {
		cu.dataPath.dec()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == ADD {
		cu.dataPath.add()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == SUB {
		cu.dataPath.sub()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == MUL {
		cu.dataPath.mul()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == DIV {
		cu.dataPath.div()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == NEG {
		cu.dataPath.neg()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == EI {
		cu.setEI(true)
		cu.tick()
	}
	if opcode.EnumIndex() == DI {
		cu.setEI(false)
		cu.tick()
	}
	cu.incrementInstructionPointer()
	cu.incrementIC()
	return
}

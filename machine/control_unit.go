package machine

import (
	"csa_3/models"
	"github.com/sirupsen/logrus"
	"os"
)

type ControlUnit struct {
	program            []models.Operation
	instructionPointer int
	instructionReg     models.Operation
	instructionCounter int
	dataPath           DataPath
	curTick            int
	halted             bool
}

func NewControlUnit(program []models.Operation, dataPath DataPath) *ControlUnit {
	return &ControlUnit{
		program:            program,
		instructionPointer: 0,
		instructionCounter: 0,
		dataPath:           dataPath,
		curTick:            0,
	}
}

func (cu *ControlUnit) printState() {
	logrus.Infof("TICK: %3d | IC: %3d | CMD: %4s | ARG: %3d | AC: %3d | DR: %3d | AR: %3d | MEM: %3d",
		cu.curTick, cu.instructionCounter, cu.instructionReg.Cmd, cu.instructionReg.Arg, cu.dataPath.accReg, cu.dataPath.dataReg, cu.dataPath.addressReg, cu.dataPath.dataMem[cu.dataPath.addressReg])
}

func (cu *ControlUnit) tick() {
	cu.printState()
	cu.curTick += 1
}

func (cu *ControlUnit) incrementIC() {
	cu.instructionCounter++
}

func (cu *ControlUnit) checkExit() {
	if cu.halted {
		logrus.Info(cu.dataPath.dataMem)
		os.Exit(0)
	}
}

func (cu *ControlUnit) latchInstructionPointer() {
}

func (cu *ControlUnit) incrementInstructionPointer() {
	cu.instructionPointer++
}

func (cu *ControlUnit) instructionFetch() {
	cu.instructionReg = cu.program[cu.instructionPointer]
	cu.tick()
}

func (cu *ControlUnit) operandFetch() {
	if cu.instructionReg.Arg != 0 {
		arg := cu.instructionReg.Arg
		if cu.instructionReg.Iam {
			cu.dataPath.latchAddressReg(arg)
			cu.tick()
			cu.dataPath.latchDataReg(DRmem, nil)
			cu.tick()
		} else {
			cu.dataPath.latchDataReg(DRir, &cu.instructionReg.Arg)
			cu.tick()
		}
	}
}

func (cu *ControlUnit) decodeExecuteCFInstruction(operation models.Operation) bool {
	if operation.Cmd == models.HLT {
		cu.halted = true
		return true
	}
	if operation.Cmd == models.JMP {
		cu.instructionPointer = cu.program[operation.Arg].Idx
		cu.tick()
		return true
	}
	if operation.Cmd == models.JZ {
		if cu.dataPath.zeroFlag() {
			cu.instructionPointer = cu.program[operation.Arg].Idx
			cu.tick()
			return true
		}
	}
	return false
}

func (cu *ControlUnit) decodeExecuteInstruction() {
	cu.instructionFetch() // 1 tick
	cu.operandFetch()     // 0 or 2 ticks
	if cu.decodeExecuteCFInstruction(cu.instructionReg) {
		return
	}
	opcode := cu.instructionReg.Cmd
	if opcode.EnumIndex() == models.LD {
		cu.dataPath.latchAcc(cu.dataPath.dataReg)
		cu.tick()
	}
	if opcode.EnumIndex() == models.ST {
		cu.dataPath.latchAddressReg(cu.dataPath.dataReg)
		cu.tick()
		cu.dataPath.latchDataReg(DRacc, nil)
		cu.tick()
		cu.dataPath.saveToMemory()
		cu.tick()
	}
	if opcode.EnumIndex() == models.ADD {
		cu.dataPath.add()
		cu.tick()
	}
	if opcode.EnumIndex() == models.SUB {
		cu.dataPath.sub()
		cu.tick()
	}
	if opcode.EnumIndex() == models.MUL {
		cu.dataPath.mul()
		cu.tick()
	}
	if opcode.EnumIndex() == models.DIV {
		cu.dataPath.div()
		cu.tick()
	}
	if opcode.EnumIndex() == models.NEG {
		cu.dataPath.neg()
		cu.tick()
	}
	cu.incrementInstructionPointer()
	return
}

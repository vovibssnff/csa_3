package machine

import (
	"csa_3/models"
	"os"
	"strconv"
)

//
//type Signal int
//
//const (
//	Input Signal = iota
//	Res
//)

type ControlUnit struct {
	program            models.MachineCode
	instructionPointer int
	instructionReg     models.Operation
	dataPath           DataPath
	curTick            int
}

func NewControlUnit(program models.MachineCode, dataPath DataPath) *ControlUnit {
	return &ControlUnit{
		program:            program,
		instructionPointer: 0,
		dataPath:           dataPath,
		curTick:            0,
	}
}

func (cu *ControlUnit) tick() {
	cu.curTick += 1
}

func (cu *ControlUnit) latchInstructionPointer() {

}

func (cu *ControlUnit) instructionFetch() {
	cu.instructionReg = cu.program.Ops[cu.instructionPointer]
	cu.tick()
}

func (cu *ControlUnit) operandFetch() {
	if cu.instructionReg.Arg != "" {
		arg, _ := strconv.Atoi(cu.instructionReg.Arg)
		cu.dataPath.latchAddressReg(arg)
		cu.tick()
		cu.dataPath.latchDataReg(DRmem)
		cu.tick()
	}
}

func (cu *ControlUnit) decodeExecuteCFInstruction(operation models.Operation) bool {
	if operation.Cmd == "HLT" {
		os.Exit(0)
	}
	if operation.Cmd == "JMP" {
		addr, _ := strconv.Atoi(operation.Arg)
		cu.instructionPointer = cu.program.Ops[addr].Inx
		cu.tick()
		return true
	}
	if operation.Cmd == "JZ" {
		if cu.dataPath.zeroFlag() {
			addr, _ := strconv.Atoi(operation.Arg)
			cu.instructionPointer = cu.program.Ops[addr].Inx
			cu.tick()
			return true
		}
	}
	return false
}

func (cu *ControlUnit) decodeExecuteInstruction() {
	cu.instructionFetch() // 1 tick
	if cu.decodeExecuteCFInstruction(cu.instructionReg) {
		return
	}
	cu.operandFetch() // 0 or 2 ticks
	opcode := cu.instructionReg.Cmd
	if opcode == "LD" {
		//TODO переписать (нельзя передавать DR тут)
		cu.dataPath.latchAcc(cu.dataPath.dataReg)
		cu.tick()
	}
	if opcode == "ST" {
		cu.dataPath.latchAddressReg(cu.dataPath.dataReg)
		cu.tick()
		cu.dataPath.latchDataReg(DRacc)
		cu.tick()
		cu.dataPath.saveToMemory()
		cu.tick()
	}
	if opcode == "ADD" {
		cu.dataPath.add()
		cu.tick()
	}
	if opcode == "SUB" {
		cu.dataPath.sub()
		cu.tick()
	}
	if opcode == "MUL" {
		cu.dataPath.mul()
		cu.tick()
	}
	if opcode == "DIV" {
		cu.dataPath.div()
		cu.tick()
	}
	if opcode == "NEG" {
		cu.dataPath.neg()
		cu.tick()
	}
	if op

}

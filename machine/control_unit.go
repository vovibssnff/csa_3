package machine

import (
	"csa_3/models"
	"github.com/sirupsen/logrus"
	"os"
)

type ControlUnit struct {
	program            []models.Operation
	instructionPointer int
	tempPointer        int
	instructionReg     models.Operation
	instructionCounter int
	dataPath           DataPath
	curTick            int
	halted             bool
	ei                 bool
	handlingInterrupt  bool
}

func NewControlUnit(program []models.Operation, dataPath DataPath) *ControlUnit {
	return &ControlUnit{
		program:            program,
		instructionPointer: 0,
		instructionCounter: 0,
		dataPath:           dataPath,
		curTick:            0,
		halted:             false,
		ei:                 false,
		handlingInterrupt:  false,
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

func (cu *ControlUnit) setEI(val bool) {
	cu.ei = val
}

func (cu *ControlUnit) checkExit() {
	if cu.halted {
		logrus.Info(cu.dataPath.dataMem)
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
	for tick, char := range cu.dataPath.portCtrl.iBuf {
		for addr, port := range cu.dataPath.portCtrl.isrMap {
			if tick == cu.curTick && port == cu.dataPath.portCtrl.iPort {
				cu.dataPath.portCtrl.interruptionRequest(cu.dataPath.intCtrl, addr)
				cu.dataPath.portCtrl.bus = char
			}
		}
	}
}

func (cu *ControlUnit) handleInterrupt() {
	if !cu.dataPath.intCtrl.interrupt || !cu.ei || cu.handlingInterrupt {
		return
	}
	cu.latchTempPointer()
	cu.dataPath.latchBufferReg()
	cu.tick()
	cu.latchInstructionPointer(cu.dataPath.intCtrl.isrAddr)
	cu.tick()
	cu.handlingInterrupt = true
}

func (cu *ControlUnit) exitInterrupt() {

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
		if cu.dataPath.zeroFLag {
			cu.instructionPointer = cu.program[operation.Arg].Idx
			cu.tick()
			return true
		}
	}
	if operation.Cmd == models.CMP {
		cu.dataPath.latchBufferReg()
		cu.dataPath.sub()
		cu.dataPath.setFlags()
		cu.tick()
		cu.dataPath.latchAcc(cu.dataPath.bufferReg)
		return true
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
		cu.dataPath.setFlags()
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
	//if opcode.EnumIndex() == models.ADD
	if opcode.EnumIndex() == models.INC {
		cu.dataPath.inc()
		cu.tick()
	}
	if opcode.EnumIndex() == models.DEC {
		cu.dataPath.dec()
		cu.tick()
	}
	if opcode.EnumIndex() == models.ADD {
		cu.dataPath.add()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == models.SUB {
		cu.dataPath.sub()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == models.MUL {
		cu.dataPath.mul()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == models.DIV {
		cu.dataPath.div()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == models.NEG {
		cu.dataPath.neg()
		cu.dataPath.setFlags()
		cu.tick()
	}
	if opcode.EnumIndex() == models.EI {
		cu.setEI(true)
		cu.tick()
	}
	if opcode.EnumIndex() == models.DI {
		cu.setEI(false)
		cu.tick()
	}
	cu.incrementInstructionPointer()
	return
}

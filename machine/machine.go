package machine

import (
	"bufio"
	"csa_3/models"
	"csa_3/translator"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type DataPath struct {
	dataMemSize     int
	dataMem         []string
	addressRegister int
	accRegister     string
	inputBuffer     []string
	outputBuffer    []string
}

type ControlUnit struct {
	program            models.MachineCode
	instructionPointer int
	dataPath           DataPath
	curTick            int
}

func NewDataPath(dataMemSize int, inputBuffer []string) *DataPath {
	return &DataPath{
		dataMemSize:     dataMemSize,
		dataMem:         make([]string, 0),
		addressRegister: 0,
		accRegister:     "",
		inputBuffer:     inputBuffer,
		outputBuffer:    make([]string, 0),
	}
}

func NewControlUnit(program models.MachineCode, dataPath DataPath) *ControlUnit {
	return &ControlUnit{
		program:            program,
		instructionPointer: 0,
		dataPath:           dataPath,
		curTick:            0,
	}
}

func (dp *DataPath) latchAcc() {
	dp.accRegister = dp.dataMem[dp.addressRegister]
}

func (dp *DataPath) wr() {

}

func (dp *DataPath) output() {
	symbol := dp.accRegister
	logrus.Debug("output: ", dp.outputBuffer, symbol)
	dp.outputBuffer = append(dp.outputBuffer, symbol)
}

func (dp *DataPath) zero() bool {
	res, _ := strconv.Atoi(dp.accRegister)
	return res == 0
}

func (cu *ControlUnit) tick() {
	cu.curTick += 1
}

func (cu *ControlUnit) latchInstructionPointer() {

}

func (cu *ControlUnit) decodeExecuteCFInstruction(operation models.Operation) bool {
	if operation.Cmd == "HLT" {
		os.Exit(0)
	}
	if operation.Cmd == "JMP" {
		cu.instructionPointer, _ = strconv.Atoi(operation.Arg)
		cu.tick()
		return true
	}
	if operation.Cmd == "JZ" {

	}
}

func simulaton(code models.MachineCode, token []string, dataMemSize int, limit int) (string, int, int) {
	dp := NewDataPath(dataMemSize, token)
	cu := NewControlUnit(code, *dp)

}

func main(i string, input string) {
	code, err := translator.Parse(i)
	if err != nil {
		logrus.Fatal(err)
	}
	inputFile, err := os.Open(input)
	if err != nil {
		logrus.Fatal(err)
	}
	defer inputFile.Close()
	var token []string
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		token = append(token, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}
	output, instrCounter, ticks := simulation(
		code,
		token,
	)
}

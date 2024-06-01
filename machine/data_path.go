package machine

import (
	"github.com/sirupsen/logrus"
)

type DataPath struct {
	dataMemSize  int
	dataMem      []int
	addressReg   int
	dataReg      int
	accReg       int
	inputBuffer  string
	outputBuffer string
}

type DRmux int

const (
	DRmem DRmux = iota
	DRacc
)

//type ALUctrl int
//
//const (
//	PASS
//)

func (mux DRmux) String() string {
	return [...]string{"DRmem", "DRacc"}[mux-1]
}

func NewDataPath(dataMemSize int, inputBuffer string) *DataPath {
	return &DataPath{
		dataMemSize:  dataMemSize,
		dataMem:      make([]int, 0),
		addressReg:   0,
		dataReg:      0,
		accReg:       0,
		inputBuffer:  inputBuffer,
		outputBuffer: "",
	}
}

func (dp *DataPath) latchAddressReg(addr int) {
	dp.addressReg = addr
	if dp.addressReg < 0 || dp.addressReg >= dp.dataMemSize {
		logrus.Fatal("Addr ", dp.addressReg, " is out of memory bounds")
	}
}

func (dp *DataPath) latchDataReg(sel DRmux) {
	if sel == DRmem {
		dp.dataReg = dp.dataMem[dp.addressReg]
	}
	if sel == DRacc {
		dp.dataReg = dp.accReg
	}
}

func (dp *DataPath) saveToMemory() {
	dp.dataMem[dp.addressReg] = dp.dataReg
}

func (dp *DataPath) latchAcc(val int) {
	dp.accReg = val
}

//func (dp *DataPath) output() {
//	symbol := dp.accReg
//	logrus.Debug("output: ", dp.outputBuffer, symbol)
//	dp.outputBuffer = append(dp.outputBuffer, symbol)
//}

func (dp *DataPath) zeroFlag() bool {
	return dp.accReg == 0
}

func (dp *DataPath) negFlag() bool {
	return dp.accReg < 0
}

func (dp *DataPath) add() {
	dp.latchAcc(dp.accReg + dp.dataReg)
}

func (dp *DataPath) sub() {
	dp.latchAcc(dp.accReg - dp.dataReg)
}

func (dp *DataPath) mul() {
	dp.latchAcc(dp.accReg * dp.dataReg)
}

func (dp *DataPath) div() {
	dp.latchAcc(dp.accReg / dp.dataReg)
}

func (dp *DataPath) neg() {
	dp.latchAcc(-dp.accReg)
}

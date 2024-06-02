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
	DRir
)

//type ALUctrl int
//
//const (
//	PASS
//)

func (mux DRmux) String() string {
	return [...]string{"DRmem", "DRacc"}[mux-1]
}

func NewDataPath(dataMem []int, inputBuffer string) *DataPath {
	return &DataPath{
		dataMemSize:  len(dataMem),
		dataMem:      dataMem,
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

// TODO solve
func (dp *DataPath) latchDataReg(sel DRmux, val *int) {
	switch sel {
	case DRmem:
		dp.dataReg = dp.dataMem[dp.addressReg]
	case DRacc:
		dp.dataReg = dp.accReg
	case DRir:
		dp.dataReg = *val
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

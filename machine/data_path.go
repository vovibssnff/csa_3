package machine

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type DRmux int

const (
	DRmem DRmux = iota
	DRacc
	DRir
)

func (mux DRmux) String() string {
	return [...]string{"DRmem", "DRacc"}[mux-1]
}

type IOPortController struct {
	iPort  int
	oPort  int
	iBuf   map[int]int
	oBuf   *string
	bus    int
	isrMap map[int]int // карта привязки хэндлеров прерываний к портам ву
}

func NewIOPortController(IPORT int, OPORT int, IBUF map[int]int, OBUF *string, ints map[int]int) *IOPortController {
	return &IOPortController{
		iPort:  IPORT,
		oPort:  OPORT,
		iBuf:   IBUF,
		oBuf:   OBUF,
		bus:    0,
		isrMap: ints,
	}
}

type DataPath struct {
	portCtrl IOPortController

	dataMemSize int
	dataMem     []int
	addressReg  int
	dataReg     int
	accReg      int
	zeroFLag    bool
	negFlag     bool
	evenFlag    bool
}

func NewDataPath(dataMem []int, ints map[int]int, inputBuffer map[int]int, out *string) *DataPath {
	return &DataPath{
		portCtrl:    *NewIOPortController(0, 1, inputBuffer, out, ints),
		dataMemSize: len(dataMem),
		dataMem:     dataMem,
		addressReg:  0,
		dataReg:     0,
		accReg:      0,
		zeroFLag:    false,
		negFlag:     false,
		evenFlag:    false,
	}
}

func (p *IOPortController) interruptionRequest(intCtrl *InterruptionController, addr int) {
	intCtrl.generateInterruption(addr)
}

func (dp *DataPath) in() {
	dp.latchAcc(dp.portCtrl.bus)
}

func (dp *DataPath) out() {
	v := fmt.Sprint(*dp.portCtrl.oBuf, string(rune(dp.accReg)))
	*dp.portCtrl.oBuf = v
}

func (ic *InterruptionController) generateInterruption(addr int) {
	ic.interrupt = true
	ic.isrAddr = addr
}

func (ic *InterruptionController) unsetInterruption() {
	ic.interrupt = false
	ic.isrAddr = 0
}

func (dp *DataPath) latchAddressReg(addr int) {
	dp.addressReg = addr
	if dp.addressReg < 0 || dp.addressReg >= dp.dataMemSize {
		logrus.Fatal("Addr ", dp.addressReg, " is out of memory bounds")
	}
}

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

func (dp *DataPath) setFlags() {
	dp.zeroFLag = dp.accReg == 0
	dp.negFlag = dp.accReg < 0
	dp.evenFlag = dp.accReg%2 == 0
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

func (dp *DataPath) inc() {
	dp.latchAcc(dp.accReg + 1)
}

func (dp *DataPath) dec() {
	dp.latchAcc(dp.accReg - 1)
}

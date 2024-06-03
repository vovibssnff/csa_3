package models

type Opcode int

const (
	LD = iota
	ST
	CMP
	JZ
	JMP
	OUT
	MUL
	DIV
	ADD
	SUB
	INC
	DEC
	NEG
	HLT
	IRET
	IN
	EI
	DI
)

var Opcodes = [...]string{
	"LD", "ST",
	"CMP", "JZ", "JMP",
	"OUT", "MUL", "DIV", "ADD", "SUB",
	"INC", "DEC", "NEG", "HLT", "IRET", "IN", "EI", "DI",
}

func (o Opcode) String() string {
	return Opcodes[o]
}

func (o Opcode) EnumIndex() int {
	return int(o)
}

type AddrMode int

const (
	DIRECT AddrMode = iota
	DEFAULT
	RELATIVE
)

func (m AddrMode) EnumIndex() int {
	return int(m)
}

type Operation struct {
	Idx      int      `json:"idx"`
	Cmd      Opcode   `json:"cmd"`
	Arg      int      `json:"arg"`
	AddrMode AddrMode `json:"adr"` // indirect addressing mode
}

type DataMemUnit struct {
	Idx int    `json:"idx"`
	Key string `json:"-"`
	Sec string `json:"-"`
	Val int    `json:"val"`
}

type Section struct {
	Name string
	Idx  int
}

type MachineCode struct {
	Data []DataMemUnit `json:"data"`
	Ints map[int]int   `json:"ints"`
	Ops  []Operation   `json:"code"`
}

type Assembly struct {
	DataSection []DataMemUnit
	Interrupts  map[string]int
	Ops         []string
	Sections    []Section
}

func (o Opcode) ToString() string {
	return [...]string{}[o]
}

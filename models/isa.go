package models

type Opcode int

const (
	//mem
	LD = iota
	ST
	//nav
	CMP
	JZ
	JMP
	HLT
	//io
	EI
	DI
	IN
	OUT
	IRET
	//math
	INC
	DEC
	MUL
	DIV
	ADD
	SUB
	NEG
)

var Opcodes = [...]string{
	"LD", "ST",
	"CMP", "JZ", "JMP", "HLT",
	"EI", "DI", "IN", "OUT", "IRET",
	"INC", "DEC", "MUL", "DIV", "ADD", "SUB", "NEG",
}

func (o Opcode) String() string {
	return Opcodes[o]
}

func (o Opcode) EnumIndex() int {
	return int(o)
}

type Operation struct {
	Idx int    `json:"idx"`
	Cmd Opcode `json:"cmd"`
	Arg int    `json:"arg"`
	Iam bool   `json:"iam"` // indirect addressing mode
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

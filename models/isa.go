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
	INR
	IRET
	OUT
	//math
	INC
	DEC
	MUL
	DIV
	ADD
	SUB
	NEG
)

var Opcodes = [...]string{"LD", "ST", "CMP", "JZ", "JMP", "HLT", "INR", "IRET", "OUT", "INC", "DEC", "MUL", "DIV", "ADD", "SUB",
	"NEG"}

func (o Opcode) String() string {
	return Opcodes[o]
}

func (o Opcode) EnumIndex() int {
	return int(o)
}

type Data map[string]string

type Operation struct {
	Inx int    `json:"inx"`
	Cmd Opcode `json:"cmd"`
	Arg int    `json:"arg"`
	Rel bool   `json:"rel"`
}

type DataMemUnit struct {
	Inx   int    `json:"mem"`
	Key   string `json:"-"`
	Sec   string `json:"-"`
	Value int    `json:"val"`
}

type Section struct {
	Name string
	Inx  int
}

type MachineCode struct {
	Data []DataMemUnit `json:"data"`
	Ops  []Operation   `json:"ops"`
}

type Assembly struct {
	DataSection []DataMemUnit
	Ops         []string
	Sections    []Section
}

func (o Opcode) ToString() string {
	return [...]string{}[o]
}

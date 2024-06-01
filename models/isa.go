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
	IN
	OUTR
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

var Opcodes = [...]string{"right", "left", "mv", "st", "ld", "inc", "dec", "mul", "div", "add", "neg", "cmp", "jz", "jmp", "hlt", "out", "in"}

type Data map[string]string

type Operation struct {
	Inx int    `json:"inx"`
	Cmd string `json:"cmd"`
	Arg string `json:"arg"`
	Lit bool   `json:"lit"`
}

type KeyValuePair struct {
	Inx   int    `json:"mem"`
	Key   string `json:"-"`
	Sec   string `json:"-"`
	Value string `json:"val"`
}

type Section struct {
	Name string
	Inx  int
}

type MachineCode struct {
	Data []KeyValuePair `json:"data"`
	Ops  []Operation    `json:"ops"`
}

type Assembly struct {
	DataSection []KeyValuePair
	Ops         []string
	Sections    []Section
}

func (o Opcode) ToString() string {
	return [...]string{}[o]
}

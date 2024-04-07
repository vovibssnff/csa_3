package models

type Opcode int

const (
	MV = iota
	ST
	LD
	INC
	DEC
	MUL
	DIV
	ADD
	NEG
	CMP
	JZ
	JMP
	HLT
	OUT
)

var Opcodes = [...]string{"mv", "st", "ld", "inc", "dec", "mul", "div", "add", "neg", "cmp", "jz", "jmp", "hlt", "out"}

type Val struct {
	reg bool
	val string
}

type Operation struct {
	inx int
	cmd string
	src Val
	dst Val
}

type MachineCode struct {
	data map[int]int
	ops  []Operation
}

func newVal(isReg bool, val string) *Val {
	return &Val{
		reg: isReg,
		val: val,
	}
}

func newOperation(inx int, op Opcode, src Val, dst Val) *Operation {
	return &Operation{
		inx: inx,
		cmd: op.ToString(),
		src: src,
		dst: dst,
	}
}

func newMachineCode(m map[int]int, ops []Operation) *MachineCode {
	return &MachineCode{
		data: m,
		ops:  ops,
	}
}

func (o Opcode) ToString() string {
	return [...]string{}[o]
}

package machine

type ALU struct {
	left  int
	right int
	out   int
}

func NewALU() *ALU {
	return &ALU{
		left:  0,
		right: 0,
		out:   0,
	}
}

//func (a *ALU) add()

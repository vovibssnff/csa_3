package main

import (
	"csa_3/machine"
	"csa_3/translator"
	"flag"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	//testFlag := flag.Bool("test", false, "test mode")
	translateFlag := flag.Bool("t", false, "translate code from .basm file to .xml machine code file")
	executeFlag := flag.Bool("e", false, "execute code from .json file")
	inputFile := flag.String("in", "", "input file of .basm extension")
	paramFile := flag.String("conf", "", "input params file")
	outputFile := flag.String("out", "", "output file")
	//golden := flag.String("golden", "", "expected file")
	flag.Parse()

	// Dereference the pointers to get the actual flag values
	if *translateFlag {
		translator.Translate(*inputFile, *outputFile)
	}

	if *executeFlag {
		machine.Main(*inputFile, *paramFile, *outputFile)

	}

	//if *testFlag {
	//	translator.Translate(*inputFile, *outputFile)
	//	t := *testing.T
	//	TestGoldenOutput()
	//}
}

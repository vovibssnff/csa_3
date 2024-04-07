package main

import (
	"csa_3/translator"
	"flag"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	translateFlag := flag.Bool("t", false, "translate code from .basm file to .xml machine code file")
	inputFile := flag.String("i", "", "input file of .basm extension")
	outputFile := flag.String("o", "", ".xml output file")
	flag.Parse()

	logrus.Info("Gonna parse other vars")

	// Dereference the pointers to get the actual flag values
	if *translateFlag {
		translator.Translate(*inputFile, *outputFile)
	}
}

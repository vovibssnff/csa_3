package machine

import (
	"bufio"
	"csa_3/models"
	"csa_3/translator"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func simulation(code models.MachineCode, token string, limit int) (string, int, int) {
	var data []int
	for _, i := range code.Data {
		data = append(data, i.Val)
	}
	dp := NewDataPath(data, token)
	cu := NewControlUnit(code.Ops, *dp)
	for cu.instructionCounter < limit {
		cu.decodeExecuteInstruction()
		cu.incrementIC()
		cu.checkExit()
	}
	if cu.instructionCounter >= limit {
		logrus.Fatal("Operation limit exceeded")
	}
	logrus.Info("Output buffer: ", dp.outputBuffer)
	return fmt.Sprint(dp.outputBuffer), cu.instructionCounter, cu.curTick
}

func Main(i string, input string) {
	code, err := translator.Parse(i)
	if err != nil {
		logrus.Fatal(err)
	}
	inputFile, err := os.Open(input)
	if err != nil {
		logrus.Fatal(err)
	}
	defer inputFile.Close()

	var token string
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		token += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}

	output, instrCounter, ticks := simulation(
		*code,
		token,
		1000,
	)
	logrus.Info("Output: ", output)
	logrus.Info("instrCounter: ", instrCounter, "ticks: ", ticks)
}

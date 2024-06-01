package machine

import (
	"bufio"
	"csa_3/models"
	"csa_3/translator"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func simulaton(code models.MachineCode, token string, dataMemSize int, limit int) (string, int, int) {
	dp := NewDataPath(code.Data, token)
	cu := NewControlUnit(code, *dp)
	instrCounter := 0
	for instrCounter < limit {
		cu.decodeExecuteInstruction()
		instrCounter++
	}
	if instrCounter >= limit {
		logrus.Fatal("Limit exceeded")
	}
	logrus.Info("Output buffer: ", dp.outputBuffer)
	return fmt.Sprint(dp.outputBuffer), instrCounter, cu.curTick
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

	output, instrCounter, ticks := simulaton(
		*code,
		token,
		100,
		1000,
	)
	logrus.Info("Output: ", output)
	logrus.Info("instrCounter: ", instrCounter, "ticks: ", ticks)
}

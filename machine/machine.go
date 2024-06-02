package machine

import (
	"bufio"
	"csa_3/models"
	"csa_3/translator"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func fileToTokens(f *os.File) map[int]int {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanRunes)

	// Read the entire file content
	var content string
	for scanner.Scan() {
		content += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}

	pairs := strings.Split(content, "), (")
	dict := make(map[int]int)

	for _, pair := range pairs {
		pair = strings.Trim(pair, "() ")
		parts := strings.Split(pair, ", ")
		if len(parts) != 2 {
			logrus.Fatal("Invalid format in file content")
		}
		number, err := strconv.Atoi(parts[0])
		if err != nil {
			logrus.Fatal(err)
		}
		char := rune(parts[1][1]) // Extracting the character inside single quotes
		dict[number] = int(char)
	}
	return dict
}

func simulation(code models.MachineCode, tokens map[int]int, limit int) (string, int, int) {
	var data []int
	for _, i := range code.Data {
		data = append(data, i.Val)
	}
	dp := NewDataPath(data, code.Ints, tokens)
	cu := NewControlUnit(code.Ops, *dp)
	for cu.instructionCounter < limit {
		cu.checkInterrupt()
		cu.handleInterrupt()
		cu.decodeExecuteInstruction()
		cu.incrementIC()
		cu.checkExit()
	}
	if cu.instructionCounter >= limit {
		logrus.Fatal("Operation limit exceeded")
	}
	//logrus.Info("Output buffer: ", dp.outputBuffer)
	return fmt.Sprint(dp.portCtrl.oBuf), cu.instructionCounter, cu.curTick
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

	tokens := fileToTokens(inputFile)

	output, instrCounter, ticks := simulation(
		*code,
		tokens,
		1000,
	)
	logrus.Info("Output: ", output)
	logrus.Info("instrCounter: ", instrCounter, "ticks: ", ticks)
}

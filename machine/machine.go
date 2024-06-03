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

func parseFileToMap(file *os.File) (map[int]int, error) {
	var AsciiZero uint8 = 48
	result := make(map[int]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "()")
		parts := strings.Split(line, ", ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format")
		}

		key, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid integer key: %v", err)
		}

		value := int(parts[1][1])
		result[key] = value
		if parts[1][1] == AsciiZero {
			result[key] = 0
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return result, nil
}

func simulation(code models.MachineCode, tokens map[int]int, limit int, out *string) (int, int) {
	var data []int
	for _, i := range code.Data {
		data = append(data, i.Val)
	}
	dp := NewDataPath(data, code.Ints, tokens, out)
	cu := NewControlUnit(code.Ops, *dp)
	for cu.instructionCounter < limit {
		cu.decodeExecuteInstruction()
		cu.checkInterrupt()
		cu.handleInterrupt()
		if cu.halted {
			break
		}
	}
	if cu.instructionCounter >= limit {
		logrus.Fatal("Operation limit exceeded")
	}
	logrus.Infof("Output buffer: %v", *dp.portCtrl.oBuf)
	return cu.instructionCounter, cu.curTick
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

	tokens, err := parseFileToMap(inputFile)
	out := ""

	instrCounter, ticks := simulation(
		*code,
		tokens,
		1000,
		&out,
	)
	//logrus.Infof("Output: %v", *output)
	logrus.Info(out)
	logrus.Info("instrCounter: ", instrCounter, " ticks: ", ticks)
}

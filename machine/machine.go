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

func simulation(code models.MachineCode, tokens map[int]int, limit int, out *string) (string, int, int) {
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
	return cu.logs, cu.instructionCounter, cu.curTick
}

func Main(i string, input string, log string) {
	code, err := translator.Parse(i)
	if err != nil {
		logrus.Fatal(err)
	}
	inputFile, _ := os.Open(input)
	defer inputFile.Close()

	logFile, err := os.OpenFile(log, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Fatal("Failed to open log file: ", err)
	}
	defer logFile.Close()

	tokens, _ := parseFileToMap(inputFile)
	out := ""

	logs, instrCounter, ticks := simulation(
		*code,
		tokens,
		1000,
		&out,
	)
	logs = logs + "Output buffer: " + out + "\ninstrCounter: " + strconv.Itoa(instrCounter) + " ticks: " + strconv.Itoa(ticks)

	_, err = logFile.WriteString(logs)
	if err != nil {
		logrus.Fatal("Failed to write to log file: ", err)
	}
	if err := logFile.Sync(); err != nil {
		logrus.Fatal("Failed to sync log file: ", err)
	}
}

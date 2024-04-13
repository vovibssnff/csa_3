package translator

import (
	"csa_3/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func ParseAssemblyCode(filename string) (models.Assembly, error) {
	var dataSection []models.KeyValuePair
	ops := make([]string, 0)
	var sections []models.Section

	content, err := os.ReadFile(filename)
	if err != nil {
		return models.Assembly{}, err
	}

	lines := strings.Split(string(content), "\n")
	var currentSection string
	inx := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ".") {
			currentSection = line
			continue
		}

		if currentSection == ".data" {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"")
				dataSection = append(dataSection, models.KeyValuePair{Inx: inx, Key: key, Value: value})
				inx += 1
			}
		} else if currentSection != "" {
			sections = append(sections, models.Section{Name: currentSection, Inx: inx})
			ops = append(ops, line)
			inx += 1
		}
	}

	return models.Assembly{
		DataSection: dataSection,
		Ops:         ops,
		Sections:    sections,
	}, nil
}

func TranslateAssemblyToMachine(assembly models.Assembly) (models.MachineCode, error) {
	machine := models.MachineCode{
		Data: assembly.DataSection,
		Ops:  make([]models.Operation, len(assembly.Ops)),
	}

	for i, op := range assembly.Ops {
		parts := strings.Fields(op)
		op := parts[0]

		var arg, dev string

		if len(parts) > 1 {
			for _, v := range assembly.DataSection {
				if parts[1] == v.Key {
					arg = strconv.Itoa(v.Inx)
				}
			}
		}
		if len(parts) > 2 {
			for _, v := range assembly.Sections {
				if parts[2] == v.Name {
					arg = strconv.Itoa(v.Inx)
					break
				}
			}

			dev = strings.Trim(parts[1], "#")
		}

		machine.Ops[i] = models.Operation{
			Inx: i,
			Cmd: op,
			Arg: arg,
			Dev: dev,
		}
	}
	return machine, nil
}

func Translate(i string, o string) {
	assembly, err := ParseAssemblyCode(i)
	if err != nil {
		logrus.Error("Error parsing .basm file: ", err)
		return
	}
	logrus.Info(assembly)

	machine, err := TranslateAssemblyToMachine(assembly)
	if err != nil {
		logrus.Error("Error translating assembly to machine code: ", err)
		return
	}

	logrus.Info(machine)

	machineJSON, err := json.MarshalIndent(machine, "", "    ")
	if err != nil {
		logrus.Error("Error marshalling machine code to JSON: ", err)
		return
	}
	err = os.WriteFile(o, machineJSON, 0644)
	if err != nil {
		logrus.Error("Output file error: ", err)
	}
}

package translator

import (
	"csa_3/models"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func ParseOpcode(op string) (models.Opcode, error) {
	for i, v := range models.Opcodes {
		if v == op {
			return models.Opcode(i), nil
		}
	}
	return 0, errors.New("invalid opcode")
}

func ParseAssemblyCode(filename string) (models.Assembly, error) {
	var dataSection []models.DataMemUnit
	ops := make([]string, 0)
	var sections []models.Section

	content, err := os.ReadFile(filename)
	if err != nil {
		return models.Assembly{}, err
	}

	lines := strings.Split(string(content), "\n")
	var currentSection string
	inx := 0
	dataInx := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ".") {
			currentSection = line
			continue
		}

		parts := strings.Split(line, "=")
		if currentSection == ".data" {
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if strings.HasPrefix(value, "\"") && strings.Contains(value, "\"") {
					lastCommaInx := strings.LastIndex(value, ",")
					lit := strings.Trim(strings.TrimSpace(value[:lastCommaInx]), "\"")
					for _, char := range lit {
						dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: int(char)})
						inx += 1
					}
					// Add the null terminator
					dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: 0})
				} else {
					// Store the numeric value
					v, _ := strconv.Atoi(strings.TrimSpace(value))
					dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: v})
				}
				inx += 1
			}
			dataInx = inx
		} else if currentSection != "" {
			parts := strings.SplitN(line, " ", 3)
			sections = append(sections, models.Section{Name: currentSection, Idx: inx})
			if len(parts) > 1 && strings.HasPrefix(parts[1], "\"") && strings.Contains(parts[1], "\"") {
				lit := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				lit = strings.ReplaceAll(lit, `"`, "")
				lit = strings.ReplaceAll(lit, ",", "")
				ops = append(ops, parts[0]+" "+strconv.Itoa(dataInx))
				for _, char := range lit {
					dataSection = append(dataSection, models.DataMemUnit{Idx: dataInx, Key: "", Val: int(char)})
					dataInx += 1
				}
				// Add the null terminator
				dataSection = append(dataSection, models.DataMemUnit{Idx: dataInx, Key: "", Val: 0})
				dataInx += 1
			} else if len(parts) > 1 {
				arg := parts[1]
				if num, err := strconv.Atoi(arg); err == nil {
					// It's a numeric literal
					ops = append(ops, parts[0]+" "+strconv.Itoa(num))
					//dataSection = append(dataSection, models.DataMemUnit{Idx: dataInx, Key: "", Val: num})
					//ops = append(ops, parts[0]+" "+strconv.Itoa(dataInx))
					//dataInx += 1
				} else {
					ops = append(ops, line)
				}
			} else {
				ops = append(ops, line)
			}
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
		op, err := ParseOpcode(parts[0])
		if err != nil {
			logrus.Fatal(err, parts[0])
		}
		var arg int
		rel := false

		// arg commands
		if len(parts) > 1 {

			// literal check
			a, err := strconv.Atoi(parts[1])
			if parts[1][0] == '"' || err == nil {
				arg = a
			}

			// section check
			if parts[1][0] == '.' || err == nil {
				for i, sec := range assembly.Sections {
					//logrus.Info(parts[1][1 : len(parts[1])-1])
					if parts[1] == sec.Name {
						arg = i
						break
					}
				}
			}

			// relative addr check
			if parts[1][0] == '(' && parts[1][len(parts[1])-1] == ')' {
				rel = true
				for _, v := range assembly.DataSection {
					if parts[1][1:len(parts[1])-1] == v.Key {
						arg = v.Idx
					}
				}
			}

			// arg name check
			for _, v := range assembly.DataSection {
				if parts[1] == v.Key {
					arg = v.Idx
				}
			}

		}

		machine.Ops[i] = models.Operation{
			Idx: i,
			Cmd: op,
			Arg: arg,
			Iam: rel,
		}
	}
	return machine, nil
}

func Parse(filename string) (*models.MachineCode, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		logrus.Error("Error reading JSON file: ", err)
		return nil, err
	}

	var machineCode models.MachineCode
	err = json.Unmarshal(fileContent, &machineCode)
	if err != nil {
		logrus.Error("Error unmarshalling JSON to machine code: ", err)
		return nil, err
	}

	return &machineCode, nil
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

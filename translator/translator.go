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
	ints := make(map[string]int)
	addressMap := make(map[string]int)

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

		parts := strings.Split(line, "=")
		if currentSection == ".data" {
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if strings.HasPrefix(value, "\"") && strings.Contains(value, "\"") {
					// Handle string literals
					lastCommaInx := strings.LastIndex(value, ",")
					lit := strings.Trim(strings.TrimSpace(value[:lastCommaInx]), "\"")
					for _, char := range lit {
						dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: int(char)})
						inx += 1
					}
					// Add the null terminator
					dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: 0})
					inx += 1
					addressMap[key] = inx - len(lit) - 1 // Store the starting address of the string
				} else if num, err := strconv.Atoi(value); err == nil {
					// Handle numeric values
					dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: num})
					inx += 1
					addressMap[key] = inx - 1 // Store the address of the numeric value
				} else if addr, ok := addressMap[value]; ok {
					// Handle address of another variable
					dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: key, Val: addr})
					inx += 1
					addressMap[key] = addr // Store the same address as the referenced variable
				}
			}
		} else if currentSection == ".int" {
			parts = strings.Split(parts[0], " ")
			if len(parts) == 2 {
				key, _ := strconv.Atoi(strings.Trim(strings.TrimSpace(parts[0]), "#"))
				value := strings.TrimSpace(parts[1])
				ints[value] = key
			}
		} else if currentSection != "" {
			parts := strings.SplitN(line, " ", 3)
			sections = append(sections, models.Section{Name: currentSection, Idx: inx})
			if len(parts) > 1 && strings.HasPrefix(parts[1], "\"") && strings.Contains(parts[1], "\"") {
				lit := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				lit = strings.ReplaceAll(lit, `"`, "")
				lit = strings.ReplaceAll(lit, ",", "")
				ops = append(ops, parts[0]+" "+strconv.Itoa(inx))
				for _, char := range lit {
					dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: "", Val: int(char)})
					inx += 1
				}
				// Add the null terminator
				dataSection = append(dataSection, models.DataMemUnit{Idx: inx, Key: "", Val: 0})
				inx += 1
			} else if len(parts) > 1 {
				arg := parts[1]
				if num, err := strconv.Atoi(arg); err == nil {
					// It's a numeric literal
					ops = append(ops, parts[0]+" "+strconv.Itoa(num))
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
		Interrupts:  ints,
		Sections:    sections,
	}, nil
}

func TranslateAssemblyToMachine(assembly models.Assembly) (models.MachineCode, error) {
	machine := models.MachineCode{
		Data: assembly.DataSection,
		Ops:  make([]models.Operation, len(assembly.Ops)),
	}
	// section resolving in .int section
	ints := make(map[int]int)
	for key, val := range assembly.Interrupts {
		for i, sec := range assembly.Sections {
			if key == sec.Name {
				ints[i] = val
				break
			}
		}
	}
	machine.Ints = ints

	for i, op := range assembly.Ops {
		parts := strings.Fields(op)
		op, err := ParseOpcode(parts[0])
		if err != nil {
			logrus.Fatal(err, parts[0])
		}
		var arg int
		addrMode := models.DIRECT

		// arg commands
		if len(parts) > 1 {
			arg, _ = strconv.Atoi(strings.Trim(parts[1], "#"))
			// literal check
			addrMode = models.DIRECT
			a, err := strconv.Atoi(parts[1])
			if parts[1][0] == '"' || err == nil {
				arg = a
			}

			// section check
			if parts[1][0] == '.' || err == nil {
				for i, sec := range assembly.Sections {
					if parts[1] == sec.Name {
						addrMode = models.DEFAULT
						arg = i
						break
					}
				}
			}

			// relative addr check
			if parts[1][0] == '(' && parts[1][len(parts[1])-1] == ')' {
				for _, v := range assembly.DataSection {
					if parts[1][1:len(parts[1])-1] == v.Key {
						if parts[0] == "ST" {
							addrMode = models.DEFAULT
						} else {
							addrMode = models.RELATIVE
						}
						arg = v.Idx
					}
				}
			}

			// arg name check
			for _, v := range assembly.DataSection {
				if parts[1] == v.Key {
					if parts[0] == "ST" {
						addrMode = models.DIRECT
					} else {
						addrMode = models.DEFAULT
					}
					arg = v.Idx
				}
			}
		}

		machine.Ops[i] = models.Operation{
			Idx:      i,
			Cmd:      op,
			Arg:      arg,
			AddrMode: addrMode,
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
	logrus.Info("done")
}

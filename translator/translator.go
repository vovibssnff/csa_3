package translator

import (
	"bufio"
	"csa_3/models"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func readFile(i string) ([]string, []string) {
	in, err := os.Open(i)
	if err != nil {
		logrus.Error(err)
		logrus.Exit(1)
	}
	sc := bufio.NewScanner(in)
	var data []string
	var ops []string
	var cur []string
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch line {
		case ".data":
			cur = data
		case ".ops":
			data = cur
			cur = ops
		default:
			cur = append(cur, line)
		}
	}
	if err := sc.Err(); err != nil {
		logrus.Error(err)
		logrus.Exit(1)
	}
	return data, ops
}

func validOpcode(op string) bool {
	for _, a := range models.Opcodes {
		if a == op {
			return true
		}
	}
	return false
}

//func parseData()

func toMachineCode(data []string) {
	for _, line := range data {
		str := strings.Split(line, " ")
		if validOpcode(str[0]) {

		}
	}
}

func Translate(i string, o string) {
	logrus.Info(i, " ", o)
	data, ops := readFile(i)
	logrus.Info(data[1])
}

package translator_test

import (
	"csa_3/machine"
	"gotest.tools/v3/golden"
	"os"
	"testing"
)

func TestTranslateAndSimulate(t *testing.T) {
	testCases := []struct {
		inputFileName   string
		machineFileName string
		configFile      string
		outputFileName  string
		goldenFileName  string
	}{
		{"testdata/hello.basm", "testdata/hello.json", "", "testdata/log.txt", "golden_hello.txt"},
		{"testdata/cat.basm", "testdata/cat.json", "testdata/in_cat.txt", "testdata/log.txt", "golden_cat.txt"},
		{"testdata/hello_user.basm", "testdata/hello_user.json", "testdata/in_hello_user.txt", "testdata/log.txt", "golden_hello_user.txt"},
		{"testdata/prob2.basm", "testdata/prob2.json", "", "testdata/log.txt", "golden_prob2.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.inputFileName, func(t *testing.T) {
			machine.Translate(tc.inputFileName, tc.machineFileName)
			machine.Main(tc.machineFileName, tc.configFile, tc.outputFileName)
			output, err := os.ReadFile(tc.outputFileName)
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}
			golden.Assert(t, string(output), tc.goldenFileName)
		})
	}
}

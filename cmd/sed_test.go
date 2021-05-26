package cmd

import (
	"io"
	"os/exec"
	"testing"
)

type SedTest struct {
	input   string
	output  string
	regex   string
	replace string
}

func TestSedArgs_Cmd(t *testing.T) {
	tests := []SedTest{
		{
			input:   `# Test\nString to replace`,
			output:  `# Test\n## String to replace`,
			regex:   "String to replace",
			replace: "## String to replace",
		},
		{
			input:   "# ToChange/ToKeep",
			output:  "# Changed/ToKeep",
			regex:   `\w+/(\w+)`,
			replace: "Changed/$1",
		},
	}

	for _, tes := range tests {
		cmd := exec.Command("go", "run", "../main.go", "sed", "-r", tes.replace, "-m", tes.regex)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal("Can't open the stdin pipe")
		}

		_, err = stdin.Write([]byte(tes.input))
		if err != nil {
			t.Fatal("Can't write to the stdin pipe")
		}
		stdin.Close()

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal("Can't open the stdout pipe")
		}

		if err = cmd.Start(); err != nil {
			t.Fatal("problem", err)
		}
		out, _ := io.ReadAll(stdout)
		if string(out) != tes.output {
			t.Fatalf("invalid output, wanted %s, got %s", tes.output, string(out))
		}
		stdout.Close()
		cmd.Wait()
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type pattern struct {
	Flags    string   `json:"flags-sed,omitempty"`
	Pattern  string   `json:"pattern,omitempty"`
	Patterns []string `json:"patterns,omitempty"`
	Engine   string   `json:"engine,omitempty"`
}

func main() {

	var listMode bool
	flag.BoolVar(&listMode, "list", false, "list available patterns")

	var dumpMode bool
	flag.BoolVar(&dumpMode, "dump", false, "prints the sed command rather than executing it")

	var replaceValue string
	flag.StringVar(&replaceValue, "replace","" ,"value to replace for FUZZ. By default search empty params")

	var newValue string
	flag.StringVar(&newValue, "new-value","FUZZ" ,"new value to replace.")

	flag.Parse()

	if listMode {
		pats, err := getPatterns()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		fmt.Println(strings.Join(pats, "\n"))
		return
	}

	patName := flag.Arg(0)
	files := flag.Arg(1)
	if files == "" {
		files = "."
	}

	patDir, err := getPatternDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to open user's pattern directory")
		return
	}

	filename := filepath.Join(patDir, patName+".json")
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "no such pattern")
		return
	}
	defer f.Close()

	pat := pattern{}
	dec := json.NewDecoder(f)
	err = dec.Decode(&pat)

	if err != nil {
		fmt.Fprintf(os.Stderr, "pattern file '%s' is malformed: %s\n", filename, err)
		return
	}

	if pat.Pattern == "" {
		// check for multiple patterns
		if len(pat.Patterns) == 0 {
			fmt.Fprintf(os.Stderr, "pattern file '%s' contains no pattern(s)\n", filename)
			return
		}

		pat.Pattern =  ""
		for _,s:= range pat.Patterns{
			pat.Pattern +=  "s/("+string(s)+")"+ replaceValue +"/\\1"+newValue+"/gi; "
		}
	}

	if dumpMode {
		fmt.Printf("sed %v %q %v\n", pat.Flags, pat.Pattern, files)

	} else {
		var cmd *exec.Cmd
		operator := "sed"
		if pat.Engine != "" {
			operator = pat.Engine
		}

		if stdinIsPipe() {
			cmd = exec.Command(operator, pat.Flags, pat.Pattern)
		} else {
			cmd = exec.Command(operator, pat.Flags, pat.Pattern, files)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}

func getPatternDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	path := filepath.Join(usr.HomeDir, ".config/gf")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// .config/gf exists
		return path, nil
	}
	return filepath.Join(usr.HomeDir, ".gf"), nil
}

func getPatterns() ([]string, error) {
	out := []string{}

	patDir, err := getPatternDir()
	if err != nil {
		return out, fmt.Errorf("failed to determine pattern directory: %s", err)
	}
	_ = patDir

	files, err := filepath.Glob(patDir + "/*.json")
	if err != nil {
		return out, err
	}

	for _, f := range files {
		f = f[len(patDir)+1 : len(f)-5]
		out = append(out, f)
	}

	return out, nil
}

func stdinIsPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

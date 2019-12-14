package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var verbose bool

func main() {
	shell := os.Getenv("SHELL")

	if shell == "" {
		log.Fatal("Shell not found, $SHELL is empty")
	}

	var pipefile string
	var fromfile string

	flag.StringVar(&pipefile, "pipe", "", "pipe definition file")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.StringVar(&fromfile, "from", "", "read from file instead of stdin")
	flag.Parse()

	if pipefile == "" {
		flag.Usage()
		os.Exit(1)
	}

	commands := parseFile(pipefile)

	if len(commands) == 0 {
		log.Fatalf("No commands found in %s", pipefile)
	}

	apply(fromfile, commands, shell)
}

func parseFile(path string) []string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("Failed to open pipe file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var commands []string

	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			continue // skip empty line
		}

		if strings.HasPrefix(trimmed, "#") {
			continue // skip comment
		}

		commands = append(commands, trimmed)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from pipe file: %v", err)
	}

	return commands
}

func apply(from string, commands []string, shell string) {
	logInfof("$$ %s", shell)

	var buffer []byte

	if from == "" {
		buffer = readStdin()
	} else {
		buffer = readFile(from)
	}

	logInfof("<< stdin")

	for _, command := range commands {
		logInfof("~~ %s", command)

		// parts := strings.Fields(command)
		// cmd := exec.Command(parts[0], parts[1:]...)

		cmd := exec.Command(shell, "-c", command)

		stdin, err := cmd.StdinPipe()

		if err != nil {
			log.Fatalf("Error accessing stdin of command: %v", err)
		}

		_, err = stdin.Write(buffer)

		if err != nil {
			log.Fatalf("Error piping to command: %v", err)
		}

		err = stdin.Close()

		if err != nil {
			log.Fatalf("Error closing stdin of command: %v", err)
		}

		stdout, err := cmd.StdoutPipe()

		if err != nil {
			log.Fatalf("Error accessing stdout of command: %v", err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatalf("Error starting command: %v", err)
		}

		buffer, err = ioutil.ReadAll(stdout)

		if err != nil {
			log.Fatalf("Error reading from stdout: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			log.Fatalf("!! %s", err)
		}
	}

	logInfof(">> stdout")

	writeStdout(buffer)
}

func readStdin() []byte {
	b, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}

	return b
}

func readFile(filepath string) []byte {
	b, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatalf("Error reading from file: %v", err)
	}

	return b
}

func writeStdout(b []byte) {
	_, err := os.Stdout.Write(b)

	if err != nil {
		log.Fatalf("Error writing to stdout: %v", err)
	}
}

func logInfof(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}

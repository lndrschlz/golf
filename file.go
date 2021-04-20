package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	funcProcessingPrefix = " ### "
)

func processFile(filename string, initialize string, verbose bool) string {
	output := ""
	file, err := os.Open(filename)
	checkFail(err, "File "+filename+" not openable!")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	golf := NewGolf(getGoPath())
	if verbose {
		fmt.Printf("Initialize: '%s'\n", initialize)
	}
	_, err = golf.interpreter.Eval(initialize)
	checkFail(err, "Initialization syntax error: '"+initialize+"'")
	lineCounter := 0
	nextLines := 0
	repeatFunction := ""
	if verbose {
		fmt.Printf("File: %s\n", filename)
	}
	for scanner.Scan() {
		line := scanner.Text()
		lineCounter++
		function, arg := getGoLineFunction(line)
		if len(function) == 0 && len(repeatFunction) > 0 {
			function = repeatFunction
			arg = line
			if nextLines > 0 {
				nextLines -= 1
			}
		} else {
			nextLines = 0
		}
		if len(function) > 0 {
			result, ok, repeat, next, err := golf.eval(function, arg)
			if err == nil {
				if verbose {
					fmt.Printf("[%4.0d] %s\n", lineCounter, line)
				}
			} else {
				processFail(filename, lineCounter, line, function, verbose, err)
			}
			if !repeat && nextLines == 0 && repeatFunction == "" && next > 0 {
				nextLines = int(next)
			}
			if repeat || nextLines > 0 {
				repeatFunction = function
			} else {
				nextLines = 0
				repeatFunction = ""
			}
			if ok {
				line = processResult(result, verbose, line)
			}
		}
		output += line + "\n"
	}
	checkFail(scanner.Err(), "File "+filename+" not readable!")
	return output
}

func processResult(result string, verbose bool, line string) string {
	processed := fmt.Sprintf("%s", result)
	if verbose {
		fmt.Printf("     ➥ %s\n", processed)
	}
	line = fmt.Sprintf("%s", processed)
	return line
}

func processFail(filename string, lineCounter int, line string, function string, verbose bool, err error) {
	if !verbose {
		printTitle()
		fmt.Printf("File: %s\n", filename)
	}
	fmt.Printf("[%4.0d] %s\n", lineCounter, line)
	println()
	fmt.Printf(" Code: %s\n", function)
	checkFail(err, fmt.Sprintf("Check the golf syntax in %s:%d", filename, lineCounter))
}

func getGoPath() string {
	return os.Getenv("GOPATH")
}

func checkFile(filename string) {
	if !FileExists(filename) {
		fail("File not found: " + filename)
	}
}

func saveFile(dst string, content string) {
	out, err := os.Create(dst)
	checkFail(err, "File could not be created: "+dst)
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = out.WriteString(content)
	checkFail(err, "File could not be written: "+dst)
	err = out.Sync()
	checkFail(err, "File could not be persisted: "+dst)
	return
}

func getGoLineFunction(line string) (function string, arg string) {
	function = ""
	prefix := funcProcessingPrefix
	functionIdx := strings.Index(line, prefix)
	if functionIdx >= 0 {
		function = line[functionIdx+len(prefix):]
		function = strings.TrimSpace(function)
		arg = line[0:functionIdx]
	}
	return
}

// Exists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

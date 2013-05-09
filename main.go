package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"regexp"
)

func stringTrim(s string) string {
	return strings.Trim(s, "\n\r \t")
}

func checkLine(line string, lineType string) {
	if lineType == "index" {
		_, err := strconv.ParseInt(line, 10, 64)
		if err != nil { panic("Expected subtitle index, but got: " + line + ": " + err.Error()) }
		return
	}
	
	if lineType == "interval" {
		pieces := strings.Split(line, "-->")
		if len(pieces) != 2 { panic("Expected subtitle interval, but got: " + line) }
		return
	}
	
	if lineType == "text" {
		if len(line) == 0 { panic("Expected subtitle text, but got empty line.") }
		return
	}
}

func filterLine(line string) string {
	var reg *regexp.Regexp
	
	// [Noise] Senator, we're making
	reg, _ = regexp.Compile("\\[.+\\]")
	line = reg.ReplaceAllString(line, "")
	
	// - MATELOT: Very good, Lieutenant.
	// TODO: handle accents
	reg, _ = regexp.Compile("[\\s\\t\\n\\r\\-]*[A-Z\\s0-9]+[\\s\\t\\n\\r\\-]*:")
	line = reg.ReplaceAllString(line, "")
	
	line = stringTrim(line)
	
	// If we removed everything but non-alphabetical characters, return an empty line
	if line == "-" { return "" }
	return line
}

func main() {
	filePath := "example.srt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("Could not read file: " + filePath + ": " + err.Error())
	}
	newline := "\n"
	lines := strings.Split(string(content), newline)
	
	output := ""
	lineType := "start"
	for i := 0; i < len(lines); i++ {
		line := stringTrim(lines[i])
		
		if lineType == "start" {
			if line == "" { continue }
			lineType = "index"
			checkLine(line, lineType)
		} else if lineType == "index" {
			lineType = "interval"
			checkLine(line, lineType)
		} else if lineType == "interval" {
			lineType = "text"
			checkLine(line, lineType)
		} else if lineType == "text" {
			if line == "" {
				lineType = "start"
			}
		}
		
		if lineType == "text" {
			line = filterLine(line)
			if line == "" { continue }
		}
	
		output += line + newline
	}
	
	fmt.Println(output)
}

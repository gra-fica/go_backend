package main

import (
	"os"
	"strings"
	"fmt"
)

type SqlParser struct {
	formats map[string]string
}

func newSqlParser() *SqlParser {
	return &SqlParser{
		formats: make(map[string]string),
	}
}

func (s *SqlParser) addFormat(name string, format string) {
	s.formats[name] = format
}

func (s* SqlParser) AddFromFile(path string) (err error) {
	// Read file
	file_bytes, err := os.ReadFile(path);
	if err != nil {
		return
	}

	file := string(file_bytes)

    // windows bullshit 
    file = strings.Replace(file, "\r", "", -1)

	lines := strings.Split(file, "\n")

	type data struct {
		name   string
		beg    int
		end	  *int
	}

	// Parse file
	// first find all the "-- @format" lines
    // it is the first place where the query starts
	formats := []data{};
	for i, line := range lines {
		if strings.Contains(line, "-- @") {
			// Parse format
			d:= data{
				name: strings.Split(line, "-- @")[1],
				beg: i,
			}

			formats = append(formats, d)
		}
	}


	// Then find the end of each format
    // The start of the format is where another one begins (or a comment)
	for i, f := range formats {
		end := f.beg + 1
		for end < len(lines) {
			if strings.Contains(lines[end], "--") {
				break
			}

			end++
		}

		formats[i].end = &end
	}

	// add each format
	for _, f := range formats {
		// Parse format
		format := ""
		for i := f.beg + 1; i < *f.end; i++ {
			format += lines[i] + " "
		}

		s.addFormat(f.name, format)
	}

	return
}

// implement interface Formatter
func (s *SqlParser) Format(f fmt.State, verb rune) {
	for name, format := range s.formats {
		fmt.Fprintf(f, "Format %s:\n%s\n\n", name, format)
	}
}

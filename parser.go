package main
import (
	"fmt"
	"os"
	"strings"
)

type SqlParser struct {
	formats map[string]Parser
}

type Parser struct {
	format string
	total  int
}

func newParser(format string) (p Parser, err error) {
	p = Parser{
		format: format,
		total: 0,
	}

	correct := 0
	for _, c := range format {
		if c == '{' {
			p.total++
			correct++
		}
		if c == '}' {
			correct--
		}

		if correct != 0 {
			err = fmt.Errorf("Invalid format")
		}
	}

	return
}

func (p *Parser) parse(params []string) (out string, error error) {
	if len(params) != p.total {
		error = fmt.Errorf("Invalid number of parameters");
		return
	}

	index := 0
	inside := false;
	for param := range p.format {
		if p.format[param] == '{' {
			inside = true;
		} else if p.format[param] == '}' {
			inside = false;
		}

		if inside{
			out += string(p.format[param])
		} else {
			out += params[index]
			index++
		}
	}
	return
}

func newSqlParser() *SqlParser {
	return &SqlParser{
		formats: make(map[string]Parser),
	}
}

func (s *SqlParser) addFormat(name string, format string) error {
	p, err := newParser(format)
	if err != nil {
		return err
	}

	s.formats[name] = p
	return nil
}

func (s* SqlParser) addFromFile(path string) (err error) {
	// Read file
	file_bytes, err := os.ReadFile(path);
	if err != nil {
		return
	}

	file := string(file_bytes)
	lines := strings.Split(file, "\n")

	formats:= make(map[int]string, 0)
	for i, line := range lines {
		if strings.Contains(line, "-- @") {
			// Parse format
			format := strings.Split(line, "-- @")[1]
			formats[i] = format
		}
	}

	for format := range formats {
		fmt.Printf("Format: %s\n", formats[format])
	}

	// Parse file
	// first find all the "-- @format" lines


	return
}

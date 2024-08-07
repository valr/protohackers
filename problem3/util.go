package problem3

import (
	"bufio"
	"strings"
)

func readLine(reader *bufio.Reader) (line string, err error) {
	if line, err = reader.ReadString('\n'); err != nil {
		return line, err
	}
	return strings.TrimSuffix(line, "\n"), nil
}

func writeLine(writer *bufio.Writer, line string) (err error) {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	if _, err = writer.WriteString(line); err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return nil
}

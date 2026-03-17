package problem5

import (
	"bufio"
	"strings"
)

func readLine(reader *bufio.Reader) (line string, err error) {
	line, err = reader.ReadString('\n')
	if err != nil {
		return line, err
	}
	return strings.TrimSuffix(line, "\n"), nil
}

func writeLine(writer *bufio.Writer, line string) (err error) {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	_, err = writer.WriteString(line)
	if err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return nil
}

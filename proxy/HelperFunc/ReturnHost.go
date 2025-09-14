package helperfunc

import (
	"bufio"
	"strings"
)

func ReturnHost(reader *bufio.Reader) string {
	for {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			// no data left, stop
			break
		}
		start := strings.Index(line, "Host: ")
		if start != -1 {
			// Return trimmed host value without trailing line breaks
			return strings.TrimSpace(line[start+6:])
		}

		if err != nil {
			// error occurred but no Host found yet, break anyway
			break
		}
	}
	return "dont contain host"
}

package tester

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

type Solver struct {
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
	buffer bytes.Buffer
}

func (s *Solver) Turn(input string, output interface{}) (Turn, error) {
	s.buffer.Reset()
	stdout := io.TeeReader(s.stdout, &s.buffer)
	starterAt := time.Now()
	if _, err := fmt.Fprintln(s.stdin, input); err != nil {
		return Turn{}, fmt.Errorf("write state: %w", err)
	}
	if output != nil {
		if _, err := fmt.Fscan(stdout, output); err != nil {
			return Turn{}, fmt.Errorf("read command from %q: %w", s.buffer.String(), err)
		}
	}
	finishedAt := time.Now()
	turnLogs, err := io.ReadAll(s.stderr)
	return Turn{
		Time:   finishedAt.Sub(starterAt),
		Input:  input,
		Output: strings.Trim(s.buffer.String(), "\n\r"),
		Logs:   strings.Trim(string(turnLogs), "\n\r"),
	}, err
}

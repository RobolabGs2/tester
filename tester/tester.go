package tester

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RobolabGs2/tester/exec"
)

type Report struct {
	Name    string
	Success bool
	Image   string
	Turns   []Turn
	Summary string
}

type TestCase interface {
	Title() string
	Run(solver *Solver) (Report, error)
}

type Turn struct {
	Time   time.Duration
	Info   string
	Input  string
	Output string
	Logs   string
}

type SyncBuffer struct {
	mut sync.Mutex
	buf bytes.Buffer
}

func (r *SyncBuffer) Write(b []byte) (int, error) {
	r.mut.Lock()
	defer r.mut.Unlock()
	return r.buf.Write(b)
}

func (r *SyncBuffer) Read(b []byte) (int, error) {
	r.mut.Lock()
	defer r.mut.Unlock()
	return r.buf.Read(b)
}

func RunTestForSolverBinary(ctx context.Context, solverPath string, test TestCase) (Report, error) {
	solverCmd := exec.MakeCmd(solverPath)
	stdin, err := solverCmd.StdinPipe()
	if err != nil {
		return Report{}, fmt.Errorf("open stdin: %w", err)
	}
	stdout, err := solverCmd.StdoutPipe()
	if err != nil {
		return Report{}, fmt.Errorf("open stdout: %w", err)
	}
	logs := new(SyncBuffer)
	solverCmd.Stderr = logs
	if err := solverCmd.Start(); err != nil {
		return Report{}, fmt.Errorf("can't start: %w", err)
	}
	localCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-localCtx.Done()
		if errors.Is(localCtx.Err(), context.DeadlineExceeded) {
			log.Println("Kill process after timeout", exec.KillProcess(solverCmd))
		}
	}()
	report, err := test.Run(&Solver{stdin: stdin, stdout: stdout, stderr: logs})
	cancel()
	if deadlineExceeded(ctx) {
		report.Summary = "Timeout!"
	}
	if err := exec.KillProcess(solverCmd); !deadlineExceeded(ctx) && err != nil {
		log.Println("Failed to kill process with solver:", err)
	}
	_ = solverCmd.Wait()
	return report, err
}

func deadlineExceeded(ctx context.Context) bool {
	return errors.Is(ctx.Err(), context.DeadlineExceeded)
}

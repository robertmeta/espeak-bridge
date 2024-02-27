package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Worker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	mu     sync.Mutex
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) startProcess() {
	w.mu.Lock()
	cmd := exec.Command("espeak-ng", "--stdin")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("WTF")
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	cmd.Start()

	    go func() {
       io.Copy(os.Stdout, stdout) // Continuously read
    }()



	w.cmd = cmd
	w.stdin = stdin
	w.stdout = stdout
	w.mu.Unlock()
}

func (w *Worker) stopProcess() {
	w.mu.Lock()
	w.cmd.Process.Kill()
	w.stdin.Close()
	w.mu.Unlock()
}

func (w *Worker) restartProcess() {
	w.stopProcess()
	w.startProcess()
}

type WorkerPool struct {
	workers    []*Worker
	currentIdx int
	mu         sync.Mutex
}

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{}
	for i := 0; i < size; i++ {
		worker := NewWorker()
		worker.startProcess() // Start each worker process immediately
		pool.workers = append(pool.workers, worker)
	}
	return pool
}

// Speak uses the current worker to speak the text.
func (p *WorkerPool) Speak(text string) {
	p.mu.Lock()
	currentWorker := p.workers[p.currentIdx]
	p.mu.Unlock()

	currentWorker.stdin.Write([]byte(text + "\n\n"))
}

// Stop the current worker and instantly switch to the next one.
func (p *WorkerPool) StopAndSwitch() {
	p.mu.Lock()

	// Stop the current worker
	currentWorker := p.workers[p.currentIdx]
	currentWorker.stopProcess()

	// Move to the next worker and start it if not already running
	p.currentIdx = (p.currentIdx + 1) % len(p.workers)
	nextWorker := p.workers[p.currentIdx]
	nextWorker.startProcess() // Ensure the next worker is ready

	// Restart the stopped worker in the background
	go currentWorker.restartProcess()
	p.mu.Unlock()
}

func main() {
	pool := NewWorkerPool(3) // Including the main worker and 2 extras.
	time.Sleep(1 * time.Second)
	pool.Speak("First message")
	time.Sleep(1 * time.Second)
	pool.Speak("First part2")
	pool.Speak("First part3")
	pool.Speak("First part4")
	time.Sleep(1 * time.Second)

	pool.StopAndSwitch()
	pool.Speak("Second message after switching")
	time.Sleep(1 * time.Second)

	pool.StopAndSwitch()
	pool.Speak("3rd message after switching")
}

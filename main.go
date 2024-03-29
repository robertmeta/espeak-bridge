package main

import "github.com/robertmeta/espeak-bridge/espeak"
import "github.com/robertmeta/espeak-bridge/tones"
import "time"
import "os"
import "fmt"

func main() {
	tempDir := os.TempDir()
	tonePath, err := tones.GetOrMakeTone(tempDir, 500, 440)
	if err != nil {
		panic(err)
	}
	fmt.Println(tonePath)

	pool := espeak.NewWorkerPool(7)
	pool.Speak("First message")
	pool.Speak("First part2")
	pool.Speak("First part3")
	pool.Speak("First part4")
	time.Sleep(3 * time.Second)

	pool.StopAndSwitch()
	pool.Speak("Second message after switching")
	time.Sleep(1 * time.Second)

	pool.StopAndSwitch()
	pool.Speak("3rd message after switching")
	time.Sleep(1 * time.Second)

	pool.StopAndSwitch()
	pool.Speak("4th message after switching")
	time.Sleep(1 * time.Second)

	pool.StopAndSwitch()
	pool.Speak("5th message after switching")

	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")
	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")
	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")
	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")
	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")
	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")
	pool.StopAndSwitch()
	pool.Speak("shoudln't hear me message after switching")

	pool.StopAndSwitch()
	pool.Speak("final message after switching")
}

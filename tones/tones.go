package tones

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// generateTone generates a tone using SoX and returns the path of the generated file.
// It checks if the file already exists to avoid regeneration.
func GetOrMakeTone(directory string, durationMs int, frequency int) (string, error) {
	// Convert duration from milliseconds to seconds for SoX command
	durationSec := float64(durationMs) / 1000.0


	// Construct the filename
	filename := fmt.Sprintf("%d_%d.ogg", durationMs, frequency)
	filepath := filepath.Join(directory, filename)

	// Check if the file already exists
	if _, err := os.Stat(filepath); err == nil {
		fmt.Println("File already exists:", filepath)
		return filepath, nil // File exists, return path
	} else if !os.IsNotExist(err) {
		return "", err // Some other error occurred
	}

	// File does not exist, generate the tone
	volume := "0.9"
	cmd := exec.Command("sox", "-n", filepath, "synth", fmt.Sprintf("%f",
		durationSec), "sine", fmt.Sprintf("%d", frequency),
		"vol", volume)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	fmt.Println("Generated tone:", filepath)
	return filepath, nil
}

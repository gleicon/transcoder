// Package ffmpeg provides functionality for audio/video processing using FFmpeg.
package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Progress represents the progress of an FFmpeg operation
type Progress struct {
	Percent float64
}

// FFmpeg represents an FFmpeg processor
type FFmpeg struct {
	Cmd *exec.Cmd // For testing purposes
}

// New creates a new FFmpeg processor
func New() (*FFmpeg, error) {
	// Check if ffmpeg is installed
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg command not found. Please install it using: brew install ffmpeg")
	}

	// Initialize without a command - we'll create new commands for each operation
	return &FFmpeg{}, nil
}

// NewWithCmd creates a new FFmpeg processor with the given command
func NewWithCmd(cmd *exec.Cmd) (*FFmpeg, error) {
	if cmd == nil {
		return nil, fmt.Errorf("command is required")
	}
	return &FFmpeg{
		Cmd: cmd,
	}, nil
}

// Close releases the FFmpeg resources (no-op for command-line wrapper)
func (f *FFmpeg) Close() {}

// ExtractAudio extracts audio from a video file
func (f *FFmpeg) ExtractAudio(ctx context.Context, input, output string) error {
	// Validate input file
	if _, err := os.Stat(input); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", input)
		}
		return fmt.Errorf("error checking input file: %v", err)
	}

	// Ensure output directory exists
	if err := EnsureOutputDir(output); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Build ffmpeg command
	args := []string{
		"-i", input,
		"-vn",                  // No video
		"-acodec", "pcm_s16le", // PCM 16-bit
		"-ar", "16000", // 16kHz sample rate
		"-ac", "1", // Mono audio
		"-y", // Overwrite output file
		output,
	}

	// Create a new command for this operation
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract audio: %v", err)
	}

	return nil
}

// ChangeSpeed changes the playback speed of a video file
func (f *FFmpeg) ChangeSpeed(ctx context.Context, input, output string, speed float64) error {
	// Validate speed
	if speed <= 0 {
		return fmt.Errorf("speed must be greater than 0, got %f", speed)
	}

	// Validate input file
	if _, err := os.Stat(input); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", input)
		}
		return fmt.Errorf("error checking input file: %v", err)
	}

	// Ensure output directory exists
	if err := EnsureOutputDir(output); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Build ffmpeg command
	args := []string{
		"-i", input,
		"-filter:v", fmt.Sprintf("setpts=PTS/%f", speed),
		"-y", // Overwrite output file
		output,
	}

	// Create a new command for this operation
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to change video speed: %v", err)
	}

	return nil
}

// EnsureOutputDir ensures the output directory exists
func EnsureOutputDir(output string) error {
	dir := filepath.Dir(output)
	if dir != "." {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

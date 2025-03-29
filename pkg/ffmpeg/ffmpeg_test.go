package ffmpeg

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// mockCommand is a helper function to create a mock command for testing
func mockCommand(t *testing.T, name string, args []string) *exec.Cmd {
	cmd := exec.Command("echo", "mock output")
	cmd.Args = append([]string{name}, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "valid",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New()
			if err != nil {
				t.Fatalf("Failed to create FFmpeg processor: %v", err)
			}
			if f == nil {
				t.Error("Expected non-nil FFmpeg processor")
			}
		})
	}
}

func createTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create a simple WAV file with a 1-second sine wave
	cmd := exec.Command("ffmpeg",
		"-f", "lavfi",
		"-i", "sine=frequency=440:duration=1",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "1",
		"-y",
		path,
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
}

func TestExtractAudio(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		output      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid_input",
			input:   "test.wav",
			output:  "output.wav",
			wantErr: false,
		},
		{
			name:        "invalid_input",
			input:       "nonexistent.wav",
			output:      "output.wav",
			wantErr:     true,
			errContains: "input file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test files
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, tt.input)
			outputPath := filepath.Join(tmpDir, tt.output)

			// Create input file for valid test case
			if tt.name == "valid_input" {
				createTestFile(t, inputPath)
			}

			// Create FFmpeg processor
			ffmpeg, err := New()
			if err != nil {
				t.Fatalf("Failed to create FFmpeg processor: %v", err)
			}
			defer ffmpeg.Close()

			// Run test
			err = ffmpeg.ExtractAudio(context.Background(), inputPath, outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractAudio() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ExtractAudio() error = %v, want error containing %v", err, tt.errContains)
				}
			}

			// Verify output file was created for successful cases
			if !tt.wantErr {
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Output file not created: %v", err)
				}
			}
		})
	}
}

func TestChangeSpeed(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		output      string
		speed       float64
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid_input",
			input:   "test.wav",
			output:  "output.wav",
			speed:   1.5,
			wantErr: false,
		},
		{
			name:        "invalid_input",
			input:       "nonexistent.wav",
			output:      "output.wav",
			speed:       1.5,
			wantErr:     true,
			errContains: "input file not found",
		},
		{
			name:        "zero_speed",
			input:       "test.wav",
			output:      "output.wav",
			speed:       0,
			wantErr:     true,
			errContains: "speed must be greater than 0",
		},
		{
			name:        "negative_speed",
			input:       "test.wav",
			output:      "output.wav",
			speed:       -1,
			wantErr:     true,
			errContains: "speed must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test files
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, tt.input)
			outputPath := filepath.Join(tmpDir, tt.output)

			// Create input file for valid test case
			if tt.name == "valid_input" {
				createTestFile(t, inputPath)
			}

			// Create FFmpeg processor
			ffmpeg, err := New()
			if err != nil {
				t.Fatalf("Failed to create FFmpeg processor: %v", err)
			}
			defer ffmpeg.Close()

			// Run test
			err = ffmpeg.ChangeSpeed(context.Background(), inputPath, outputPath, tt.speed)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeSpeed() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ChangeSpeed() error = %v, want error containing %v", err, tt.errContains)
				}
			}

			// Verify output file was created for successful cases
			if !tt.wantErr {
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Output file not created: %v", err)
				}
			}
		})
	}
}

func TestEnsureOutputDir(t *testing.T) {
	// Test with current directory
	err := EnsureOutputDir("file.txt")
	if err != nil {
		t.Errorf("EnsureOutputDir() error = %v", err)
	}

	// Test with new directory
	testDir := filepath.Join(os.TempDir(), "ffmpeg_test", "output")
	defer os.RemoveAll(filepath.Join(os.TempDir(), "ffmpeg_test"))

	err = EnsureOutputDir(filepath.Join(testDir, "file.txt"))
	if err != nil {
		t.Errorf("EnsureOutputDir() error = %v", err)
	}

	// Check if directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

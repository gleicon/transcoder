package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func TestExtractAudio(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create a temporary test file
	testFile := filepath.Join(tmpDir, "test.mp4")
	if err := os.WriteFile(testFile, []byte("test video data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a temporary output file
	outputFile := filepath.Join(tmpDir, "output.wav")

	tests := []struct {
		name    string
		input   string
		output  string
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   testFile,
			output:  outputFile,
			wantErr: false,
		},
		{
			name:    "invalid input",
			input:   "nonexistent.mp4",
			output:  outputFile,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FFmpeg{}
			f.Cmd = mockCommand(t, "ffmpeg", []string{
				"-i", tt.input,
				"-vn",
				"-acodec", "pcm_s16le",
				"-ar", "16000",
				"-ac", "1",
				"-y",
				tt.output,
			})

			err := f.ExtractAudio(context.Background(), tt.input, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractAudio() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeSpeed(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create a temporary test file
	testFile := filepath.Join(tmpDir, "test.mp4")
	if err := os.WriteFile(testFile, []byte("test video data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a temporary output file
	outputFile := filepath.Join(tmpDir, "output.mp4")

	tests := []struct {
		name    string
		input   string
		output  string
		speed   float64
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   testFile,
			output:  outputFile,
			speed:   2.0,
			wantErr: false,
		},
		{
			name:    "invalid input",
			input:   "nonexistent.mp4",
			output:  outputFile,
			speed:   2.0,
			wantErr: true,
		},
		{
			name:    "zero speed",
			input:   testFile,
			output:  outputFile,
			speed:   0.0,
			wantErr: true,
		},
		{
			name:    "negative speed",
			input:   testFile,
			output:  outputFile,
			speed:   -1.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create FFmpeg processor with mock command
			f := &FFmpeg{}
			f.Cmd = mockCommand(t, "ffmpeg", []string{
				"-i", tt.input,
				"-filter:v", fmt.Sprintf("setpts=PTS/%f", tt.speed),
				"-y",
				tt.output,
			})

			err := f.ChangeSpeed(context.Background(), tt.input, tt.output, tt.speed)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeSpeed() error = %v, wantErr %v", err, tt.wantErr)
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

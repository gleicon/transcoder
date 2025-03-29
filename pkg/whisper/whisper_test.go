package whisper

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
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create a temporary model file
	modelFile := filepath.Join(tmpDir, "mock-model.bin")
	if err := os.WriteFile(modelFile, []byte("mock model data"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				ModelPath: modelFile,
				Device:    "cpu",
				Threads:   4,
				Language:  "auto",
			},
			wantErr: false,
		},
		{
			name: "invalid model path",
			config: Config{
				ModelPath: "nonexistent.bin",
				Device:    "cpu",
				Threads:   4,
				Language:  "auto",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && w == nil {
				t.Error("Expected non-nil Whisper processor")
			}
		})
	}
}

func TestTranscribe(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create a temporary test file
	testFile := filepath.Join(tmpDir, "test.wav")
	if err := os.WriteFile(testFile, []byte("test audio data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a temporary output file
	outputFile := filepath.Join(tmpDir, "output.srt")

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
			input:   "nonexistent.wav",
			output:  outputFile,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Whisper{
				config: DefaultConfig(),
			}
			w.Cmd = mockCommand(t, "whisper-cli", []string{
				"-m", w.config.ModelPath,
				"-osrt",
				"-of", strings.TrimSuffix(tt.output, ".srt"),
				"-f", tt.input,
			})

			err := w.Transcribe(context.Background(), tt.input, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transcribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTranscribeWithTranslation(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create a temporary test file
	testFile := filepath.Join(tmpDir, "test.wav")
	if err := os.WriteFile(testFile, []byte("test audio data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a temporary output file
	outputFile := filepath.Join(tmpDir, "output.srt")

	tests := []struct {
		name       string
		input      string
		output     string
		targetLang string
		wantErr    bool
	}{
		{
			name:       "valid input",
			input:      testFile,
			output:     outputFile,
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "invalid input",
			input:      "nonexistent.wav",
			output:     outputFile,
			targetLang: "en",
			wantErr:    true,
		},
		{
			name:       "empty target language",
			input:      testFile,
			output:     outputFile,
			targetLang: "",
			wantErr:    false,
		},
		{
			name:       "invalid target language",
			input:      testFile,
			output:     outputFile,
			targetLang: "invalid",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Whisper{
				config: DefaultConfig(),
			}
			w.Cmd = mockCommand(t, "whisper-cli", []string{
				"-m", w.config.ModelPath,
				"-osrt",
				"-tr",
				"-of", strings.TrimSuffix(tt.output, ".srt"),
				"-f", tt.input,
			})

			err := w.TranscribeWithTranslation(context.Background(), tt.input, tt.output, tt.targetLang)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranscribeWithTranslation() error = %v, wantErr %v", err, tt.wantErr)
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
	testDir := filepath.Join(os.TempDir(), "whisper_test", "output")
	defer os.RemoveAll(filepath.Join(os.TempDir(), "whisper_test"))

	err = EnsureOutputDir(filepath.Join(testDir, "file.txt"))
	if err != nil {
		t.Errorf("EnsureOutputDir() error = %v", err)
	}

	// Check if directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

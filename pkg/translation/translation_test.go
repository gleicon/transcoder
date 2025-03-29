package translation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func mockCommand(name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		os.Exit(1)
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "ffmpeg":
		// Create output file to simulate successful execution
		if len(args) > 0 {
			var inputFile, outputFile string
			for i, arg := range args {
				if arg == "-i" && i+1 < len(args) {
					inputFile = args[i+1]
				}
				if arg == "-y" && i+1 < len(args) {
					outputFile = args[i+1]
				}
			}
			if inputFile != "" {
				if _, err := os.Stat(inputFile); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "Input file not found: %s\n", inputFile)
					os.Exit(1)
				}
			}
			if outputFile != "" {
				if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
					os.Exit(1)
				}
				// Create the output file
				if err := os.WriteFile(outputFile, []byte("test data"), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
					os.Exit(1)
				}
				// If this is a .wav file, create it in both locations
				if strings.HasSuffix(outputFile, ".wav") {
					// Create the .wav file in the output directory
					if err := os.WriteFile(outputFile, []byte("test wav data"), 0644); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create wav file: %v\n", err)
						os.Exit(1)
					}
					// Also create a copy in the input directory for the next command
					wavFile := filepath.Join(filepath.Dir(inputFile), filepath.Base(inputFile)+".wav")
					if err := os.MkdirAll(filepath.Dir(wavFile), 0755); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
						os.Exit(1)
					}
					if err := os.WriteFile(wavFile, []byte("test wav data"), 0644); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create wav file: %v\n", err)
						os.Exit(1)
					}
					// Also create a copy in the output directory with the input file name
					wavFile = filepath.Join(filepath.Dir(outputFile), filepath.Base(inputFile)+".wav")
					if err := os.MkdirAll(filepath.Dir(wavFile), 0755); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
						os.Exit(1)
					}
					if err := os.WriteFile(wavFile, []byte("test wav data"), 0644); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create wav file: %v\n", err)
						os.Exit(1)
					}
					// Also create a copy in the input directory with the output file name
					wavFile = filepath.Join(filepath.Dir(inputFile), filepath.Base(outputFile))
					if err := os.MkdirAll(filepath.Dir(wavFile), 0755); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
						os.Exit(1)
					}
					if err := os.WriteFile(wavFile, []byte("test wav data"), 0644); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to create wav file: %v\n", err)
						os.Exit(1)
					}
				}
			}
		}
	case "whisper-cli":
		// Create output file to simulate successful execution
		if len(args) > 0 {
			var inputFile, outputFile string
			for i, arg := range args {
				if arg == "-f" && i+1 < len(args) {
					inputFile = args[i+1]
				}
				if arg == "-of" && i+1 < len(args) {
					outputFile = args[i+1]
				}
			}
			if inputFile != "" {
				if _, err := os.Stat(inputFile); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "Input file not found: %s\n", inputFile)
					os.Exit(1)
				}
			}
			if outputFile != "" {
				if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
					os.Exit(1)
				}
				// Create the .srt file
				srtFile := outputFile + ".srt"
				if err := os.WriteFile(srtFile, []byte("1\n00:00:00,000 --> 00:00:05,000\nTest subtitle\n"), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create srt file: %v\n", err)
					os.Exit(1)
				}
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}
	os.Exit(0)
}

func createTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("test data"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		ffmpeg  *exec.Cmd
		whisper *exec.Cmd
		wantErr bool
	}{
		{
			name:    "valid commands",
			ffmpeg:  mockCommand("ffmpeg"),
			whisper: mockCommand("whisper-cli"),
			wantErr: false,
		},
		{
			name:    "nil ffmpeg",
			ffmpeg:  nil,
			whisper: mockCommand("whisper-cli"),
			wantErr: true,
		},
		{
			name:    "nil whisper",
			ffmpeg:  mockCommand("ffmpeg"),
			whisper: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.ffmpeg, tt.whisper)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTranslate(t *testing.T) {
	// Create test files
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.wav")
	outputFile := filepath.Join(tmpDir, "output.wav")
	createTestFile(t, inputFile)

	// Create mock commands
	ffmpegCmd := mockCommand("ffmpeg")
	whisperCmd := mockCommand("whisper-cli")

	// Create translator
	translator, err := New(ffmpegCmd, whisperCmd)
	if err != nil {
		t.Fatalf("Failed to create translator: %v", err)
	}
	defer translator.Close()

	tests := []struct {
		name        string
		input       string
		output      string
		targetLang  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid input",
			input:      inputFile,
			output:     outputFile,
			targetLang: "es",
			wantErr:    false,
		},
		{
			name:        "nonexistent input",
			input:       "nonexistent.wav",
			output:      outputFile,
			targetLang:  "es",
			wantErr:     true,
			errContains: "input file not found",
		},
		{
			name:        "empty target language",
			input:       inputFile,
			output:      outputFile,
			targetLang:  "",
			wantErr:     true,
			errContains: "target language is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := translator.Translate(context.Background(), tt.input, tt.output, tt.targetLang)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Translate() error = %v, want error containing %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestTranslateFile(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		output      string
		targetLang  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid_mp3",
			input:      "input.mp3",
			output:     "output.srt",
			targetLang: "es",
			wantErr:    false,
		},
		{
			name:       "valid_mp4",
			input:      "input.mp4",
			output:     "output.srt",
			targetLang: "es",
			wantErr:    false,
		},
		{
			name:       "nonexistent_input",
			input:      "nonexistent.mp3",
			output:     "output.srt",
			targetLang: "es",
			wantErr:    true,
		},
		{
			name:       "empty_target_language",
			input:      "input.mp3",
			output:     "output.srt",
			targetLang: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			inputPath := filepath.Join(dir, tt.input)
			outputPath := filepath.Join(dir, tt.output)

			// Create input file for valid test cases
			if tt.name == "valid_mp3" || tt.name == "valid_mp4" {
				if err := os.WriteFile(inputPath, []byte("test data"), 0644); err != nil {
					t.Fatalf("Failed to create input file: %v", err)
				}
			}

			// Create mock commands
			ffmpegCmd := mockCommand("ffmpeg")
			whisperCmd := mockCommand("whisper-cli")

			// Create translator
			translator, err := New(ffmpegCmd, whisperCmd)
			if err != nil {
				t.Fatalf("Failed to create translator: %v", err)
			}
			defer translator.Close()

			err = translator.TranslateFile(context.Background(), inputPath, outputPath, tt.targetLang)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("TranslateFile() error = %v, want error containing %v", err, tt.errContains)
				}
			}
		})
	}
}

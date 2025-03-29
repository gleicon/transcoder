package whisper

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config holds the configuration for the Whisper processor
type Config struct {
	ModelPath string
	Device    string // "cpu", "cuda", "metal"
	Threads   int
	Language  string
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}

	// Default model path in the user's home directory
	defaultModel := filepath.Join(homeDir, ".cache", "whisper", "base.bin")

	return Config{
		ModelPath: defaultModel,
		Device:    "cpu",
		Threads:   4,
		Language:  "auto",
	}
}

// Whisper represents a Whisper processor
type Whisper struct {
	config Config
	Cmd    *exec.Cmd // For testing purposes
}

// New creates a new Whisper processor with the given configuration
func New(config Config) (*Whisper, error) {
	// Check if whisper-cli is installed
	if _, err := exec.LookPath("whisper-cli"); err != nil {
		return nil, fmt.Errorf("whisper-cli command not found. Please install it using: brew install whisper-cpp")
	}

	// Validate model path
	if _, err := os.Stat(config.ModelPath); err != nil {
		return nil, fmt.Errorf("model file not found at %s: %v", config.ModelPath, err)
	}

	return &Whisper{
		config: config,
	}, nil
}

// NewWithCmd creates a new Whisper processor with the given command and configuration
func NewWithCmd(cmd *exec.Cmd, config Config) (*Whisper, error) {
	if cmd == nil {
		return nil, fmt.Errorf("command is required")
	}

	// Validate model path
	if _, err := os.Stat(config.ModelPath); err != nil {
		return nil, fmt.Errorf("model file not found at %s: %v", config.ModelPath, err)
	}

	return &Whisper{
		config: config,
		Cmd:    cmd,
	}, nil
}

// Close releases the Whisper resources (no-op for command-line wrapper)
func (w *Whisper) Close() {}

// Transcribe transcribes an audio file to SRT format
func (w *Whisper) Transcribe(ctx context.Context, input, output string) error {
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

	// Build whisper-cli command
	args := []string{
		"-m", w.config.ModelPath,
		"-osrt",
		"-of", strings.TrimSuffix(output, ".srt"),
	}

	if w.config.Language != "" && w.config.Language != "auto" {
		args = append(args, "-l", w.config.Language)
	}

	if w.config.Threads > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", w.config.Threads))
	}

	args = append(args, "-f", input)

	// Create command with context
	if w.Cmd == nil {
		w.Cmd = exec.CommandContext(ctx, "whisper-cli", args...)
		w.Cmd.Stdout = os.Stdout
		w.Cmd.Stderr = os.Stderr
	}

	// Run the command
	if err := w.Cmd.Run(); err != nil {
		return fmt.Errorf("failed to transcribe audio: %v", err)
	}

	return nil
}

// TranscribeWithTranslation transcribes and translates an audio file
func (w *Whisper) TranscribeWithTranslation(ctx context.Context, input, output, targetLang string) error {
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

	// Validate target language if provided
	if targetLang != "" && targetLang != "auto" {
		// List of supported language codes
		supportedLangs := []string{
			"en", "zh", "de", "es", "ru", "ko", "fr", "ja", "pt", "tr", "pl", "ca", "nl", "ar", "sv", "it", "id", "hi", "fi", "vi", "he", "uk", "el", "ms", "cs", "ro", "da", "hu", "ta", "no", "th", "ur", "hr", "bg", "lt", "la", "mi", "ml", "cy", "sk", "te", "fa", "lv", "bn", "sr", "az", "sl", "kn", "et", "mk", "br", "eu", "is", "hy", "ne", "mn", "bs", "kk", "sq", "sw", "gl", "mr", "pa", "si", "km", "sn", "yo", "so", "af", "oc", "ka", "be", "tg", "sd", "gu", "am", "yi", "lo", "uz", "fo", "ht", "ps", "tk", "nn", "mt", "sa", "lb", "my", "bo", "tl", "mg", "as", "tt", "haw", "ln", "ha", "ba", "jw", "su",
		}
		isValid := false
		for _, lang := range supportedLangs {
			if lang == targetLang {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("unsupported language code: %s", targetLang)
		}
	}

	// Build whisper-cli command
	args := []string{
		"-m", w.config.ModelPath,
		"-osrt",
		"-tr", // Enable translation
		"-of", strings.TrimSuffix(output, ".srt"),
	}

	// Only add language if it's not empty
	if targetLang != "" {
		args = append(args, "-l", targetLang)
	}

	if w.config.Threads > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", w.config.Threads))
	}

	args = append(args, "-f", input)

	// Create command with context
	if w.Cmd == nil {
		w.Cmd = exec.CommandContext(ctx, "whisper-cli", args...)
		w.Cmd.Stdout = os.Stdout
		w.Cmd.Stderr = os.Stderr
	}

	// Run the command
	if err := w.Cmd.Run(); err != nil {
		return fmt.Errorf("failed to translate audio: %v", err)
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

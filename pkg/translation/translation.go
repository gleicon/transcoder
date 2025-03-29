package translation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gleicon/transcoder/pkg/ffmpeg"
	"github.com/gleicon/transcoder/pkg/whisper"
)

// Translator represents a translator that can transcribe and translate audio
type Translator struct {
	whisperProcessor *whisper.Whisper
	ffmpegProcessor  *ffmpeg.FFmpeg
}

// New creates a new translator with the given FFmpeg and Whisper commands
func New(ffmpegCmd, whisperCmd *exec.Cmd) (*Translator, error) {
	if ffmpegCmd == nil {
		return nil, fmt.Errorf("FFmpeg command is required")
	}
	if whisperCmd == nil {
		return nil, fmt.Errorf("Whisper command is required")
	}

	// Create FFmpeg processor
	ffmpegProcessor, err := ffmpeg.NewWithCmd(ffmpegCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create FFmpeg processor: %v", err)
	}

	// Create Whisper processor with default config
	whisperProcessor, err := whisper.NewWithCmd(whisperCmd, whisper.DefaultConfig())
	if err != nil {
		ffmpegProcessor.Close()
		return nil, fmt.Errorf("failed to create Whisper processor: %v", err)
	}

	return &Translator{
		whisperProcessor: whisperProcessor,
		ffmpegProcessor:  ffmpegProcessor,
	}, nil
}

// NewWithCmd creates a new translator with default commands
func NewWithCmd() (*Translator, error) {
	ffmpegCmd := exec.Command("ffmpeg")
	whisperCmd := exec.Command("whisper-cli")
	return New(ffmpegCmd, whisperCmd)
}

// Close releases the translator's resources
func (t *Translator) Close() {
	if t.whisperProcessor != nil {
		t.whisperProcessor.Close()
	}
	if t.ffmpegProcessor != nil {
		t.ffmpegProcessor.Close()
	}
}

// Translate transcribes and translates an audio file
func (t *Translator) Translate(ctx context.Context, input, output, targetLang string) error {
	if targetLang == "" {
		return fmt.Errorf("target language is required")
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// If the input is not a WAV file, convert it
	var audioFile string
	if strings.ToLower(filepath.Ext(input)) != ".wav" {
		audioFile = strings.TrimSuffix(output, ".srt") + ".wav"
		if err := t.ffmpegProcessor.ExtractAudio(ctx, input, audioFile); err != nil {
			return fmt.Errorf("failed to extract audio: %v", err)
		}
		defer os.Remove(audioFile)
	} else {
		audioFile = input
	}

	// Transcribe and translate audio
	if err := t.whisperProcessor.TranscribeWithTranslation(ctx, audioFile, output, targetLang); err != nil {
		return fmt.Errorf("failed to translate audio: %w", err)
	}

	return nil
}

// TranslateFile transcribes and translates a video file
func (t *Translator) TranslateFile(ctx context.Context, input, output, targetLang string) error {
	if targetLang == "" {
		return fmt.Errorf("target language is required")
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", input)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Extract audio from input file
	audioFile := filepath.Join(filepath.Dir(output), filepath.Base(input)+".wav")
	if err := t.ffmpegProcessor.ExtractAudio(ctx, input, audioFile); err != nil {
		return fmt.Errorf("failed to extract audio: %w", err)
	}

	// Transcribe and translate audio
	if err := t.whisperProcessor.TranscribeWithTranslation(ctx, audioFile, output, targetLang); err != nil {
		return fmt.Errorf("failed to transcribe and translate audio: %w", err)
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

// FFmpegProcessor returns the FFmpeg processor
func (t *Translator) FFmpegProcessor() *ffmpeg.FFmpeg {
	return t.ffmpegProcessor
}

// WhisperProcessor returns the Whisper processor
func (t *Translator) WhisperProcessor() *whisper.Whisper {
	return t.whisperProcessor
}

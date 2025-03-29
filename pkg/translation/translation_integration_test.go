package translation

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestTranslateWithRealFiles(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Get the test file's directory
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)

	// Get the project root (two levels up from pkg/translation)
	projectRoot := filepath.Join(testDir, "..", "..")

	// Get the model path
	modelPath := filepath.Join(projectRoot, "models", "ggml-base.en.bin")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Fatalf("Whisper model not found at %s. Please download it first.", modelPath)
	}

	testDataDir := filepath.Join(projectRoot, "testdata")
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Fatalf("Test data directory not found: %v", err)
	}

	// Create temporary directory for test outputs
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create a context with timeout for the entire test
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	tests := []struct {
		name        string
		input       string
		output      string
		targetLang  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "original_mp3",
			input:      filepath.Join(testDataDir, "Minoru_Arakawa_–_Nintendo_–_Gameboy,_interviewed_by_Maximilian_Schönherr_1990.mp3"),
			output:     filepath.Join(outputDir, "output_original"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "wav_44khz",
			input:      filepath.Join(testDataDir, "nintendo_interview_44khz.wav"),
			output:     filepath.Join(outputDir, "output_44khz"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "wav_16khz",
			input:      filepath.Join(testDataDir, "nintendo_interview_16khz.wav"),
			output:     filepath.Join(outputDir, "output_16khz"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "m4a_format",
			input:      filepath.Join(testDataDir, "nintendo_interview.m4a"),
			output:     filepath.Join(outputDir, "output_m4a"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "first_30s",
			input:      filepath.Join(testDataDir, "nintendo_interview_first_30s.mp3"),
			output:     filepath.Join(outputDir, "output_first_30s"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "last_30s",
			input:      filepath.Join(testDataDir, "nintendo_interview_last_30s.mp3"),
			output:     filepath.Join(outputDir, "output_last_30s"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "faster_speed",
			input:      filepath.Join(testDataDir, "nintendo_interview_1.5x.mp3"),
			output:     filepath.Join(outputDir, "output_1.5x"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "slower_speed",
			input:      filepath.Join(testDataDir, "nintendo_interview_0.75x.mp3"),
			output:     filepath.Join(outputDir, "output_0.75x"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "mono_audio",
			input:      filepath.Join(testDataDir, "nintendo_interview_mono.mp3"),
			output:     filepath.Join(outputDir, "output_mono"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "stereo_audio",
			input:      filepath.Join(testDataDir, "nintendo_interview_stereo.mp3"),
			output:     filepath.Join(outputDir, "output_stereo"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "high_bitrate",
			input:      filepath.Join(testDataDir, "nintendo_interview_320k.mp3"),
			output:     filepath.Join(outputDir, "output_320k"),
			targetLang: "en",
			wantErr:    false,
		},
		{
			name:       "low_bitrate",
			input:      filepath.Join(testDataDir, "nintendo_interview_128k.mp3"),
			output:     filepath.Join(outputDir, "output_128k"),
			targetLang: "en",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt // Create new variable for closure
		t.Run(tt.name, func(t *testing.T) {
			// Create a new context for this test case
			testCtx, cancel := context.WithTimeout(ctx, testTimeout)
			defer cancel()

			// Step 1: Verify input file exists
			if _, err := os.Stat(tt.input); os.IsNotExist(err) {
				t.Fatalf("Input file not found: %v", err)
			}

			// Create output directory
			if err := os.MkdirAll(filepath.Dir(tt.output), 0755); err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}

			// Step 2: Extract audio using FFmpeg
			wavFile := tt.output + ".wav"
			var ffmpegStdout, ffmpegStderr bytes.Buffer
			ffmpegCmd := exec.CommandContext(testCtx, "ffmpeg",
				"-i", tt.input,
				"-vn",
				"-acodec", "pcm_s16le",
				"-ar", "16000",
				"-ac", "1",
				"-y",
				wavFile,
			)
			ffmpegCmd.Stdout = &ffmpegStdout
			ffmpegCmd.Stderr = &ffmpegStderr

			if err := ffmpegCmd.Run(); err != nil {
				t.Fatalf("FFmpeg failed: %v\nFFmpeg stderr: %s", err, ffmpegStderr.String())
			}

			// Verify WAV file was created
			if _, err := os.Stat(wavFile); os.IsNotExist(err) {
				t.Fatalf("FFmpeg output file not created: %v\nFFmpeg stderr: %s", err, ffmpegStderr.String())
			}

			// Step 3: Run Whisper on the WAV file
			var whisperStdout, whisperStderr bytes.Buffer
			whisperCmd := exec.CommandContext(testCtx, "whisper-cli",
				"-m", modelPath,
				"-osrt",
				"-of", tt.output,
				"-f", wavFile,
			)
			whisperCmd.Stdout = &whisperStdout
			whisperCmd.Stderr = &whisperStderr

			if err := whisperCmd.Run(); err != nil {
				t.Fatalf("Whisper failed: %v\nWhisper stderr: %s", err, whisperStderr.String())
			}

			// Verify SRT file was created
			srtFile := tt.output + ".srt"
			if _, err := os.Stat(srtFile); os.IsNotExist(err) {
				t.Fatalf("Whisper output file not created: %v\nWhisper stderr: %s", err, whisperStderr.String())
			}

			// Log success
			t.Logf("Successfully processed %s", tt.name)
		})
	}
}

// testTimeout is the maximum time to wait for a test case
const testTimeout = 5 * time.Minute

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gleicon/transcoder/pkg/translation"
)

func isVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mp4" || ext == ".avi" || ext == ".mkv" || ext == ".mov"
}

func isAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mp3" || ext == ".wav" || ext == ".ogg" || ext == ".flac"
}

func processVideo(ctx context.Context, translator *translation.Translator, input, output, targetLang string, speed float64) error {
	if speed <= 0 {
		return fmt.Errorf("speed must be greater than 0")
	}

	// Extract audio from video
	audioFile := filepath.Join(filepath.Dir(output), filepath.Base(input)+".wav")
	if err := translator.FFmpegProcessor().ExtractAudio(ctx, input, audioFile); err != nil {
		return fmt.Errorf("failed to extract audio: %w", err)
	}

	// Translate audio
	translatedAudioFile := filepath.Join(filepath.Dir(output), filepath.Base(input)+".translated.wav")
	if err := translator.Translate(ctx, audioFile, translatedAudioFile, targetLang); err != nil {
		return fmt.Errorf("failed to translate audio: %w", err)
	}

	// Change video speed
	if err := translator.FFmpegProcessor().ChangeSpeed(ctx, input, output, speed); err != nil {
		return fmt.Errorf("failed to change video speed: %w", err)
	}

	return nil
}

func processAudio(ctx context.Context, translator *translation.Translator, input, output, targetLang string) error {
	if targetLang == "" {
		return fmt.Errorf("target language is required")
	}

	// Translate audio
	if err := translator.Translate(ctx, input, output, targetLang); err != nil {
		return fmt.Errorf("failed to transcribe and translate audio: %w", err)
	}

	return nil
}

func main() {
	input := flag.String("input", "", "Input file path")
	output := flag.String("output", "", "Output file path")
	targetLang := flag.String("lang", "", "Target language for translation")
	speed := flag.Float64("speed", 1.0, "Speed factor for video (default: 1.0)")
	flag.Parse()

	if *input == "" || *output == "" {
		log.Fatal("Input and output file paths are required")
	}

	if *speed <= 0 {
		log.Fatal("Speed must be greater than 0")
	}

	if _, err := os.Stat(*input); os.IsNotExist(err) {
		log.Fatal("Input file not found")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create translator
	translator, err := translation.NewWithCmd()
	if err != nil {
		log.Fatalf("Failed to create translator: %v", err)
	}
	defer translator.Close()

	ctx := context.Background()

	// Process file based on type
	if isVideoFile(*input) {
		if err := processVideo(ctx, translator, *input, *output, *targetLang, *speed); err != nil {
			log.Fatal(err)
		}
	} else if isAudioFile(*input) {
		if err := processAudio(ctx, translator, *input, *output, *targetLang); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unsupported file type")
	}

	fmt.Println("Processing completed successfully!")
}

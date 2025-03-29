#!/bin/bash

# Input file
INPUT_FILE="Minoru_Arakawa_–_Nintendo_–_Gameboy,_interviewed_by_Maximilian_Schönherr_1990.mp3"

# Create different formats
ffmpeg -i "$INPUT_FILE" -acodec pcm_s16le -ar 44100 -ac 2 "nintendo_interview_44khz.wav"
ffmpeg -i "$INPUT_FILE" -acodec pcm_s16le -ar 16000 -ac 1 "nintendo_interview_16khz.wav"
ffmpeg -i "$INPUT_FILE" -acodec libvorbis "nintendo_interview.ogg"
ffmpeg -i "$INPUT_FILE" -acodec aac "nintendo_interview.m4a"

# Create different durations (first 30 seconds, last 30 seconds)
ffmpeg -i "$INPUT_FILE" -t 30 -acodec copy "nintendo_interview_first_30s.mp3"
ffmpeg -i "$INPUT_FILE" -ss -30 -acodec copy "nintendo_interview_last_30s.mp3"

# Create different speeds
ffmpeg -i "$INPUT_FILE" -filter:a "atempo=1.5" -acodec libmp3lame "nintendo_interview_1.5x.mp3"
ffmpeg -i "$INPUT_FILE" -filter:a "atempo=0.75" -acodec libmp3lame "nintendo_interview_0.75x.mp3"

# Create mono and stereo versions
ffmpeg -i "$INPUT_FILE" -ac 1 -acodec libmp3lame "nintendo_interview_mono.mp3"
ffmpeg -i "$INPUT_FILE" -ac 2 -acodec libmp3lame "nintendo_interview_stereo.mp3"

# Create different bitrates
ffmpeg -i "$INPUT_FILE" -ab 128k -acodec libmp3lame "nintendo_interview_128k.mp3"
ffmpeg -i "$INPUT_FILE" -ab 320k -acodec libmp3lame "nintendo_interview_320k.mp3"

echo "Generated additional audio samples from $INPUT_FILE" 
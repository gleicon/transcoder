#!/bin/bash

# Create testdata directory if it doesn't exist
mkdir -p testdata

# Function to generate a WAV file with specific parameters
generate_wav() {
    local output=$1
    local duration=$2
    local sample_rate=16000
    local channels=1
    local bits=16
    
    ffmpeg -f lavfi -i "sine=frequency=440:duration=$duration" \
           -ar $sample_rate \
           -ac $channels \
           -acodec pcm_s16le \
           "$output"
}

# Function to generate a video file
generate_video() {
    local output=$1
    local duration=$2
    local resolution=$3
    
    ffmpeg -f lavfi -i "testsrc=duration=$duration:size=$resolution:rate=30" \
           -f lavfi -i "sine=frequency=440:duration=$duration" \
           -c:v libx264 \
           -c:a aac \
           "$output"
}

# Generate WAV files
echo "Generating WAV files..."
generate_wav "testdata/sample_16khz.wav" 5
generate_wav "testdata/sample_16khz_short.wav" 2
generate_wav "testdata/sample_16khz_long.wav" 10

# Generate video files
echo "Generating video files..."
generate_video "testdata/sample.mp4" 5 "1280x720"
generate_video "testdata/sample_short.mp4" 2 "1280x720"
generate_video "testdata/sample_long.mp4" 10 "1920x1080"

# Generate other audio formats
echo "Generating other audio formats..."
ffmpeg -i "testdata/sample_16khz.wav" -c:a libmp3lame "testdata/sample.mp3"
ffmpeg -i "testdata/sample_16khz.wav" -c:a aac "testdata/sample.m4a"
ffmpeg -i "testdata/sample_16khz.wav" -c:a libvorbis "testdata/sample.ogg"

# Generate other video formats
echo "Generating other video formats..."
ffmpeg -i "testdata/sample.mp4" -c:v libx264 -c:a aac "testdata/sample.avi"
ffmpeg -i "testdata/sample.mp4" -c:v libx264 -c:a aac "testdata/sample.mov"
ffmpeg -i "testdata/sample.mp4" -c:v libx264 -c:a aac "testdata/sample.mkv"

echo "Test files generated successfully!" 
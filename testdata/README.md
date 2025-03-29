# Test Data Directory

This directory contains test files used for testing the Transcoder application. All files are carefully selected to meet specific requirements for each test case.
The audio files are based on https://commons.wikimedia.org/wiki/File:Minoru_Arakawa_%E2%80%93_Nintendo_%E2%80%93_Gameboy,_interviewed_by_Maximilian_Sch%C3%B6nherr_1990.MP3

## Original file
- 
*
* 
## Audio Files

### WAV Files (16kHz, 16-bit, mono)
- `sample_16khz.wav`: A 16kHz WAV file that meets Whisper's requirements
- `sample_16khz_short.wav`: A short 16kHz WAV file for quick tests
- `sample_16khz_long.wav`: A longer 16kHz WAV file for performance tests

### Other Audio Formats
- `sample.mp3`: MP3 file for format conversion tests
- `sample.m4a`: M4A file for format conversion tests
- `sample.ogg`: OGG file for format conversion tests

## Video Files

### MP4 Files
- `sample.mp4`: Standard MP4 video file
- `sample_short.mp4`: Short MP4 video for quick tests
- `sample_long.mp4`: Longer MP4 video for performance tests

### Other Video Formats
- `sample.avi`: AVI file for format support tests
- `sample.mov`: MOV file for format support tests
- `sample.mkv`: MKV file for format support tests

## Requirements

### WAV Files for Whisper
- Sample rate: 16kHz
- Bit depth: 16-bit
- Channels: Mono
- Format: PCM

### Video Files
- Resolution: 720p or 1080p
- Codec: H.264
- Audio: AAC
- Duration: Various lengths for different test cases

## Usage

These files are used in the test suite to verify:
1. Audio extraction from video files
2. Audio format conversion
3. Transcription and translation
4. Video speed adjustment
5. File format support

## Maintenance

When adding new test files:
1. Ensure they meet the specified requirements
2. Update this README with file details
3. Keep file sizes reasonable for repository management
4. Use descriptive names that indicate the file's purpose 

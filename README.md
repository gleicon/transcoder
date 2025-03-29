# Transcoder

A command-line tool for video and audio processing using FFmpeg, Whisper, and local translation capabilities.

## Features

- Transcribe audio from video files using Whisper
- Translate audio to different languages
- Change video playback speed
- Support for various video and audio formats

## System Requirements

- macOS or Linux
- FFmpeg
- Whisper (installed via Homebrew)
- Go 1.23 or later

## Installation

### macOS

1. Install Homebrew if you haven't already:
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. Install FFmpeg and Whisper:
   ```bash
   brew install ffmpeg whisper-cpp
   ```

3. Clone the repository:
   ```bash
   git clone https://github.com/gleicon/transcoder.git
   cd transcoder
   ```

4. Build the project:
   ```bash
   make setup
   ```

### Linux

1. Install FFmpeg and Whisper:
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install ffmpeg
   # Install whisper-cpp from source
   git clone https://github.com/ggerganov/whisper.cpp.git
   cd whisper.cpp
   make
   sudo make install

   # Fedora
   sudo dnf install ffmpeg
   # Install whisper-cpp from source
   git clone https://github.com/ggerganov/whisper.cpp.git
   cd whisper.cpp
   make
   sudo make install
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/gleicon/transcoder.git
   cd transcoder
   ```

3. Build the project:
   ```bash
   make setup
   ```

## Usage

The tool supports both video and audio files. For video files, it can extract audio, transcribe it, and optionally change the video speed. For audio files, it can transcribe and translate them directly.

### Process Video

```bash
transcoder -i input.mp4 -o output.mp4 -lang en -speed 1.5
```

Options:
- `-i`: Input video file
- `-o`: Output video file
- `-lang`: Target language for translation (e.g., en, es, fr)
- `-speed`: Playback speed multiplier (e.g., 1.5 for 50% faster)

### Process Audio

```bash
transcoder -i input.wav -o output.srt -lang en
```

Options:
- `-i`: Input audio file
- `-o`: Output SRT file
- `-lang`: Target language for translation (e.g., en, es, fr)

## Supported File Types

### Video Files
- MP4 (.mp4)
- AVI (.avi)
- MOV (.mov)
- MKV (.mkv)

### Audio Files
- WAV (.wav)
- MP3 (.mp3)
- M4A (.m4a)
- OGG (.ogg)

## Troubleshooting

### Missing Dependencies

If you encounter errors about missing dependencies:

1. For FFmpeg:
   ```bash
   # macOS
   brew install ffmpeg

   # Ubuntu/Debian
   sudo apt install ffmpeg

   # Fedora
   sudo dnf install ffmpeg
   ```

2. For Whisper:
   ```bash
   # macOS
   brew install whisper-cpp

   # Linux
   git clone https://github.com/ggerganov/whisper.cpp.git
   cd whisper.cpp
   make
   sudo make install
   ```

### Common Issues

1. **Whisper not found**: Make sure whisper-cli is installed and available in your PATH:
   ```bash
   which whisper-cli
   ```

2. **FFmpeg not found**: Ensure FFmpeg is installed and available in your PATH:
   ```bash
   which ffmpeg
   ```

3. **Permission issues**: If you encounter permission errors, make sure you have write access to the output directory.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
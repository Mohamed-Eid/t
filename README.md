# t

A simple and efficient task runner for Go projects, similar to Make but with YAML configuration.

## 🚀 Installation

### Option 1: Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/Mohamed-Eid/t/releases):

- **Windows**: `t-windows-amd64.zip`
- **Linux**: `t-linux-amd64.tar.gz`
- **macOS Intel**: `t-darwin-amd64.tar.gz`
- **macOS Apple Silicon**: `t-darwin-arm64.tar.gz`

Extract the archive and move the binary to your PATH:

#### Windows

```powershell
# Extract the zip file and move to a directory in your PATH
Move-Item t.exe C:\Windows\System32\
```

#### Linux/macOS

```bash
# Extract and install
tar -xzf t-*.tar.gz
sudo mv t /usr/local/bin/
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/Mohamed-Eid/t.git
cd t

# Build the binary
go build -o t .

# Install globally (optional)
go install .
```

## 📖 Usage

### Basic Usage

```bash
t <task-name>
```

### Configuration

Create a `tasks.yaml` file in your project root:

```yaml
tasks:
  build:
    description: "Build the application"
    deps: [clean]
    command: "go build -o app ."

  test:
    description: "Run tests"
    command: "go test ./..."

  clean:
    description: "Clean build artifacts"
    command: "rm -f app"

  dev:
    description: "Run in development mode"
    command: "go run ."
```

### Example Commands

```bash
# Build your project
t build

# Run tests
t test

# Clean artifacts
t clean

# Start development server
t dev
```

## 🔧 Features

- ✅ Simple YAML configuration
- ✅ Cross-platform support (Windows, Linux, macOS)
- ✅ Fast execution
- ✅ Easy to use and configure
- ✅ Lightweight alternative to Make

## 📁 Project Structure

```
your-project/
├── tasks.yaml          # Task configuration
├── main.go
├── go.mod
└── other files...
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🐛 Issues

If you encounter any issues or have suggestions, please [open an issue](https://github.com/Mohamed-Eid/t/issues).

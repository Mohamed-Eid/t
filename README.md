# t

A simple and efficient task runner for your projects, similar to Make but with YAML configuration.

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
# Extract the zip file to a directory
Expand-Archive -Path "t-windows-amd64.zip" -DestinationPath "C:\tools\t"

# Add to PATH (run as Administrator)
$env:PATH += ";C:\tools\t"

# Or copy to System32 (not recommended for production)
# Copy-Item "C:\tools\t\t.exe" "C:\Windows\System32\"
```

#### Linux/macOS

```bash
# Extract and install
tar -xzf t-*.tar.gz

# Make executable and move to PATH
chmod +x t
sudo mv t /usr/local/bin/

# Verify installation
t --help || echo "Ready to use! Create a tasks.yaml file first."
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

Create a `tasks.yaml` file in your project root with the correct format:

```yaml
version: "1"

vars:
  BIN: "bin/app"
  SRC: "main.go"

tasks:
  build:
    desc: "Build the application"
    deps: [clean]
    cmds:
      - "go build -o {{.BIN}} {{.SRC}}"

  test:
    desc: "Run tests"
    cmds:
      - "go test ./..."

  clean:
    desc: "Clean build artifacts"
    cmds:
      - "rm -f {{.BIN}}"
      - "rm -rf bin/"

  dev:
    desc: "Run in development mode"
    cmds:
      - "go run {{.SRC}}"

  install:
    desc: "Install dependencies"
    cmds:
      - "go mod download"
      - "go mod tidy"
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

# Install dependencies
t install
```

## 🔧 Configuration Format

- **`version`**: Configuration version (currently "1")
- **`vars`**: Variables that can be used in commands with `{{.VARIABLE_NAME}}`
- **`tasks`**: Available tasks with the following properties:
  - **`desc`**: Task description
  - **`deps`**: List of dependencies (tasks to run first)
  - **`cmds`**: List of commands to execute

## 🚨 Troubleshooting

### Error: "open tasks.yaml: The system cannot find the file specified"

This error occurs when you run `t` in a directory without a `tasks.yaml` file.

**Solutions:**

1. **Create a `tasks.yaml` file** in your project root using the format above
2. **Run `t` from the correct directory** where `tasks.yaml` exists
3. **Check if the file exists**: `ls tasks.yaml` (Linux/macOS) or `dir tasks.yaml` (Windows)

### Quick Start Template

Create this basic `tasks.yaml` in your project:

```yaml
version: "1"

tasks:
  hello:
    desc: "Hello world task"
    cmds:
      - "echo Hello from t task runner!"
```

Then run: `t hello`

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

# t

A lightweight and efficient task runner for your projects, similar to Make but with YAML configuration. Built with Go for cross-platform compatibility and speed.

## ✨ Features

- 🚀 **Simple YAML configuration** - Easy to read and write
- 🔗 **Task dependencies** - Automatic dependency resolution
- 🔄 **Variable substitution** - Reusable configuration with variables
- 🌍 **Cross-platform** - Works on Windows, Linux, and macOS
- ⚡ **Fast execution** - Built with Go for performance
- 🛠️ **Built-in commands** - Project initialization and task listing
- 🔒 **No conflicts** - Tool commands use `:` prefix to avoid task name conflicts

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

# Verify installation
t --help
```

#### Linux/macOS

```bash
# Extract and install
tar -xzf t-*.tar.gz

# Make executable and move to PATH
chmod +x t
sudo mv t /usr/local/bin/

# Verify installation
t --help
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

## 📖 Quick Start

### 1. Initialize a new project

```bash
# Create a tasks.yaml with default tasks
t :init
```

### 2. List available tasks

```bash
# Show all defined tasks
t :list
```

### 3. Run tasks

```bash
# Run any task by name
t build
t test
t clean
```

## 🔧 Usage

### Tool Commands (`:` prefix)

```bash
t :init         # Initialize tasks.yaml with defaults
t :list         # List all available tasks
t :ls           # Alias for :list
t :version      # Show version information
t --help        # Show help information
```

### User Tasks (no prefix)

```bash
t <task-name>   # Run any task defined in tasks.yaml
t build         # Example: run build task
t test          # Example: run test task
```

### Configuration

The `tasks.yaml` file defines your project tasks:

```yaml
version: "1"

vars:
  APP_NAME: "myapp"
  BUILD_DIR: "bin"
  VERSION: "1.0.0"

tasks:
  build:
    desc: "Build the application"
    deps: [clean]
    cmds:
      - "New-Item -ItemType Directory -Force -Path {{.BUILD_DIR}}"
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o {{.BUILD_DIR}}/{{.APP_NAME}}.exe .'

  test:
    desc: "Run tests"
    cmds:
      - "go test ./..."

  clean:
    desc: "Clean build artifacts"
    cmds:
      - "go clean"
      - "Remove-Item -Recurse -Force {{.BUILD_DIR}} -ErrorAction SilentlyContinue"

  dev:
    desc: "Run in development mode"
    cmds:
      - "go run ."

  install:
    desc: "Install dependencies"
    cmds:
      - "go mod download"
      - "go mod tidy"

  lint:
    desc: "Run linter and formatter"
    cmds:
      - "go fmt ./..."
      - "go vet ./..."

  release:
    desc: "Build optimized release binary"
    deps: [test, lint]
    cmds:
      - "New-Item -ItemType Directory -Force -Path dist"
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o dist/{{.APP_NAME}}.exe .'
```

## � Configuration Reference

### Structure

- **`version`**: Configuration version (currently "1")
- **`vars`**: Variables that can be used in commands with `{{.VARIABLE_NAME}}`
- **`tasks`**: Available tasks with the following properties:
  - **`desc`**: Task description (shown in `:list`)
  - **`deps`**: List of dependencies (tasks to run first)
  - **`cmds`**: List of commands to execute

### Variables

Use variables in commands for reusability:

```yaml
vars:
  APP_NAME: "myapp"
  VERSION: "1.0.0"

tasks:
  build:
    cmds:
      - "go build -o {{.APP_NAME}} ."
      - "echo Built {{.APP_NAME}} version {{.VERSION}}"
```

### Dependencies

Tasks can depend on other tasks:

```yaml
tasks:
  release:
    deps: [test, lint, build] # Runs test, lint, then build before release
    cmds:
      - "echo Ready for release!"
```

## 🚨 Troubleshooting

### Error: "tasks.yaml not found in current directory"

This error occurs when you run `t` in a directory without a `tasks.yaml` file.

**Solutions:**

1. **Create a `tasks.yaml` file**: Run `t :init` to create one with defaults
2. **Run from the correct directory**: Navigate to where your `tasks.yaml` exists
3. **Check if the file exists**:
   - Windows: `dir tasks.yaml`
   - Linux/macOS: `ls tasks.yaml`

### Error: "task <name> not found"

The task name doesn't exist in your `tasks.yaml` file.

**Solutions:**

1. **Check available tasks**: Run `t :list` to see all defined tasks
2. **Check spelling**: Ensure the task name matches exactly
3. **Check YAML syntax**: Ensure your `tasks.yaml` is valid

### Commands not working on Windows

If you get "command not found" errors on Windows:

1. **PowerShell commands**: Use PowerShell syntax (e.g., `New-Item`, `Remove-Item`)
2. **Cross-platform commands**: Use Go commands or tools available on all platforms

## 🏗️ Project Structure

```
your-project/
├── tasks.yaml          # Task configuration (required)
├── src/                # Your source code
├── docs/               # Documentation
├── tests/              # Test files
└── ...                 # Other project files
```

## 🎯 Examples

### Go Project

```yaml
version: "1"

vars:
  APP_NAME: "myapp"
  VERSION: "1.0.0"

tasks:
  build:
    desc: "Build the Go application"
    deps: [test]
    cmds:
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o {{.APP_NAME}} .'

  test:
    desc: "Run Go tests"
    cmds:
      - "go test ./..."
      - "go vet ./..."

  clean:
    desc: "Clean build artifacts"
    cmds:
      - "go clean"
      - "Remove-Item -Force {{.APP_NAME}}.exe -ErrorAction SilentlyContinue"

  run:
    desc: "Run the application"
    deps: [build]
    cmds:
      - "./{{.APP_NAME}}"
```

### Node.js Project

```yaml
version: "1"

vars:
  NODE_ENV: "development"

tasks:
  install:
    desc: "Install npm dependencies"
    cmds:
      - "npm install"

  build:
    desc: "Build the project"
    deps: [install]
    cmds:
      - "npm run build"

  test:
    desc: "Run tests"
    cmds:
      - "npm test"

  dev:
    desc: "Start development server"
    deps: [install]
    cmds:
      - "npm run dev"

  lint:
    desc: "Lint code"
    cmds:
      - "npm run lint"
      - "npm run format"
```

### Docker Project

```yaml
version: "1"

vars:
  IMAGE_NAME: "myapp"
  TAG: "latest"

tasks:
  build:
    desc: "Build Docker image"
    cmds:
      - "docker build -t {{.IMAGE_NAME}}:{{.TAG}} ."

  run:
    desc: "Run Docker container"
    deps: [build]
    cmds:
      - "docker run -p 8080:8080 {{.IMAGE_NAME}}:{{.TAG}}"

  push:
    desc: "Push image to registry"
    deps: [build]
    cmds:
      - "docker push {{.IMAGE_NAME}}:{{.TAG}}"

  clean:
    desc: "Clean Docker artifacts"
    cmds:
      - "docker image prune -f"
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Make](https://www.gnu.org/software/make/) and [Task](https://taskfile.dev/)
- Built with [Cobra](https://cobra.dev/) for CLI functionality
- Uses [Go](https://golang.org/) for cross-platform compatibility

## 🐛 Issues & Support

If you encounter any issues or have suggestions:

- 📋 [Open an issue](https://github.com/Mohamed-Eid/t/issues)
- 💬 [Start a discussion](https://github.com/Mohamed-Eid/t/discussions)
- 📧 Email: medoeid50@gmail.com

---

**Made with ❤️ by [Mohamed Eid](https://github.com/Mohamed-Eid)**

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

# t

A lightweight and efficient task runner for your projects, similar to Make but with YAML configuration. Built with Go for cross-platform compatibility and speed.

## ✨ Features

- 🚀 **Simple YAML configuration** - Easy to read and write
- 🔗 **Task dependencies** - Automatic dependency resolution
- ⚡ **Parallel execution** - Dependencies run concurrently for faster builds
- 🔄 **Variable substitution** - Reusable configuration with variables
- 🌍 **Cross-platform** - Works on Windows, Linux, and macOS
- 🏃 **Fast execution** - Built with Go and Goroutines for maximum performance
- 🛠️ **Built-in commands** - Project initialization and task listing
- 🔒 **No conflicts** - Tool commands use `:` prefix to avoid task name conflicts
- ⏱️ **Timing information** - Track execution time and performance
- 🧵 **Thread-safe** - Concurrent execution without race conditions
- 🔄 **Detached execution** - Run long-living tasks in background
- 📝 **Process management** - Track, monitor, and control background tasks
- 🤔 **Interactive tasks** - Prompt users for input during execution

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

# Run with timing information to see parallel execution
t :parallel build

# Example output:
# ⏱️  Starting task 'build' at 15:04:05.123
# 🔧 Running task: format
# 🔧 Running task: vet
# 🔧 Running task: deps
# ✅ done (all three tasks run in parallel!)
# 🔧 Running task: build
# ✅ done
# 🎉 Task 'build' completed successfully in 2.1s!
```

## 🔧 Usage

### Tool Commands (`:` prefix)

```bash
t :init         # Initialize tasks.yaml with defaults
t :list         # List all available tasks
t :ls           # Alias for :list
t :parallel     # Run task with timing information
t :time         # Alias for :parallel
t :detach       # Run task in background (detached mode)
t :d            # Alias for :detach (short form)
t :bg           # Alias for :detach (background)
t :ps           # List running detached tasks
t :p            # Alias for :ps (short form)
t :processes    # Alias for :ps (descriptive)
t :stop         # Stop a running detached task
t :kill         # Alias for :stop (forceful)
t :s            # Alias for :stop (short form)
t :logs         # View logs of a detached task
t :log          # Alias for :logs (singular)
t :l            # Alias for :logs (short form)
t :tail         # Alias for :logs (tail-like)
t :version      # Show version information
t --help        # Show help information
```

### User Tasks (no prefix)

```bash
t <task-name>   # Run any task defined in tasks.yaml
t build         # Example: run build task
t test          # Example: run test task

# Performance commands
t :parallel <task-name>  # Run task with detailed timing information
t :time <task-name>      # Alias for :parallel (short form)
```

## 🔗 Quick Reference

### Command Aliases

| Full Command  | Short Aliases                    | Purpose                  |
| ------------- | -------------------------------- | ------------------------ |
| `t :detach`   | `:d`, `:bg`, `:background`       | Start task in background |
| `t :ps`       | `:p`, `:processes`, `:status`    | List running tasks       |
| `t :stop`     | `:s`, `:kill`, `:terminate`      | Stop running task        |
| `t :logs`     | `:l`, `:log`, `:tail`            | View task logs           |
| `t :parallel` | `:time`, `:timing`, `:benchmark` | Run with timing          |
| `t :list`     | `:ls`                            | List available tasks     |

### Common Workflows

```bash
# Quick development server workflow
t :d serve                    # Start server in background
t :p                          # Check if it's running
t :l serve -f                 # Follow logs
t :s serve                    # Stop when done

# Build with performance monitoring
t :time build                 # Build with timing info

# Interactive task workflows
t echo                        # Prompt for message and echo it
t greet                       # Interactive greeting with name
t commit                      # Interactive git commit
t deploy                      # Interactive deployment with confirmation

# Background task management
t :bg watch                   # Start file watcher
t :bg serve                   # Start dev server
t :processes                  # List all running
t :kill watch                 # Stop file watcher
t :terminate serve            # Stop dev server
```

### Configuration

The `tasks.yaml` file defines your project tasks with support for parallel execution:

```yaml
version: "1"

vars:
  APP_NAME: "myapp"
  BUILD_DIR: "bin"
  VERSION: "1.0.0"

tasks:
  # Independent tasks that can run in parallel
  format:
    desc: "Format code"
    cmds:
      - "go fmt ./..."

  vet:
    desc: "Vet code"
    cmds:
      - "go vet ./..."

  deps:
    desc: "Download dependencies"
    cmds:
      - "go mod download"

  # Tasks with dependencies - these will run in parallel when possible
  test:
    desc: "Run tests"
    deps: [format, vet] # format and vet run in parallel
    cmds:
      - "go test ./..."

  lint:
    desc: "Run linter"
    deps: [format, vet] # format and vet run in parallel
    cmds:
      - "echo Linting completed"

  build:
    desc: "Build the application"
    deps: [test, deps] # test and deps can run in parallel
    cmds:
      - "New-Item -ItemType Directory -Force -Path {{.BUILD_DIR}}"
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o {{.BUILD_DIR}}/{{.APP_NAME}}.exe .'

  release:
    desc: "Build optimized release binary"
    deps: [test, lint] # test and lint run in parallel
    cmds:
      - "New-Item -ItemType Directory -Force -Path dist"
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o dist/{{.APP_NAME}}.exe .'
```

## 🔧 Configuration Reference

### Structure

- **`version`**: Configuration version (currently "1")
- **`vars`**: Variables that can be used in commands with `{{.VARIABLE_NAME}}`
- **`tasks`**: Available tasks with the following properties:
  - **`desc`**: Task description (shown in `:list`)
  - **`deps`**: List of dependencies (tasks to run first) - **runs in parallel when possible**
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

Tasks can depend on other tasks, and **t automatically runs dependencies in parallel** when possible:

```yaml
tasks:
  release:
    deps: [test, lint, build] # Runs test, lint, and build in parallel when possible
    cmds:
      - "echo Ready for release!"
```

## ⚡ Parallel Execution

**t** automatically detects which tasks can run in parallel and executes them concurrently using Goroutines:

### How It Works

1. **Dependency Analysis**: t analyzes the dependency graph to find tasks that can run simultaneously
2. **Concurrent Execution**: Independent tasks run in parallel using Goroutines
3. **Thread-Safe**: Uses `sync.RWMutex` to prevent race conditions
4. **Optimal Performance**: Reduces total execution time significantly

### Example

```yaml
tasks:
  # These three tasks are independent and will run in parallel
  format:
    desc: "Format code"
    cmds: ["go fmt ./..."]

  vet:
    desc: "Vet code"
    cmds: ["go vet ./..."]

  deps:
    desc: "Download deps"
    cmds: ["go mod download"]

  # This task depends on all three above - they run in parallel first
  build:
    desc: "Build app"
    deps: [format, vet, deps] # All three run concurrently!
    cmds: ["go build ."]
```

### Performance Comparison

**Without Parallel Execution:**

```
format (2s) → vet (2s) → deps (1s) → build (1s) = 6 seconds total
```

**With Parallel Execution:**

```
format (2s) ┐
vet (2s)    ├─→ build (1s) = 3 seconds total
deps (1s)   ┘
```

### Monitoring Performance

Use the `:parallel` command to see timing information:

```bash
# Run with timing information
t :parallel build    # or t :time build, t :timing build

# Example output:
# ⏱️  Starting task 'build' at 15:04:05.123
# 🔧 Running task: format
# 🔧 Running task: vet
# 🔧 Running task: deps
# ... (tasks run concurrently)
# 🎉 Task 'build' completed successfully in 3.2s!
```

# 🎉 Task 'build' completed successfully in 3.2s!

````

## 🔄 Detached Execution

**t** supports running long-living tasks in the background, perfect for development servers, file watchers, and other persistent processes.

### Background Tasks

Start any task in detached mode:

```bash
# Start a development server in the background
t :detach serve    # or t :d serve, t :bg serve

# Start a file watcher
t :detach watch    # or t :d watch

# Start any long-running task
t :detach long-task
```

### Process Tree Management

The detach feature properly handles **process trees and child processes**:

- **Windows**: Uses `taskkill /T` to terminate the entire process tree
- **Unix/Linux**: Attempts to kill the process group first, then individual processes
- **Child Process Cleanup**: When you stop a detached task like `php artisan serve`, all child processes are properly terminated

This ensures that commands like `php artisan serve`, `npm run dev`, or any server that spawns child processes won't leave orphaned processes running when stopped.

**Example use cases:**
- `php artisan serve` - PHP development server
- `npm run dev` - Node.js development server
- `python manage.py runserver` - Django development server
- `docker-compose up` - Docker services with multiple containers
- Any long-running server or daemon process that spawns children`

### Process Management

Monitor and control background tasks:

```bash
# List all running detached tasks
t :ps              # or t :p, t :processes, t :status

# View live logs (follow mode)
t :logs serve --follow    # or t :log serve -f, t :l serve -f, t :tail serve -f

# View recent logs
t :logs serve      # or t :log serve, t :l serve

# Stop a task by name
t :stop serve      # or t :kill serve, t :s serve

# Stop a task by PID
t :stop 12345      # or t :kill 12345
```

### Example Output

```bash
$ t :detach serve
🚀 Starting detached task: serve
✅ Task 'serve' started in background (PID: 12345)
📝 Logs: .t-logs/serve-20250809-071236.log
🛑 Stop with: t :stop serve (or PID 12345)

$ t :ps
🔧 Running detached tasks (1):

  📋 Task: serve
     🆔 PID: 12345
     ⏰ Running for: 2m30s
     📝 Log file: .t-logs/serve-20250809-071236.log
     🛑 Stop with: t :stop serve

$ t :logs serve -f
📝 Logs for task 'serve':
📡 Following logs (Press Ctrl+C to exit)...
─────────────────────────────────────────────
[15:04:05] Server is running...
[15:04:10] Server is running...
[15:04:15] Server is running...
```

### Log Management

All detached tasks automatically log their output:

- **Log Directory**: `.t-logs/`
- **Log Format**: `<task-name>-<timestamp>.log`
- **Auto-cleanup**: Process tracking files are removed when tasks stop

### Perfect For

- 🌐 **Development servers** (`php artisan serve`, `npm run dev`)
- 👀 **File watchers** (`npm run watch`, `sass --watch`)
- 🏗️ **Build processes** (`webpack --watch`)
- 🐳 **Docker containers** (`docker-compose up`)
- 🧪 **Test runners** (`jest --watchAll`)

### Example Tasks

```yaml
tasks:
  serve:
    desc: "Start development server"
    cmds:
      - "echo Starting server..."
      - "php artisan serve --host=0.0.0.0 --port=8000"

  watch:
    desc: "Watch files for changes"
    cmds:
      - "npm run watch"

  docker:
    desc: "Start Docker services"
    cmds:
      - "docker-compose up"
```

## 🤔 Interactive Tasks

**t** supports interactive tasks that prompt users for input during execution, perfect for commands that need dynamic values like commit messages, deployment targets, or user preferences.

### Defining Interactive Tasks

```yaml
tasks:
  echo:
    desc: "Echo a message with user input"
    interactive:
      message:
        message: "Enter the message to echo"
        required: true
    cmds:
      - "echo $message"

  greet:
    desc: "Greet someone with custom name"
    interactive:
      name:
        message: "Enter your name"
        required: true
      greeting:
        message: "Enter greeting"
        default: "Hello"
        required: false
    cmds:
      - "echo $greeting $name!"

  commit:
    desc: "Git commit with interactive message"
    interactive:
      message:
        message: "Enter commit message"
        required: true
      files:
        message: "Enter files to add (or leave empty for all)"
        default: "."
        required: false
    cmds:
      - "git add $files"
      - 'git commit -m "$message"'
```

### Interactive Configuration

- **`message`**: The prompt text shown to the user
- **`required`**: Whether the input is mandatory (true/false)
- **`default`**: Default value used if user provides no input

### Variable Syntax

Use `$variable_name` in commands to reference interactive inputs:

```yaml
tasks:
  deploy:
    interactive:
      env:
        message: "Target environment (dev/prod)"
        required: true
    cmds:
      - "kubectl apply -f k8s/$env/"
      - "echo Deployed to $env environment"
```

### Example Usage

```bash
$ t echo
🤔 Task 'echo' requires interactive input:

📝 Enter the message to echo (required): Hello World!
✅ message: Hello World!

➡️  echo Hello World!
Hello World!
✅ done
🎉 Task 'echo' completed successfully!

$ t greet
🤔 Task 'greet' requires interactive input:

📝 Enter your name (required): John
✅ name: John
📝 Enter greeting [Hello]: Hi there
✅ greeting: Hi there

➡️  echo Hi there John!
Hi there John!
✅ done
🎉 Task 'greet' completed successfully!
```

### Perfect For

- 📝 **Git operations** - Interactive commit messages, branch names
- 🚀 **Deployments** - Environment selection, confirmation prompts
- 🔧 **Configuration** - Dynamic settings, user preferences
- 📦 **Package management** - Version selection, dependency updates
- 🧪 **Testing** - Test environment selection, test data input

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
  # Independent tasks - these will run in parallel
  format:
    desc: "Format Go code"
    cmds:
      - "go fmt ./..."

  vet:
    desc: "Vet Go code"
    cmds:
      - "go vet ./..."

  deps:
    desc: "Download dependencies"
    cmds:
      - "go mod download"
      - "go mod tidy"

  # Dependent tasks - dependencies run in parallel when possible
  test:
    desc: "Run Go tests"
    deps: [format, vet] # format and vet run in parallel
    cmds:
      - "go test ./..."

  build:
    desc: "Build the Go application"
    deps: [test, deps] # test and deps can run in parallel
    cmds:
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o {{.APP_NAME}} .'

  clean:
    desc: "Clean build artifacts"
    cmds:
      - "go clean"
      - "Remove-Item -Force {{.APP_NAME}}.exe -ErrorAction SilentlyContinue"

  release:
    desc: "Create release build"
    deps: [test, build] # test and build dependencies handled optimally
    cmds:
      - "New-Item -ItemType Directory -Force -Path dist"
      - 'go build -ldflags="-s -w -X main.Version={{.VERSION}}" -o dist/{{.APP_NAME}} .'
```

### Node.js Project

```yaml
version: "1"

vars:
  NODE_ENV: "development"

tasks:
  # Independent setup tasks
  install:
    desc: "Install npm dependencies"
    cmds:
      - "npm install"

  lint:
    desc: "Lint code"
    cmds:
      - "npm run lint"

  format:
    desc: "Format code"
    cmds:
      - "npm run format"

  # Quality checks that can run in parallel
  quality:
    desc: "Run quality checks"
    deps: [lint, format] # lint and format run in parallel
    cmds:
      - "echo Quality checks completed"

  # Build with optimized dependencies
  build:
    desc: "Build the project"
    deps: [install, quality] # install and quality can run in parallel
    cmds:
      - "npm run build"

  test:
    desc: "Run tests"
    deps: [install]
    cmds:
      - "npm test"

  dev:
    desc: "Start development server"
    deps: [install]
    cmds:
      - "npm run dev"

  # Production build with all checks
  production:
    desc: "Production build"
    deps: [test, build] # test and build run optimally
    cmds:
      - "npm run build:prod"
```

build:
desc: "Build the project"
deps: [install]
cmds: - "npm run build"

test:
desc: "Run tests"
cmds: - "npm test"

dev:
desc: "Start development server"
deps: [install]
cmds: - "npm run dev"

lint:
desc: "Lint code"
cmds: - "npm run lint" - "npm run format"

````

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

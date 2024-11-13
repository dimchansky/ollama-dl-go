# ollama-dl (Go Version)

A fast and efficient Go implementation for downloading models from the [Ollama Library](https://ollama.com/library). This tool mirrors the functionality of the [original Python version](https://github.com/akx/ollama-dl) but offers the performance benefits of a statically-compiled Go binary. With `ollama-dl` written in Go, you get the advantage of native concurrency and a single executable without additional Python dependencies.

## 🚀 Key Features
- **Concurrent Downloads**: Downloads multiple model layers simultaneously using Go's goroutines.
- **Resumable Downloads**: Supports partial downloads using HTTP range requests, allowing you to resume interrupted downloads.
- **Progress Display**: Provides a live progress bar to keep you informed of the download status.
- **Simple CLI**: Easy to use, with minimal setup required.

## 📦 Installation

1. **Build from source**:
   Ensure you have [Go 1.22+](https://golang.org/dl/) installed, then run:

   ```bash
   git clone https://github.com/dimchansky/ollama-dl-go.git
   cd ollama-dl-go
   go build -o ollama-dl

2. **Download the binary (coming soon)**:
   Precompiled binaries will be available for Windows, macOS, and Linux.

## 🛠 Example Usage

To download the default latest version of a model (e.g., Meta Llama 3.2):

```
$ ./ollama-dl llama3.2
library-llama3.2-latest/license-fcc5a6bec9da.txt 100% |█████████████████| (7.7/7.7 kB, 6.5 MB/s)
library-llama3.2-latest/template-966de95ca8a6.txt 100% |████████████████| (1.4/1.4 kB, 1.8 MB/s)
library-llama3.2-latest/params-56bb8bd477a5.json 100% |████████████████████| (96/96 B, 257 kB/s)
library-llama3.2-latest/license-a70ff7e570d9.txt 100% |██████████████████| (6.0/6.0 kB, 12 MB/s)
library-llama3.2-latest/model-dde5aa3fc5ff.gguf 100% |██████████████████| (2.0/2.0 GB, 78 MB/s)
Download complete
```

To download a specific version (e.g., llama3.2:3b):

```
$ ./ollama-dl llama3.2:3b
library-llama3.2-3b/params-56bb8bd477a5.json 100% |████████████████████████| (96/96 B, 107 kB/s)
library-llama3.2-3b/template-966de95ca8a6.txt 100% |████████████████████| (1.4/1.4 kB, 3.0 MB/s)
library-llama3.2-3b/license-fcc5a6bec9da.txt 100% |█████████████████████| (7.7/7.7 kB, 3.6 MB/s)
library-llama3.2-3b/license-a70ff7e570d9.txt 100% |█████████████████████| (6.0/6.0 kB, 785 kB/s)
library-llama3.2-3b/model-dde5aa3fc5ff.gguf 100% |██████████████████████| (2.0/2.0 GB, 74 MB/s)
Download complete
```

## 🔥 Why Use the Go Version?

-	Speed: Go’s concurrency model and lightweight binaries ensure fast and reliable downloads.
-	No Dependencies: Unlike the Python version, there is no need for a Python environment or virtualenv. Just a single executable.
-	Cross-Platform: Works seamlessly on Windows, macOS, and Linux.

## 🛡 Requirements

- Go 1.22 or higher (for building from source)
- Internet connection (obviously, for downloading the models)

## 🧑‍💻 Using with llama.cpp

Once you have downloaded a model, you can use it directly with llama.cpp:
```
$ llama-cli -m library-llama3.2-3b/model-dde5aa3fc5ff.gguf -p "We're no strangers to love"
```

## 🛠 Development

1. Clone the repository:

   ```bash
   git clone https://github.com/dimchansky/ollama-dl-go.git
   cd ollama-dl-go
   ```

2. Build and test:
   
   ```bash
   go build -o ollama-dl
   ./ollama-dl llama3.2:3b
   ```

3. Run the tests:

   ```bash
   go test ./...
   ```

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## ❤️ Contributions

Contributions are welcome! Feel free to open issues or submit pull requests to improve the functionality or add new features.

## 📞 Support

If you encounter any problems or have questions, please open an issue on GitHub or contact the maintainers.
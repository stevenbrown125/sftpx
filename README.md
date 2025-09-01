# SFTPX

A lightweight Go application that monitors a local folder for newly added **files or subfolders** and automatically uploads them (files up to 4GB) to a remote SFTP server.
The app runs as a standalone binary (`.exe` on Windows, or native builds on macOS/Linux) and is configurable via a JSON file.

---

## 🚀 Features

* Watches a directory for **new files and folders**.
* Configurable delays to let files “settle” before upload.
* Concurrent uploads with configurable worker count.
* Uploads files over SFTP (handles large files efficiently).
* Logs activity and errors to timestamped files in a configurable log directory.
* Behavior is controlled by `configs/config.json` (no rebuild needed).
* Supports both **password** and **SSH key authentication**.

---

## 🛠 Requirements

* **Go 1.21+** (to build from source)

  * [Download Go](https://go.dev/dl/)
* **Windows (x64)**, macOS, or Linux for running the published binary
* Network access to your SFTP server
* An SFTP **user account** on the target server with:

  * Either a password **or** an SSH private key (`id_rsa` or similar)
  * Permissions to create directories and upload files in the configured `remoteDir`
  * If using key authentication, the public key must be installed in the server’s `authorized_keys`

---

## 📂 Project Structure

```
sftpx/
├── cmd/
│   └── sftpx/            # Main entry point
├── internal/
│   ├── config/           # Config loader
│   ├── watcher/          # File/folder watcher logic
│   ├── sftp/             # SFTP helpers
│   └── logger/           # Logging setup
├── configs/
│   └── config.json       # App configuration (local, git-ignored)
├── dist/                 # Build output (git-ignored)
├── go.mod
├── .gitignore
└── README.md
```

---

## 📦 Dependencies

The project uses these Go libraries:

```bash
go get github.com/fsnotify/fsnotify
go get github.com/pkg/sftp
go get golang.org/x/crypto/ssh
```

---

## ⚙️ Config Example

```json
{
  "watchDir": "C:\watched",
  "remoteDir": "/upload",
  "sftp": {
    "host": "sftp.example.com",
    "port": 22,
    "user": "username",
    "password": "password"
  },
  "logDir": "./logs",
  "logFile": "sftpx.log",
  "delaySeconds": 5,
  "workers": 4
}
```

* `watchDir` → Local folder to monitor.
* `remoteDir` → Remote path on the SFTP server (POSIX style).
* `logDir` → Directory for logs. Created automatically if missing.
* `logFile` → Base log filename. Actual file will include a timestamp (e.g., `sftpx-2025-09-01_09-30-00.log`).
* `delaySeconds` → Wait time before uploading files/folders.
* `workers` → Number of concurrent uploads. Higher values improve throughput on fast networks/servers but can overwhelm slower servers or networks. Typical range: `2–8`. Start low and increase if you have many files and a high‑capacity SFTP server.
* Use `password` **or** `privateKeyPath`+`passphrase` for authentication.

### SSH Key Authentication Example

```json
{
  "watchDir": "C:\\watched",
  "remoteDir": "/upload",
  "sftp": {
    "host": "sftp.example.com",
    "port": 22,
    "user": "username",
    "privateKeyPath": "C:\\keys\\id_rsa",
    "passphrase": ""
  },
  "logDir": "./logs",
  "logFile": "sftpx.log",
  "delaySeconds": 15,
  "workers": 4
}
```

---

## 🏗 Building

A `Makefile` is provided to simplify builds and organize output into the `/dist` directory.

From the project root:

```bash
# Clean and rebuild for your host OS
make build

# Build for Windows 64-bit
make build-windows

# Build for Linux 64-bit
make build-linux

# Remove build artifacts
make clean
```

This will place binaries under `dist/` (e.g. `dist/sftpx` or `dist/sftpx.exe`).

---

## ▶️ Running

On macOS/Linux:

```bash
./dist/sftpx
```

On Windows:

```powershell
.\dist\sftpx.exe
```

The app will:

* Watch the configured folder (`watchDir`)
* Upload files/folders to the SFTP destination (`remoteDir`)
* Write logs into the configured `logDir` with timestamped filenames

---

## 🧪 Local Testing

You can test SFTPX locally using Docker.

### Password-based user

```bash
docker run --platform linux/amd64 -p 2222:22 -d \
  atmoz/sftp foo:pass::::upload
```

Then connect manually:

```bash
sftp -P 2222 foo@127.0.0.1
```

(password: `pass`)

Update your `config.json`:

```json
"sftp": {
  "host": "127.0.0.1",
  "port": 2222,
  "user": "foo",
  "password": "pass"
}
```

### Key-based user

Generate a key pair:

```bash
ssh-keygen -t rsa -b 4096 -f ./id_rsa_sftpx -N ""
```

Run the container with the public key:

```bash
docker run --platform linux/amd64 -p 2222:22 -d \
  -v $(pwd)/id_rsa_sftpx.pub:/home/foo/.ssh/keys/id_rsa.pub:ro \
  atmoz/sftp foo::::upload
```

Connect manually:

```bash
sftp -i ./id_rsa_sftpx -P 2222 foo@127.0.0.1
```

Update your `config.json`:

```json
"sftp": {
  "host": "127.0.0.1",
  "port": 2222,
  "user": "foo",
  "privateKeyPath": "./id_rsa_sftpx",
  "passphrase": ""
}
```

---

## 🐞 Issues & Support

If you encounter bugs, unexpected behavior, or have feature requests:

1. Check the [existing issues](../../issues) to see if it’s already reported.
2. If not, open a new issue and include:

   * OS and version (Windows, macOS, Linux)
   * SFTPX version (from the release tag)
   * Steps to reproduce the problem
   * Relevant configuration (without sensitive data)
   * Logs from `./logs/` showing the error

Clear, detailed reports help us resolve issues faster.

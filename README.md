# 🧩 Mocker — JSON-Driven Mock API Server

**Mocker** is a lightweight command-line tool that helps frontend or mobile developers **mock APIs instantly** — without waiting for the backend to be ready.  

Define your mock routes and responses in a simple JSON file, and Mocker automatically serves them over HTTP.

---

## 📚 Table of Contents

1. [🚀 Features](#-features)
2. [🧰 Installation](#-installation)
   - [macOS (Apple Silicon / ARM64)](#macos-apple-silicon--arm64)
   - [macOS (Intel / AMD64)](#macos-intel--amd64)
   - [Linux](#linux)
   - [Windows](#windows)
3. [🧱 Self-Build (for developers)](#-self-build-for-developers)
4. [💡 Usage](#-usage)
   - [🧭 CLI Options](#-cli-options)
5. [🔄 Updating Mocker](#-updating-mocker)
6. [🗑️ Uninstalling Mocker](#-uninstalling-mocker)
7. [📄 Example Config File](#-example-config-file)
8. [🧪 Example Session](#-example-session)
9. [🧱 Binary Releases](#-binary-releases)
10. [🧾 License](#-license)
11. [💬 Credits](#-credits)

---

## 🚀 Features

- 🔧 **Simple JSON-based configuration**
- ⚡ **Instant API mocking** — no setup, no database, no dependencies
- 🧠 **Custom routes, methods, status codes, and responses**
- 📦 **Single binary** — portable across macOS, Linux, and Windows
- 💬 **CLI options** for version, update, uninstall, and sample config generation
- 💾 **Auto-generate a sample config** with `--download=example.json`
- 🔄 **Self-updating support** via `--update`
- 🗑️ **Easy uninstall** via `--uninstall`

---

## 🧰 Installation

### **macOS (Apple Silicon / ARM64)**
```bash
curl -L -o mocker \
  https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-macos
chmod +x mocker
sudo mv mocker /usr/local/bin/mocker
````

### **macOS (Intel / AMD64)**

```bash
curl -L -o mocker \
  https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-macos-amd
chmod +x mocker
sudo mv mocker /usr/local/bin/mocker
```

### **Linux**

```bash
curl -L -o mocker \
  https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-linux
chmod +x mocker
sudo mv mocker /usr/local/bin/mocker
```

### **Windows**

1. Download [`mocker-windows.exe`](https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-windows.exe)
2. Move it to `C:\Tools\mocker\mocker.exe`
3. Add `C:\Tools\mocker` to your **PATH**
4. Open a new Command Prompt or PowerShell:

   ```powershell
   mocker.exe --path .\example.json
   ```

---

## 🧱 Self-Build (for developers)

You can also build `mocker` locally using Go (v1.20+ recommended):

```bash
git clone https://github.com/Yuddhaa/mocker.git
cd mocker

# Build for your own OS
go build -o mocker main.go
```

### Cross-compile for other systems

| Target OS             | Architecture | Command                                                                |
| --------------------- | ------------ | ---------------------------------------------------------------------- |
| Linux                 | amd64        | `GOOS=linux GOARCH=amd64 go build -o bin/mocker-linux main.go`         |
| macOS (Intel)         | amd64        | `GOOS=darwin GOARCH=amd64 go build -o bin/mocker-macos-amd main.go`    |
| macOS (Apple Silicon) | arm64        | `GOOS=darwin GOARCH=arm64 go build -o bin/mocker-macos main.go`        |
| Windows               | amd64        | `GOOS=windows GOARCH=amd64 go build -o bin/mocker-windows.exe main.go` |

---

## 💡 Usage

Run Mocker with your JSON config:

```bash
mocker --path=./example.json
```

Example output:

```
GET /users set
POST /login set
server is up and running at port: 8080
```

---

### 🧭 CLI Options

| Flag                                 | Description                                       |
| ------------------------------------ | ------------------------------------------------- |
| `--path <file>`                      | Path to your JSON config (default: `./example.json`) |
| `--version`                          | Print version info                                |
| `--help`                             | Show all flags and usage                          |
| `--download=<filename>`              | Generate example JSON config and exit             |
| `--update`                           | Update to the latest version                      |
| `--update --download_version=v1.1.0` | Update to a specific version                      |
| `--uninstall`                        | Uninstall Mocker from your system                 |

---

## 🔄 Updating Mocker

To update to the latest version:

```bash
mocker --update
```

To update to a specific version:

```bash
mocker --update --download_version=v1.1.0
```

---

## 🗑️ Uninstalling Mocker

If you installed via curl/mv:

```bash
sudo rm /usr/local/bin/mocker
```

If using the built-in flag:

```bash
mocker --uninstall
```


---

## 📄 Example Config File

Mocker uses a simple, human-readable **JSON configuration file** that defines:

* The port on which to run the mock server, and
* A list of HTTP routes with their methods, responses, and status codes.

You can generate this example automatically by running:

```bash
mocker --download=example.json
```

### 🧠 Config Structure Overview

| Key          | Type                 | Required | Description                                                                                        |
| ------------ | -------------------- | -------- | -------------------------------------------------------------------------------------------------- |
| **`port`**   | `string`  | ✅ Yes    | The TCP port where Mocker will start the HTTP server. Example: `"8080"`.                 |
| **`routes`** | `array (of route object)`  | ✅ Yes    | List of mock routes. Each object inside defines an API endpoint with a method, path, and response. |

Each **route** object supports the following fields:

| Key                   | Type                            | Required | Description                                                                                         |
| --------------------- | ------------------------------- | -------- | --------------------------------------------------------------------------------------------------- |
| **`path`**            | `string`                        | ✅ Yes    | The URL path to handle (e.g. `/api/users`). You can include path parameters like `/api/users/{id}`. |
| **`method`**          | `string`                        | ✅ Yes    | The HTTP method to match (`GET`, `POST`, `PATCH`, `PUT`, `DELETE`, etc.). Case-insensitive.         |
| **`response`**        | `object`                        | ✅ Yes    | Defines what Mocker returns when this route is called.                                              |
| **`response.status`** | `number`                        | ✅ Yes    | The HTTP status code to return (e.g. `200`, `201`, `404`, etc.).                                    |
| **`response.body`**   | `any` (i.e. `object` or `array` or `string`) | ✅ Yes    | The JSON body to send back in the response. Can be any valid JSON value.                            |

---

### 📘 Example Configuration

```json
{
  "port": "6969",
  "routes": [
    {
      "path": "/api/users",
      "method": "GET",
      "response": {
        "status": 200,
        "body": {
          "users": ["alice", "bob", "charlie"]
        }
      }
    },
    {
      "path": "/api/users",
      "method": "POST",
      "response": {
        "status": 201,
        "body": {
          "message": "User created successfully"
        }
      }
    },
    {
      "path": "/api/users/{id}",
      "method": "PATCH",
      "response": {
        "status": 200,
        "body": {
          "id": "{id}",
          "updated": true,
          "role": "admin"
        }
      }
    }
  ]
}
```

---

### 💡 Notes

* **Port must be a string.**
  
* **Dynamic path parameters:**
  You can use `{variable}` segments in your `path` (e.g. `/api/users/{id}`),
  and Mocker will match any value there -> **It supports dynamic routes to be mocked**.

* **Allowed response types: Any valid JSON**

  * Object `{}` — most common for structured JSON
  * Array `[]` — useful for lists
  * String `"message"` — for simple text responses
  * Number `123` or Boolean `true` — also supported

* **Methods supported:**
  `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`, `HEAD` (case-insensitive)

* **Status codes:**
  You can return any standard HTTP status code. **Must be a number**. (e.g. `200`, `201`, `400`, `401`, `404`, `500`).

* **Content type:**
  Mocker automatically sets `Content-Type: application/json` for all responses.

---


## 🧪 Example Session

```bash
mocker --download=example.json
# ✅ Example configuration downloaded to: example.json
# 🚀 Run it with: mocker --path=example.json

mocker --path=example.json
# GET /api/users set
# POST /api/users set
# PATCH /api/users/{id} set
# server is up and running at port: 6969
```

Now visit:

* `GET http://localhost:6969/api/users`
* `POST http://localhost:6969/api/users`
* `PATCH http://localhost:6969/api/users/123`

---

## 🧱 Binary Releases

| OS                    | Architecture | Binary                                                                                              |
| --------------------- | ------------ | --------------------------------------------------------------------------------------------------- |
| macOS (Apple Silicon) | arm64        | [mocker-macos](https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-macos)             |
| macOS (Intel)         | amd64        | [mocker-macos-amd](https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-macos-amd)     |
| Linux                 | amd64        | [mocker-linux](https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-linux)             |
| Windows               | amd64        | [mocker-windows.exe](https://github.com/Yuddhaa/mocker/releases/latest/download/mocker-windows.exe) |


---

## 🧾 License

MIT License © 2025 [Yuddhaa](https://github.com/Yuddhaa)

---

## 💬 Credits

Built in Go by [Yuddhaa](https://github.com/Yuddhaa)
Powering faster frontend development by removing backend blockers.






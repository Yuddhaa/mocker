// Mocker is a small CLI tool to spin up a mock HTTP server from a JSON config.
// It is meant to be used by frontend/mobile developers while the real backend
// is still being developed.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
)

// routesType represents a single mocked API route defined in the JSON config.
//
// Each route specifies:
//   - The HTTP method (e.g. GET, POST, PATCH)
//   - The request path (e.g. /api/users)
//   - The response object containing a status code and body
//
// Example JSON fragment:
//
//	{
//	  "method": "GET",
//	  "path": "/api/users",
//	  "response": {
//	    "status": 200,
//	    "body": {
//	      "users": ["alice", "bob", "charlie"]
//	    }
//	  }
//	}
type routesType struct {
	Method   string   `json:"method"`   // HTTP method to match (GET, POST, PATCH, etc.)
	Path     string   `json:"path"`     // HTTP path to match (supports static or parameterized paths like /api/users/{id})
	Response response `json:"response"` // Response definition containing status and body
}

// response defines the structure of the HTTP response returned for a mock route.
//
// Example JSON fragment:
//
//	"response": {
//	  "status": 201,
//	  "body": {
//	    "message": "User created successfully"
//	  }
//	}
type response struct {
	Status int `json:"status"` // HTTP status code to return (e.g. 200, 201, 404)
	Body   any `json:"body"`   // JSON body to return — can be object, array, string, number, or boolean
}

// inputType represents the top-level JSON configuration used by Mocker.
//
// Example JSON:
//
//	{
//	  "port": "8080",
//	  "routes": [ ... ]
//	}
type inputType struct {
	Port   string       `json:"port"`   // Port on which the mock server listens
	Routes []routesType `json:"routes"` // List of routes to configure
}

// appVersion is the version string printed by the --version flag.
var appVersion = "v1.0.0"

// main is the entry point of the Mocker CLI.
//
// It handles:
//   - CLI flags (version, help, path, download, update, uninstall)
//   - Optional generation of an example config
//   - Reading and parsing the JSON configuration
//   - Wiring up HTTP routes using chi
//   - Starting the HTTP server
func main() {
	var input inputType

	// CLI flags
	versionFlag := flag.Bool("version", false, "Print version info")
	path := flag.String("path", "./example.json", "path of the test json file")
	helpFlag := flag.Bool("help", false, "Show help message")
	downloadPath := flag.String("download", "", "Generate an example config file (usage: -download=example.json)")
	update := flag.Bool("update", false, "update to either latest or specific version")
	version := flag.String("download_verison", "", "mention specific version to be updated to")
	uninstall := flag.Bool("uninstall", false, "Can be used to uninstall mocker")
	flag.Parse()

	// Handle uninstall flow first so nothing else runs.
	if *uninstall {
		uninstallMocker()
		return
	}

	// Handle self-update flow before anything else.
	if *update {
		updateMocker(version)
		return
	}

	// Show help if requested.
	if *helpFlag {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	// Print version and exit.
	if *versionFlag {
		fmt.Println("Mocker", appVersion)
		return
	}

	// Generate an example JSON config and exit.
	if *downloadPath != "" {
		err := os.WriteFile(*downloadPath, []byte(exampleConfig), 0o644)
		if err != nil {
			fmt.Printf("Error creating example file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Example configuration downloaded to: %s\n", *downloadPath)
		fmt.Printf("🚀 Run it with: mocker --path=%s\n", *downloadPath)
		return // Exit so we don't start the server
	}

	// Read the JSON config from the provided path.
	data, err := os.ReadFile(*path)
	if err != nil {
		log.Fatalf("error in reading the file, err: %s", err.Error())
	}

	// Parse the JSON into inputType.
	if err := json.Unmarshal(data, &input); err != nil {
		log.Fatalf("error in Unmarshal of the JSON, err: %s", err.Error())
	}

	// Set up router and create handlers for each configured route.
	router := chi.NewRouter()
	for _, route := range input.Routes {
		v := route // copy to avoid closure capturing issues
		router.Method(strings.ToUpper(v.Method), v.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%v %v was called\n", r.Method, v.Path)
			if err := respondWithJSON(w, v.Response.Status, v.Response.Body); err != nil {
				log.Fatalf("err in responding with json, Error: %s\n", err.Error())
			}
		}))
		fmt.Printf("%v %v set\n", v.Method, v.Path)
	}

	// Start the HTTP server.
	fmt.Println("server is up and running at port: ", input.Port)
	log.Fatal(http.ListenAndServe(":"+input.Port, router))
}

// respondWithJSON marshals the given payload into JSON and writes it to the
// HTTP response with the given status code.
//
// It returns an error if the JSON marshaling fails.
func respondWithJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	return err
}

// updateMocker downloads and replaces the currently running Mocker binary
// with either the latest release or a specific version from GitHub.
//
// Behavior:
//   - If version == "" or "latest", the latest GitHub release is used.
//   - Otherwise, the given version string is used (e.g. "v1.1.0").
//   - On macOS/Linux, the running binary is replaced in-place.
//   - On Windows, the new binary is downloaded but must be manually swapped.
func updateMocker(version *string) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Determine which binary asset name to download based on OS + arch.
	var binaryName string
	switch osName {
	case "darwin":
		if arch == "arm64" {
			binaryName = "mocker-macos"
		} else {
			binaryName = "mocker-macos-amd"
		}
	case "linux":
		binaryName = "mocker-linux"
	case "windows":
		binaryName = "mocker-windows.exe"
	default:
		log.Fatalf("unsupported OS: %s/%s", osName, arch)
	}

	// Choose download URL.
	var url string
	if *version == "" || *version == "latest" {
		// Use the latest release asset.
		url = fmt.Sprintf("https://github.com/Yuddhaa/mocker/releases/latest/download/%s", binaryName)
		fmt.Println("🔄 Updating to the latest version...")
	} else {
		// Use a specific tagged release.
		url = fmt.Sprintf("https://github.com/Yuddhaa/mocker/releases/download/%s/%s", *version, binaryName)
		fmt.Printf("🔄 Updating to version %s...\n", *version)
	}

	// Perform the HTTP GET request.
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to download binary, HTTP %d", resp.StatusCode)
	}

	// Temporary file to hold the new binary before replacement.
	tmpFile := "mocker_new"
	if osName == "windows" {
		tmpFile += ".exe"
	}

	out, err := os.Create(tmpFile)
	if err != nil {
		log.Fatalf("failed to create temp file: %v", err)
	}
	defer out.Close()

	// Stream the response body into the temporary file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("failed to write binary: %v", err)
	}

	// Set executable permissions (best-effort, mainly for Unix-like systems).
	if err := out.Chmod(0o755); err != nil {
		log.Printf("warning: could not set permissions: %v", err)
	}

	// Find the path of the currently running executable.
	currentPath, err := os.Executable()
	fmt.Println("currentPath", currentPath)
	if err != nil {
		log.Fatalf("cannot locate current executable: %v", err)
	}

	// Replace the current binary.
	if osName == "windows" {
		// Windows cannot replace the running .exe directly.
		fmt.Printf("✅ Downloaded new version to %s\n", tmpFile)
		fmt.Println("⚠️ On Windows, please manually replace the old binary with the new one.")
	} else {
		err = os.Rename(tmpFile, currentPath)
		if err != nil {
			// Check specifically for permission errors
			if os.IsPermission(err) {
				fmt.Println("🔒 Permission denied while trying to replace the binary.")
				fmt.Printf("Current location: %s\n", currentPath)
				fmt.Printf("New binary:       %s\n", tmpFile)
				fmt.Println()
				fmt.Print("Would you like Mocker to try replacing it using sudo? (y/N): ")

				var confirm string
				fmt.Scanln(&confirm)
				confirm = strings.ToLower(strings.TrimSpace(confirm))

				if confirm != "y" && confirm != "yes" {
					fmt.Println("🟡 Skipped automatic sudo replacement.")
					fmt.Println("You can manually run:")
					fmt.Printf("  sudo mv %s %s\n", tmpFile, currentPath)
					return
				}

				fmt.Println("🔧 Attempting to elevate privileges with sudo...")

				// Try to move using sudo
				cmd := exec.Command("sudo", "mv", tmpFile, currentPath)

				// Connect standard IO so the user can see the sudo password prompt
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					// Clean up temp file if sudo fails
					_ = os.Remove(tmpFile)
					log.Fatalf("❌ Failed to update even with sudo: %v.\nPlease try manually running:\n  sudo mv %s %s\n", err, tmpFile, currentPath)
				}

				fmt.Println("✅ Successfully updated Mocker!")
				return
			}

			// Handle other types of rename errors
			log.Fatalf("❌ Failed to replace binary: %v", err)
		}

		fmt.Println("✅ Successfully updated Mocker!")
	}
}

// uninstallMocker attempts to remove the currently running Mocker binary
// from the filesystem.
//
// Behavior:
//   - Prompts the user for confirmation.
//   - On Unix-like systems, deletes the executable in-place.
//   - On Windows, prints instructions for manual deletion because a running
//     .exe cannot delete itself.
func uninstallMocker() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Unable to locate current executable:", err)
		return
	}

	fmt.Printf("⚠️  This will remove Mocker from your system.\nLocation: %s\n", execPath)
	fmt.Print("Are you sure you want to uninstall Mocker? (y/N): ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("🟡 Uninstall cancelled.")
		return
	}

	switch runtime.GOOS {
	case "windows":
		fmt.Println("⚠️ On Windows, an application cannot delete itself while running.")
		fmt.Println("Please manually delete this file:")
		fmt.Println(execPath)
		return

	default:
		err := os.Remove(execPath)
		if err != nil {
			// Handle permission errors gracefully and optionally elevate with sudo.
			if os.IsPermission(err) {
				fmt.Println("❌ Permission denied while trying to uninstall Mocker.")
				fmt.Printf("Binary location: %s\n", execPath)
				fmt.Println()
				fmt.Print("Would you like Mocker to try removing it using sudo? (y/N): ")

				var confirm string
				fmt.Scanln(&confirm)
				confirm = strings.ToLower(strings.TrimSpace(confirm))

				if confirm != "y" && confirm != "yes" {
					fmt.Println("🟡 Skipped automatic sudo removal.")
					fmt.Println("You can manually run:")
					fmt.Printf("  sudo rm %s\n", execPath)
					return
				}

				fmt.Println("🔧 Attempting to elevate privileges with sudo...")

				cmd := exec.Command("sudo", "rm", execPath)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					fmt.Printf("❌ Failed to uninstall with sudo: %v\n", err)
					fmt.Println("You can try manually running:")
					fmt.Printf("  sudo rm %s\n", execPath)
					return
				}

				fmt.Println("✅ Successfully uninstalled Mocker.")
				return
			}

			// Non-permission-related errors.
			fmt.Printf("❌ Failed to uninstall: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully uninstalled Mocker.")
	}
}

// exampleConfig is the built-in example configuration used when the user
// passes the --download flag.
//
// Users can run:
//
//	mocker --download=example.json
//	mocker --path=example.json
const exampleConfig = `{
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
}`

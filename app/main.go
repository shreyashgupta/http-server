package main

import (
	"net"
	"os"

	httpserver "github.com/codecrafters-io/http-server-starter-go/http-server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var directory = ""

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path) // Go 1.16+
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {

	args := os.Args[1:] // skip program name

	if len(args) == 2 {
		directory = args[1]
	}

	mux := httpserver.NewMux()
	mux.Get("/", defaultHandler)
	mux.Get("/echo/{echo_str}", echoHandler)
	mux.Get("/user-agent", userAgentHandler)
	mux.Get("/files/{file_name}", fileHandler)
	mux.Post("/files/{file_name}", filePostHandler)

	server := httpserver.NewHttpServer(mux, "0.0.0.0:4221")

	server.Serve()
}

func defaultHandler(_ *httpserver.Request, w *httpserver.Writer) {
	w.Write()
}

func echoHandler(r *httpserver.Request, w *httpserver.Writer) {
	echoStr := r.GetCapture("echo_str")
	w.SetHeader("Content-Type", "text/plain")
	w.SetContent(echoStr)
	w.Write()
}

func userAgentHandler(r *httpserver.Request, w *httpserver.Writer) {
	userAgent, _ := r.GetHeader("User-Agent")
	w.SetHeader("Content-Type", "text/plain")
	w.SetContent(userAgent)
	w.Write()
}

func fileHandler(r *httpserver.Request, w *httpserver.Writer) {
	absFilePath := directory + r.GetCapture("file_name")
	fileContent, err := readFile(absFilePath)
	if err != nil {
		w.SetStatusCode(httpserver.HTTP_NOT_FOUND)
		w.Write()
		return
	}
	w.SetHeader("Content-Type", "application/octet-stream")
	w.SetContent(fileContent)
	w.Write()
}

func writeToFile(path string, data []byte) error {
	// Create file (or truncate if it already exists)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func filePostHandler(r *httpserver.Request, w *httpserver.Writer) {
	absFilePath := directory + r.GetCapture("file_name")
	err := writeToFile(absFilePath, r.GetBodyData())
	if err != nil {
		w.SetStatusCode(httpserver.HTTP_BAD_REQUEST)
		w.Write()
		return
	}

	w.SetStatusCode(httpserver.HTTP_CREATED)
	w.Write()
}

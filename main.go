package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

var version = "2.0.0"

type PageData struct {
	Runtime      string
	RuntimeClass string
	Hostname     string
	PodName      string
	Namespace    string
	Version      string
}

func detectRuntime() (runtime, runtimeClass string) {
	// Check for Kubernetes environment variables
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return "Kubernetes", "kubernetes"
	}

	// Check for pod name (usually set in K8s deployments)
	if os.Getenv("POD_NAME") != "" || os.Getenv("HOSTNAME") != "" {
		if strings.HasPrefix(os.Getenv("HOSTNAME"), "heartflow-demo-") {
			return "Kubernetes", "kubernetes"
		}
	}

	// Check if running under systemd
	if os.Getenv("INVOCATION_ID") != "" {
		return "systemd", "systemd"
	}

	// Fallback check for systemd by checking parent process
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		// Running on a systemd system, check if we're a service
		if os.Getppid() == 1 {
			return "systemd", "systemd"
		}
	}

	// Default fallback
	return "systemd", "systemd"
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}

	runtime, runtimeClass := detectRuntime()
	hostname, _ := os.Hostname()

	data := PageData{
		Runtime:      runtime,
		RuntimeClass: runtimeClass,
		Hostname:     hostname,
		PodName:      getEnvOrDefault("POD_NAME", getEnvOrDefault("HOSTNAME", hostname)),
		Namespace:    getEnvOrDefault("POD_NAMESPACE", "N/A"),
		Version:      version,
	}

	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func main() {
	port := getEnvOrDefault("PORT", "8080")

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/healthz", healthHandler)

	log.Printf("HeartFlow Demo Service starting on port %s", port)
	log.Printf("Runtime: %s", func() string { r, _ := detectRuntime(); return r }())

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

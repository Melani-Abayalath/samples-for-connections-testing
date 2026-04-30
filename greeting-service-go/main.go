package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/greeter/greet", greet)
	serverMux.HandleFunc("/greeter/env", getEnvVars)

	serverPort := 9094
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: serverMux,
	}
	go func() {
		log.Printf("Starting HTTP Greeter on port %d\n", serverPort)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP ListenAndServe error: %v", err)
		}
		log.Println("HTTP server stopped serving new requests.")
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Shutdown complete.")
}

func greet(w http.ResponseWriter, r *http.Request) {
	// Choreo injects these at deploy time based on the project-level connection in component.yaml
	serviceURL   := os.Getenv("CHOREO_PROJECT_LEVL_BALLERINA_GREETING_CON_SERVICEURL") + "/greeting"
	tokenURL     := os.Getenv("CHOREO_PROJECT_LEVL_BALLERINA_GREETING_CON_TOKENURL")
	clientID     := os.Getenv("CHOREO_PROJECT_LEVL_BALLERINA_GREETING_CON_CONSUMERKEY")
	clientSecret := os.Getenv("CHOREO_PROJECT_LEVL_BALLERINA_GREETING_CON_CONSUMERSECRET")
	choreoapikey = os:getEnv("CHOREO_PROJECT_LEVL_BALLERINA_GREETING_CON_CHOREOAPIKEY");

	fmt.Printf("serviceURL: %s\n", serviceURL)
	fmt.Printf("Client ID: %s\n", clientID)

	clientCredsConfig := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	client := clientCredsConfig.Client(context.Background())
	response, err := client.Get(serviceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Server returned non-200 status: %d %s", response.StatusCode, response.Status), response.StatusCode)
		return
	}

	_, err = io.Copy(w, response.Body)
	if err != nil {
		log.Printf("Error writing response body to client: %v\n", err)
	}
}

func getEnvVars(w http.ResponseWriter, r *http.Request) {
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		pair := splitEnv(env)
		envVars[pair[0]] = pair[1]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(envVars); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func splitEnv(env string) [2]string {
	var pair [2]string
	for i, char := range env {
		if char == '=' {
			pair[0] = env[:i]
			pair[1] = env[i+1:]
			break
		}
	}
	return pair
}

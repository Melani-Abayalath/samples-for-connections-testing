/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
        "io"
        "golang.org/x/oauth2/clientcredentials"
	"encoding/json"
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
	<-stopCh // Wait for shutdown signal

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Shutdown complete.")
}

// func greet(w http.ResponseWriter, r *http.Request) {
// 	// Read environment variables
// 	serviceURL := os.Getenv("CHOREO_GQL_CON_5_SERVICEURL") // NOTE: no suffix; this is the GraphQL endpoint
// 	tokenURL := os.Getenv("CHOREO_GQL_CON_5_TOKENURL")
// 	clientSecret := os.Getenv("CHOREO_GQL_CON_5_CONSUMERSECRET")
// 	clientID := os.Getenv("CHOREO_GQL_CON_5_CONSUMERKEY")
// 	apiKey := os.Getenv("CHOREO_GQL_CON_5_CHOREOAPIKEY")

// 	fmt.Printf("Client ID: %s\n", clientID)
// 	fmt.Printf("serviceURL: %s\n", serviceURL)
// 	fmt.Printf("apiKey (set?): %t\n", apiKey != "")

// 	// OAuth2 client credentials
// 	var clientCredsConfig = clientcredentials.Config{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		TokenURL:     tokenURL,
// 	}
// 	httpClient := clientCredsConfig.Client(context.Background())

// 	// GraphQL request payload
// 	type gqlRequest struct {
// 		Query string `json:"query"`
// 	}
// 	type gqlError struct {
// 		Message string `json:"message"`
// 	}
// 	type gqlResponse struct {
// 		Data struct {
// 			Greeting string `json:"greeting"`
// 		} `json:"data"`
// 		Errors []gqlError `json:"errors"`
// 	}

// 	payload := gqlRequest{
// 		Query: `query GreetWorld { greeting(name: "Hansiiiii") }`,
// 	}
// 	bodyBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to marshal GraphQL request: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Build POST request
// 	req, err := http.NewRequest(http.MethodPost, serviceURL, bytes.NewReader(bodyBytes))
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	// If your gateway expects an API key, set it (adjust header name if your gateway uses a different one)
// 	if apiKey != "" {
// 		req.Header.Set("x-api-key", apiKey)
// 	}

// 	// Execute
// 	resp, err := httpClient.Do(req)
// 	if err != nil {
//         fmt.Printf("Eror while calling consumer")
// 		http.Error(w, fmt.Sprintf("Error making GraphQL request: %v", err), http.StatusBadGateway)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Read and decode
// 	respBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusBadGateway)
// 		return
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		http.Error(w, fmt.Sprintf("Upstream status %d: %s", resp.StatusCode, string(respBody)), resp.StatusCode)
// 		return
// 	}

// 	var gqlResp gqlResponse
// 	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to decode GraphQL response: %v", err), http.StatusBadGateway)
// 		return
// 	}
// 	if len(gqlResp.Errors) > 0 {
// 		// Aggregate GraphQL errors
// 		msg := "GraphQL errors:"
// 		for _, e := range gqlResp.Errors {
// 			msg += " " + e.Message + ";"
// 		}
// 		http.Error(w, msg, http.StatusBadGateway)
// 		return
// 	}

// 	// Return the greeting (plain text). Change to JSON if you prefer.
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	if _, err := w.Write([]byte(gqlResp.Data.Greeting)); err != nil {
// 		log.Printf("Error writing response to client: %v\n", err)
// 	}
// }


func greet(w http.ResponseWriter, r *http.Request) {
	// Read environment variables
	serviceURL := os.Getenv("CHOREO_TESTCON6_SERVICEURL")+"/greeting"
	tokenURL := os.Getenv("CHOREO_TESTCON6_TOKENURL")
	clientSecret := os.Getenv("CHOREO_TESTCON6_CONSUMERSECRET")
	clientID := os.Getenv("CHOREO_TESTCON6_CONSUMERKEY")
	apiKey := os.Getenv("CHOREO_TESTCON6_CHOREOAPIKEY")
	// Log the client ID for debugging
	fmt.Printf("Client ID: %s\n", clientID)
	fmt.Printf("serviceURL: %s\n", serviceURL)
	fmt.Printf("apiKey: %s\n", apiKey)

	// // Prepare a JSON response with the environment variables
	// envVars := map[string]string{
	// 	"SVC_URL":         serviceURL,
	// 	"TOKEN_URL":       tokenURL,
	// 	"CONSUMER_SECRET": clientSecret,
	// 	"CONSUMER_KEY":    clientID,
	// }

	// // Set the content type to JSON and write the response
	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(envVars); err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	// }
	
	var clientCredsConfig = clientcredentials.Config{
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
	// Check the response status code
	if response.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Server returned non-200 status: %d %s", response.StatusCode, response.Status), response.StatusCode)
		return
	}
	// Copy the response body directly to the client
	_, err = io.Copy(w, response.Body)
	if err != nil {
		// If there's an error writing the response body to the client, log it
		log.Printf("Error writing response body to client: %v\n", err)
	}
}

// Returns all environment variables as a JSON response
func getEnvVars(w http.ResponseWriter, r *http.Request) {
	envVars := make(map[string]string)

	// Retrieve all environment variables
	for _, env := range os.Environ() {
		pair := splitEnv(env)
		envVars[pair[0]] = pair[1]
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the environment variables as JSON and send the response
	if err := json.NewEncoder(w).Encode(envVars); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

// Helper function to split environment variables into key-value pairs
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Schema struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ValidationRequest struct {
	SchemaName string      `json:"schema_name"`
	Data       interface{} `json:"data"`
}

type ValidationResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

func main() {
	log.Println("Schema Service starting...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "schema-service",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Schema management endpoints
	r.Route("/schemas", func(r chi.Router) {
		r.Get("/", handleGetSchemas)
		r.Post("/", handleCreateSchema)
		r.Get("/{name}", handleGetSchema)
		r.Put("/{name}", handleUpdateSchema)
		r.Delete("/{name}", handleDeleteSchema)
	})

	// Validation endpoints
	r.Route("/validate", func(r chi.Router) {
		r.Post("/", handleValidateData)
		r.Post("/{schema}", handleValidateWithSchema)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Schema Service listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Failed to start Schema Service:", err)
	}
}

// Placeholder handlers - implement these based on your requirements
func handleGetSchemas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	schemas := []Schema{
		{
			Name:        "user",
			Version:     "1.0.0",
			Description: "User schema validation",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "message",
			Version:     "1.0.0",
			Description: "Message schema validation",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"schemas": schemas,
		"total":   len(schemas),
	})
}

func handleCreateSchema(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Schema created successfully",
	})
}

func handleGetSchema(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "name")
	w.Header().Set("Content-Type", "application/json")

	// Return a sample schema
	schema := Schema{
		Name:        schemaName,
		Version:     "1.0.0",
		Description: fmt.Sprintf("Schema for %s", schemaName),
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":   map[string]string{"type": "string"},
				"name": map[string]string{"type": "string"},
			},
			"required": []string{"id", "name"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	json.NewEncoder(w).Encode(schema)
}

func handleUpdateSchema(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Schema updated successfully",
	})
}

func handleDeleteSchema(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Schema deleted successfully",
	})
}

func handleValidateData(w http.ResponseWriter, r *http.Request) {
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationResponse{
			Valid:  false,
			Errors: []string{"Invalid JSON payload"},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Basic validation - always return valid for now
	json.NewEncoder(w).Encode(ValidationResponse{
		Valid:  true,
		Errors: nil,
	})
}

func handleValidateWithSchema(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schema")

	var data interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationResponse{
			Valid:  false,
			Errors: []string{"Invalid JSON payload"},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Basic validation - always return valid for now
	json.NewEncoder(w).Encode(ValidationResponse{
		Valid:  true,
		Errors: nil,
	})

	log.Printf("Validated data against schema: %s", schemaName)
}

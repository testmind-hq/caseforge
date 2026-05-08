// internal/sandbox/handler.go
package sandbox

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// resourceTypeFromPath extracts the primary resource name from an OpenAPI path.
// "/pets/{petId}" → "pets", "/users/{userId}/posts/{postId}" → "users_posts"
func resourceTypeFromPath(path string) string {
	var segments []string
	for _, seg := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		if seg != "" && !strings.HasPrefix(seg, "{") {
			segments = append(segments, seg)
		}
	}
	if len(segments) == 0 {
		return "resource"
	}
	return strings.Join(segments, "_")
}

// idParamName returns the name of the first path parameter in an operation, or "id".
func idParamName(op *spec.Operation) string {
	for _, p := range op.Parameters {
		if p.In == "path" {
			return p.Name
		}
	}
	return "id"
}

// hasPathParam reports whether the operation has any path parameters.
func hasPathParam(op *spec.Operation) bool {
	for _, p := range op.Parameters {
		if p.In == "path" {
			return true
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg, "status": status})
}

// makeHandler returns an http.HandlerFunc for the given operation.
func makeHandler(op *spec.Operation, store StateStore, gen *ResponseGenerator, logger *slog.Logger) http.HandlerFunc {
	resourceType := resourceTypeFromPath(op.Path)
	idParam := idParamName(op)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := strings.ToUpper(op.Method)

		switch method {
		case "DELETE":
			id := r.PathValue(idParam)
			if _, ok := store.Read(resourceType, id); !ok {
				writeError(w, http.StatusNotFound, "resource not found")
				logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 404, "latency", time.Since(start))
				return
			}
			store.Delete(resourceType, id)
			w.WriteHeader(http.StatusNoContent)
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 204, "latency", time.Since(start))

		case "GET":
			if hasPathParam(op) {
				id := r.PathValue(idParam)
				obj, ok := store.Read(resourceType, id)
				if !ok {
					writeError(w, http.StatusNotFound, "resource not found")
					logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 404, "latency", time.Since(start))
					return
				}
				writeJSON(w, http.StatusOK, obj)
				logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 200, "latency", time.Since(start))
			} else {
				list := store.List(resourceType)
				writeJSON(w, http.StatusOK, list)
				logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 200, "latency", time.Since(start))
			}

		default: // POST, PUT, PATCH
			if r.Header.Get("Content-Type") != "" && !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
				logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", 415, "latency", time.Since(start))
				return
			}
			pathParams := map[string]string{}
			if hasPathParam(op) {
				pathParams[idParam] = r.PathValue(idParam)
			}
			status, body := gen.Generate(op, pathParams)
			id := extractID(body, op.Path)
			if id == "" {
				id = fmt.Sprintf("%v", body["id"])
			}
			if method == "POST" {
				store.Write(resourceType, id, body)
			}
			w.Header().Set("X-Sandbox-ID", id)
			writeJSON(w, status, body)
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", status, "latency", time.Since(start))
		}
	}
}

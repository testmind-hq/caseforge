package datagen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParsePostmanCollection_ExtractsBodyFields(t *testing.T) {
	col := map[string]any{
		"info": map[string]any{
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"item": []any{
			map[string]any{
				"name": "Create user",
				"request": map[string]any{
					"method": "POST",
					"body": map[string]any{
						"mode": "raw",
						"raw":  `{"name": "Alice", "age": 30}`,
					},
				},
			},
		},
	}
	data, _ := json.Marshal(col)
	tmp := t.TempDir()
	path := filepath.Join(tmp, "collection.json")
	os.WriteFile(path, data, 0644)

	pool, err := ParsePostmanCollection(path)
	if err != nil {
		t.Fatalf("ParsePostmanCollection: %v", err)
	}
	if v, ok := pool.ValueFor("name"); !ok || v != "Alice" {
		t.Errorf("name = %v %v, want Alice true", v, ok)
	}
	if v, ok := pool.ValueFor("age"); !ok {
		t.Errorf("age not found: %v %v", v, ok)
	}
}

func TestParsePostmanCollection_RecursesFolders(t *testing.T) {
	col := map[string]any{
		"item": []any{
			map[string]any{
				"name": "Users folder",
				"item": []any{
					map[string]any{
						"name": "Create user",
						"request": map[string]any{
							"method": "POST",
							"body": map[string]any{
								"mode": "raw",
								"raw":  `{"email": "test@example.com"}`,
							},
						},
					},
				},
			},
		},
	}
	data, _ := json.Marshal(col)
	tmp := t.TempDir()
	path := filepath.Join(tmp, "collection.json")
	os.WriteFile(path, data, 0644)

	pool, err := ParsePostmanCollection(path)
	if err != nil {
		t.Fatalf("ParsePostmanCollection: %v", err)
	}
	if v, ok := pool.ValueFor("email"); !ok || v != "test@example.com" {
		t.Errorf("email = %v %v, want test@example.com true", v, ok)
	}
}

func TestParsePostmanCollection_SkipsNonJSON(t *testing.T) {
	col := map[string]any{
		"item": []any{
			map[string]any{
				"name": "Upload",
				"request": map[string]any{
					"method": "POST",
					"body": map[string]any{
						"mode": "formdata",
						"raw":  "",
					},
				},
			},
		},
	}
	data, _ := json.Marshal(col)
	tmp := t.TempDir()
	path := filepath.Join(tmp, "collection.json")
	os.WriteFile(path, data, 0644)

	pool, err := ParsePostmanCollection(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool.Len() != 0 {
		t.Errorf("expected empty pool for non-JSON body, got %d fields", pool.Len())
	}
}

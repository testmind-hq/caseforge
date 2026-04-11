package datagen

import (
	"path/filepath"
	"testing"
)

func TestDataPool_AddAndValueFor(t *testing.T) {
	p := NewDataPool()
	p.Add("name", "Alice")
	p.Add("age", float64(30))

	if v, ok := p.ValueFor("name"); !ok || v != "Alice" {
		t.Errorf("ValueFor name = %v %v, want Alice true", v, ok)
	}
	if v, ok := p.ValueFor("age"); !ok || v != float64(30) {
		t.Errorf("ValueFor age = %v %v, want 30 true", v, ok)
	}
	if _, ok := p.ValueFor("unknown"); ok {
		t.Error("expected false for unknown field")
	}
}

func TestDataPool_CaseInsensitive(t *testing.T) {
	p := NewDataPool()
	p.Add("UserName", "Bob")
	if v, ok := p.ValueFor("username"); !ok || v != "Bob" {
		t.Errorf("case-insensitive lookup failed: %v %v", v, ok)
	}
}

func TestDataPool_SaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "pool.json")

	p := NewDataPool()
	p.Add("email", "test@example.com")
	p.Add("count", float64(5))
	if err := p.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p2, err := LoadDataPool(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if v, ok := p2.ValueFor("email"); !ok || v != "test@example.com" {
		t.Errorf("after round-trip email = %v %v", v, ok)
	}
}

func TestDataPool_Merge(t *testing.T) {
	p1 := NewDataPool()
	p1.Add("x", "a")
	p2 := NewDataPool()
	p2.Add("y", "b")
	p1.Merge(p2)
	if _, ok := p1.ValueFor("y"); !ok {
		t.Error("merge failed: y not found in p1 after merge")
	}
}

func TestDataPool_Len(t *testing.T) {
	p := NewDataPool()
	if p.Len() != 0 {
		t.Errorf("empty pool Len = %d, want 0", p.Len())
	}
	p.Add("a", 1)
	p.Add("b", 2)
	if p.Len() != 2 {
		t.Errorf("Len = %d, want 2", p.Len())
	}
}

package collection

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
)

func defaultCfg(name string) Config {
	return Config{
		Name:           name,
		Dim:            16,
		Metric:         MetricL2,
		M:              8,
		EfConstruction: 50,
	}
}

func TestCreate_Success(t *testing.T) {
	m := NewManager()
	cfg := defaultCfg("c1")
	if err := m.Create(cfg); err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := m.Get("c1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Config != cfg {
		t.Errorf("Config mismatch: got %+v, want %+v", got.Config, cfg)
	}
	if got.Index == nil {
		t.Errorf("Index is nil")
	}
}

func TestCreate_Duplicate(t *testing.T) {
	m := NewManager()
	if err := m.Create(defaultCfg("c1")); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	err := m.Create(defaultCfg("c1"))
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second Create: got %v, want ErrAlreadyExists", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	m := NewManager()
	_, err := m.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Get: got %v, want ErrNotFound", err)
	}
}

func TestDrop(t *testing.T) {
	m := NewManager()
	if err := m.Create(defaultCfg("c1")); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := m.Drop("c1"); err != nil {
		t.Fatalf("Drop: %v", err)
	}
	_, err := m.Get("c1")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Get after Drop: got %v, want ErrNotFound", err)
	}
}

func TestDrop_NotFound(t *testing.T) {
	m := NewManager()
	err := m.Drop("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Drop: got %v, want ErrNotFound", err)
	}
}

func TestList(t *testing.T) {
	m := NewManager()
	if got := m.List(); len(got) != 0 {
		t.Errorf("empty List: got %v, want []", got)
	}
	want := []string{"a", "b", "c"}
	for _, n := range want {
		if err := m.Create(defaultCfg(n)); err != nil {
			t.Fatalf("Create %s: %v", n, err)
		}
	}
	got := m.List()
	sort.Strings(got)
	if len(got) != len(want) {
		t.Fatalf("List len: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List[%d]: got %s, want %s", i, got[i], want[i])
		}
	}
}

func TestConcurrentCreate(t *testing.T) {
	m := NewManager()
	const N = 100
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			if err := m.Create(defaultCfg(fmt.Sprintf("c%d", i))); err != nil {
				t.Errorf("Create c%d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()
	if got := len(m.List()); got != N {
		t.Errorf("List len after concurrent Create: got %d, want %d", got, N)
	}
}

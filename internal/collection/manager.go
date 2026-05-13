package collection

import (
	"errors"
	"sync"

	"github.com/zhangchuqi1998/velosearch/internal/hnsw"
)

var (
	ErrAlreadyExists = errors.New("collection already exists")
	ErrNotFound      = errors.New("collection not found")
	ErrDimMismatch   = errors.New("vector dimension does not match collection")
)

type Collection struct {
	Config Config
	Index  *hnsw.Index
}

type Manager struct {
	mu   sync.RWMutex
	cols map[string]*Collection
}

func NewManager() *Manager {
	return &Manager{cols: make(map[string]*Collection)}
}

func (m *Manager) Create(cfg Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cols[cfg.Name]; ok {
		return ErrAlreadyExists
	}
	idx := hnsw.NewIndex(cfg.Dim, cfg.M, cfg.EfConstruction, cfg.Metric.DistanceFunc())
	m.cols[cfg.Name] = &Collection{Config: cfg, Index: idx}
	return nil
}

func (m *Manager) Get(name string) (*Collection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.cols[name]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (m *Manager) Drop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cols[name]; !ok {
		return ErrNotFound
	}
	delete(m.cols, name)
	return nil
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.cols))
	for n := range m.cols {
		names = append(names, n)
	}
	return names
}

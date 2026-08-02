package source

import (
	"os"
	"slices"
	"sync"
)

// Manager handles the lifecycle and tracking of source files across the workspace.
// It acts as the central repository for code state and is fully thread-safe
// for concurrent access by Language Server background routines.
type Manager struct {
	mu    sync.RWMutex
	files map[string]*SourceText
}

// NewManager creates and initializes a new thread-safe file manager.
func NewManager() *Manager {
	return &Manager{
		files: make(map[string]*SourceText),
	}
}

// Update registers a new file or overwrites an existing one with new contents.
// It clones the text buffer to safely isolate it from external shared memory (e.g. LSP buffers).
func (m *Manager) Update(path string, version int32, text []byte) *SourceText {
	m.mu.Lock()
	defer m.mu.Unlock()

	file := NewSourceText(path, version, slices.Clone(text))
	m.files[path] = file
	return file
}

// Delete removes a file from the manager's tracking system.
// This is commonly triggered by textDocument/didClose events in an LSP.
func (m *Manager) Delete(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.files, path)
}

// SourceText safely retrieves a tracked file by its absolute or relative path.
// The boolean return value indicates whether the file was successfully found.
func (m *Manager) File(path string) (*SourceText, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.files[path]
	return f, ok
}

// Open reads a file directly from the filesystem, tracks it, and returns the instance.
func (m *Manager) Open(path string) (*SourceText, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	file := NewSourceText(path, 0, text)
	m.files[path] = file
	return file, nil
}

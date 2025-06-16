package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

type Storage[T any] interface {
	Get(id int64) *T
	Set(id int64, value *T) error
}

var _ Storage[any] = (*MemoryStorage[any])(nil)

type MemoryStorage[T any] struct {
	mutex sync.RWMutex
	file  *os.File
	data  map[int64]*T
}

func NewMemory[T any](filename string) (*MemoryStorage[T], error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		slog.Error("failed to open file", "file", filename, "error", err)
		return nil, fmt.Errorf("opening file %s: %w", filename, err)
	}

	var data map[int64]*T

	dec := gob.NewDecoder(file)
	if err := dec.Decode(&data); err != nil {
		slog.Debug("file has no data")
		data = make(map[int64]*T)
	}

	return &MemoryStorage[T]{
		file: file,
		data: data,
	}, nil
}

var ErrNoFileToSave = errors.New("no file to save data")

func (c *MemoryStorage[T]) Get(id int64) *T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value := c.data[id]

	if value == nil {
		slog.Debug("no data found for id", "id", id)
		return new(T)
	}

	return value
}

func (c *MemoryStorage[T]) Set(id int64, value *T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// slog.Debug("setting data for id", "id", id, "value", value)

	if c.data == nil {
		c.data = make(map[int64]*T)
	}

	c.data[id] = value

	if c.file == nil {
		slog.Error("no file to save")
		return ErrNoFileToSave
	}

	c.file.Seek(0, io.SeekStart)
	enc := gob.NewEncoder(c.file)
	if err := enc.Encode(c.data); err != nil {
		slog.Error("failed to encode data", "error", err)
		return fmt.Errorf("saving data to file: %w", err)
	}

	return nil
}

package platform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spaolacci/murmur3"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

type Ticker struct {
	ticks      []string
	timestamps []int64
}

func NewTicker() *Ticker {
	return &Ticker{}
}

func (t *Ticker) Tick(name string) {
	t.ticks = append(t.ticks, name)
	t.timestamps = append(t.timestamps, time.Now().UnixNano())
}

func (t *Ticker) String() string {
	var buf bytes.Buffer
	buf.WriteString("ticker")
	var prev int64
	for i := range t.ticks {
		n, ts := t.ticks[i], t.timestamps[i]
		if prev == 0 {
			prev = ts
		}
		buf.WriteString(fmt.Sprintf("|%s:%0.2fms", n, float32(ts-prev)/float32(time.Millisecond.Nanoseconds())))
		prev = ts
	}
	return buf.String()
}

// OpenFile is a general function to take a context and file path, then return
// a ReadSeekCloser and error.
func OpenFile(_ context.Context, path string) (ReadSeekCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Hash is a function to take a set of bytes and return a hash using murmur3.Sum64.
func Hash(data []byte) int {
	return int(murmur3.Sum64(data))
}

// Hash32 is a function to take a set of bytes and return a hash using murmur3.Sum32.
func Hash32(data []byte) int {
	return int(murmur3.Sum32(data))
}

// ReadFile implementation for open-source.
func ReadFile(_ context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Glob implementation for open-source.
func Glob(_ context.Context, pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

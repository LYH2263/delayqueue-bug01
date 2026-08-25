package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (j *Journal) Rotate(newPath string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, rec := range j.queue {
		b, err := json.Marshal(rec)
		if err != nil {
			return err
		}
		if _, err := j.w.Write(b); err != nil {
			return err
		}
		if err := j.w.WriteByte('\n'); err != nil {
			return err
		}
	}
	j.queue = nil
	if err := j.w.Flush(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("open rotated journal: %w", err)
	}
	old := j.f
	if err := old.Close(); err != nil {
		_ = f.Close()
		return err
	}
	j.path = newPath
	j.f = f
	j.w = bufio.NewWriter(f)
	return nil
}

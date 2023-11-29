package pace

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	r := &recorder{}
	w := NewWriter(r)
	w.SetRateLimit(1250, 100)
	random := rand.Reader
	size := 5000
	data := make([]byte, size)
	n, err := random.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, size, n)
	start := time.Now()
	n, err = w.Write(data)
	duration := time.Since(start)
	rate := float64(size) / float64(duration.Seconds()) * 8 / 1000.0
	fmt.Printf("size: %v, duration: %v, rate: %v kbit/s\n", size, duration, rate)
	assert.NoError(t, err)
	assert.Equal(t, size, n)
	assert.Equal(t, r.data, data)
}

type recorder struct {
	data      []byte
	callCount int
}

func (r *recorder) Write(b []byte) (int, error) {
	r.callCount++
	r.data = append(r.data, b...)
	return len(b), nil
}

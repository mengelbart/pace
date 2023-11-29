package pace

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type Writer struct {
	w   io.Writer
	l   *rate.Limiter
	ctx context.Context
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:   w,
		ctx: context.Background(),
	}
}

func (w *Writer) SetRateLimit(bytesPerSec float64, burst int) {
	w.l = rate.NewLimiter(rate.Limit(bytesPerSec), burst)
}

func (w *Writer) Write(b []byte) (int, error) {
	if w.l == nil {
		return w.w.Write(b)
	}
	burst := w.l.Burst()
	n := 0
	for i := 0; i < len(b); i += burst {
		end := i + burst
		if end > len(b) {
			end = len(b)
		}
		tokens := end - i
		if err := w.l.WaitN(w.ctx, tokens); err != nil {
			return len(b), err
		}
		x, err := w.w.Write(b[i:end])
		if err != nil {
			return n, err
		}
		n += x
	}
	return n, nil
}

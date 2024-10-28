package daily

import "os"

type Option func(h *DailyHandler)

func WithMaxAge(maxAge int) Option {
	return func(h *DailyHandler) {
		h.maxAge = maxAge
	}
}

func WithCompress(compress bool) Option {
	return func(h *DailyHandler) {
		h.compress = compress
	}
}

func WithFileMode(mode os.FileMode) Option {
	return func(h *DailyHandler) {
		h.mode = mode
	}
}

package model

import (
	"time"
)

type (
	Activity struct {
		Title       string
		Description string
		Link        string
		UpdatedAt   time.Time
	}

	ServiceActivity interface {
		CategoryCandidates() []string
		Activity() *Activity
	}
)

package model

import (
	"time"
)

type (
	Activity struct {
		Path        []string
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

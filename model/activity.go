package model

type (
	Activity struct {
		Title       string
		Description string
		Link        string
		Meta        []string
	}

	ServiceActivity interface {
		CategoryCandidates() []string
		Activity() *Activity
	}
)

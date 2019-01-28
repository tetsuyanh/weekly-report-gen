package reporter

import (
	"fmt"
	"io"
)

const (
	markdownDateFormat = "01/02/2006"
)

type (
	Markdown struct{}

	MarkdownConf struct{}
)

func NewMarkdown() Reporter {
	return &Markdown{}
}

// implemented Reporter
func (md *Markdown) Report(catActs CategorizedActivities, w io.Writer) error {
	for cat, acts := range catActs {
		if len(acts) == 0 {
			continue
		}
		if _, e := w.Write([]byte(fmt.Sprintf("- %s\n", cat))); e != nil {
			return e
		}
		for _, act := range acts {
			w.Write([]byte(fmt.Sprintf("  - [%s](%s) %s on %s\n", act.Title, act.Link, act.Description, act.UpdatedAt.Format(markdownDateFormat))))
		}
	}
	return nil
}

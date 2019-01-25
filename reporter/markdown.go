package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

type (
	Markdown struct{}
)

func NewMarkdown() Reporter {
	return &Markdown{}
}

// implemented Reporter
func (md *Markdown) Report(r *model.Report, w io.Writer) error {
	for cat, acts := range r.CategorizedActivity {
		if len(acts) == 0 {
			continue
		}
		if _, e := w.Write([]byte(fmt.Sprintf("- %s\n", cat))); e != nil {
			return e
		}
		for _, act := range acts {
			w.Write([]byte(fmt.Sprintf("  - [%s](%s): %s\n", act.Title, act.Link, strings.Join(act.Meta, ", "))))
		}
	}
	return nil
}

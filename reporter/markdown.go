package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/tetsuyanh/weekly-report-gen/categorizer"
)

const (
	markdownDateFormat = "01/02"
)

type (
	Markdown struct {
		conf *MarkdownConf
	}

	MarkdownConf struct {
		ShowPath bool
	}
)

func NewMarkdown(conf *MarkdownConf) Reporter {
	return &Markdown{
		conf: conf,
	}
}

// implemented Reporter
func (md *Markdown) Report(catActs categorizer.CategorizedActivities, w io.Writer) error {
	for cat, acts := range catActs {
		if len(acts) == 0 {
			continue
		}
		if _, e := w.Write([]byte(fmt.Sprintf("- %s\n", cat))); e != nil {
			return e
		}
		for _, act := range acts {
			path := ""
			if md.conf.ShowPath {
				path = strings.Join(act.Path, "/") + "/"
			}
			w.Write([]byte(fmt.Sprintf("  - [%s%s](%s) %s on %s\n", path, act.Title, act.Link, act.Description, act.UpdatedAt.Format(markdownDateFormat))))
		}
	}
	return nil
}

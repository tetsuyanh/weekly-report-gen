package reporter

import (
	"fmt"
	"io"

	"github.com/tetsuyanh/weekly-report-gen/categorizer"
)

const (
	ReportTypeMarkdown = "markdown"
)

var (
	ReportTypes = []string{
		ReportTypeMarkdown,
	}
)

type (
	Conf struct {
		MarkdownConf MarkdownConf
	}

	Reporter interface {
		Report(catActs categorizer.CategorizedActivities, w io.Writer) error
	}
)

func NewReporter(repType string) (Reporter, error) {
	var repo Reporter
	switch repType {
	case ReportTypeMarkdown:
		repo = NewMarkdown()
	default:
		return nil, fmt.Errorf("unknown report type: %s\n", repType)
	}
	return repo, nil
}

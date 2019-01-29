package reporter

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	ReporterMarkdown = "markdown"
)

var (
	Reporters = []string{
		ReporterMarkdown,
	}
)

type (
	Conf struct {
		Categorizer  CategorizerConf
		MarkdownConf MarkdownConf
	}

	Reporter interface {
		Report(catActs CategorizedActivities, w io.Writer) error
	}
)

func ReportActivities(conf *Conf, repId string, srvActs []model.ServiceActivity, w io.Writer) error {
	cater := NewCategorizer(&conf.Categorizer)

	var rep Reporter
	switch repId {
	case ReporterMarkdown:
		rep = NewMarkdown()
	default:
		return fmt.Errorf("unknown reporter: %s\n", repId)
	}

	if err := rep.Report(cater.Categorize(srvActs), w); err != nil {
		return errors.Wrap(err, "rep.Report")
	}

	return nil
}

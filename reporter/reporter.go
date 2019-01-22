package reporter

import (
	"io"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

type (
	Reporter interface {
		Report(repo *model.Report, w io.Writer) error
	}
)

package service

import (
	"time"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

type (
	Service interface {
		CollectServiceActivity(begin, end *time.Time) ([]model.ServiceActivity, error)
	}
)

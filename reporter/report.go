package reporter

import (
	"strings"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	CategoryOthers = "others"
)

type (
	CategoryMap map[string][]string

	ReportBuilder struct {
		CategoryMap CategoryMap
	}
)

func NewReportBuilder(cm CategoryMap) *ReportBuilder {
	rb := &ReportBuilder{
		CategoryMap: cm,
	}
	return rb
}

func (rb *ReportBuilder) Build(sas []model.ServiceActivity) (*model.Report, error) {
	r := &model.Report{
		CategorizedActivity: make(map[string][]*model.Activity, 0),
	}
	// conf category
	for k, _ := range rb.CategoryMap {
		r.CategorizedActivity[k] = make([]*model.Activity, 0)
	}
	// others category
	r.CategorizedActivity[CategoryOthers] = make([]*model.Activity, 0)

	for _, sa := range sas {
		notFound := true
		for _, c := range sa.CategoryCandidates() {
			if cat, ok := rb.findCategory(c); ok {
				r.CategorizedActivity[cat] = append(r.CategorizedActivity[cat], sa.Activity())
				notFound = false
				break
			}
		}
		if notFound {
			r.CategorizedActivity[CategoryOthers] = append(r.CategorizedActivity[CategoryOthers], sa.Activity())
		}
	}

	return r, nil
}

func (rb *ReportBuilder) findCategory(text string) (string, bool) {
	for cat, vals := range rb.CategoryMap {
		for _, val := range vals {
			if strings.Contains(text, val) {
				return cat, true
			}
		}
	}
	return "", false
}

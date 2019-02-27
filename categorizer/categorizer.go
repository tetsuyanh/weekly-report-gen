package categorizer

import (
	"sort"
	"strings"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	CategoryOthers = "others"
)

type (
	Conf struct {
		CategoryMap CategoryMap
	}

	CategoryMap map[string][]string

	Categorizer struct {
		catMap  CategoryMap
		catActs CategorizedActivities
	}

	CategorizedActivities map[string][]*model.Activity
)

func NewCategorizer(conf *Conf) *Categorizer {
	catActs := make(map[string][]*model.Activity, 0)
	for k, _ := range conf.CategoryMap {
		// specified category key
		catActs[k] = make([]*model.Activity, 0)
	}
	// not categorized key
	catActs[CategoryOthers] = make([]*model.Activity, 0)

	return &Categorizer{
		catMap:  conf.CategoryMap,
		catActs: catActs,
	}
}

func (c *Categorizer) Categorize(srvActs []model.ServiceActivity) CategorizedActivities {
	for _, srvAct := range srvActs {
		act := srvAct.Activity()
		notFound := true
		for _, candi := range srvAct.CategoryCandidates() {
			if categ, ok := c.findCategory(candi); ok {
				c.catActs[categ] = append(c.catActs[categ], act)
				notFound = false
				break
			}
		}
		if notFound {
			c.catActs[CategoryOthers] = append(c.catActs[CategoryOthers], act)
		}
	}

	for _, acts := range c.catActs {
		sort.Slice(acts, func(i int, j int) bool {
			return acts[i].UpdatedAt.Before(acts[j].UpdatedAt)
		})
	}

	return c.catActs
}

func (c *Categorizer) findCategory(text string) (string, bool) {
	for cat, vals := range c.catMap {
		for _, val := range vals {
			if strings.Contains(text, val) {
				return cat, true
			}
		}
	}
	return "", false
}

package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tetsuyanh/esa-go/esa"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	esaPageMax    = 30
	esaDateFormat = "2006-01-02"
)

// esa.io, https://esa.io

type (
	Esa struct {
		cli  esa.Client
		user string
	}

	EsaConf struct {
		Enable   bool
		APIToken string
		Team     string
		User     string
	}

	EsaActivity struct {
		Cans []string
		Act  *model.Activity
	}
)

func NewEsa(c *EsaConf) Service {
	apiToken := c.APIToken
	team := c.Team
	user := c.User
	if !c.Enable || apiToken == "" || team == "" || user == "" {
		// not using
		return nil
	}

	return &Esa{
		cli:  esa.NewClient(c.APIToken, c.Team),
		user: c.User,
	}
}

// implemented Service
func (e *Esa) CollectServiceActivity(begin, end *time.Time) ([]model.ServiceActivity, error) {

	param := esa.ListPostsParam{
		Q:       fmt.Sprintf("user:%s updated:>=%s updated:<%s", e.user, begin.Format(esaDateFormat), end.Format(esaDateFormat)),
		Include: []esa.ListPostsParamInclude{},
		Sort:    esa.ListPostsParamSortUpdated,
		Order:   esa.ASC,
	}
	// TODO: recursively retrieve page
	resp, err := e.cli.ListPosts(context.Background(), param, 0, esaPageMax)
	if err != nil {
		return nil, err
	}

	acts := make([]model.ServiceActivity, 0)
	for _, p := range resp.Posts {
		a := &EsaActivity{
			Cans: categoryCandidate(&p),
			Act: &model.Activity{
				Title:       p.FullName,
				Description: "",
				Link:        p.URL,
				Meta:        meta(&p),
			},
		}
		acts = append(acts, a)
	}
	return acts, nil
}

func categoryCandidate(p *esa.Post) []string {
	// 1st candidates are tree, order by splited
	cans := strings.Split(p.Category, "/")
	// 2nd candidates are tags, order by response array
	cans = append(cans, p.Tags...)
	return cans
}

func meta(p *esa.Post) []string {
	var action string
	if p.CreatedAt == p.UpdatedAt {
		action = "Created"
	} else {
		action = "Updated"
	}
	return []string{fmt.Sprintf("esa post %s at %s", action, p.UpdatedAt.Format(esaDateFormat))}
}

// implemented model.ServiceActivity
func (ea *EsaActivity) CategoryCandidates() []string {
	return ea.Cans
}

// implemented model.ServiceActivity
func (ea *EsaActivity) Activity() *model.Activity {
	return ea.Act
}

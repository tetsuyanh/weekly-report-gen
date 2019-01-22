package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tetsuyanh/esa-go/esa"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	pageMax    = 30
	dateFormat = "2006-01-02"
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
		Q:       fmt.Sprintf("user:%s updated:>=%s updated:<%s", e.user, begin.Format(dateFormat), end.Format(dateFormat)),
		Include: []esa.ListPostsParamInclude{},
		Sort:    esa.ListPostsParamSortUpdated,
		Order:   esa.ASC,
	}
	// TODO: recursively retrieve page
	resp, err := e.cli.ListPosts(context.Background(), param, 0, pageMax)
	if err != nil {
		return nil, err
	}

	acts := make([]model.ServiceActivity, 0)
	for _, p := range resp.Posts {
		log.Infof("full: %s\n", p.FullName)
		a := &EsaActivity{
			Cans: categoryCandidate(&p),
			Act: &model.Activity{
				Title:       p.Name,
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
	cans := strings.Split(p.Category, "/")
	cans = append(cans, p.Tags...)
	return cans
}

func meta(p *esa.Post) []string {
	// tree layer
	return []string{fmt.Sprintf("esa %s", p.Message)}
}

// implemented model.ServiceActivity
func (ea *EsaActivity) CategoryCandidates() []string {
	return ea.Cans
}

// implemented model.ServiceActivity
func (ea *EsaActivity) Activity() *model.Activity {
	return ea.Act
}

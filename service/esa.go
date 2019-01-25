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
		conf *EsaConf
	}

	EsaConf struct {
		Enable   bool
		APIToken string
		Team     string
		User     string
	}

	EsaActivity struct {
		Post *esa.Post
	}
)

func NewEsa(conf *EsaConf) Service {
	if !conf.Enable {
		return nil
	}

	return &Esa{
		cli:  esa.NewClient(conf.APIToken, conf.Team),
		conf: conf,
	}
}

// implemented Service
func (e *Esa) CollectServiceActivity(begin, end *time.Time) ([]model.ServiceActivity, error) {

	param := esa.ListPostsParam{
		Q:       fmt.Sprintf("user:%s updated:>=%s updated:<%s", e.conf.User, begin.Format(esaDateFormat), end.Format(esaDateFormat)),
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
		post := p
		acts = append(acts, &EsaActivity{Post: &post})
	}
	return acts, nil
}

// implemented model.ServiceActivity
func (ea *EsaActivity) CategoryCandidates() []string {
	cans := []string{}
	cans = append(cans, strings.Split(ea.Post.Category, "/")...)
	// order dependeds on array
	cans = append(cans, ea.Post.Tags...)
	return cans
}

// implemented model.ServiceActivity
func (ea *EsaActivity) Activity() *model.Activity {
	a := &model.Activity{}

	a.Title = ea.Post.FullName
	a.Link = ea.Post.URL

	var action string
	if ea.Post.CreatedAt == ea.Post.UpdatedAt {
		action = "Created"
	} else {
		action = "Updated"
	}
	a.Meta = []string{fmt.Sprintf("esa post %s at %s", action, ea.Post.UpdatedAt.Format(esaDateFormat))}

	return a
}

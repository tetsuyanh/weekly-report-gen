package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tetsuyanh/go-github/github"
	"golang.org/x/oauth2"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	githubTaskURLHeader = "https://app.github.com/0"
	githubPageMax       = 30
	githubDateFormat    = "2006-01-02"
)

type (
	Github struct {
		cli  *github.Client
		conf *GithubConf
	}

	GithubConf struct {
		Enable       bool
		APIToken     string
		Organization string
		Assignee     string
	}

	GithubToken struct {
		APIToken string
	}

	GithubActivity struct {
		Cans []string
		Act  *model.Activity
	}
)

func NewGithub(conf *GithubConf) Service {
	if !conf.Enable {
		// not using
		return nil
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.APIToken},
	)
	return &Github{
		cli:  github.NewClient(oauth2.NewClient(ctx, ts)),
		conf: conf,
	}
}

// implemented Service
func (g *Github) CollectServiceActivity(begin, end *time.Time) ([]model.ServiceActivity, error) {

	opt := &github.IssueListOptions{
		Filter: "assigned",
		State:  "closed",
		Since:  *begin,
	}
	is, _, e := g.cli.Issues.ListByOrg(context.Background(), g.conf.Organization, opt)
	if e != nil {
		return nil, e
	}

	acts := make([]model.ServiceActivity, 0)
	for _, i := range is {
		// skip out of range
		if i.ClosedAt.After(*end) {
			continue
		}
		// log.Infof("%s, %v", i.GetTitle(), i)
		act := &GithubActivity{
			Cans: makeCategoryCandidateByGithub(i),
			Act:  makeActivityByGithub(i),
		}
		acts = append(acts, act)
	}
	return acts, nil
}

func makeCategoryCandidateByGithub(i *github.Issue) []string {
	cans := []string{}
	// 1st candidate is repo name
	if i.Repository != nil {
		cans = append(cans, *i.Repository.Name)
	}
	// 2nd candidate is milestone name
	if i.Milestone != nil {
		cans = append(cans, *i.Milestone.Title)
	}
	return cans
}

func makeActivityByGithub(i *github.Issue) *model.Activity {
	// log.Infof("issue: %v\n", i)

	titles := []string{}
	if i.Repository != nil {
		titles = append(titles, *i.Repository.Name)
	}
	if i.Milestone != nil {
		titles = append(titles, *i.Milestone.Title)
	}
	titles = append(titles, *i.Title)

	link := i.URL

	labels := []string{}
	for _, l := range i.Labels {
		labels = append(labels, *l.Name)
	}
	meta := []string{fmt.Sprintf("github issue#%d(%s) Closed at %s", *i.Number, strings.Join(labels, ","), i.ClosedAt.Format(githubDateFormat))}

	return &model.Activity{
		Title:       strings.Join(titles, "/"),
		Description: "",
		Link:        *link,
		Meta:        meta,
	}
}

// implemented model.ServiceActivity
func (ga *GithubActivity) CategoryCandidates() []string {
	return ga.Cans
}

// implemented model.ServiceActivity
func (ga *GithubActivity) Activity() *model.Activity {
	return ga.Act
}

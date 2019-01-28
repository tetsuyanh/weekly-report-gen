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

	GithubActivity struct {
		Issue *github.Issue
	}
)

func NewGithub(conf *GithubConf) Service {
	if !conf.Enable {
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
		issue := i
		// skip out of range
		if issue.ClosedAt.After(*end) {
			continue
		}
		// skip other's
		if *issue.Assignee.Login != g.conf.Assignee {
			continue
		}
		acts = append(acts, &GithubActivity{Issue: issue})
	}
	return acts, nil
}

// implemented model.ServiceActivity
func (ga *GithubActivity) CategoryCandidates() []string {
	cans := []string{}
	if ga.Issue.Repository != nil {
		cans = append(cans, *ga.Issue.Repository.Name)
	}
	if ga.Issue.Milestone != nil {
		cans = append(cans, *ga.Issue.Milestone.Title)
	}
	return cans
}

// implemented model.ServiceActivity
func (ga *GithubActivity) Activity() *model.Activity {
	a := &model.Activity{}

	titles := []string{}
	if ga.Issue.Repository != nil {
		titles = append(titles, *ga.Issue.Repository.Name)
	}
	if ga.Issue.Milestone != nil {
		titles = append(titles, *ga.Issue.Milestone.Title)
	}
	titles = append(titles, *ga.Issue.Title)
	a.Title = strings.Join(titles, "/")

	a.Link = *ga.Issue.URL

	labels := []string{}
	for _, l := range ga.Issue.Labels {
		labels = append(labels, *l.Name)
	}
	a.Description = fmt.Sprintf("github issue#%d(%s) Closed", *ga.Issue.Number, strings.Join(labels, ","))
	a.UpdatedAt = *ga.Issue.ClosedAt

	return a
}

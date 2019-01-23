package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tetsuyanh/go-asana/asana"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

const (
	asanaTaskURLHeader = "https://app.asana.com/0"
	asanaPageMax       = 30
	asanaDateFormat    = "2006-01-02"
)

type (
	Asana struct {
		cli  *asana.Client
		conf *AsanaConf
	}

	AsanaConf struct {
		Enable    bool
		APIToken  string
		Workspace string
		Assignee  string
	}

	AsanaToken struct {
		APIToken string
	}

	AsanaActivity struct {
		Cans []string
		Act  *model.Activity
	}

	AuthClient struct {
		Token string
	}
)

func (at *AuthClient) Do(req *http.Request) (*http.Response, error) {
	cli := http.DefaultClient
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", at.Token))
	return cli.Do(req)
}

func NewAsana(conf *AsanaConf) Service {
	if !conf.Enable {
		// not using
		return nil
	}

	return &Asana{
		cli:  asana.NewClient(&AuthClient{Token: conf.APIToken}),
		conf: conf,
	}
}

// implemented Service
func (a *Asana) CollectServiceActivity(begin, end *time.Time) ([]model.ServiceActivity, error) {
	opt := &asana.Filter{
		CompletedSince: begin.Format(asanaDateFormat),
		ModifiedSince:  begin.Format(asanaDateFormat),
		Workspace:      a.conf.Workspace,
		Assignee:       a.conf.Assignee,
		OptExpand:      []string{"completed", "completed_at", "modified_at", "name", "assignee", "projects", "memberships"},
	}
	ts, err := a.cli.ListTasks(context.Background(), opt)
	if err != nil {
		return nil, err
	}

	acts := make([]model.ServiceActivity, 0)
	for _, t := range ts {
		// skip out of range
		if t.ModifiedAt.After(*end) || t.CompletedAt.After(*end) {
			continue
		}
		// skip another's assignee
		if t.Assignee.GID != a.conf.Assignee {
			continue
		}
		a := makeActivityByAsana(&t)
		if a == nil {
			continue
		}
		act := &AsanaActivity{
			Cans: makeCategoryCandidateByAsana(&t),
			Act:  a,
		}
		acts = append(acts, act)
	}
	return acts, nil
}

func makeCategoryCandidateByAsana(t *asana.Task) []string {
	cans := []string{}
	// 1st candidates are team name, order dependeds on array.
	for _, p := range t.Projects {
		cans = append(cans, p.Team.Name)
	}
	// 2nd candidates are project name, order dependeds on array.
	for _, p := range t.Projects {
		cans = append(cans, p.Name)
	}
	// 3rd candidates are section name, order dependeds on array.
	// for _, p := range t.M.Projects {
	// 	for _, m := range p.Me
	// 	cans = append(cans, p.Name)
	// }
	return cans
}

func makeActivityByAsana(t *asana.Task) *model.Activity {
	var prj *asana.Project
	var sec *asana.Section
	// Prioritize membership has section
	for _, m := range t.Memberships {
		if m.Section != nil {
			prj = &m.Project
			sec = m.Section
		}
	}
	if prj == nil {
		// pick head project
		for _, m := range t.Memberships {
			prj = &m.Project
		}
	}
	// no project
	if prj == nil {
		return nil
	}

	var team *asana.Team
	for _, p := range t.Projects {
		if p.GID == prj.GID {
			team = &p.Team
		}
	}

	title := ""
	if team != nil {
		title = fmt.Sprintf("%s/", team.Name)
	}
	if prj != nil {
		title = fmt.Sprintf("%s%s/", title, prj.Name)
	}
	if sec != nil {
		title = fmt.Sprintf("%s%s/", title, sec.Name)
	}
	title = fmt.Sprintf("%s%s", title, t.Name)

	var link string
	link = fmt.Sprintf("%s/%s/%s/f", asanaTaskURLHeader, prj.GID, t.GID)

	var meta []string
	if t.Completed {
		meta = append(meta, fmt.Sprintf("asana task Completed at %s", t.CompletedAt.Format(asanaDateFormat)))
	} else {
		meta = append(meta, fmt.Sprintf("asana task Modified at %s", t.ModifiedAt.Format(asanaDateFormat)))
	}

	return &model.Activity{
		Title:       title,
		Description: "",
		Link:        link,
		Meta:        meta,
	}
}

// implemented model.ServiceActivity
func (ea *AsanaActivity) CategoryCandidates() []string {
	return ea.Cans
}

// implemented model.ServiceActivity
func (ea *AsanaActivity) Activity() *model.Activity {
	return ea.Act
}

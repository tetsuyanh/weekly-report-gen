package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tetsuyanh/go-asana/asana"
	"golang.org/x/oauth2"

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

	AsanaActivity struct {
		Task *asana.Task
	}
)

func NewAsana(conf *AsanaConf) Service {
	if !conf.Enable {
		return nil
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.APIToken},
	)
	return &Asana{
		cli:  asana.NewClient(oauth2.NewClient(ctx, ts)),
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
		task := t
		// skip out of range
		if task.ModifiedAt.After(*end) || task.CompletedAt.After(*end) {
			continue
		}
		// skip other's
		if task.Assignee.GID != a.conf.Assignee {
			continue
		}
		acts = append(acts, &AsanaActivity{Task: &task})
	}
	return acts, nil
}

// implemented model.ServiceActivity
func (ea *AsanaActivity) CategoryCandidates() []string {
	cans := []string{}
	// order dependeds on array
	for _, p := range ea.Task.Projects {
		cans = append(cans, p.Team.Name)
	}
	// order dependeds on array
	for _, p := range ea.Task.Projects {
		cans = append(cans, p.Name)
	}
	// order dependeds on array
	for _, m := range ea.Task.Memberships {
		if m.Section != nil {
			cans = append(cans, m.Section.Name)
			break
		}
	}
	return cans
}

// implemented model.ServiceActivity
func (aa *AsanaActivity) Activity() *model.Activity {
	team, prj, sec := getRelations(aa)
	a := &model.Activity{}

	titles := []string{}
	if team != nil {
		titles = append(titles, team.Name)
	}
	if prj != nil {
		titles = append(titles, prj.Name)
	}
	if sec != nil {
		titles = append(titles, sec.Name)
	}
	titles = append(titles, aa.Task.Name)
	a.Title = strings.Join(titles, "/")

	// '0' will be redirected
	var prjGID string
	if prj != nil {
		prjGID = prj.GID
	} else {
		prjGID = "0"
	}
	a.Link = fmt.Sprintf("%s/%s/%s/f", asanaTaskURLHeader, prjGID, aa.Task.GID)

	var meta string
	if aa.Task.Completed {
		meta = fmt.Sprintf("asana task Completed at %s", aa.Task.CompletedAt.Format(asanaDateFormat))
	} else {
		meta = fmt.Sprintf("asana task Modified at %s", aa.Task.ModifiedAt.Format(asanaDateFormat))
	}
	a.Meta = []string{meta}

	return a
}

// no project
func getRelations(aa *AsanaActivity) (*asana.Team, *asana.Project, *asana.Section) {
	var team *asana.Team
	var prj *asana.Project
	var sec *asana.Section

	// prioritize membership has section
	for _, m := range aa.Task.Memberships {
		if m.Section != nil {
			prj = &m.Project
			sec = m.Section
		}
	}

	// find project
	if prj == nil {
		// pick haad project
		for _, m := range aa.Task.Memberships {
			prj = &m.Project
		}
	}

	// find team
	if prj != nil {
		for _, p := range aa.Task.Projects {
			if p.GID == prj.GID {
				team = &p.Team
			}
		}
	}
	return team, prj, sec
}

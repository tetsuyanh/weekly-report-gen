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
		Enable      bool
		APIToken    string
		WorkspaceID string
		UserID      string
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
		Workspace:      a.conf.WorkspaceID,
		Assignee:       a.conf.UserID,
		OptExpand:      []string{"completed", "completed_at", "modified_at", "name", "assignee", "projects", "memberships", "parent"},
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
		if task.Assignee.GID != a.conf.UserID {
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
	if ea.Task.ParentTask != nil {
		cans = append(cans, ea.Task.ParentTask.Name)
	}
	return cans
}

// implemented model.ServiceActivity
func (aa *AsanaActivity) Activity() *model.Activity {
	a := &model.Activity{}

	a.Title = buildTitle(aa.Task)

	// path '0' will be redirected
	a.Link = fmt.Sprintf("%s/0/%s/f", asanaTaskURLHeader, aa.Task.GID)

	if aa.Task.Completed {
		a.Description = "asana task Completed"
		a.UpdatedAt = aa.Task.CompletedAt
	} else {
		a.Description = "asana task Modified"
		a.UpdatedAt = aa.Task.ModifiedAt
	}

	return a
}

func buildTitle(t *asana.Task) string {
	origin := t
	if t.ParentTask != nil {
		origin = t.ParentTask
	}

	team, prj, sec := getRelations(origin)
	layers := []string{}
	if team != nil {
		layers = append(layers, team.Name)
	}
	if prj != nil {
		layers = append(layers, prj.Name)
	}
	if sec != nil {
		layers = append(layers, sec.Name)
	}
	if t.ParentTask != nil {
		layers = append(layers, t.ParentTask.Name)
	}
	layers = append(layers, t.Name)

	return strings.Join(layers, "/")
}

// no project
func getRelations(t *asana.Task) (*asana.Team, *asana.Project, *asana.Section) {
	var team *asana.Team
	var prj *asana.Project
	var sec *asana.Section

	// prioritize membership has section
	for _, m := range t.Memberships {
		if m.Section != nil {
			prj = &m.Project
			sec = m.Section
		}
	}

	// find project
	if prj == nil {
		// pick haad project
		for _, m := range t.Memberships {
			prj = &m.Project
		}
	}

	// find team
	if prj != nil {
		for _, p := range t.Projects {
			if p.GID == prj.GID {
				team = &p.Team
			}
		}
	}
	return team, prj, sec
}

package service

import (
	"time"

	"github.com/pkg/errors"

	"github.com/tetsuyanh/weekly-report-gen/model"
)

type (
	Service interface {
		CollectServiceActivity(begin, end *time.Time) ([]model.ServiceActivity, error)
	}

	Conf struct {
		Esa    EsaConf
		Asana  AsanaConf
		Github GithubConf
	}
)

func CollectActivities(conf *Conf, begin, end *time.Time) ([]model.ServiceActivity, error) {
	srvs := []Service{
		NewEsa(&conf.Esa),
		NewAsana(&conf.Asana),
		NewGithub(&conf.Github),
	}

	acts := make([]model.ServiceActivity, 0)
	for _, srv := range srvs {
		if srv == nil {
			continue
		}
		a, err := srv.CollectServiceActivity(begin, end)
		if err != nil {
			return nil, errors.Wrap(err, "srv.CollectServiceActivity")
		}
		acts = append(acts, a...)
	}

	return acts, nil
}

package main

import (
	"os"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tetsuyanh/weekly-report-gen/model"
	"github.com/tetsuyanh/weekly-report-gen/reporter"
	"github.com/tetsuyanh/weekly-report-gen/service"
)

type (
	Conf struct {
		Out         string
		CategoryMap reporter.CategoryMap
		Esa         service.EsaConf
	}
)

var conf Conf
var rb *reporter.ReportBuilder

func main() {
	log.Info("gen start\n")

	// period 'end' bigging of tomorrow
	end := beginningOfShiftDay(time.Now(), 1)
	// period 'begin' 7 days ago of 'end'
	begin := end.AddDate(0, 0, -7)
	log.Infof("week: %s - %s", begin, end)

	sas, e := collectActivities(&begin, &end)
	if e != nil {
		log.Error(e)
		return
	}

	r, e := rb.Build(sas)
	if e != nil {
		log.Error(e)
		return
	}

	md := reporter.NewMarkdown()
	if e := md.Report(r, os.Stdout); e != nil {
		log.Error(e)
		return
	}

	log.Info("gen finish\n")
}

func collectActivities(begin, end *time.Time) ([]model.ServiceActivity, error) {
	var conf Conf
	// fixed path and filename
	viper.AddConfigPath("./")
	viper.SetConfigName("conf")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	sas := make([]model.ServiceActivity, 0)
	if esa := service.NewEsa(&conf.Esa); esa != nil {
		esaActs, err := esa.CollectServiceActivity(begin, end)
		if err != nil {
			return nil, err
		}
		sas = append(sas, esaActs...)
	}

	rb = reporter.NewReportBuilder(conf.CategoryMap)

	return sas, nil
}

func beginningOfShiftDay(base time.Time, shift int) time.Time {
	d := base.Add(time.Duration(shift) * 24 * time.Hour)
	t := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
	return t
}

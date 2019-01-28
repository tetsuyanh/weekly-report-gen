package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tetsuyanh/weekly-report-gen/reporter"
	"github.com/tetsuyanh/weekly-report-gen/service"
)

var rootCmd = &cobra.Command{
	Use:   "weekly-report-gen",
	Short: "weekly-report-gen is a weekly report generator",
	Long: `a
                b`,
	Run: func(cmd *cobra.Command, args []string) {
		path := ""

		conf, err := LoadConf(&path)
		if err != nil {
			log.Error(err)
			return
		}

		// period 'end' bigging of tomorrow
		end := beginningOfShiftDay(time.Now(), 1)
		// period 'begin' 7 days ago of 'end'
		begin := end.AddDate(0, 0, -7)
		log.Infof("week: %s - %s", begin, end)

		srvActs, e := service.CollectActivities(&conf.Service, &begin, &end)
		if e != nil {
			log.Error(e)
			return
		}

		repId := reporter.ReporterMarkdown
		if e := reporter.ReportActivities(&conf.Reporter, repId, srvActs, os.Stdout); e != nil {
			log.Error(e)
			return
		}

		log.Info("gen finish\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func beginningOfShiftDay(base time.Time, shift int) time.Time {
	d := base.Add(time.Duration(shift) * 24 * time.Hour)
	t := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
	return t
}

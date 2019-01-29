package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tetsuyanh/weekly-report-gen/reporter"
	"github.com/tetsuyanh/weekly-report-gen/service"
)

type (
	Conf struct {
		Version  string
		Root     RootConf
		Reporter reporter.Conf
		Service  service.Conf
	}

	RootConf struct {
		Date     string
		Reporter string
		Out      string
		Writer   io.Writer
	}
)

const (
	DateFormat = "2006-01-02"

	// if break conf compatibility, have to update
	// keep format v{major}.{middle}.{minor}
	ConfVersionCompatibility = "v0.1.0"
)

var (
	// expect Makefile build
	version string

	// expect command line option
	argConf     string
	argDate     string
	argReporter string
	argOut      string
)

var rootCmd = &cobra.Command{
	Use:   "weekly-report-gen",
	Short: "weekly-report-gen is a weekly report generator",
	Long: `a
                b`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("gen started.\n")

		conf, err := loadConf(argConf)
		if err != nil {
			log.Error(err)
			return
		}

		begin, end, err := weekPeriod(conf.Root.Date)
		if err != nil {
			log.Error(err)
		}
		log.Infof("week from %s to %s", begin, end)

		srvActs, err := service.CollectActivities(&conf.Service, begin, end)
		if err != nil {
			log.Error(err)
			return
		}

		if e := reporter.ReportActivities(&conf.Reporter, conf.Root.Reporter, srvActs, conf.Root.Writer); e != nil {
			log.Error(e)
			return
		}

		log.Info("gen finished.\n")
	},
}

func Execute() {
	// default working directory
	rootCmd.Flags().StringVarP(&argConf, "conf", "c", "./conf", "config filepath without extension")

	rootCmd.Flags().StringVarP(&argDate, "date", "d", "", "the end date of week, e.g. '2019-01-28'")
	rootCmd.Flags().StringVarP(&argReporter, "reporter", "r", "", fmt.Sprintf("reporter name [%s], default 'markdown'", strings.Join(reporter.Reporters, ",")))
	rootCmd.Flags().StringVarP(&argOut, "out", "o", "", "output filepath, default os.Stdout")

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func loadConf(path string) (*Conf, error) {
	var conf Conf

	viper.AddConfigPath(filepath.Dir(path))
	viper.SetConfigName(filepath.Base(path))
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	if conf.Version < ConfVersionCompatibility {
		return nil, fmt.Errorf("require conf %s", ConfVersionCompatibility)
	}

	// override if argument is specified
	if argDate != "" {
		conf.Root.Date = argDate
	}
	if argReporter != "" {
		conf.Root.Reporter = argReporter
	}
	if argOut != "" {
		conf.Root.Out = argOut
	}

	// set up, default and Writer
	if conf.Root.Date == "" {
		conf.Root.Date = time.Now().Format(DateFormat)
	}
	if conf.Root.Reporter == "" {
		conf.Root.Reporter = reporter.ReporterMarkdown
	}
	if conf.Root.Out != "" {
		var err error
		if conf.Root.Writer, err = os.Create(conf.Root.Out); err != nil {
			return nil, errors.Wrap(err, "os.Create")
		}
	} else {
		conf.Root.Writer = os.Stdout
	}

	return &conf, nil
}

func weekPeriod(endDate string) (*time.Time, *time.Time, error) {
	endTime, err := time.Parse(DateFormat, endDate)
	if err != nil {
		return nil, nil, errors.Wrap(err, "time.Parse")
	}
	// period 'end' bigging of next day
	end := beginningOfShiftDay(endTime, 1)
	// period 'begin' 7 days ago of 'end'
	begin := end.AddDate(0, 0, -7)

	return &begin, &end, nil
}

func beginningOfShiftDay(base time.Time, shift int) time.Time {
	d := base.Add(time.Duration(shift) * 24 * time.Hour)
	t := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
	return t
}

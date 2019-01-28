package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/tetsuyanh/weekly-report-gen/reporter"
	"github.com/tetsuyanh/weekly-report-gen/service"
)

type (
	Conf struct {
		Out      string
		Reporter reporter.Conf
		Service  service.Conf
	}
)

func LoadConf(path *string) (*Conf, error) {
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

	return &conf, nil
}

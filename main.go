package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/marpaia/graphite-golang"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugins-go-library/sensu"
)

type HandlerConfig struct {
	sensu.PluginConfig
	Prefix   string
	Port     uint64
	Host     string
	NoPrefix bool
}

const (
	prefix      = "prefix"
	port        = "port"
	host        = "host"
	noPrefix    = "no-prefix"
	defaultPort = 2003
)

var (
	stdin *os.File

	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-graphite-handler",
			Short:    "The Sensu Go Graphite for sending metrics to Carbon/Graphite",
			Keyspace: "sensu.io/plugins/graphite/config",
		},
	}

	graphiteConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      prefix,
			Argument:  prefix,
			Shorthand: "P",
			Default:   "sensu",
			Usage:     "The prefix to use in graphite for these metrics",
			Value:     &config.Prefix,
		},
		{
			Path:      port,
			Argument:  port,
			Shorthand: "p",
			Default:   uint64(defaultPort),
			Usage:     "The port number to which to connect on the graphite server",
			Value:     &config.Port,
		},
		{
			Path:      host,
			Argument:  host,
			Shorthand: "H",
			Default:   "localhost",
			Usage:     "The hostname or address of the graphite server",
			Value:     &config.Host,
		},
		{
			Path:      noPrefix,
			Argument:  noPrefix,
			Shorthand: "n",
			Default:   false,
			Usage:     "Do not include *any* prefixes, use the bare metrics.point.name",
			Value:     &config.NoPrefix,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, graphiteConfigOptions, checkArgs, sendMetrics)
	goHandler.Execute()
}

func checkArgs(_ *corev2.Event) error {

	if config.NoPrefix {
		config.Prefix = ""
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	return nil
}

func sendMetrics(event *corev2.Event) error {

	var (
		metrics        []graphite.Metric
		tmp_point_name string
	)

	Graphite, err := graphite.NewGraphiteWithMetricPrefix(config.Host, int(config.Port), config.Prefix)
	if err != nil {
		return err
	}

	for _, point := range event.Metrics.Points {
		if config.NoPrefix {
			tmpvalue := fmt.Sprintf("%f", point.Value)
			metrics = append(metrics, graphite.NewMetric(point.Name, tmpvalue, point.Timestamp))
		} else {
			// Deal with special cases, such as disk checks that return file system paths as the name
			// Graphite places these on disk using the name, so using any slashes would cause confusion and lost metrics
			if point.Name == "/" {
				tmp_point_name = "root"
			} else {
				tmp_point_name = strings.Replace(point.Name, "/", "_", -1)
			}
			tmpname := fmt.Sprintf("%s.%s_%s", strings.Replace(event.Entity.Name, ".", "_", 1), event.Check.Name, tmp_point_name)
			tmpvalue := fmt.Sprintf("%f", point.Value)
			metrics = append(metrics, graphite.NewMetric(tmpname, tmpvalue, point.Timestamp))
		}
	}

	if err = Graphite.SendMetrics(metrics); err != nil {
		return err
	}

	return Graphite.Disconnect()
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/marpaia/graphite-golang"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

type HandlerConfig struct {
	sensu.PluginConfig
	Prefix           string
	PrefixSource     bool
	AnnotationPrefix string
	Port             uint64
	Host             string
	NoPrefix         bool
}

const (
	// flags
	prefix           = "prefix"
	prefixSource     = "prefix-source"
	annotationPrefix = "annotation-prefix"
	port             = "port"
	host             = "host"
	noPrefix         = "no-prefix"

	// defaults
	defaultPrefix           = "sensu"
	defaultPort             = 2003
	defaultHost             = "127.0.0.1"
	defaultAnnotationPrefix = "sensu.annotations"
)

var (
	stdin        *os.File
	metricPrefix string

	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-graphite-handler",
			Short:    "The Sensu Go Graphite for sending metrics to Carbon/Graphite",
			Keyspace: "sensu.io/plugins/graphite/config",
		},
	}

	graphiteConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      prefixSource,
			Argument:  prefixSource,
			Shorthand: "s",
			Default:   false,
			Usage:     "if true, prepends the sensu entity name (source) as a string to each metric in graphite",
			Value:     &config.PrefixSource,
		},
		{
			Path:      prefix,
			Argument:  prefix,
			Shorthand: "P",
			Default:   defaultPrefix,
			Usage:     "the string to be prepended to each metric in graphite",
			Value:     &config.Prefix,
		},
		{
			Path:      annotationPrefix,
			Argument:  annotationPrefix,
			Shorthand: "a",
			Default:   defaultAnnotationPrefix,
			Usage:     "the string to be prepended to each annotation in graphite",
			Value:     &config.AnnotationPrefix,
		},
		{
			Path:      port,
			Argument:  port,
			Shorthand: "p",
			Default:   uint64(defaultPort),
			Usage:     "the port number to which to connect on the graphite server",
			Value:     &config.Port,
		},
		{
			Path:      host,
			Argument:  host,
			Shorthand: "H",
			Default:   defaultHost,
			Usage:     "the hostname or address of the graphite server",
			Value:     &config.Host,
		},
		{
			Path:      noPrefix,
			Argument:  noPrefix,
			Shorthand: "n",
			Default:   false,
			Usage:     "unsets the default prefix value, use the bare metrics.point.name",
			Value:     &config.NoPrefix,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, graphiteConfigOptions, checkArgs, sendMetrics)
	goHandler.Execute()
}

func checkArgs(event *corev2.Event) error {
	if !event.HasMetrics() {
		return fmt.Errorf("event does not contain metrics")
	}

	prefixSource := strings.Replace(event.Entity.Name, ".", "_", 1)
	if config.Prefix != "" && !config.NoPrefix {
		// --prefix is set, --no-prefix is not set
		metricPrefix = config.Prefix
		if config.PrefixSource {
			// --prefix and --prefix-source are both set
			metricPrefix = fmt.Sprintf("%s.%s", metricPrefix, prefixSource)
		}
	} else if config.PrefixSource {
		// --prefix-source and --no-prefix are both set
		metricPrefix = prefixSource
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	return nil
}

func sendMetrics(event *corev2.Event) error {
	var (
		metrics      []graphite.Metric
		tmpPointName string
		tmpname      string
		check        string
	)

	entity := strings.Replace(event.Entity.Name, ".", "_", 1)
	g, err := graphite.NewGraphite(config.Host, int(config.Port))
	if err != nil {
		return err
	}

	for _, point := range event.Metrics.Points {
		if metricPrefix == "" {
			metrics = append(metrics, graphite.NewMetric(point.Name, fmt.Sprintf("%f", point.Value), point.Timestamp))
		} else {
			// Deal with special cases, such as disk checks that return file system paths as the name
			// Graphite places these on disk using the name, so using any slashes would cause confusion and lost metrics
			if point.Name == "/" {
				tmpPointName = "root"
			} else {
				tmpPointName = strings.Replace(point.Name, "/", "_", -1)
			}

			if event.HasCheck() {
				check = event.Check.Name
			} else {
				check = "statsd"
			}
			tmpname = fmt.Sprintf("%s.%s.%s_%s", metricPrefix, entity, check, tmpPointName)
			metrics = append(metrics, graphite.NewMetric(tmpname, fmt.Sprintf("%f", point.Value), point.Timestamp))
		}
	}

	if event.HasCheck() {
		annotationPath := fmt.Sprintf("%s.events.%s.%s", config.AnnotationPrefix, entity, event.Check.Name)
		metrics = append(metrics, graphite.NewMetric(fmt.Sprintf("%s.action.%s", annotationPath, event.Check.State), "1", event.Timestamp))
		metrics = append(metrics, graphite.NewMetric(fmt.Sprintf("%s.status", annotationPath), fmt.Sprintf("%d", event.Check.Status), event.Timestamp))
	}

	if err = g.SendMetrics(metrics); err != nil {
		return err
	}

	return g.Disconnect()
}

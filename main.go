package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/marpaia/graphite-golang"
	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

var (
	prefix string
	port   int
	host   string
	stdin  *os.File
)

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "sensu-graphite-handler",
		Short: "a graphite handler built for use with sensu",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&host, "host", "H", "localhost", "the hostname or address of the graphite server")
	cmd.Flags().IntVarP(&port, "port", "p", 2003, "the port number to which to connect on the graphite server")
	cmd.Flags().StringVarP(&prefix, "prefix", "P", "sensu", "the prefix to use in graphite for these metrics")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stdin data: %s", err)
	}

	if err = event.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %s", err)
	}

	if !event.HasMetrics() {
		return fmt.Errorf("event does not contain metrics")
	}

	return sendMetrics(event)
}

func sendMetrics(event *types.Event) error {

	var (
		metrics        []graphite.Metric
		tmp_point_name string
	)

	Graphite, err := graphite.NewGraphiteWithMetricPrefix(host, port, prefix)
	if err != nil {
		return err
	}

	for _, point := range event.Metrics.Points {
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

	if err = Graphite.SendMetrics(metrics); err != nil {
		return err
	}

	return Graphite.Disconnect()
}

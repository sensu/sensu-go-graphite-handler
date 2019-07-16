# Sensu Graphite Handler

The Sensu Graphite Handler is a [Sensu Event Handler][3] that sends metrics to
the time series database [Graphite][2]. [Sensu][1] can collect metrics using
check output metric extraction or the StatsD listener. Those collected metrics
pass through the event pipeline, allowing Sensu to deliver the metrics to the
configured metric event handlers. This Graphite handler will allow you to
store, instrument, and visualize the metric data from Sensu.

## Installation

Download the latest version of the sensu-graphite-handler from [releases][4],
or create an executable script from this source.

From the local path of the sensu-graphite-handler repository:
```
go build -o /usr/local/bin/sensu-graphite-handler main.go
```

## Configuration

Example Sensu Go handler definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "graphite"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-graphite-handler",
        "timeout": 10,
        "filters": [
            "has_metrics"
        ]
    }
}
```

Example Sensu Go check definition:

```json
{
    "api_version": "core/v2",
    "type": "CheckConfig",
    "metadata": {
        "namespace": "default",
        "name": "dummy-app-prometheus"
    },
    "spec": {
        "command": "sensu-prometheus-collector -exporter-url http://localhost:8080/metrics",
        "subscriptions":[
            "dummy"
        ],
        "publish": true,
        "interval": 10,
        "output_metric_format": "influxdb_line",
        "output_metric_handlers": [
            "graphite"
        ]
    }
}
```

That's right, you can collect different types of metrics (ex. Influx,
Graphite, OpenTSDB, Nagios, etc.), Sensu will extract and transform
them, and this handler will populate them into your Graphite.


## Usage Examples

Help:
```
Usage:
  sensu-graphite-handler [flags]

Flags:
  -h, --help            help for sensu-graphite-handler
  -H, --host string     the hostname or address of the graphite server (default "localhost")
  -p, --port int        the port number to which to connect on the graphite server (default 2003)
  -n, --no-prefix       do not include *any* prefixes, use the bare metrics.point.name
  -P, --prefix string   the prefix to use in graphite for these metrics (default "sensu")
```

[1]: https://github.com/sensu/sensu-go
[2]: https://graphiteapp.org
[3]: https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/nixwiz/sensu-graphite-handler/releases

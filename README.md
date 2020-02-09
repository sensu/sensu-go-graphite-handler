[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/sensu-go-graphite-handler)
![Go Test](https://github.com/nixwiz/sensu-go-graphite-handler/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/sensu-go-graphite-handler/workflows/goreleaser/badge.svg)

## Sensu Go Graphite Handler Plugin

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Sensu Go](#sensu-go)
    - [Asset registration](#asset-registration)
    - [Asset definition](#asset-definition)
    - [Check definition](#check-definition)
    - [Handler definition](#handler-definition)
  - [Sensu Core](#sensu-core)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

### Overview

The Sensu Graphite Handler is a [Sensu Event Handler][3] that sends metrics to the time series database [Graphite][2]. [Sensu][1] can collect metrics using check output metric extraction or the StatsD listener. Those collected metrics
pass through the event pipeline, allowing Sensu to deliver the metrics to the configured metric event handlers. This Graphite handler will allow you to store, instrument, and visualize the metric data from Sensu.

## Files

N/A

## Usage examples

### Help

```
Usage:
  sensu-go-graphite-handler [flags]
  sensu-go-graphite-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -a, --annotation-prefix string   the string to be prepended to each annotation in graphite (default "sensu.annotations")
  -h, --help                       help for sensu-go-graphite-handler
  -H, --host string                the hostname or address of the graphite server (default "127.0.0.1")
  -n, --no-prefix                  unsets the default prefix value, use the bare metrics.point.name
  -p, --port uint                  the port number to which to connect on the graphite server (default 2003)
  -P, --prefix string              the string to be prepended to each metric in graphite (default "sensu")
  -s, --prefix-source              if true, prepends the sensu entity name (source) as a string to each metric in graphite
```

## Configuration
### Sensu Go
#### Asset registration

Assets are the best way to make use of this plugin. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 or later, you can use the following command to add the asset: 

`sensuctl asset add nixwiz/sensu-go-graphite-handler`

If you're using an earlier version of sensuctl, you can download the asset definition from [this project's Bonsai asset index page][5] or [releases][4] or create an executable script from this source.

From the local path of the sensu-go-graphite-handler repository:
```
go build -o /usr/local/bin/sensu-go-graphite-handler main.go
```

#### Asset definition

```yaml
---
type: Asset
api_version: core/v2
metadata:
  name: sensu-go-graphite-handler
spec:
  url: https://assets.bonsai.sensu.io/793026667633e5cb3e60ba1d063eb5a38ac9cd6b/sensu-go-graphite-handler_0.3.0_linux_amd64.tar.gz
  sha512: af738d13865fdce508fc0c4457ef7473c01639cc92da98590d842eb535db0b51bccdef5c310adf0135b5e3b3677487fe7a1b4370ae3028367bc8117c3fb1824c
```

#### Check definition

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

#### Handler definition

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
        "command": "sensu-go-graphite-handler",
        "timeout": 10,
        "filters": [
            "has_metrics"
        ]
    }
}
```

That's right, you can collect different types of metrics (ex. Influx, Graphite, OpenTSDB, Nagios, etc.), Sensu will extract and transform them, and this handler will populate them into your Graphite.

### Sensu Core

N/A

## Installation from source

### Sensu Go

See the instructions above for [asset registration][7].

### Sensu Core

Install and setup plugins on [Sensu Core][6].

## Additional notes

N/A

## Contributing

N/A

[1]: https://github.com/sensu/sensu-go
[2]: https://graphiteapp.org
[3]: https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/nixwiz/sensu-go-graphite-handler/releases
[5]: https://bonsai.sensu.io/assets/nixwiz/sensu-go-graphite-handler
[6]: https://docs.sensu.io/sensu-core/latest/installation/installing-plugins/
[7]: #asset-registration

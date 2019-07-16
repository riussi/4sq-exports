# 4sq-exports

Simple CLI-tool to export your Foursquare checkins into a KML-file.



## Usage

1. First you need to run the authorise command to get an access token by logging into your Foursquare account.
2. Then you use the access token with the checkins command to get your checkins list from the API.

```bash
$ 4sq-exports
Export your data from Foursquare API

Usage:
  4sq-exports [command]

Available Commands:
  authorise   authorise the app
  checkins    get your foursquare checkins as KML
  help        Help about any command
  version     shows the application version

Flags:
      --config string   config file (default is $HOME/.4sq-exports.yaml)
  -h, --help            help for 4sq-exports

Use "4sq-exports [command] --help" for more information about a command.
```
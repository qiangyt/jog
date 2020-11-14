# jog [![Build Status](https://travis-ci.org/qiangyt/jog.svg?branch=master)](https://travis-ci.org/qiangyt/jog)
Command line tool to on-the-fly convert and view structured(JSON) log as regular flat line format. Jog supports follow mode like 'tail -f', as well as filtering by log level and time range.

## Background

Structured log, AKA. JSON line log, is great for log collectors but hard to read by developers themselves during local development. Jog helps to on-the-fly convert those structured JSON log to traditional flat line log. It then decreases the need to have environment-specific output formats - for ex. we don't need any more to configure JSON log for test / production but flat line log for local development.

Extra feature includes filtering by log level, by time ranage, helpful for daily local development as well.

## Features

   Feature request is welcomed, for ex. new JSON log format. Submit issue for that please.

   - [x] Support various of formats out-of-box and without customization.
         Verified includes (submit issue for new one):
      - [x] Logstash
      - [x] GOLANG Uber zap (https://github.com/uber-go/zap)
      - [x] Node.js Bunyan (https://github.com/trentm/node-bunyan)
      - [ ] Node.js Winston (https://github.com/winstonjs/winston)
      - [ ] Logrus
      - [ ] AWS CloudWatch Logs

   - [x] Follow mode like `tail -f`, with optional beginning from latest specified lines like `tail -n`.
         (see example #1 and #2)

   - [x] Read from stdin (stream) or local file

   - [ ] Straightforward filtering:
      - [x] by log level (see example #6)
      - [x] by absolute time range (see example #7)
      - [x] by relative time range (see example #8)
      - [ ] show surrounding logs

   - [x] output the raw JSON but then able to apply filters (see example #9)

   - [x] Support JSON log mixed with non-JSON text, includes:
      - [x] Mixed with regular flat log lines, for ex., springboot banner, and PM2 banner
      - [x] Extra non-JSON prefix, followed by JSON log, for ex., docker-compose multi-services log

   - [x] Supports nested escaped JSON value (escaped by `\"...\"`)

   - [x] Compressed logger name - only first letters of package names are outputed

   - [x] Print line number as line prefix

   - [x] Customizable although usually you no need it.
         Run `jog -t` to export default configuration, or see [./static_files/DefaultConfiguration.yml](./static_files/DefaultConfiguration.yml)
      - [x] Output pattern
      - [x] Hightlight startup line
      - [x] Colorization
      - [x] Print unknown fields as 'others'
      - [x] For fields that not explictly in output pattern, print as 'others'
      - [x] Show/hide fields

   - [x] A GOLANG application, so single across-platform executable binary, support Mac OSX, Windows, Linux.

## Usage:
  Download the executable binary (https://github.com/qiangyt/jog/releases/) to $PATH. For ex., for Mac OSX and Linux,

  ```shell
     curl -L https://github.com/qiangyt/jog/releases/download/v1.0.0-alpha/jog.$(echo `uname -s` | tr A-Z a-z) -o /usr/local/bin/jog
     chmod +x /usr/local/bin/jog
  ```

   * View a local JSON log file: `jog sample.log`

     Or follows begining from latest 20 lines: `jog -n 20 -f sample.log`

   * Follow stdin stream, for ex. docker: `docker logs -f my.container | ./jog -n 20`

   * Check full usage: `jog -h`

      ```
      Usage:
        jog  [option...]  <your JSON log file path>
           or
        <stdin stream>  |  jog  [option...]

      Examples:
	     1) follow with last 10 lines:         jog -f app-20200701-1.log
	     2) follow with specified lines:       jog -n 100 -f app-20200701-1.log
	     3) with specified config file:        jog -c another.jog.yml app-20200701-1.log
	     4) view docker-compose log:           docker-compose logs | jog
	     5) print the default template:        jog -t
	     6) only shows WARN & ERROR level:     jog -l warn -l error app-20200701-1.log
	     7) shows with timestamp range:        jog --after 2020-7-1 --before 2020-7-3 app-20200701-1.log
	     8) natural timestamp range:           jog --after "1 week" --before "2 days" app-20200701-1.log
	     9) output raw JSON and apply time range filter:      jog --after "1 week" --before "2 days" app-20200701-1.log --json
	     10) disable colorization:             jog -cs colorization=false app-20200701-1.log
	     11) view apache log, non-JSON log     jog -g COMMONAPACHELOG example_logs/grok_apache.log")

      Options:
        -a,  --after <timestamp>                                    'after' time filter. Auto-detect the timestamp format; can be natural datetime
        -b,  --before <timestamp>                                   'before' time filter. Auto-detect the timestamp format; can be natural datetime
        -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml
        -cs, --config-set <config item path>=<config item value>    Set value to specified config item
        -cg, --config-get <config item path>                        Get value to specified config item
        -d,  --debug                                                Print more error detail
        -f,  --follow                                               Follow mode - follow log output
        -g,  --grok <grok pattern name>                             For non-json log line. The default patterns are saved in ...
        -j,  --json                                                 Output the raw JSON but then able to apply filters
        -h,  --help                                                 Display this information
        -l,  --level <level value>                                  Filter by log level. For ex. --level warn
        -n,  --lines <number of tail lines>                         Number of tail lines. 10 by default, for follow mode
             --reset-grok-library-dir                               Save default GROK patterns to [%s]
        -t,  --template                                             Print a configuration YAML file template
        -V,  --version                                              Display app version information
     ```

## Build

   *  Install GOLANG version >= 1.13

   *  In current directory, run `./build.sh`

## Status

   Not yet ready for first release, still keep refactoring and fixing and adding new features. I create pre-release even for single bug fix or small feature. I won't test much before version 1.0 is ready.
   Just feel free to use it since it wouldn't affect something.

## TODO

   * version 1.0 TODO
     - unit test coverage: >= 80%
     - manual test suite
     - read from native Docker containers log
     - SSH remote host to read log directly

## License

[MIT](/LICENSE)

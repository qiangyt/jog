# jog [![Build Status](https://travis-ci.org/qiangyt/jog.svg?branch=master)](https://travis-ci.org/qiangyt/jog)
Command line tool to view structured (JSON) log, as regular flat line format


## Background

Structured log, AKA. JSON line log, is great for log collectors but hard to read by developers themselves, usually during local development. Jog helps to on-the-fly convert those structured JSON log to traditional space-separated flat line log, friendly for developers. It then removes the need to maintenain different output format for different environments - for ex. we don't need any more to configure JSON log for test / production but flat line log for local development.

## Features

   Feature request is welcomed, for ex. new JSON log format. Submit issue for that please.

   - [x] Automatically detect various of formats, by customizeable field aliases. For ex.:
      - [x] Logstash
      - [x] GOLANG Uber zap
      - [x] Node.js Bunyan (https://github.com/trentm/node-bunyan)

   - [x] Follow mode like `tail -f`, with optional beginning from latest specified lines like `tail -n`.

   - [x] Read from stdin (stream) or local file

   - [ ] Straightforard filtering:
      - [x] by logger level
      - [ ] by time range, relative or absolute

   - [x] Support JSON log mixed with non-JSON text, includes:
      - [x] Mixed with regular flat log lines, for ex., springboot banner, and PM2 banner
      - [x] Extra non-JSON prefix, followed by JSON log, for ex., docker-compose multi-services log

   - [x] Supports nested escaped JSON value (escaped by `\"...\"`)

   - [x] Compressed logger name - only first letters of package names are outputed

   - [x] Print line number as line prefix

   - [x] Customization. But I think most-likely you no need customization. Anyway, run `jog -t` to see configuration example.
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
     sudo curl -L https://github.com/qiangyt/jog/releases/download/v0.9.17/jog.$(echo `uname -s` | tr A-Z a-z) -o /usr/local/bin/jog
     sudo chmod +x /usr/local/bin/jog
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
        1) follow with last 10 lines:         jog -f app-20200701-1.log")
        2) follow with specified lines:       jog -n 100 -f app-20200701-1.log")
        3) with specified config file:        jog -c another.jog.yml app-20200701-1.log")
        4) view docker-compose log:           docker-compose logs | jog")
        5) print the default template:        jog -t")
        6) only shows WARN & ERROR level:     jog -l warn -l error app-20200701-1.log")
        7) with WARN level foreground color set to RED: jog -cs fields.level.enums.WARN.color=FgRed app-20200701-1.log")
        8) view the WARN level config item:   jog -cg fields.level.enums.WARN")
        9) disable colorization:              jog -cs colorization=false app-20200701-1.log")

      Options:
        -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml \n")
        -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
        -cg, --config-get <config item path>                        Get value to specified config item \n")
        -d,  --debug                                                Print more error detail\n")
        -f,  --follow                                               Follow mode - follow log output\n")
        -h,  --help                                                 Display this information\n")
        -l,  --level <level value>                                  Filter by log level. For ex. --level warn \n")
        -n,  --lines <number of tail lines>                         Number of tail lines. 10 by default, for follow mode\n")
        -t,  --template                                             Print a config YAML file template\n")
        -V,  --version                                              Display app version information\n")
     ```

## Build

   *  Install GOLANG version >= 1.13

   *  In current directory, run `./build.sh`

## Status

   Not yet ready for firt release, still keep refactoring and fixing and adding new features. I create pre-release even for single bug fix or small feature. I won't test much before version 1.0 is ready.
   Just feel free to use it since it wouldn't affect something.

## License

[MIT](/LICENSE)

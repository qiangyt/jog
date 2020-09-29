# jog [![Build Status](https://travis-ci.org/qiangyt/jog.svg?branch=master)](https://travis-ci.org/qiangyt/jog)
Command line tool to view structured (JSON) log, as regular flat line format


## Background

Structured log, AKA. JSON line log, is great for log collectors but hard to read by developers themselves, usually during local development. Jog helps to on-the-fly convert those structured JSON log to traditional space-separated flat line log, friendly for developers. It then removes the effort to maintenain different output format for different environments (for ex. JSON log for test / production, but flat line log for local development).

## Features

   Feature request is welcomed, for ex. new JSON log format. Submit issue for that please.

   - [x] Detect format, automatically

   - [x] Most-likely know unknow format via customizeable dictionary

   - [ ] Built-in supports as many as possible formats:

      - [x] Logstash
      - [x] Uber zap
      - [x] Bunyan (https://github.com/trentm/node-bunyan)
      - [x] Customizable format. Run `jog -t` to see configuration example.

   - [x] Support Mac OSX, Windows, Linux

   - [x] Supports customized fields

   - [x] Supports nested JSON fields

   - [x] Customizable output pattern

   - [x] Customizable output colorization

   - [x] Hightlight startup line

   - [x]  Support JSON log mixed with non-JSON log line (for ex., springboot banner) - just directly print them

   - [ ] Able to directly read many sources:
      - [x] stdin & stream
      - [x] local file
      - [x] docker-compose log
      - [x] docker log
      - [x] OMS-docker (https://github.com/microsoft/OMS-docker)
      - [ ] aggregate multiple log

   - [x]  Friendly to multi-containers log outputted by docker-compose

   - [x]  Compressed logger name - only first letters of package names are outputed

   - [ ]  Filtering by level and field

## Usage:
  Download the executable binary (https://github.com/qiangyt/jog/releases/) to $PATH. For ex., for Mac OSX and Linux,

  ```shell
     sudo curl -L https://github.com/qiangyt/jog/releases/download/v0.9.13/jog.$(echo `uname -s` | tr A-Z a-z) -o /usr/local/bin/jog
     sudo chmod +x /usr/local/bin/jog
  ```

   * View a local JSON log file: `jog sample.log`

     And, with specified number of lines: `jog -n 20 -f sample.log`

   * Follow stdin stream, for ex. docker: `docker logs -f my.container | ./jog`

     Also, with specified number of lines: `docker logs -f my.container | ./jog -n 20`

   * From stdin steam: `tail -f sample.log | ./jog`

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
         6) with WARN level foreground color set to RED: jog -cs fields.level.enums.WARN.color=FgRed app-20200701-1.log
         7) view the WARN level config item:   jog -cg fields.level.enums.WARN
         8) disable colorization:              jog -cs colorization=false app-20200701-1.log

      Options:
         -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml
         -cs, --config-set <config item path>=<config item value>    Set value to specified config item
         -cg, --config-get <config item path>                        Get value to specified config item
         -f,  --follow                                               Follow mode - follow log output
         -n,  --lines <number of tail lines>                         Number of tail lines. 10 by default, for follow mode
         -t,  --template                                             Print a config YAML file template
         -h,  --help                                                 Display this information
         -V,  --version                                              Display app version information
         -d,  --debug                                                Print more error detail
     ```

## Build

   *  Install GOLANG version >= 1.13

   *  In current directory, run `./build.sh`

## License

[MIT](/LICENSE)

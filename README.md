# jog
Convert and view structured (JSON) log


## Background

Structured log, AKA. JSON line log, is great for log collectors but hard to read for developers themselves, usually during local development. This tool helps to on-the-fly convert those structured JSON log to traditional space-separated line log, friendly for developers. It then removes the effort to maintenain different output format for different environments (for ex. JSON log for test / production, but line log for local development).

## Features

   Feature request is welcomed, for ex. new JSON log format. Submit issue for that please.

   - [ ] Supports as many as possible formats:

      - [x] Logstash
      - [ ] Uber zap
      - [ ] Bunyan
      - [x] Actually you could define your own format. Run `jog -t` to see configuration example. Most-likely your JSON log format is already supported, automatically

   - [x] Supports customized fields

   - [x] Customizable output pattern

   - [x] Customizable output colorization

   - [x] Hightlight startup line

   - [x]  Support JSON log mixed with non-JSON log line (for ex., springboot banner) - just directly print them

   - [ ] Able to directly read many sources:
      - [x] stdin & stream
      - [x] local file
      - [ ] remote file (HTTP/HTTPs/FTP/SFTP)
      - [ ] k8s log
      - [ ] docker-compose log
      - [ ] docker log
      - [ ] aggregate multiple log

   - [ ]  Friendly to multi-containers log outputted by docker-compose

   - [x]  Compressed logger name - only first letters of package names are outputed

   - [ ]  Filtering, both command line and embedded Web UI


## Executable binaries:

   https://github.com/qiangyt/jog/releases/

## Usage:
  Copy the downloaded binary to $PATH. For ex.

  ```shell
     cp jog /usr/local/bin/
  ```

   * View a local JSON log file: `jog sample.log`

   * From stdin: `cat sample.log | ./jog`

   * From stdin steam: `tail -f sample.log | ./jog`

   * Check full usage: `jog -h`

## Build

   *  Install GOLANG

   *  In current directory, run `./build.sh`

## License

[MIT](/LICENSE)

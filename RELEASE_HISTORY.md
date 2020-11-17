# targeting v1.0.0 (TODO)
1. User manual - English
2. Unit test coverage: 100%
3. New feature: field path

# v1.0.0-rc-2 (TODO)
1. Unit test coverage: >= 80%
2. Fix: ...
3. New feature: initial arbitrary field filter

# v1.0.0-rc-1 (TODO)
1. Unit test coverage: >= 50%
2. Fix: ...
3. Enhancement: Log json log

# v1.0.0-belta (TODO)
1. User manual - Chineses
2. Enhancement: improve grok support
3. Enhancement: merge multiple configuration
4. Replace dynamic object design with https://github.com/imdario/mergo & https://github.com/mitchellh/mapstructure

# v1.0.0-alpha
1. New feature: initial support flat (non-JSON) log, using logstash GROK to parse
2. New feature: new configuration parameter `print-format`, for ex., able to tabularize log output - by @frudolph77
3. Fix: replace tab in log lines by four spaces - by @frudolph77
4. Remove `before` and `after` configuration parameters. Uses `print-format` instead
5. A bit refactoring

# v0.9.20 (2020-10-27)
1. New feature: output the raw JSON but then able to apply filters (see example #9)

# v0.9.19 (2020-10-24)
1. New feature: specify the time range filter by natural date time, for ex. `1 day`, `2 hour`

# v0.9.18 (2020-10-23)
1. New feature: filtering by absolute time range (option `--before` and `--after`)

# v0.9.17 (2020-10-22)
1. New feature: filtering by logger level (option `-l`)

Comments:
  The filtering feature initially works. So far it only supports filtering by logger level. Time range filtering should be ready soon (v0.9.18?). Then feature freezed before v1.0.0 ready.

# v0.9.16 (2020-10-21)
1. Enhancement: if a field doesnot explicitly appears in output print pattern, I name such field as `implicit field`,
   and now implicit field will be printed in `${others}`.
   This is incompatible behavior - for old versions, implicit fields are hidden. Now begins from this version, to hide
   the implicit field, must set `print` attribute as `false`. There're examples in default configuration template.
2. Refine default template: some fields are printed with different color and style, some fields are hidden if be implicit.
3. Fix: should not print a implicit field if its print attribute is false.
4. Fix: for bunyan logger, it logger field takes 'id' as name, but then not printed.
5. Enhancement: ${others} fields are sorted by alphabet order.
6. Enhancement: detect logback error stacktrace.

# v0.9.15.1
1. Fix: regression by regression by https://github.com/qiangyt/jog/commit/cea3edbb5f6c19079e21688d657a85a5587d4394

# v0.9.15
1. Fix: failed to load default configuration file due to a stupid error that takes path as yaml.
   Thanks @https://github.com/nseba for reporting and @https://github.com/frudolph77 for reproducing.

# v0.9.14
1. Enhancement: support -n option to specifiy amount of tail lines, like `tail -f`
2. Refine README and help messages

# v0.9.13
  skipped due to mis-operation

# v0.9.12
1. Fix: array value should not be ignored
2. Refactor the code for ongoing data type support
3. Initial time data type format

# v0.9.11
1. Enhancement: support -f argument for local file (follow mode like `tail -f`)

# v0.9.10
1. Fix: typo in README and help text
2. Update README.md to remove deprecated feature plan

# v0.9.9
1. Refactor the command line argument handling
2. Enhancement: use github.com/mitchellh/go-homedir to detect user home directory, which works better in Windows
3. Enhancement: Update dependencies
4. Enhancement: Remove a debug code
5. Fix: embedded stacktrace is not printed

# v0.9.8
1. Feature: able to enable/disable colorization (see example #7 running `jog -h`)

# v0.9.7
1. Feature: support to remove escaped double quote as possible as we can, for ex., for log outputted by OMS-docker (https://github.com/microsoft/OMS-docker)

# v0.9.6
1. Feature: field compress feature now supports white-list

# v0.9.5
1. Feature: Finish field prefix compress
2. Fix: bunyan logs with stacktrace should not be ignored
3. Enhancement: add "critical" as alias of ERROR log level
4. Enhancement: Internal log improvements: move from current folder to ${HOME}/.jog/, increase max size to 100M, ...

# v0.9.4
1. Fix: NPD if the field value is nil
2. Verified it works on Windows

# v0.9.3
1. Feature: Add command line option to get/set configuration items

# v0.9.2
1. Enhancement: Re-design the configuration so that we can define customized arbitrarily

# v0.9.1
  N/A

# v0.9.0
  First release

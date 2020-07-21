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

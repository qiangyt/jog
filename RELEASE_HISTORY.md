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

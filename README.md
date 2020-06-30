# j2log
Convert and view structured (JSON) log

UNDER EARLY DEVELOPMENT. DON'T USE FOR NOW.

## Feature request is welcomed, for ex. new JSON log format. Submit issue for that please.

## Planned Features

- [ ] Supported many formats:

   - [x] logstash
   - [ ] zap
   - [ ] bunyan
   - [ ] will add more

- [ ] Customizable colorization (not yet supports windows)

- [ ] able to directly read many sources:
   - [ ] stdin
   - [ ] local file
   - [ ] remote file (HTTP/HTTPs/FTP/SFTP)
   - [ ] k8s log
   - [ ] docker-compose log
   - [ ] docker log
   - [ ] aggregate multiple log

- [ ]  Friendly to multi-containers log outputted by docker-compose

- [ ]  Customizable output format

- [ ]  Compressed logger name - only first letters of package names are outputed

- [x]  Support JSON log mixed with non-JSON log line (for ex., springboot banner) - directly print them

- [ ]  Filtering, both command line and embedded Web UI

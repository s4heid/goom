[![Go Report Card](https://goreportcard.com/badge/s4heid/goom)](https://goreportcard.com/report/s4heid/goom)
[![Actions Status](https://github.com/s4heid/goom/workflows/.github/workflows/goom.yml/badge.svg)](https://github.com/s4heid/goom/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/s4heid/goom/blob/master/LICENSE)

# goom

A simple golang CLI application for opening url's associated with an alias in
the web browser.

## Configuration

Create a configuration file `~/.goom.yml`. The following properties need to be
specified:

* `url` - the target url. Templating syntax can be used to build the target url
  with properties of a room.
* `rooms` - list of rooms.
  - `id` - identifier of a room.
  - `alias` - alias for a room that should be unique.
  - `name` *(optional)* - descriptive name of a room.

Currently supported input source formats are [YAML](https://yaml.org) and
[JSON](https://www.json.org/)

**Example Config**:

```yaml
---
url: https://zoom.us/j/{{.Id}}
rooms:
- id: 0123456789
  name: Team Standup
  alias: daily
- id: 9876543210
  name: John Doe
  alias: jd
```

If you want to open the url that is associated with the room of John Doe, use
goom's `open` command and pass the alias `jd` as a command line argument:

```sh
$ goom open jd
Opening "https://zoom.us/j/9876543210" in the browser...
```

The interpolated url `https://zoom.us/j/9876543210` will be opened in a new
window of your default web browser.

## License

[Apache License, Version 2.0](LICENSE)

# gozoom

Hop into a zoom rooms from the command line.

## Usage

1. Create the binary in the `bin` folder of your $GOPATH:
    ```sh
    go install
    ```
1. Create a configuration file `~/.gozoom/config.yml`.

**Example Config**:

```yaml
room: https://zoom.us/j
people:
- id: 0123456789
  name: Team Standup
  alias: daily
- id: 9876543210
  name: John Doe
  alias: jd
```


Now, you can use the `gozoom` command and the alias of the room as command line argument to
jump into the zoom of the corresponding person:

```sh
gozoom jd
```


## License

[Apache License](LICENSE)


## corectl measure ls

Print a list of all generic measures in the current app

### Synopsis

Print a list of all generic measures in the current app

```
corectl measure ls [flags]
```

### Examples

```
corectl measure ls
```

### Options

```
  -h, --help   help for ls
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl measure](corectl_measure.md)	 - Explore and manage measures


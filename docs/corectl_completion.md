## corectl completion

Generate auto completion scripts

### Synopsis

Generate a shell completion script for the specified shell (bash or zsh). The shell script must be evaluated to provide
interactive completion. This can be done by sourcing it in your ~/.bashrc or ~/.zshrc file. 
Note that bash-completion is required and needs to be installed on your system.

```
corectl completion <shell> [flags]
```

### Examples

```
   Add the following to your ~/.bashrc or ~/.zshrc file

   . <(corectl completion zsh)

   or

   . <(corectl completion bash)
```

### Options

```
  -h, --help   help for completion
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

* [corectl](corectl.md)	 - 


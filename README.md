# gappwrapp
Lightweight process wrapper (for Linux/Mac) to execute commands with user-defined side-effects. *gappwrapp* = **G**o-**App**lication-**Wrapp**er.

Being developed with personal use-cases in mind, trying to improve my golang skills along the way. **Just starting out - no code here yet :-)**

Here's a list of things I'd like gappwrapp to do in the future (most are not implemented yet):

- Capture status and duration of program, export them as prometheus metrics to a file
- Timeout a command after a given time (alternative to the `timeout` command)
- Pipe stderr and stdout to files (shorthands for `> /var/log/myfile.log`)
- Prepend timestamps to stdout/stderr (similar to `ts` or `gomon`)
- Pass along environment variables and signals to target command

This is how I think you'd use gappwrapp:

```sh
# Execute python scrip myscript.py using gappwrapp
$ gappwrapp --config etc/gappwrapp-myscript.toml -- python myscript.py --foo 123 --bar 456

$ cat gappwrapp-myscript.toml
[general]
timeout=120

[stdout]
timestamp=true
timestamp-format="%Y %M %d"
logfile=/var/log/mylogfile.log

[stderr]
timestamp=true
timestamp-format="%Y %M %d"
logfile=/var/log/mylogfile.log

[record]
duration=yes
status=yes

[export]
format=prometheus
metric_prefix="myscript_"
file=/opt/node_exporter/myfile.prom
```

# Development

```sh
# Run during development
go run gapprwapp.go -- ls -la .
# IMPORTANT: When using 'go run', go will print the exit code of the program and not report the correct statuscode back to bash
# See here: https://stackoverflow.com/questions/26893774/how-to-disable-exit-status-1-when-executing-os-exit1

# Build binary
go build -o bin/darwin/gappwrapp gappwrapp.go

# Run binary
bin/darwin/gappwrapp

# Run tests
go test -v
```
# Check NG

This is a simple port of Mathias Kettner's awesome check-mk agent program to Golang.

It tries to be compatible with the output of MKs agent and can even use the same script directory but it doesn't
require the use of xinetd since it is self-contained with it's own TCP server.  It works in the exact same way
but runs on port 5665 by default.

## Building

Is easy.

    go get github.com/penguinpowernz/check-ng
    go build github.com/penguinpowernz/check-ng/cmd/check-ng

## Usage

This will cause the agent to dump the output directly to stdout so you could potentially if you wanted to use xinetd.

    check-ng -dump

This will start a TCP server on port 5665:

    check-ng
    nc localhost 5665       # get the output via TCP

You can also specify the directory to run the scripts from (by default it uses `/usr/lib/check_mk_agent/local/`):

    check-ng -dump -dir /var/lib/check-ng/scripts

You can specify the host and port to listen on:

    check-ng -host 0.0.0.0 -port 6556

You can also specify to use UDP instead:

    check-ng -udp -port 4545

Or even HTTP:

    check-ng -http -port 8080
    curl localhost:8080/ -o raw.txt
    curl localhost:8080/tree -o sections_and_lines.json
    curl localhost:8080/tree/df -o df_lines.json

## Your Own Scripts

You can create your own plugins, in the same was as you would for check-mk-agent.  Simply make an executable script
and put it in the directory you run `check-ng` from (default `/usr/lib/check_mk_agent/local`):

```
#!/bin/bash

# this script puts the database info out for check-mk to serve

echo "<<<databases>>>"  # must include this as it becomes the title
ls -l /data/*.db        # the output for this script
```

Of course it can be written in any language, or can even be a binary as long as it includes the header line: `<<<title>>>` and
you may opt to output JSON or YAML to make it easier for processes that read your output to parse.

## Todo

- [ ] add ability to change port
- [ ] add compatability mode to be fully backwards compatible
- [ ] complete the tree endpoints in the API
- [ ] multiple directory support
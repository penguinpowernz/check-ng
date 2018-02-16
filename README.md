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

## Todo

- [ ] add ability to change port
- [ ] add compatability mode to be fully backwards compatible
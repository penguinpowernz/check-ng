package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	checkmk "github.com/penguinpowernz/check-mk-ng"
)

const (
	CONN_HOST   = "localhost"
	CONN_PORT   = "5665"
	CONN_TYPE   = "tcp"
	DEFAULT_DIR = "/usr/lib/check_mk_agent/local"
)

func main() {
	var dir string
	var dump bool

	flag.StringVar(&dir, "dir", DEFAULT_DIR, "")
	flag.BoolVar(&dump, "dump", false, "")
	flag.Parse()

	log.Println("Using dir:", dir)

	log.SetOutput(os.Stderr)

	if dump {
		log.Println("dumping...")
		if err := Write(dir, os.Stdout); err != nil {
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Send a response back to person contacting us.
		if err := Write(dir, conn); err != nil {
			log.Printf("ERROR: %s", err)
		}
		// Close the connection when you're done with it.
		conn.Close()
	}
}

func Write(dir string, w io.Writer) (err error) {
	err = checkmk.WriteHeader(dir, w)
	if err != nil {
		return
	}

	err = checkmk.WriteDefaultProcOutput(w)
	if err != nil {
		return
	}

	err = checkmk.WriteDefaultCommandsOutput(w)
	if err != nil {
		return
	}

	err = checkmk.WriteScripts(dir, w)
	if err != nil {
		return
	}

	return nil
}

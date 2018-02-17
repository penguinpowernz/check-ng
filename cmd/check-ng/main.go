package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	checkmk "github.com/penguinpowernz/check-ng"
	"github.com/penguinpowernz/check-ng/client"
)

const defaultDir = "/usr/lib/check_mk_agent/local"

func main() {
	var dir, host, port string
	var dump, http, udp bool

	flag.StringVar(&dir, "dir", defaultDir, "")
	flag.StringVar(&host, "host", "localhost", "")
	flag.StringVar(&port, "port", "5665", "")
	flag.BoolVar(&dump, "dump", false, "")
	flag.BoolVar(&http, "http", false, "")
	flag.BoolVar(&udp, "udp", false, "")
	flag.Parse()

	log.SetOutput(os.Stderr)
	log.Println("Using dir:", dir)

	switch {
	case http:
		runHTTP(dir, host, port)
	case dump:
		dumpOut(dir)
	default:
		runRaw(udp, dir, host, port)
	}
}

func dumpOut(dir string) {
	log.Println("dumping...")
	if err := write(dir, os.Stdout); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runHTTP(dir, host, port string) {
	api := gin.Default()

	api.GET("/", func(c *gin.Context) {
		c.Status(200)
		write(dir, c.Writer)
	})

	api.GET("/tree", func(c *gin.Context) {
		cl := client.New()
		err := write(dir, cl)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, cl.Tree())
	})

	api.GET("/tree/:section", func(c *gin.Context) {
		cl := client.New()
		err := write(dir, cl)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		sect, found := cl.Tree()[c.Param("section")]

		if !found {
			c.AbortWithStatus(404)
			return
		}

		c.JSON(200, sect)
	})

	api.Run(host + ":" + port)
}

func runRaw(udp bool, dir, host, port string) {
	// Listen for incoming connections.
	connType := "tcp"
	if udp {
		connType = "udp"
	}

	l, err := net.Listen(connType, host+":"+port)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + host + ":" + port)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Send a response back to person contacting us.
		if err := write(dir, conn); err != nil {
			log.Printf("ERROR: %s", err)
		}
		// Close the connection when you're done with it.
		conn.Close()
	}
}

func write(dir string, w io.Writer) (err error) {
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

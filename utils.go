package check_ng

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func WriteDefaultProcOutput(w io.Writer) error {
	cat("/proc/meminfo").WithTitle("mem").Write(w)
	cat("/proc/uptime").WithTitle("uptime").Write(w)
	cat("/proc/loadavg").WithTitle("loadavg").Write(w)
	cat("/proc/diskstats").WithTitle("diskstat").Write(w)
	cat("/proc/mdstat").WithTitle("md").Write(w)
	return nil
}

func WriteHeader(dir string, w io.Writer) error {
	lines := []string{}
	lines = append(lines, "<<<check-mk>>>")
	lines = append(lines, "Version: 0.0.1-ng")
	lines = append(lines, "AgentOS: linux")
	hn, _ := os.Hostname()
	lines = append(lines, fmt.Sprintf("Hostname: %s", hn))
	lines = append(lines, fmt.Sprintf("LocalDirectory: %s", dir))

	_, err := w.Write([]byte(strings.Join(lines, "\n") + "\n"))
	return err
}

func WriteDefaultCommandsOutput(w io.Writer) error {
	command("/bin/mount").WithTitle("mounts").Write(w)
	command("/bin/df").WithTitle("df").Write(w)
	command("/bin/df", "-i").WithTitle("df-inodes").Write(w)
	return nil
}

func WriteScripts(dir string, w io.Writer) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no scripts found in %s", dir)
	}

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}

		path := filepath.Join(dir, fi.Name())
		log.Printf("running script %s", path)

		if _, err := script(path).Write(w); err != nil {
			log.Printf("Failed to run %s", err)
		}
	}

	return nil
}

type Command struct {
	Cmd       *exec.Cmd
	Title     string
	OmitTitle bool
}

func (c *Command) WithTitle(t string) *Command {
	c.Title = t
	return c
}

func (c *Command) Write(w io.Writer) (bw int, err error) {
	if !c.OmitTitle {
		if c.Title == "" {
			c.Title = c.Cmd.Args[0]
		}

		bw, err = w.Write([]byte(fmt.Sprintf("<<<%s>>>\n", c.Title)))
		if err != nil {
			return
		}
	}

	c.Cmd.Stdout = w
	err = c.Cmd.Run()
	if err != nil {
		return bw, err
	}

	return
}

func command(script string, args ...string) *Command {
	cmd := exec.Command(script, args...)
	cmd.Stderr = os.Stderr
	return &Command{Cmd: cmd}
}

func script(script string, args ...string) *Command {
	cmd := exec.Command(script, args...)
	cmd.Stderr = os.Stderr
	return &Command{Cmd: cmd, OmitTitle: true}
}

type File struct {
	Path  string
	Title string
}

func (f *File) WithTitle(t string) *File {
	f.Title = t
	return f
}

func cat(fn string) *File {
	return &File{Path: fn}
}

func (f *File) Write(w io.Writer) (bw int, err error) {
	if f.Title == "" {
		f.Title = filepath.Base(f.Path)
	}

	bw, err = w.Write([]byte(fmt.Sprintf("<<<%s>>>\n", f.Title)))
	if err != nil {
		return
	}

	fh, err := os.Open(f.Path)
	if err != nil {
		return
	}

	bw64, err := io.Copy(w, fh)
	bw += int(bw64)

	return
}

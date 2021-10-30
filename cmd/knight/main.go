package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"time"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/parser"
	"github.com/polyscone/knight/value"
)

var (
	version = "unset"
	branch  = "unset"
	commit  = "unset"
	built   = "unset"
	tags    = "unset"
	target  = "unset"
	race    = "unset"
)

var (
	expression   = flag.String("e", "", "An expression to evaluate")
	filename     = flag.String("f", "", "A path to a file to run")
	profile      = flag.String("p", "", "The name of a profile to record")
	versionQuery = flag.Bool("version", false, "Display binary version information")
)

func main() {
	flag.Parse()

	if *versionQuery || flag.Arg(0) == "version" {
		if t, err := time.Parse(time.RFC3339, built); err == nil {
			built = t.Local().Format(time.UnixDate)
		}

		var message string

		message += fmt.Sprintln("Version:      ", version)
		message += fmt.Sprintln("Branch:       ", branch)
		message += fmt.Sprintln("Commit:       ", commit)
		message += fmt.Sprintln("Built:        ", built)
		message += fmt.Sprintln("Tags:         ", tags)
		message += fmt.Sprintln("Go version:   ", strings.TrimPrefix(runtime.Version(), "go"))
		message += fmt.Sprintln("OS/Arch:      ", runtime.GOOS+"/"+runtime.GOARCH)
		message += fmt.Sprintln("Target:       ", target)
		message += fmt.Sprintln("Race detector:", race)

		fmt.Print(message)

		return
	}

	if (*expression == "" && *filename == "") || (*expression != "" && *filename != "") {
		flag.Usage()

		os.Exit(1)
	}

	//nolint:gocritic // it's ok for defers to not run in profile code
	switch *profile {
	case "cpu":
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	case "mem":
		f, err := os.Create("mem.pprof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		runtime.MemProfileRate = 4096

		//nolint:errcheck // any errors here will be apparent so we don't check
		defer pprof.Lookup("heap").WriteTo(f, 0)
	case "trace":
		f, err := os.Create("trace.out")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := trace.Start(f); err != nil {
			log.Fatal(err)
		}
		defer trace.Stop()
	default:
		if *profile != "" {
			log.Fatalf("unknown profile type %q", *profile)
		}
	}

	var b []byte
	if *expression != "" {
		b = []byte(*expression)
	} else {
		var err error
		if b, err = os.ReadFile(*filename); err != nil {
			fmt.Println(err)

			os.Exit(1)
		}
	}

	l := lexer.New()
	p := parser.New(l)
	g := value.NewGlobalStore()
	program, err := p.Parse(g, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	if _, err := interpreter.New(g, p).Execute(program); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

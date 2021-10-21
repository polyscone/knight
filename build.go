//go:build ignore

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	buildTags = []string{}
	testTags  = append([]string{"test"}, buildTags...)
)

var (
	debug    = flag.Bool("debug", false, "Enable symbol table/DWARF generation and disable optimisations/inlining")
	race     = flag.Bool("race", false, "Enable data race detection")
	reckless = flag.Bool("reckless", false, "Enable reckless behaviour")
)

func main() {
	flag.Parse()

	fmt.Println("Building:")

	generate()
	test()
	lint()
	build()
}

func generate() {
	fmt.Print("    go generate... ")

	out, err := exec.Command("go", "generate", "./...").CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			fmt.Println(string(out))
		}

		os.Exit(1)
	}

	fmt.Println("ok")
}

func test() {
	fmt.Print("    go test... ")

	out, err := exec.Command("go", "test", "-race", "-tags", strings.Join(testTags, " "), "./...").CombinedOutput()
	if err != nil {
		fmt.Println()

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		os.Exit(1)
	}

	fmt.Println("ok")
}

func lint() {
	fmt.Print("    golangci-lint... ")

	if lint, err := exec.LookPath("golangci-lint"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Println("not found, please install")
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	} else {
		out, err := exec.Command(lint, "run", "--color", "always", "--build-tags", strings.Join(testTags, " ")).CombinedOutput()
		if err != nil {
			fmt.Println()

			if len(out) > 0 {
				fmt.Println(string(out))
			}

			os.Exit(1)
		}

		fmt.Println("ok")
	}
}

func build() {
	fmt.Print("    go build... ")

	if *reckless {
		buildTags = append(buildTags, "reckless")
	}

	cmd := "go"
	args := []string{"build", "-v", "-tags", strings.Join(buildTags, " ")}
	gcflags := []string{}
	ldflags := []string{
		fmt.Sprintf("-X 'main.version=%v'", closestTag()),
		fmt.Sprintf("-X 'main.branch=%v'", commitBranch()),
		fmt.Sprintf("-X 'main.commit=%v'", commitHash()),
		fmt.Sprintf("-X 'main.built=%v'", time.Now().UTC().Format(time.RFC3339)),
	}

	if len(buildTags) == 0 {
		ldflags = append(ldflags, "-X 'main.tags=None'")
	} else {
		ldflags = append(ldflags, fmt.Sprintf("-X 'main.tags=%v'", strings.Join(buildTags, " ")))
	}

	if *debug {
		// -N disables all optimisations
		// -l disables inlining
		// See: go tool compile --help
		gcflags = append(gcflags, "all=-N -l")

		ldflags = append(ldflags, "-X 'main.target=Debug'")
	} else {
		args = append(args, "-trimpath")

		ldflags = append(ldflags, "-X 'main.target=Release'")

		// -s disables the symbol table
		// -w disables DWARF generation
		// See: go tool link --help
		ldflags = append(ldflags, "-s")
		ldflags = append(ldflags, "-w")
	}

	if *race {
		args = append(args, "-race")

		ldflags = append(ldflags, "-X 'main.race=Enabled'")
	} else {
		ldflags = append(ldflags, "-X 'main.race=Disabled'")
	}

	if len(gcflags) > 0 {
		args = append(args, "-gcflags", strings.Join(gcflags, " "))
	}

	if len(ldflags) > 0 {
		args = append(args, "-ldflags", strings.Join(ldflags, " "))
	}

	if out, err := exec.Command(cmd, args...).CombinedOutput(); err != nil {
		fmt.Println()

		if len(out) > 0 {
			fmt.Println(string(out))
		}

		os.Exit(1)
	}

	fmt.Print("ok ")

	var info []string

	if *debug {
		info = append(info, "debug")
	} else {
		info = append(info, "release")
	}

	if *race {
		info = append(info, "race")
	}

	fmt.Printf("(%v)\n", strings.Join(info, "/"))
}

func closestTag() string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "describe", "--long", "--tags").CombinedOutput()
	if err != nil {
		return "unknown"
	}

	parts := strings.Split(string(out), "-")
	tag := strings.Join(parts[:len(parts)-2], "-")
	commitsAhead, _ := strconv.Atoi(parts[len(parts)-2])
	commitHash := tagCommitHash(tag)

	version := strings.TrimPrefix(strings.TrimSpace(tag), "v")

	var additional []string
	if commitsAhead > 0 {
		noun := "commit"
		if commitsAhead > 1 {
			noun = "commits"
		}

		additional = append(additional, fmt.Sprintf("%v %v ahead of %v", commitsAhead, noun, commitHash))
	}
	if hasUncommittedChanges() {
		additional = append(additional, "built with uncommitted changes")
	}
	if len(additional) > 0 {
		version = fmt.Sprintf("%v (%v)", version, strings.Join(additional, "; "))
	}

	return version
}

func tagCommitHash(tag string) string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "show-ref", "-d", "--tags", tag).CombinedOutput()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, "^{}") {
			parts := strings.Split(line, " ")

			return parts[0][:7]
		}
	}

	return "unknown"
}

func hasUncommittedChanges() bool {
	git, err := exec.LookPath("git")
	if err != nil {
		return false
	}

	out, err := exec.Command(git, "status", "-su").CombinedOutput()
	if err != nil {
		return false
	}

	return len(out) > 0
}

func commitBranch() string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "branch", "--show-current").CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}

func commitHash() string {
	git, err := exec.LookPath("git")
	if err != nil {
		return "unavailable"
	}

	out, err := exec.Command(git, "rev-list", "-1", "HEAD").CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}

// Copyright (C) 2019 The testecho Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Command testecho assists test cases which assert the subject handles process execution as expected. Its flags allow selection of stdout, stderr, exit code, etc.
//
// Usage:
//
//   testecho --help
//
// Print "out" to standard output:
//
//   testecho --stdout out
//
// Same as above but also print "err" to standard error:
//
//   testecho --stdout out --stderr err
//
// Same as above but also exit with code 7 instead of 0:
//
//   testecho --stdout out --stderr err --code 7
//
// Same as above but also sleep for 5 seconds after printing:
//
//   testecho --stdout out --stderr err --code 7 --sleep 5
//
// Spawn another testecho proecss, print its PID, and then sleep "forever" (10000 seconds):
//
//   testecho --spawn
//
// Same as above but also print "err" to standard error:
//
//   testecho --spawn --stderr err
//
// Print standard input:
//
//   echo "out" | testecho
//
// Same as above but also print "err" to standard error:
//
//   echo "out" | testecho --stderr err
//
// Same as above but also exit with code 7 instead of 0:
//
//   echo "out" | testecho --stderr err --code 7
//
// Same as above but also sleep for 5 seconds after printing
//
//   echo "out" | testecho --stderr err --code 7 --sleep 5
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	tp_os "github.com/codeactual/testecho/internal/third_party/stackexchange/os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/codeactual/testecho/internal/cage/cli/handler"
	handler_cobra "github.com/codeactual/testecho/internal/cage/cli/handler/cobra"
	cage_reflect "github.com/codeactual/testecho/internal/cage/reflect"
)

func main() {
	err := handler_cobra.NewHandler(&Handler{
		Session: &handler.DefaultSession{},
	}).Execute()
	if err != nil {
		panic(errors.WithStack(err))
	}
}

const (
	spawnSleepSecs = "10000" // assumption: effectively forever, longer than use case timeouts
)

// Handler defines the sub-command flags and logic.
type Handler struct {
	handler.Session

	Code   int `usage:"Exit with this code"`
	Sleep  int `usage:"Sleep for this number of seconds before exiting (but after printing any selected messages"`
	Spawn  bool
	Stderr string `usage:"Print message to standard error"`
	Stdout string `usage:"Print message to standard output (if stdin is empty)"`
}

// Init defines the command, its environment variable prefix, etc.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Init() handler_cobra.Init {
	return handler_cobra.Init{
		Cmd: &cobra.Command{
			Use:   "testecho",
			Short: "Control this process's standard output/error, exit code, and run duration",
		},
		EnvPrefix: "TESTECHO",
	}
}

// BindFlags binds the flags to Handler fields.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) BindFlags(cmd *cobra.Command) []string {
	cmd.Flags().IntVarP(&h.Code, "code", "", 0, cage_reflect.GetFieldTag(*h, "Code", "usage"))
	cmd.Flags().IntVarP(&h.Sleep, "sleep", "", 0, cage_reflect.GetFieldTag(*h, "Sleep", "usage"))
	cmd.Flags().StringVarP(&h.Stderr, "stderr", "", "", cage_reflect.GetFieldTag(*h, "Stderr", "usage"))
	cmd.Flags().StringVarP(&h.Stdout, "stdout", "", "", cage_reflect.GetFieldTag(*h, "Stdout", "usage"))
	cmd.Flags().BoolVarP(&h.Spawn, "spawn", "", false, "Spawn a second testecho process, print its PID to standard output, block while it sleeps for "+spawnSleepSecs+"s")
	return []string{}
}

// Run performs the sub-command logic.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Run(ctx context.Context, input handler.Input) {
	if len(input.Args) > 0 {
		fmt.Fprintln(h.Err(), "Received unexpected arguments, see --help")
		os.Exit(1)
	}

	if h.Code < 0 {
		fmt.Fprintln(h.Err(), "--code cannot be less than 0")
		os.Exit(1)
	}

	spawnPath, execErr := os.Executable()
	h.ExitOnErrShort(execErr, "failed to spawn another instance", 1)

	var wg sync.WaitGroup

	if h.Spawn {
		wg.Add(1)

		go func() {
			cmd := exec.Command(spawnPath, "--sleep", spawnSleepSecs)
			if startErr := cmd.Start(); startErr != nil {
				panic(startErr)
			}
			fmt.Fprintf(h.Out(), "%d", cmd.Process.Pid)
			if waitErr := cmd.Wait(); waitErr != nil {
				panic(waitErr)
			}

			wg.Done()
		}()
	} else {
		pipeStdin, pipeErr := tp_os.IsPipeStdin()
		if pipeErr != nil {
			panic(pipeErr)
		}

		if pipeStdin {
			if h.Stdout != "" {
				panic(errors.New("received both --stdout and stdin"))
			}

			// For some reason if --sleep is e.g. 1, this will be the last string written to stdout.
			// These related hacks had no effect: h.Out().{Sync,Close} on signal or before Sleep,
			// Println("]"), multi-second Sleep to give time for sync, etc.
			fmt.Fprint(h.Out(), "stdin [")

			stdinScan := bufio.NewScanner(os.Stdin)
			for stdinScan.Scan() {
				fmt.Fprint(h.Out(), stdinScan.Text())
			}
			fmt.Fprint(h.Out(), "]")
			if scanErr := stdinScan.Err(); scanErr != nil {
				panic(scanErr)
			}
		} else {
			fmt.Fprint(h.Out(), h.Stdout)
		}
	}

	fmt.Fprint(h.Err(), h.Stderr)

	if h.Sleep > 0 {
		time.Sleep(time.Duration(h.Sleep) * time.Second)
	}

	wg.Wait()

	os.Exit(h.Code)
}

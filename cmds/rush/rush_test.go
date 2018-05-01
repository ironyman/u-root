// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestRush(t *testing.T) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "TestExit")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for i, tt := range []struct {
		stdin  string // input
		stdout string // output (regular expression)
		stderr string // output (regular expression)
		ret    int    // output
	}{
		// TODO: Create a `-c` flag for rush so stdout does not contain
		// prompts, or have the prompt be derived from $PS1.
		{
			stdin:  "exit\n",
			stdout: "% ",
			stderr: "",
			ret:    0,
		},
		{
			stdin:  "exit 77\n",
			stdout: "% ",
			stderr: "",
			ret:    77,
		},
		{
			stdin:  "exit 1 2 3\n",
			stdout: "% % ",
			stderr: "Too many arguments\n",
			ret:    0,
		},
		{
			stdin:  "exit abcd\n",
			stdout: "% % ",
			stderr: "Non numeric argument\n",
			ret:    0,
		},
		{
			stdin:  "time cd .\n",
			stdout: "% % ",
			stderr: `real 0.0\d\d\n`,
			ret:    0,
		},
		{
			stdin:  "time sleep 0.25\n",
			stdout: "% % ",
			stderr: `real \d+.\d{3}\nuser \d+.\d{3}\nsys \d+.\d{3}\n`,
			ret:    0,
		},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			cmd := testutil.Command(t)
			cmd.Stdin = strings.NewReader(tt.stdin)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			// Check return code
			if err := testutil.IsExitCode(cmd.Run(), tt.ret); err != nil {
				t.Error(err)
			}

			// Check stdout
			strout := string(stdout.Bytes())
			if !regexp.MustCompile("^" + tt.stdout + "$").MatchString(strout) {
				t.Errorf("Want: %#v; Got: %#v", tt.stdout, strout)
			}

			// Check stderr
			strerr := string(stderr.Bytes())
			if !regexp.MustCompile("^" + tt.stderr + "$").MatchString(strerr) {
				t.Errorf("Want: %#v; Got: %#v", tt.stderr, strerr)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

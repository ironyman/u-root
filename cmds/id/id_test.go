// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var logPrefixLength = len("2009/11/10 23:00:00 ")

func TestInvocation(t *testing.T) {
	for _, test := range []struct {
		args     []string
		out      string
		exitCode int
	}{
		{
			args:     []string{"-n"},
			out:      "id: cannot print only names in default format\n",
			exitCode: 1,
		}, {
			args:     []string{"-G", "-g"},
			out:      "id: cannot print \"only\" of more than one choice\n",
			exitCode: 1,
		}, {
			args:     []string{"-G", "-u"},
			out:      "id: cannot print \"only\" of more than one choice\n",
			exitCode: 1,
		}, {
			args:     []string{"-g", "-u"},
			out:      "id: cannot print \"only\" of more than one choice\n",
			exitCode: 1,
		}, {
			args:     []string{"-g", "-u", "-G"},
			out:      "id: cannot print \"only\" of more than one choice\n",
			exitCode: 1,
		},
	} {
		c := testutil.Command(t, test.args...)
		stderr := &bytes.Buffer{}
		c.Stderr = stderr
		err := c.Run()
		if err := testutil.IsExitCode(err, test.exitCode); err != nil {
			t.Error(err)
		}

		e := stderr.String()
		// Ignore the date and time because we're using log.Fatalf.
		if e[logPrefixLength:] != test.out {
			t.Errorf("id %s failed: got %q, want %q", test.args, e, test.out)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

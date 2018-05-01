// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestUniq(t *testing.T) {
	var (
		put1 string = "test\ntest\ngo\ngo\ngo\ncoool\ncoool\ncool\nlegaal\ntest\n"
		put2 string = "u-root\nuniq\nron\nron\nteam\nbinaries\ntest\n\n\n\n\n\n"
	)

	tmpDir, err := ioutil.TempDir("", "UniqTest")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	for i, tt := range []struct {
		in       string
		out      string
		exitCode int
		args     []string
	}{
		{
			in:       put1,
			out:      "test\ngo\ncoool\ncool\nlegaal\ntest\n",
			exitCode: 0,
		},
		{
			in:       put1,
			out:      "2\ttest\n3\tgo\n2\tcoool\n1\tcool\n1\tlegaal\n1\ttest\n",
			exitCode: 0,
			args:     []string{"-c"},
		},
		{
			in:       put1,
			out:      "cool\nlegaal\ntest\n",
			exitCode: 0,
			args:     []string{"-u"},
		},
		{
			in:       put1,
			out:      "test\ngo\ncoool\n",
			exitCode: 0,
			args:     []string{"-d"},
		},
		{
			in:       put2,
			out:      "u-root\nuniq\nron\nteam\nbinaries\ntest\n\n",
			exitCode: 0,
		},
		{
			in:       put2,
			out:      "1\tu-root\n1\tuniq\n2\tron\n1\tteam\n1\tbinaries\n1\ttest\n5\t\n",
			exitCode: 0,
			args:     []string{"-c"}},
		{
			in:       put2,
			out:      "u-root\nuniq\nteam\nbinaries\ntest\n",
			exitCode: 0,
			args:     []string{"-u"},
		},
		{
			in:       put2,
			out:      "ron\n\n",
			exitCode: 0,
			args:     []string{"-d"},
		},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			c := testutil.Command(t, tt.args...)
			c.Stdin = bytes.NewReader([]byte(tt.in))

			o, err := c.CombinedOutput()
			if err := testutil.IsExitCode(err, tt.exitCode); err != nil {
				t.Fatal(err)
			}

			if string(o) != tt.out {
				t.Errorf("uniq %v < %v: got %v, want %v", tt.args, tt.in, string(o), tt.out)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

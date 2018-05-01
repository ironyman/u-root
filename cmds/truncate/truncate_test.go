// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var truncateTests = []struct {
	name  string
	flags []string

	// Exit status of the truncate process.
	ret int

	// if set, a temporary file will be created before the test (used for -c)
	genFile bool

	// if set, we expect that the file will exist after the test
	fileExistsAfter bool

	// -1 to signal we don't care for size test, early continue
	size int64

	// only used when genFile is true
	initSize int64
}{
	{
		name:  "Without args",
		flags: []string{},
		ret:   1,
	},
	{
		name:  "Invalid, valid args, but -s is missing",
		flags: []string{"-c"},
		ret:   1,
	},
	{
		name:  "Invalid, invalid flag",
		flags: []string{"-x"},
		ret:   2,
	},
	{
		name:            "Valid, file does not exist",
		flags:           []string{"-s", "0"},
		ret:             0,
		genFile:         false,
		fileExistsAfter: true,
		size:            0,
	},
	{
		name:            "Valid, file does exist and is smaller",
		flags:           []string{"-s", "1"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        0,
		size:            1,
	},
	{
		name:            "Valid, file does exist and is bigger",
		flags:           []string{"-s", "1"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        2,
		size:            1,
	},
	{
		name:            "Valid, file does exist grow",
		flags:           []string{"-s", "+3K"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        2,
		size:            2 + 3*1024,
	},
	{
		name:            "Valid, file does exist shrink",
		flags:           []string{"-s", "-3"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        5,
		size:            2,
	},
	{
		name:            "Valid, file does exist shrink lower than 0",
		flags:           []string{"-s", "-3M"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        2,
		size:            0,
	},
	{
		name:            "Weird GNU behavior that this actual error is ignored",
		flags:           []string{"-c", "-s", "2"},
		ret:             0,
		genFile:         false,
		fileExistsAfter: false,
		size:            -1,
	},
	{
		name:            "Existing one",
		flags:           []string{"-c", "-s", "3"},
		ret:             0,
		genFile:         true,
		fileExistsAfter: true,
		initSize:        0,
		size:            3,
	},
}

// TestTruncate implements a table-driven test.
func TestTruncate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "truncate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for i, test := range truncateTests {
		t.Run(test.name, func(t *testing.T) {
			testfile := filepath.Join(tmpDir, fmt.Sprintf("txt%d", i))
			if test.genFile {
				data := make([]byte, test.initSize)
				if err := ioutil.WriteFile(testfile, data, 0600); err != nil {
					t.Fatal(err)
				}
			}

			cmd := testutil.Command(t, append(test.flags, testfile)...)
			if err := testutil.IsExitCode(cmd.Run(), test.ret); err != nil {
				t.Fatal(err)
			}

			st, err := os.Stat(testfile)
			if err != nil && test.fileExistsAfter {
				t.Fatalf("Expected %s to exist, but os.Stat() retuned error: %v\n", testfile, err)
			}

			if err != nil {
				return
			}
			if s := st.Size(); test.size != -1 && s != test.size {
				t.Fatalf("Expected that %s has size: %d, but it has size: %d\n", testfile, test.size, s)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

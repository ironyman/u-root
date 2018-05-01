// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type file struct {
	name string
	a    string
	val  []byte
	o    string
	e    string
	x    int // XXX wrong for Plan 9 and Harvey
}

func TestValidate(t *testing.T) {
	var data = []byte(`127.0.0.1	localhost
127.0.1.1	akaros
192.168.28.16	ak
192.168.28.131	uroot

# The following lines are desirable for IPv6 capable hosts
::1     localhost ip6-localhost ip6-loopback
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`)

	tmpDir, err := ioutil.TempDir("", "validate")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := ioutil.WriteFile(filepath.Join(tmpDir, "hosts"), data, 0444); err != nil {
		t.Fatalf("Can't set up data file: %v", err)
	}

	for _, v := range []file{
		// TODO: what kind of table-driven test only has one test? what the fuck?
		{
			name: "hosts.sha1",
			val:  []byte("3f397a3b3a7450075da91b078afa35b794cf6088  hosts"),
			o:    "SHA1\n",
		},
	} {
		if err := ioutil.WriteFile(filepath.Join(tmpDir, v.name), v.val, 0444); err != nil {
			t.Fatalf("Can't set up hash file: %v", err)
		}

		c := testutil.Command(t, filepath.Join(tmpDir, v.name), filepath.Join(tmpDir, "hosts"))
		o, err := c.Output()
		if err != nil {
			t.Fatalf("Can't start StdoutPipe: %v", err)
		}

		if err := testutil.IsExitCode(err, v.x); err != nil {
			t.Error(err)
			continue
		}

		if string(o) != v.o {
			t.Errorf("Validate %v hosts %v (%v): want stdout: %v, got %v)", v.a, v.name, string(v.val), v.o, string(o))
			continue
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

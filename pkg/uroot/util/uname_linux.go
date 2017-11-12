// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"

	"golang.org/x/sys/unix"
)

func toString(d []byte) string {
	return string(d[:bytes.IndexByte(d[:], 0)])
}

// uname does a uname and returns a uroot.Utsname
func Uname() (*Utsname, error) {
	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		return nil, err
	}
	return &Utsname{
		Sysname:    toString(u.Sysname[:]),
		Nodename:   toString(u.Nodename[:]),
		Release:    toString(u.Release[:]),
		Version:    toString(u.Version[:]),
		Machine:    toString(u.Machine[:]),
		Domainname: toString(u.Domainname[:]),
	}, nil
}

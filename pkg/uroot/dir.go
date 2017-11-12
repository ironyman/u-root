// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
	"github.com/u-root/u-root/pkg/ramfs"
)

// DirArchiver is an Archiver that dumps the desired files into a directory.
type DirArchiver struct{}

func (da DirArchiver) OpenOutputFile(path string, 

func (da DirArchiver) Archive(opts ArchiveOpts) error {
	if err := da.arch(opts); err != nil {
		return err
	}
	log.Printf("Archive was unpacked into %q", opts.TempDir)
	return nil
}

func (DirArchiver) arch(opts ArchiveOpts) error {
	ca := CPIOArchiver{
		Format: "newc",
	}

	if fi, err := opts.OutputFile.

	if err := ca.Archive(opts); err != nil {
		return err
	}

	archiver, err := cpio.Format(ca.Format)
	if err != nil {
		return err
	}

	rr := archiver.Reader(opts.OutputFile)
	for {
		rec, err := rr.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("internal error reading records: %v", err)
		}
		if err := cpio.CreateFileInRoot(rec, opts.TempDir); err != nil {
			return fmt.Errorf("Creating %q failed: %v", rec.Name, err)
		}
	}
	return nil
}

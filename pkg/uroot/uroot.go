// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
)

var (
	builders = map[string]Build{
		"source": SourceBuild,
		"bb":     BBBuild,
	}
	archivers = map[string]Archiver{
		"cpio": CPIOArchiver{
			Format: "newc",
		},
	}
)

// BuildOpts are arguments to the Build function.
type BuildOpts struct {
	// Env is the Go environment to use to compile and link packages.
	Env golang.Environ

	// Packages are the Go package import paths to compile.
	//
	// E.g. cmd/go or github.com/u-root/u-root/cmds/ls.
	Packages []string

	// TempDir is a temporary directory where the compilation mode compiled
	// binaries can be placed.
	TempDir string
}

// Build uses the given options to build Go packages and returns a list of
// files to be included in an initramfs archive.
type Build func(BuildOpts) (ArchiveFiles, error)

// ArchiveOpts are the options for building the initramfs archive.
type ArchiveOpts struct {
	// ArchiveFiles are the files to be included.
	ArchiveFiles

	// OutputFile is the file to write to.
	OutputFile *os.File

	// BaseArchive is an existing archive to add files to.
	//
	// BaseArchive may be nil.
	BaseArchive *os.File

	// UseExistingInit determines whether the init from BaseArchive is used
	// or not, if BaseArchive is specified.
	//
	// If this is false, the "init" file in BaseArchive will be renamed
	// "inito" in the output archive.
	UseExistingInit bool
}

// Archiver is an archive format that builds an archive using a given set of
// files.
type Archiver interface {
	// Archive builds an archive file.
	Archive(ArchiveOpts) error

	// DefaultExtension is the default file extension of the archive format.
	DefaultExtension() string
}

// GetBuilder returns the Build function for the named build mode.
func GetBuilder(name string) (Build, error) {
	build, ok := builders[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find builder %q", name)
	}
	return build, nil
}

// GetArchiver returns the archive mode for the named archive.
func GetArchiver(name string) (Archiver, error) {
	archiver, ok := archivers[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find archival format %q", name)
	}
	return archiver, nil
}

// ArchiveFiles are host files and records to add to
type ArchiveFiles struct {
	// Files is a map of relative archive path -> absolute host file path.
	Files map[string]string

	// Records is a map of relative archive path -> Record to use.
	//
	// TODO: While the only archive mode is cpio, this will be a
	// cpio.Record. If or when there is another archival mode, we can add a
	// similar uroot.Record type.
	Records map[string]cpio.Record
}

// NewArchiveFiles returns a new archive files map.
func NewArchiveFiles() ArchiveFiles {
	return ArchiveFiles{
		Files:   make(map[string]string),
		Records: make(map[string]cpio.Record),
	}
}

// SortedKeys returns a list of sorted paths in the archive.
func (af ArchiveFiles) SortedKeys() []string {
	keys := make([]string, 0, len(af.Files)+len(af.Records))
	for dest := range af.Files {
		keys = append(keys, dest)
	}
	for dest := range af.Records {
		keys = append(keys, dest)
	}
	sort.Sort(sort.StringSlice(keys))
	return keys
}

// AddFile adds a host file at `src` into the archive at `dest`.
func (af ArchiveFiles) AddFile(src string, dest string) error {
	if filepath.IsAbs(dest) {
		return fmt.Errorf("cannot add absolute path %q (from %q) to archive", dest, src)
	}
	if !filepath.IsAbs(src) {
		return fmt.Errorf("source file %q (-> %q) must be absolute", src, dest)
	}

	if _, ok := af.Records[dest]; ok {
		return fmt.Errorf("record for %q already exists in archive", dest)
	}
	if srcFile, ok := af.Files[dest]; ok {
		// Just a duplicate.
		if src == srcFile {
			return nil
		}
		return fmt.Errorf("archive file %q already comes from %q", dest, src)
	}

	af.Files[dest] = src
	return nil
}

// AddRecord adds a cpio.Record into the archive at `r.Name`.
func (af ArchiveFiles) AddRecord(r cpio.Record) error {
	if filepath.IsAbs(r.Name) {
		return fmt.Errorf("cannot add absolute path %q to archive", r.Name)
	}

	if _, ok := af.Files[r.Name]; ok {
		return fmt.Errorf("record for %q already exists in archive", r.Name)
	}
	if rr, ok := af.Records[r.Name]; ok {
		if rr.Info == r.Info {
			return nil
		}
		return fmt.Errorf("record for %q already exists", r.Name)
	}

	af.Records[r.Name] = r
	return nil
}

// Contains returns whether path `dest` is already contained in the archive.
func (af ArchiveFiles) Contains(dest string) bool {
	_, fok := af.Files[dest]
	_, rok := af.Records[dest]
	return fok || rok
}

// DefaultPackageImports returns a list of default u-root packages to include.
func DefaultPackageImports(env golang.Environ) ([]string, error) {
	// Find u-root directory.
	urootDir, err := env.FindPackageDir("github.com/u-root/u-root")
	if err != nil {
		return nil, fmt.Errorf("Couldn't find u-root src directory: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(urootDir, "cmds/*"))
	if err != nil {
		return nil, fmt.Errorf("couldn't find u-root cmds: %v", err)
	}
	pkgs := make([]string, 0, len(matches))
	for _, match := range matches {
		importPath, err := env.FindPackageByPath(match)
		if err == nil {
			pkgs = append(pkgs, importPath)
		}
	}
	return pkgs, nil
}


u-root
======

[![Build Status](https://travis-ci.org/u-root/u-root.svg?branch=master)](https://travis-ci.org/u-root/u-root) [![Go Report Card](https://goreportcard.com/badge/github.com/u-root/u-root)](https://goreportcard.com/report/github.com/u-root/u-root) [![GoDoc](https://godoc.org/github.com/u-root/u-root?status.svg)](https://godoc.org/github.com/u-root/u-root) [![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/u-root/blob/master/LICENSE)

# Description

u-root is a "universal root". It's a root file system with mostly Go source with the exception of 5 binaries.

u-root contains simple Go versions of many standard Linux tools, similar to
busybox. u-root can create an initramfs in two different modes:

 * source mode: Go toolchain binaries + simple shell + Go source for tools to be
                compiled on the fly by the shell.

   When you run a command that is not built, you fall through to the command
   that does a `go build` of the command, and then execs the command once it is
   built. From that point on, when you run the command, you get the one in
   tmpfs. This is fast.

 * bb mode: One busybox-like binary comprised of all the Go tools you ask to
            include.

   In this mode, u-root copies and rewrites the source of the tools you asked to
   include to be able to compile everything into one busybox-like binary.

That's the interesting part. This set of utilities is all Go, and mostly source.

# Usage

Make sure your Go version is the latest (>=1.9). Make sure your `GOPATH` is set
up correctly.

Download and install u-root:

```shell
go get github.com/u-root/u-root
```

You can now use the u-root command to build an initramfs. Here are some
examples:

```shell
# Build a bb-mode cpio initramfs of all the Go cmds in ./cmds/...
u-root --build=bb

# Generate a cpio archive named initramfs.cpio.
u-root --format=cpio --build=source -o initramfs.cpio

# Generate a bb-mode archive with only these given commands.
u-root --format=cpio --build=bb ./cmds/{ls,ip,dhclient,wget,tcz,cat}
```

`--format=cpio` and `--build=source` are the defaults. The default set of
packages to include is all packages in `github.com/u-root/u-root/cmds/...`.

In addition to using paths to specify Go source packages to include, you may
also use Go package import paths (e.g. "golang.org/x/tools/imports") to include
commands. Only the `main` package and its dependencies in those source
directories will be included. For example:

```shell
# Both are required for the elvish shell.
go get github.com/boltdb/bolt
go get github.com/elves/elvish
u-root --build=bb ./cmds/\* github.com/elves/elvish
```

Side note: `elvish` is a nicer shell than our default shell `rush`; and also
written in Go.

You can build the initramfs built by u-root into the kernel via the
`CONFIG_INITRAMFS_SOURCE` config variable or you can load it separately via an
option in for example Grub or the QEMU command line or coreboot config variable.

A good way to test the initramfs generated by u-root is with qemu:

```shell
qemu-system-x86_64 -kernel path/to/kernel -initrd /tmp/initramfs.linux_amd64.cpio
```

Note that you do not have to build a special kernel on your own, it is
sufficient to use an existing one. Usually you can find one in `/boot`.

You may also include additional files in the initramfs using the `--files` flag.
As example for Debian, you want to add two kernel modules for testing, executing
your currently booted kernel:

```shell
u-root -files "$HOME/hello.ko $HOME/hello2.ko"
qemu-system-x86_64 -kernel /boot/vmlinuz-$(uname -r) -initrd /tmp/initramfs.linux_amd64.cpio
```

## Getting Packages of TinyCore

Using the `tcz` command included in u-root, you can install tinycore linux
packages for things you want.

You can use QEMU NAT to allow you to fetch packages. Let's suppose, for
example, you want bash. Once u-root is running, you can do this:

```shell
% tcz bash
```

The tcz command computes and fetches all dependencies. If you can't get to
tinycorelinux.net, or you want package fetching to be faster, you can run your
own server for tinycore packages.

You can do this to get a local server using the u-root srvfiles command:

```shell
% srvfiles -p 80 -d path-to-local-tinycore-packages
```

Of course you have to fetch all those packages first somehow :-)

## Build an Embeddable U-root

You can build this environment into a kernel as an initramfs, and further
embed that into firmware as a coreboot payload.

In the kernel and coreboot case, you need to configure ethernet. We have a
`dhclient` command that works for both ipv4 and ipv6. Since v6 does not yet work
that well for most people, a typical invocation looks like this:

```shell
% dhclient -ipv4 -ipv6=false
```

Or, on newer linux kernels (> 4.x) boot with ip=dhcp in the command line,
assuming your kernel is configured to work that way.

# Hardware

If you want to see u-root on real hardware, this
[board](https://www.pcengines.ch/apu2.htm) is a good start.

# Contributions

See [CONTRIBUTING.md](CONTRIBUTING.md)

Improving existing commands (e.g., additional currently unsupported) flags is
very welcome. In this case it is not even required to build an initramfs, just
enter the `cmds/` directory and start coding. A list of commands that are on
the roadmap can be found [here](roadmap.md).

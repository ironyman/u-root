#!/bin/bash
# that's me
DEV_USER=cyl

sudo apt-get update

# my own tools
sudo apt install -y ranger openssh-server 
# apt-file lookupfile
sudo apt install -y apt-file
sudo apt update


# vscode
wget -q https://packages.microsoft.com/keys/microsoft.asc -O- | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main"
sudo apt update
sudo apt install -y code

# https://github.com/oweisse/u-root/tree/kexec_test/booting_windows

sudo apt install -y golang alien kpartx

# kernel build dep https://wiki.ubuntu.com/Kernel/BuildYourOwnKernel#Building_the_kernel
sudo apt-get -y install libncurses-dev flex bison openssl libssl-dev dkms libelf-dev libudev-dev libpci-dev libiberty-dev autoconf
sudo apt install -y git


sudo apt install -y qemu-system-x86-64
sudo gpasswd -a $DEV_USER kvm
# may need to reboot
# This script expect a windows install iso to be in $HOME.
# Copy an install image to this machine, run on domain joined windows machine.

# goes in ~/.profile
cat >> ~/.bashrc <<"EOF"
export WORKSPACE=/home/cyl/uroot/workspace
export EFI_WORKSPACE=/home/cyl/uroot/efi-workspace
export PATH=$PATH:${HOME}/go/bin
export GOPATH=${HOME}/go/
EOF

. ~/.profile
mkdir -p $EFI_WORKSPACE $WORKSPACE $GOPATH

go get github.com/u-root/u-root

## Installing the Modified u-root
pushd ~/go/src/github.com/u-root/u-root
git remote add oweisse git@github.com:ironyman/u-root.git  # our revised uroot repo
git fetch oweisse
git checkout -b kexec_test oweisse/kexec_test
go install
popd

## Setting up windows image
pushd ~/go/src/github.com/u-root/u-root/booting_windows
tar xvf ovmf_uefi.tar.gz -C $EFI_WORKSPACE/
mv ~/*.iso $EFI_WORKSPACE/windows_installer.iso
qemu-img create -f raw "${WORKSPACE}"/windows.img 20G
./install_windows.sh
# finish install in qemu, use default everything.
# Backup!
cp ${WORKSPACE}/{windows.img,windows.img.back}
popd

## Setting up the Kernel.
pushd ~/go/src/github.com/u-root/u-root/booting_windows
./setup.sh
./run_vm.sh rebuild_uroot rebuild_kernel
popd

## Follow instructions in Running u-Root and booting Windows

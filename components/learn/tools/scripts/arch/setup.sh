# Setup Yay
setupAruPackage https://aur.archlinux.org/yay.git

# Install base packages
pacstrap /mnt base base-devel linux linux-firmware btrfs-progs vim git

# Install packages
pacman -S --needed - < packages.txt

# Install packages from AUR
yay -S --needed - < aur.txt

# Helpers

# Install ARU Package Git in $1 
function setupAruPackage() {
  git clone $1
  dir=$(basename $1 .git)
  cd $dir
  makepkg -si
  cd ..
  rm -rf $dir
}
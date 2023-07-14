################## Package Install #####################

## OS ##
# Base Packages and Kernal (option linux-lts, linux-lts-headers)
# base and linux are only bare minimum
pacstrap -i /mnt base linux base-devel linux-headers

## Enter Distro ##
# Key Mismatch: pacman-key --populate; pacman -S archlinux-keyring;

arch-chroot /mnt

#### pacman
# -S Install/Sync
# -R Remove
# -Q Query -Qe Explicit
# -y Update

## Network ##
pacman -S --needed networkmanager wpa_supplicant wireless_tools netctl dialog

## LVM ##
# pacman -S --needed lvm2
# TODO: Add Hooks
# mkinitcpio -p linux

## Essential ##
pacman -S --needed vi git firefox

## Configuration ##
echo "\033[1;33m Setting Region Settings \033[0m \n";
echo "KEYMAP=dvorak" >> /etc/vconsole.conf
echo "aman" >> /etc/hostname

echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
locale-gen

ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
hwclock --systohc
timedatectl

## Services ##
# TODO: enable Important Services
systemctl enable NetworkManager

## Users ##
echo "\033[1;33m Set Root Password \033[0m \n";
passwd

echo "\033[1;33m Create User \033[0m \n";
useradd -m -g users -G wheel aman
passwd aman
#Update /etc/sudoers list
echo "\033[1;34m Uncomment Wheel using visudo \033[0m \n";
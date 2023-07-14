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
pacman -S --needed vi git firefox tldr zsh

## Configuration ##
echo "\033[1;33m Setting Region Settings \033[0m \n";
echo "KEYMAP=dvorak" >> /etc/vconsole.conf
echo "aman" >> /etc/hostname

echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
locale-gen

timedatectl set-timezone Asia/Kolkata
hwclock --systohc
timedatectl

## Services ##
# TODO: enable Important Services
systemctl enable NetworkManager

## Users ##
# https://wiki.archlinux.org/title/Users_and_groups
#groups - Show Groups
#usermod - Modifications
echo "\033[1;33m Set Root Password \033[0m \n";
passwd

echo "\033[1;33m Create User \033[0m \n";
useradd -m --create-home -g users -G wheel aman -s 
passwd aman
#Update /etc/sudoers list
echo "\033[1;34m Uncomment %Wheel using visudo \033[0m \n";

## Grub ##
pacman -S grub efibootmgr dosfstools os-prober mtools
#Populates /mnt/grub and /mnt/efi/EFI Folders
grub-install --target=x86_64-efi --bootloader-id=grub_uefi --recheck
cp /usr/share/locale/en\@quot/LC_MESSAGES/grub.mo /boot/grub/locale/en.mo
# TODO: Crypto /etc/default/grub
grub-mkconfig -o /boot/grub/grub.cfg
# TODO: OS Prober

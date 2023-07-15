################## CHROOT Installation #####################
echo -en "\033[1;33m Driver Installation \033[0m \n";

## Network ##
pacman -S --needed --noconfirm networkmanager wpa_supplicant wireless_tools netctl dialog

## Drivers ##
# pacman -S --needed virtualbox-guest-utils xf86-video-vmware
pacman -S --needed --noconfirm amd-ucode nvidia

## Display ##
pacman -S --needed --noconfirm xorg-server plasma-meta kde-applications

## LVM ##
# pacman -S --needed lvm2
# TODO: Add Hooks
# mkinitcpio -p linux

## Essential ##
pacman -S --needed --noconfirm vi git tldr btrfs-progs

## Grub ##
pacman -S --needed --noconfirm grub efibootmgr dosfstools os-prober mtools

################## Configuration #####################
## Local, Layouts etc ##
echo -en "\033[1;33m Performing Configuration \033[0m \n";
echo "KEYMAP=dvorak" >> /etc/vconsole.conf
hostnamectl set-hostname aman
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
locale-gen

timedatectl set-timezone Asia/Kolkata
hwclock --systohc
timedatectl

## Users ##
echo -en "\033[1;33m User Management \033[0m \n";
useradd -m -g users -G wheel aman
usermod -p `openssl passwd -1 chageme` root
usermod -p `openssl passwd -1 chageme` aman

sed -i '0,/^# %wheel ALL/s/^# //' /etc/sudoers
################## Grub Setup #####################
echo -en "\033[1;33m Grub Setup \033[0m \n";
grub-install --target=x86_64-efi --bootloader-id=grub_uefi --recheck
grub-mkconfig -o /boot/grub/grub.cfg

#Populates /mnt/grub and /mnt/efi/EFI Folders
# cp /usr/share/locale/en\@quot/LC_MESSAGES/grub.mo /boot/grub/locale/en.mo
# TODO: Crypto /etc/default/grub
# TODO: OS Prober
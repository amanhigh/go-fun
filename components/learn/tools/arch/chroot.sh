################## CHROOT Installation #####################
echo -en "\033[1;33m Driver Installation \033[0m \n";

## Display ##
read -p "Install Display (Y) ?: " confirm
if [ "$confirm" == 'Y' ]; then
    # pacman -S --needed virtualbox-guest-utils xf86-video-vmware
    # Minimal: plasma-desktop < plasma-meta < plasma | konsole < kde-applications-meta < kde-applications
    pacman -S --needed --noconfirm nvidia xorg-server sddm plasma-desktop konsole;
fi

## Network ##
pacman -S --needed --noconfirm networkmanager

## Drivers ##
pacman -S --needed --noconfirm amd-ucode ntfs-3g

## LVM ##
# pacman -S --needed lvm2

## Essential ##
pacman -S --needed --noconfirm vi git tldr btrfs-progs

## Grub ##
pacman -S --needed --noconfirm grub efibootmgr os-prober

################## Configuration #####################
## Local, Layouts etc ##
echo -en "\033[1;33m Performing Configuration \033[0m \n";
echo "KEYMAP=dvorak" >> /etc/vconsole.conf
echo "aman" > /etc/hostname
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
locale-gen
hwclock --systohc
echo "\033[1;34m Time Check: `date` \033[0m \n";
# TODO: etchosts

## Users ##
echo -en "\033[1;33m User Management \033[0m \n";
useradd -m -g users -G wheel aman
usermod -p `openssl passwd -1 changeme` root
usermod -p `openssl passwd -1 changeme` aman

sed -i '0,/^# %wheel ALL/s/^# //' /etc/sudoers
################## Grub Setup #####################
echo "\033[1;33m mkinitcpio (Hooks) \033[0m \n";
# TODO: Add Hooks
# mkinitcpio -p linux

echo -en "\033[1;33m Grub Setup \033[0m \n";
grub-install --target=x86_64-efi --bootloader-id=grub_uefi --recheck
grub-mkconfig -o /boot/grub/grub.cfg

#Populates /mnt/grub and /mnt/efi/EFI Folders
# cp /usr/share/locale/en\@quot/LC_MESSAGES/grub.mo /boot/grub/locale/en.mo
# TODO: Crypto /etc/default/grub
# TODO: OS Prober
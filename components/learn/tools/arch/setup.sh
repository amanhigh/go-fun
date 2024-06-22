################## Package Install #####################
## OS ##
# Base Packages and Kernal (option linux-lts, linux-lts-headers)
# base and linux are only bare minimum
echo -en "\033[1;33m Setup Base Packages: Confirm (y/N) ?\033[0m \n";
read base
if [ "$base" == 'y' ]; then
    pacstrap /mnt base linux base-devel linux-headers linux-firmware #Add -i for Interactive
    else
    echo -en "\033[1;34m Skipping Base \033[0m \n";
fi

## Enter Distro ##
# Key Mismatch: pacman-key --populate; pacman -S archlinux-keyring;

## Change Root ##
SCRIPT_PATH=$(cd .;pwd -P)
echo -en "\033[1;32m Changing Root: $SCRIPT_PATH \033[0m \n";
cp $SCRIPT_PATH/chroot.sh /mnt/root/setup.sh
arch-chroot /mnt chmod 755 /root/setup.sh
arch-chroot /mnt /root/setup.sh

## Exit Change Root ##
# Create Snapshot after Setup
if [ "$base" == 'y' ]; then
    btrfs subvolume snapshot /mnt /mnt/.snapshots/base;
else
    echo -en "\033[1;34m Skipping Snapshot \033[0m \n";
fi

echo -en "\033[1;32m Installation Complete \033[0m"

################## Useful Command #####################
## pacman ##
# -S Install/Sync
# -R Remove
# -Q Query -Qe Explicit
# -y Update

## Yay Install ##
# yay <search>
# yay -R <name>
# https://github.com/Jguer/yay

## Key Management ##
# sudo rm -r /etc/pacman.d/gnupg
# sudo pacman-key --init
# sudo pacman-key --populate archlinux
# sudo pacman -Sy archlinux-keyring

## Services ##
# systemctl start <svc>
# systemctl status <svc>
# systemctl enable --now <svc> (Autostart)

## Users ##
# https://wiki.archlinux.org/title/Users_and_groups
# groups - Show Groups
# usermod - Modifications (-s shell) (-p password) (-r remove)
# lslogins - Show all user logins (<name> detailed info)
# faillock --user <name> --reset (Reset Failed Attempts)

## Localization ##
# timedatectl set-timezone Asia/Kolkata
# hostnamectl set-hostname aman
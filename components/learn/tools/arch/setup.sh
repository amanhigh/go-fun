################## Package Install #####################
echo -en "\033[1;33m Base Packages \033[0m \n";
## OS ##
# Base Packages and Kernal (option linux-lts, linux-lts-headers)
# base and linux are only bare minimum
pacstrap /mnt base linux base-devel linux-headers #Add -i for Interactive

## Enter Distro ##
# Key Mismatch: pacman-key --populate; pacman -S archlinux-keyring;

## Change Root ##
SCRIPT_PATH=$(cd .;pwd -P)
echo -en "\033[1;33m Changing Root \033[0m \n";
cp $SCRIPT_PATH/chroot.sh /mnt/root/setup.sh
arch-chroot /mnt /root/setup.sh

## Exit Change Root ##

################## Useful Command #####################
## pacman ##
# -S Install/Sync
# -R Remove
# -Q Query -Qe Explicit
# -y Update

## Services ##
# systemctl start <svc>
# systemctl status <svc>
# systemctl enable --now <svc> (Autostart)

## Users ##
# https://wiki.archlinux.org/title/Users_and_groups
# groups - Show Groups
# usermod - Modifications (-s shell) (-p password) (-r remove)
# lslogins - Show all user logins (<name> detailed info)
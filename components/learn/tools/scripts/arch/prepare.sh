# Arch Linux Prepration Script
# Tutorials
# - https://wiki.archlinux.org/title/Installation_guide
# - https://www.youtube.com/watch?v=DPLnBPM4DhI
# - https://www.learnlinux.tv/arch-linux-full-installation-guide/

################## Basics #####################
# Ctrl+d to exit anywhere

# Set Keyboard
loadkeys dvorak

echo "\033[1;34m UEFI Verify Value: 64 \033[0m \n";
cat /sys/firmware/efi/fw_platform_size


# Wifi - iwctl
# device list
# station wlan0 scan
# station wlan0 connect <sid> -P <password>
# station wlan0 show

# Network Check
# ip addr show (Check inet value)

# TODO: DHCP Setup

#Internet Check
ping archlinux.org

# Time Check
timedatectl

################## Disk Setup #####################
# Disk Info: fdisk -l ; lsblk (-f) ; findmnt ; df -hl
## Setup Partitions ##
# fdisk /dev/sda
# Layout: Boot:/mnt/efi (300MB+), Swap (500MB+), Root:/mnt, Home:/home, Others:
# n - Create Partition, Size (+500M,+5G)
# d - Delete Partition
# p - Print Current Layout
# t - Set Type (EF: UEFI, 8E: LVM)

## Move/Resize Partition ##
# Clone: partclone.btrfs -c -d -s /dev/sda2 -o cloned.img
# Restore: partclone.btrfs -r -s cloned.img -o /dev/sdb1
# Block Copy: partclone.btrfs -b -s /dev/sda2 -o /dev/sdb1
# btrfstune -u /dev/sda2; lsblk -f; (Change UUID)
# mount /dev/sdb1 /mnt (Target Mount)
# btrfs filesystem resize max /mnt (Fix Size)
# arch-chroot /mnt (New Disk: Verify Size and Files)
# Redo all fstab and grub Steps

## btrfs ##
# btrfs check /dev/sdb1 (Check Disk for Errors)

# TODO: LVM Setup
# TODO: Encryption

# Enable Swap (If Required):  swapon /dev/sda2

## Format ##

# Input Partition Names
read -p "Enter EFI Partition Name: " efi
read -p "Enter Primary Partition Name: " primary
read -p "Enter Home Partition Name: " home

# Format EFI using FAT32
mkfs.fat -F32 $efi -n BOOT

# Format Primary and Home using btrfs
mkfs.btrfs $primary -L ROOT
mkfs.btrfs $home -L HOME

## Mounts ##
# findmnt or mount - Show all Mounts
# umount /mnt (-a All) (-R Recursive)

echo "\033[1;32m Mounting Drives \033[0m \n";
mount $primary /mnt
mkdir -p /mnt/boot/efi
mkdir -p /mnt/home
mount $efi /mnt/boot/efi
mount $home /mnt/home

echo "\033[1;33m Generate FsTab \033[0m \n";
mkdir -p /mnt/etc
genfstab -U /mnt >> /mnt/etc/fstab
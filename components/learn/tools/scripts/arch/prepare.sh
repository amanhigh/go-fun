# Arch Linux Prepration Script
# Tutorials
# - https://wiki.archlinux.org/title/Installation_guide
# - https://www.youtube.com/watch?v=DPLnBPM4DhI
# - https://www.learnlinux.tv/arch-linux-full-installation-guide/

################## Basics #####################
# Ctrl+d to exit anywhere
# loadkeys dvorak

echo -en "\033[1;34m UEFI Verify Value: 64 \033[0m \n";
cat /sys/firmware/efi/fw_platform_size

## Wifi - iwctl ##
# device list
# station wlan0 scan
# station wlan0 connect <sid> -P <password>
# station wlan0 show

## Network Check ##
# ip addr show (Check inet value)
ping archlinux.org
# TODO: DHCP Setup

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

# TODO: LVM Setup
# TODO: Encryption


## Format ##
# Input Partition Names
fdisk -l
read -p "Enter Disk Name (Eg. /dev/sda): " disk
read -p "Formatting $disk. Confirm ?: " confirm

boot=${disk}1
root=${disk}2

# Check if the disk value is not 'N'
if [ "$confirm" == 'Y' ]; then
  # Format EFI using FAT32
  mkfs.fat -F32 $boot -n BOOT
  mkfs.btrfs $root -L ROOT
  #swapon /dev/sda2
else
  echo 'Skipping Disk Formatting'
fi

echo -en "\033[1;33m Creating Sub Partitions \033[0m \n";
mount $root /mnt
# https://archive.kernel.org/oldwiki/btrfs.wiki.kernel.org/index.php/SysadminGuide.html#Managing_Snapshots
btrfs sub cr /mnt/@
btrfs sub cr /mnt/@home
btrfs sub cr /mnt/@log
btrfs sub cr /mnt/@snapshots

echo -en "\033[1;33m Mounting Drives \033[0m \n";
# Mount ROOT Subvolume @
umount /mnt
mount -o subvol=@ $root /mnt

# Create directory for each partitions and subvolumes:
mkdir -p /mnt/{etc,boot/efi,home,var/log}

mount -o subvol=@home $root /mnt/home
mount -o subvol=@log $root /mnt/var/log
mount $boot /mnt/boot/efi

echo -en "\033[1;33m Generate Fstab \033[0m \n";
genfstab -U /mnt >> /mnt/etc/fstab

################## Useful Command #####################
## Move/Resize Partition ##
# Clone: partclone.btrfs -c -d -s /dev/sda2 -o cloned.img
# Restore: partclone.btrfs -r -s cloned.img -o /dev/sdb1
# Block Copy: partclone.btrfs -b -s /dev/sda2 -o /dev/sdb1
# btrfstune -u /dev/sda2; lsblk -f; (Change UUID)
# mount /dev/sdb1 /mnt (Target Mount)
# btrfs filesystem resize max /mnt (Fix Size)
# arch-chroot /mnt (New Disk: Verify Size and Files)
# Refresh fstab, grub-install --recheck and grub-mkconfig -o /boot/grub/grub.cfg

## Ram Filesystem ##
# tmpfs - Stays in Swap
# mount -t ramfs -o size=2g ram_bkp /backup

## Mounts ##
# findmnt or mount - Show all Mounts
# umount /mnt (-a All) (-R Recursive)

#### BTRFS ####
# btrfs check /dev/sdb1 (Check Disk for Errors)
## Subvolumes ##
# btrfs sub li /mnt
# btrfs sub cr /mnt/mysub
# btrfs sub del /mnt/mysub
## Snapshot ##
# btrfs sub snap /mnt/mysub /mnt/mysub-backup
# mount -o subvol=mysub-backup /dev/sda2 /mnt/restore/

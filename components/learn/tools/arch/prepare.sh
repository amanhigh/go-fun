# Arch Linux Prepration Script
# ------------------ RUN ------------------
# pacman -Sy git; git clone https://github.com/amanhigh/go-fun;
# cd go-fun/components/learn/tools/arch; ./prepare.sh; ./setup.sh

################## Basics #####################
# Ctrl+d to exit anywhere
# loadkeys dvorak
echo -en "\033[1;33m Basic Chceks \033[0m \n";
echo -en "\033[1;34m UEFI Verify Value: 64 \033[0m \n";
cat /sys/firmware/efi/fw_platform_size

## Network Check ##
# ip addr show (Check inet value)
ping archlinux.org -c1
# TODO: DHCP Setup

# Time Check
timedatectl

################## Disk Setup #####################
# TODO: LVM Setup
# TODO: Encryption

## Format ##
# Input Partition Names
echo -en "\033[1;33m Disk Layout \033[0m \n";
fdisk -l
echo -en "\033[1;33m Disk Formatting \033[0m \n";
read -p "Enter Disk Name (Eg. /dev/sda): " disk
read -p "Formatting $disk. Confirm (y/N) ?: " confirm

boot=${disk}1
root=${disk}2

# Check if the disk value is not 'N'
if [ "$confirm" == 'y' ]; then
  # Format EFI using FAT32
  mkfs.fat -F32 $boot -n BOOT

  # Encrypt Root Partition
  # --type luks2 has Limited Support in Grub
  echo -en "\033[1;33m Encryption \033[0m \n";
  read -p "Encrypt $root. Confirm (y/N) ?: " confirm
  if [ "$confirm" == 'y' ]; then
    cryptsetup luksFormat $root
    cryptsetup open $root cryptroot
    root=/dev/mapper/cryptroot
  fi

  #Normal Format on Crypt Root
  mkfs.btrfs $root -L ROOT
  #swapon /dev/sda3
else
    echo -en "\033[1;33m Skipping Disk Formatting \033[0m \n";
fi

echo -en "\033[1;33m Creating Sub Partitions (Any Key to Continue) \033[0m \n";
read
mountpoint -q /mnt || mount $root /mnt
btrfs sub cr /mnt/@
btrfs sub cr /mnt/@home
btrfs sub cr /mnt/@log
btrfs sub cr /mnt/@snapshots

echo -en "\033[1;33m Mounting Drives \033[0m \n";
read
# Mount ROOT at Subvolume @
mountpoint -q /mnt && umount /mnt
mountpoint -q /mnt || mount -o subvol=@ $root /mnt

# Create directory for each partitions and subvolumes:
mkdir -p /mnt/{etc,boot/efi,home,var/log,.snapshots}

# TODO: -o defaults,noatime,discard=async,ssd,space_cache,compress=zstd,subvol=
mountpoint -q /mnt/home || mount -o subvol=@home $root /mnt/home
mountpoint -q /mnt/var/log || mount -o subvol=@log $root /mnt/var/log
mountpoint -q /mnt/boot/efi || mount $boot /mnt/boot/efi
mountpoint -q /mnt/.snapshots || mount -o subvol=@snapshots $root /mnt/.snapshots
findmnt -R -M /mnt

# cryptsetup luksHeaderBackup /dev/device --header-backup-file /mnt/backup/file.img
# Test Header cryptsetup -v --header /mnt/backup/file.img open /dev/device test
# cryptsetup luksHeaderRestore /dev/device --header-backup-file ./mnt/backup/file.img

echo -en "\033[1;33m Generate Fstab \033[0m \n";
read
genfstab -U /mnt > /mnt/etc/fstab
cat /mnt/etc/fstab

################## Useful Command #####################
# Tutorials
# - https://wiki.archlinux.org/title/Installation_guide
# - https://www.youtube.com/watch?v=DPLnBPM4DhI
# - https://www.learnlinux.tv/arch-linux-full-installation-guide/
## Commands ##
# Ttys - Ctrl + Alt + F1-F10

## Wifi - iwctl ##
# device list
# station wlan0 scan
# station wlan0 connect <sid> -P <password>
# station wlan0 show
# Network Manager
# nmtui
# nmcli device list wifi connect <ssid> password <password>
# lsusb (-s -v)

# iwd Auto Setup Wifi using <Sid> and <Password>
# iwctl --passphrase <password> station wlan0 connect <ssid>
# Ensure password is remembered on reboot


## Setup Partitions ##
# Disk Info: fdisk -l ; lsblk (-f) ; findmnt ; df -hl ; blkid
# fdisk /dev/sda
# Partition Table: GPT (g) or MBR (Backward Compaitable)
# Layout: Boot:/mnt/efi (300MB+), Swap (500MB+), Root:/mnt, Home:/home, Others:
# n - Create Partition, Size (+500M,+5G)
# d - Delete Partition
# p - Print Current Layout
# t - Set Type (EF: UEFI, 8E: LVM) / GPT (1.EFI System)

## Move/Resize Partition ##
# Clone: partclone.btrfs -c -d -s /dev/sda2 -o cloned.img
# Restore: partclone.btrfs -r -s cloned.img -o /dev/sdb1
# Block Copy: partclone.btrfs -b -s /dev/sda2 -o /dev/sdb1
# btrfstune -u /dev/sda2; lsblk -f; (Change UUID)
# mount /dev/sdb1 /mnt (Target Mount);
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
# https://archive.kernel.org/oldwiki/btrfs.wiki.kernel.org/index.php/SysadminGuide.html#Managing_Snapshots
# btrfs sub snap -r /mnt/mysub /mnt/mysub-backup
# mount -o subvol=mysub-backup /dev/sda2 /mnt/restore/

#### Encrypted External Disks #####
# https://www.youtube.com/watch?v=co5V2YmFVEE
# https://github.com/Szwendacz99/Arch-install-encrypted-btrfs
##LUKS
# cryptsetup luksOpen /dev/sda2 cryptroot
# cryptsetup luksClose cryptroot
# cryptsetup luksDum /dev/sda2
## Veracrypt
# cryptsetup --type tcrypt --veracrypt open /dev/sda1 my_decrypted_volume
## Mounting
# mkdir /mnt/my_decrypted_volume
# mount /dev/mapper/my_decrypted_volume /mnt/my_decrypted_volume
## Password Change
# see key slots, max -8 i.e. max 8 passwords can be setup for each device
# cryptsetup luksChangeKey /dev/sda2
# cryptsetup luksAddKey /dev/sda2 (Set New Password)
# cryptsetup luksRemoveKey /dev/sda2 (Remove old Password)

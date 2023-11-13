# Arch Linux Prepration Script
# ------------------ RUN ------------------
# pacman -Sy git; git clone https://github.com/amanhigh/go-fun;
# cd go-fun/components/learn/tools/arch; ./prepare.sh /dev/sdx y; ./setup.sh

if [ $# -ne 2 ]; then
  # Input Partition Names
  echo -en "\033[1;31m Usage:  $0 <Disk Name: /dev/sda> <Encrypt (y/N)> \033[0m \n";
  echo -en "\033[1;33m Disk Layout \033[0m \n";
  fdisk -l
  exit 1
fi

disk=$1
encrypt=$2

echo -en "\033[1;32m Disk: $disk, Encrpt: $encrypt \033[0m \n";

################## Basics #####################
# Ctrl+d to exit anywhere
# loadkeys dvorak
echo -en "\033[1;33m Basic Chceks \033[0m \n";
echo -en "\033[1;34m UEFI Verify Value: 64 \033[0m \n";
cat /sys/firmware/efi/fw_platform_size

## Network Check ##
# ip addr show (Check inet value)
ping archlinux.org -c1

# Time Check
timedatectl

################## Disk Setup #####################
# XXX: LVM Setup

## Format ##
boot=${disk}1
root=${disk}2

echo -en "\033[1;33m Format $disk (y/N) ?\033[0m \n";
read confirm
if [ "$confirm" == 'y' ]; then
  # Format EFI using FAT32
  mkfs.fat -F32 $boot -n BOOT
  
  if [ "$encrypt" == 'y' ]; then
    # Encrypt Root Partition
    # --type luks2 has Limited Support in Grub
    cryptsetup luksFormat --type luks1 $root
    cryptsetup open $root cryptroot
    root=/dev/mapper/cryptroot
  fi

  #Normal Format on Crypt Root
  mkfs.btrfs $root -L ROOT
  #swapon /dev/sda3

  # Subpartitions
  echo -en "\033[1;33m Creating Sub Partitions (Any Key to Continue) \033[0m \n";
  read
  mountpoint -q /mnt || mount $root /mnt
  btrfs sub cr /mnt/@
  btrfs sub cr /mnt/@home
  btrfs sub cr /mnt/@log
  btrfs sub cr /mnt/@snapshots
else
    echo -en "\033[1;34m Skipping Disk Formatting \033[0m \n";
    if [ "$encrypt" == 'y' ]; then
      [ ! -e /dev/mapper/cryptroot ] && cryptsetup open $root cryptroot; root=/dev/mapper/cryptroot
    fi
fi

echo -en "\033[1;33m Mounting Drives \033[0m \n";
read
# Mount ROOT at Subvolume @
mountpoint -q /mnt && umount /mnt
mountpoint -q /mnt || mount -o subvol=@ $root /mnt

# Create directory for each partitions and subvolumes:
mkdir -p /mnt/{root,etc,boot/efi,home,var/log,.snapshots}

MOUNT_OPT="defaults,noatime,discard=async,ssd,space_cache,compress=zstd"
mountpoint -q /mnt/home || mount -o $MOUNT_OPT,subvol=@home $root /mnt/home
mountpoint -q /mnt/var/log || mount -o $MOUNT_OPT,subvol=@log $root /mnt/var/log
mountpoint -q /mnt/.snapshots || mount -o $MOUNT_OPT,subvol=@snapshots $root /mnt/.snapshots
mountpoint -q /mnt/boot/efi || mount $boot /mnt/boot/efi
findmnt -R -M /mnt

# Crypt File
if [ "$encrypt" == 'y' ] && [ ! -f /mnt/root/crypt.keyfile ]; then
    # cryptsetup luksHeaderBackup /dev/device --header-backup-file /mnt/backup/file.img
    # Test Header cryptsetup -v --header /mnt/backup/file.img open /dev/device test
    # cryptsetup luksHeaderRestore /dev/device --header-backup-file ./mnt/backup/file.img

    echo -en "\033[1;34m Generating Crypt File \033[0m \n";
    dd bs=512 count=4 if=/dev/random of=/mnt/root/crypt.keyfile iflag=fullblock
    chmod 000 /mnt/root/crypt.keyfile
    cryptsetup -v luksAddKey ${disk}2 /mnt/root/crypt.keyfile
fi

echo -en "\033[1;33m Generate Fstab (y/N) ? \033[0m \n";
read confirm
if [ "$confirm" == 'y' ]; then
  cp /mnt/etc/fstab /mnt/etc/fstab.bkp || true
  genfstab -U /mnt > /mnt/etc/fstab
  cat /mnt/etc/fstab
  echo -en "\033[1;34m Diff Backup:/mnt/etc/fstab.bkp \033[0m \n";
  [ -e /mnt/etc/fstab.bkp ] && diff /mnt/etc/fstab /mnt/etc/fstab.bkp
else
  echo -en "\033[1;34m Skipping Fstab Generation \033[0m \n";
fi

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

## Move/Resize/Backup/Restore Partition ##
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
# mount (-a Fstab) - Mount all Partitions in /etc/fstab
# umount /mnt (-a All) (-R Recursive)

#### BTRFS ####
# btrfs check /dev/sdb1 (Check Disk for Errors)
## Subvolumes ##
# btrfs sub li /mnt
# btrfs sub cr /mnt/@snapshots
# btrfs sub del /mnt/@snapshots
# btrfs su sh @home
# Defaults
# btrfs sub get-default /mnt
# btrfs sub set-default 5 /mnt (Root Disk)

#### Snapshot ####
# https://archive.kernel.org/oldwiki/btrfs.wiki.kernel.org/index.php/SysadminGuide.html#Managing_Snapshots
## Create
# btrfs sub snap -r /mnt/@ /mnt/@snapshots/root-backup (-r Readonly) [Snapshot Create]
## Restore
# mount -o subvol=@snapshots/root-backup /dev/sda2 /mnt/root-readonly (Temp Mount)
# Writable: Refer Snapper Restore
## Backup Image ##
# btrfs send /home/.snapshots/2/snapshot/ > home.img
# btrfs recieve /restore/home < home.img
## Delete
# btrfs sub del /.snapshots/base/
# btrfs sub del --subvolid 271 / (Delete Root)
# rm -rf /mnt/@ (Change Default SubVol before this)
# btrfs filesystem du -s /.snapshots (Snapshot Size)
## Properties
# btrfs property list -ts /.snapshots/8/snapshot/
# btrfs property get -ts /.snapshots/8/snapshot/ ro
# btrfs property set -ts /.snapshots/8/snapshot/ ro false (Writeable Snapshot)

############### Snapper ############
# https://www.youtube.com/watch?v=sm_fuBeaOqE
## Explore Snapshots
# mount -o subvolid=5 /mnt (Mount Top Disk, No Sub Volume)
# cat /mnt/@snapshots/7/inf.xml
# Writable Restore
# rm -rf /mnt/@ (Change Default SubVol Before: BTRFS:Defaults Section)
# btrfs sub snap /mnt/@snapshots/7/snapshot /mnt/@

#### Encrypted External Disks #####
# https://www.youtube.com/watch?v=co5V2YmFVEE
# https://github.com/Szwendacz99/Arch-install-encrypted-btrfs
##LUKS
# cryptsetup luksOpen /dev/sda2 cryptroot
# cryptsetup luksClose cryptroot
## Veracrypt
# cryptsetup --type tcrypt --veracrypt open /dev/sda1 my_decrypted_volume
# cryptsetup tcryptClose my_decrypted_volume
## Mounting
# mkdir /mnt/my_decrypted_volume
# mount /dev/mapper/my_decrypted_volume /mnt/my_decrypted_volume
## Password Change
# see key slots, max -8 i.e. max 8 passwords can be setup for each device
# cryptsetup luksDump /dev/sda2
# cryptsetup luksChangeKey /dev/sda2
# cryptsetup luksAddKey /dev/sda2 (Set New Password)
# cryptsetup luksRemoveKey /dev/sda2 (Remove old Password)

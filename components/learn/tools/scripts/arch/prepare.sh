# Arch Linux Prepration Script

################## Basics #####################
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

# Setup Timezone as Kolkata and Verify
ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
hwclock --systohc
timedatectl

################## Disk Setup #####################

# fdisk -l See all Disks

# Setup Partitions
# fdisk /dev/sda
# Layout: Boot:/mnt/efi (300MB+), Swap (500MB+), Root:/mnt, Home:/home, Others:
# n - Create Partition, Size (+500M,+5G)
# d - Delete Partition
# p - Print Current Layout
# t - Set Type (EF: UEFI, 8E: LVM)

# TODO: LVM Setup
# TODO: Encryption


# Input Partition Names
read -p "Enter EFI Partition Name: " efi
read -p "Enter Primary Partition Name: " primary
read -p "Enter Home Partition Name: " home

# Format EFI using FAT32
mkfs.fat -F32 /dev/$efi

# Format Primary and Home using btrfs
mkfs.btrfs $primary
mkfs.btrfs $home

echo "\033[1;32m Mounting Drives \033[0m \n";
mount $primary /mnt
mount $efi /mnt/efi
mount $home /mnt/home

echo "\033[1;33m Setting Region Settings \033[0m \n";
echo "KEYMAP=dvorak" >> /etc/vconsole.conf
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
echo "arch" >> /etc/hostname
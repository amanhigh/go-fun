# Arch Linux Prepration Script

# Setup Partitions
parted /dev/sda mklabel gpt
parted /dev/sda mkpart primary fat32 1MiB 513MiB
parted /dev/sda set 1 esp on
parted /dev/sda mkpart primary ext4 513MiB 100%

# Input Partition Names
read -p "Enter EFI Partition Name: " efi
read -p "Enter Primary Partition Name: " primary
read -p "Enter Home Partition Name: " home

# Format EFI using FAT32
mkfs.fat -F32 $efi

# Format Primary and Home using btrfs
mkfs.btrfs $primary
mkfs.btrfs $home

# Setup EFI and Primary Mount
mount $primary /mnt
mount $efi /mnt/boot

# Setup Home Mount
mkdir /mnt/home
mount $home /mnt/home

# Setup Locale and Keyboard to Dvorak
echo "KEYMAP=dvorak" >> /etc/vconsole.conf
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen

# Setup Hostname
echo "arch" >> /etc/hostname

# Setup Timezone as Kolkata
ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
hwclock --systohc

# Setup Wifi
read -p "Enter Wifi SSID: " ssid
read -p "Enter Wifi Password: " password
echo "ctrl_interface=/run/wpa_supplicant" >> /etc/wpa_supplicant/wpa_supplicant.conf
echo "update_config=1" >> /etc/wpa_supplicant/wpa_supplicant.conf
echo "network={" >> /etc/wpa_supplicant/wpa_supplicant.conf
echo "  ssid=\"$ssid\"" >> /etc/wpa_supplicant/wpa_supplicant.conf
echo "  psk=\"$password\"" >> /etc/wpa_supplicant/wpa_supplicant.conf
echo "}" >> /etc/wpa_supplicant/wpa_supplicant.conf

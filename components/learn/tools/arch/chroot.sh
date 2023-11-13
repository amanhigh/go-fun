################## CHROOT Installation #####################
echo -en "\033[1;33m Software Setup \033[0m \n";
## Network ##
pacman -S --needed --noconfirm networkmanager

## Drivers ##
pacman -S --needed --noconfirm amd-ucode ntfs-3g

## LVM ##
# pacman -S --needed lvm2

## Essential ##
pacman -S --needed --noconfirm vim git tldr btrfs-progs cronie

## Grub ##
pacman -S --needed --noconfirm grub efibootmgr os-prober

## Display ##
echo -en "\033[1;33m Install Display: Confirm (y/N) ?\033[0m \n";
read confirm
if [ "$confirm" == 'y' ]; then
    # pacman -S --needed virtualbox-guest-utils xf86-video-vmware
    # Minimal: plasma-desktop < plasma-meta < plasma | konsole < kde-applications-meta < kde-applications
    pacman -S --needed --noconfirm nvidia xorg-server sddm plasma-desktop konsole dolphin firefox;
else
    echo -en "\033[1;34m Skipping Base Setup \033[0m \n";
fi

################## Configuration #####################
## Local, Layouts etc ##
echo -en "\033[1;33m Config & User Management ?: Confirm (y/N) ?\033[0m \n";
read confirm
if [ "$confirm" == 'y' ]; then
    echo "KEYMAP=dvorak" >> /etc/vconsole.conf
    echo "aman" > /etc/hostname
    echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
    ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
    locale-gen
    hwclock --systohc
    echo -en "\033[1;34m Time Check: `date` \033[0m \n";

    ## Users ##
    echo -en "\033[1;33m User Management \033[0m \n";
    useradd -m -g users -G wheel aman
    usermod -p `openssl passwd -1 changeme` root
    usermod -p `openssl passwd -1 changeme` aman

    sed -i '0,/^# %wheel ALL/s/^# //' /etc/sudoers
else
    echo -en "\033[1;34m Skipping Configuration \033[0m \n";
fi
################## Encryption #####################
echo -en "\033[1;33m Generating Encryption Config. Confirm (y/N) ?\033[0m \n";
read confirm
if [ "$confirm" == 'y' ]; then
    ## Hooks and modules
    sed -i 's/^MODULES=()$/MODULES=(btrfs)/' /etc/mkinitcpio.conf
    sed -i 's/^FILES=()$/FILES=(\/root\/crypt.keyfile)/' /etc/mkinitcpio.conf
    sed -i 's/^HOOKS=(.*)$/HOOKS=(base udev autodetect modconf kms keyboard consolefont block encrypt btrfs filesystems fsck)/' /etc/mkinitcpio.conf
    mkinitcpio -p linux
    chmod 600 /boot/initramfs-linux* #Secure embedded crypt.keyfile

    ## Grub Config ##
    sed -i '/^#.*GRUB_ENABLE_CRYPTODISK/s/^#//' /etc/default/grub
    # Set UUID of Encrypted Partition
    ID=`blkid -s UUID -o value -t TYPE=crypto_LUKS`
    sed -i "s/GRUB_CMDLINE_LINUX_DEFAULT.*/GRUB_CMDLINE_LINUX_DEFAULT=\"loglevel=3 quiet cryptdevice=UUID=$ID:cryptroot:allow-discards root=\/dev\/mapper\/cryptroot cryptkey=rootfs:\/root\/crypt.keyfile\"/" /etc/default/grub
    #sed -i '/^#.*GRUB_DISABLE_OS_PROBER/s/^#//' /etc/default/grub
else
    echo -en "\033[1;34m Skipping Encryption Config \033[0m \n";
fi
################## Grub Setup #####################
echo -en "\033[1;33m Grub Install: Confirm (y/N) ?\033[0m \n";
read confirm
if [ "$confirm" == 'y' ]; then
    echo -en "\033[1;33m Grub Setup \033[0m \n";
    grub-install --target=x86_64-efi --bootloader-id=grub_uefi --recheck
    grub-mkconfig -o /boot/grub/grub.cfg
else
    echo -en "\033[1;34m Skipping Grub Install \033[0m \n";
fi

#Populates /mnt/grub and /mnt/efi/EFI Folders
# XXX: OS Prober
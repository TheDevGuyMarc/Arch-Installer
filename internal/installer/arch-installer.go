package installer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// Creates all needed partitions for the system
func setupPartitions(disk string) {
	RunCommand("parted", disk, "mklabel", "gpt")
	// Create a new EFI partition
	RunCommand("parted", disk, "mkpart", "EFI", "fat32", "1MiB", "513MiB")
	RunCommand("parted", disk, "set", "1", "esp", "on")
	
	// Create a new swap partition with 4GiB size
	RunCommand("parted", disk, "mkpart", "SWAP", "linux-swap", "513MiB", "4609MiB")
	
	// Create a new root partition with rest size of harddrive and BTRFS
	RunCommand("parted", disk, "mkpart", "ROOT", "btrfs", "4609MiB", "100%")
}

// Encrypts the root partition of the system for enhanced security
func encryptRootPartition(disk string) {
	// Encrypt the root partition
	RunCommand("cryptsetup", "luksFormat", disk + "3")

	// Try to open the encrypted root partition
	RunCommand("cryptsetup", "open", disk + "3", "cryptroot")
}

// Create subvolumes for btrfs to enable backups and advanced features
func createBtrfsSubvolumes(rootMountPoint string) {
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@")
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@home")
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@var")
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@snapshots")
}

// Format & mount any partition that exists on the hard drive
func mountPartitions(efiPartition, swapPartition, rootPartition, rootMountPoint string) {
	RunCommand("mks.fat", "-F32", efiPartition)

	RunCommand("mkswap", swapPartition)
	RunCommand("swapon", swapPartition)

	RunCommand("mkfs.btrfs", rootPartition)
	RunCommand("mkdir", "-m" "0755", rootMountPoint)
	RunCommand("mount", rootPartition, rootMountPoint)

	createBtrfsSubvolumes(rootMountPoint)
	RunCommand("mount", "-o", "subvol=@", rootPartition, rootMountPoint)
	
	RunCommand("mkdir", "-m", "0755", rootMountPoint + "/home")
	RunCommand("mount", "-o", "subvol=@home", rootPartition, rootMountPoint + "/home")
	
	RunCommand("mkdir", "-m", "0755", rootMountPoint + "/var")
	RunCommand("mount", "-o", "subvol=@var", rootPartition, rootMountPoint + "/var")


	RunCommand("mkdir", "-m", "0755", rootMountPoint + "/.snapshots")
	RunCommand("mount", "-o", "subvol=@snapshots", rootPartition, rootMountPoint + "/.snapshots")

	RunCommand("mkdir", "-m", "0755", rootMountPoint + "/boot")
	RunCommand("mount", efiPartition, rootMountPoint + "/boot")
} 

// Install all needed packages for the base system
func installBaseSystem(rootMountPoint string) {
	RunCommand("pacstrap", rootMountPoint, "linux", "linux-firmware", "networkmanager", "base", "git", "nano")
}

// Generate the fstab file
func generateFstab(rootMountPoint string) {
	RunCommand("genfstab", "-U", rootMountPoint, ">>", rootMountPoint + "/etc/fstab")
}

// Configure the arch linux base installation
func configureBaseInstallation(config SystemConfig, rootMountPoint string) {
	// Configure System Clock to be on correct time zone
	RunCommand("chroot", rootMountPoint, "sh", "-c", "ln", "-sf", "/usr/share/zoneinfo/" + config.Timezone, "/etc/localtime")

	// Configure System clock to be synchronized
	RunCommand("chroot", rootMountPoint, "sh", "-c", "hwclock", "--systohc")

	// Configure locale and system language
	RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", config.Locale, ">>", "/etc/locale.gen")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", config.LocaleLang, ">>", "/etc/locale.gen")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "locale-gen")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", "LANG=" + config.LocaleLang, ">", "/etc/locale.conf")

	// Configure keyboard layout
	RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", "KEYMAP=" + config.Keymap, ">", "/etc/vconsole.conf")

	// Configure hostname
	RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", config.Hostname, ">", "/etc/hostname")

	// Update hosts file to have the hostname included
	RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", "-e", "127.0.0.1\tlocalhost\\n::1\tlocalhost\\n127.0.1.1\t" + config.Hostname + ".localdomain " + config.Hostname, ">>", "/etc/hosts")

	// Enable network services on next startup
	RunCommand("chroot", rootMountPoint, "sh", "-c", "systemctl", "enable", "NetworkManager.service")
}

// Install the grub bootloader and configure it
func installGrubBootLoader(config SystemConfig, rootMountPoint string) {
	// Install needed packages
	RunCommand("pacstrap", rootMountPoint, "grub", "efibootmgr")

	// Install grub on the system
	RunCommand("chroot", rootMountPoint, "sh", "-c", "grub-install", "--target=x86_64-efi", "--efi-directory=/boot")

	// Generate grub configuration files
	RunCommand("chroot", rootMountPoint, "sh", "-c", "grub-mkconfig", "-o", "/boot/grub/grub.cfg")

	// Install grub theme
	if themePath != "" {
		grubConfig := "/mnt/etc/default/grub"
		RunCommand("chroot", rootMountPoint, "sh", "-c", "echo", "GRUB_THEME=" + config.GrubThemePath, "|", "tee", "-a", grubConfig)
		RunCommand("chroot", rootMountPoint, "sh", "-c", "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	}
}

// Installs & configures wireless support (necessary for laptops mostly)
func configureWirelessSupport(rootMountPoint string) {
	RunCommand("pacstrap", rootMountPoint, "iwd")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "systemctl", "enable", "iwd")
}

func installAURHelper(helper string) {
	RunCommand("chroot", rootMountPoint, "sh", "-c", "git", "clone", "https://aur.archlinux.org/" + helper + ".git")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "cd " + helper, "&&", "makepkg", "-si", "--no-confirm")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "rm", "-rf", helper)
}

// Install sddm login manager & configures it
func installSddmLoginManager(config SystemConfig, rootMountPoint) {
	// Install AUR helper
	installAURHelper("paru")
	
	// Install & enable sddm
	RunCommand("pacstrap", rootMountPoint, "sddm")
	RunCommand("chroot", rootMountPoint, "sh", "-c", "systemctl", "enable", "sddm")

	// Configure sddm
	RunCommand("chroot", rootMountPoint, "sh", "-c", "ln", "-sf", config.Dotfilepath + "/sddm/sddm.conf", config.SddmCfgPath)

	// Install sddm theme
	RunCommand("chroot", rootMountPoint, "sh", "-c", "paru", "-S", "--no-confirm", config.SddmTheme)
}

// Runs first update on system (completes correct installation)
func runUpdatesAfterInstallation(rootMountPoint string) {
	RunCommand("chroot", rootMountPoint, "sh", "-c", "pacman", "-Syu", "--no-confirm")
}

func InstallArchBase(config SystemConfig) error {
 	// 1. Partition Discs
	// TODO: Create output for users
 	setupPartitions(config.Disk)

 	// 1.1 Encrypt Root Partition
	// TODO: Create output for users
	encryptRootPartition(config.Disk)

	// 1.2 Mount Partitions
	// TODO: Create output for users
	mountPartitions(config.Disk + 1, config.Disk + 2, "/dev/mapper/cryptroot", "/mnt")

	// 2. Install Base system packages (base, linux, linux-firmware, networkmanager, git, etc.)
	installBaseSystem(rootMountPoint)

	// 3. Generate fstab
	generateFstab(rootMountPoint)

	// 4. Base System configuration (Locale, Keyboard layout, clock, etc.)
	configureBaseInstallation(config)

	// 5. Install Grub + Configuration
	installGrubBootLoader(config, rootMountPoint)

	// 6. Install wireless support (install iwd & enable service)
	configureWirelessSupport(rootMountPoint)

	// 7. Install sddm + Configuration
	installSddmLoginManager(config, rootMountPoint)

	// 8. Update system
	runUpdatesAfterInstallation(rootMountPoint)
}
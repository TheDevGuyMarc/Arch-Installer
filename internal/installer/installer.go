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

func createBtrfsSubvolumes(rootMountPoint string) {
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@")
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@home")
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@var")
	RunCommand("brtfs", "subvolume", "create", rootMountPoint + "/@snapshots")
}

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

	// 3. Generate fstab

	// 4. Base System configuration (Locale, Keyboard layout, clock, etc.)

	// 5. Install Grub + Configuration

	// 6. Install wireless support (install iwd & enable service)

	// 7. Install sddm + Configuration

	// 8. Update system
}
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var version = "1.0.0"

type SystemConfig struct {
	Disk          string
	BootPartition string
	RootPartition string
	Hostname      string
	Timezone      string
	Locale        string
	LocaleLang    string
	Keymap				string
	Dotfilepath		string
	GrubThemePath string
	SddmCfgPath		string
	SddmTheme     string
	RootPassword  string
	Username      string
	UserPassword  string
	Packages      []string
}

// prompt helper function to handle user input with default values
func prompt(reader *bufio.Reader, question, defaultValue string) string {
	fmt.Printf("%s [%s]: ", question, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func getUserInput() SystemConfig {
	reader := bufio.NewReader(os.Stdin)

	defaults := SystemConfig{
		Disk:          "/dev/sda",
		BootPartition: "/dev/sda1",
		RootPartition: "/dev/sda2",
		Hostname:      "archlinux",
		Timezone:      "Europe/Berlin",
		Locale:        "en_US.UTF-8",
		LocaleLang:    "de_DE.UTF-8",
		Keymap:				 "de",
		Dotfilepath		 "/mnt/dotfiles",
		GrubThemePath: "",
		SddmCfgPath:   "/etc/sddm.conf",
		SddmTheme:     "",
		RootPassword:  "rootpassword",
		Username:      "g405t",
		UserPassword:  "userpassword",
		Packages: []string{
			"linux", "linux-firmware", "grub", "efibootmgr", "sddm", "xorg", "vim", "git", "base",
		},
	}

	fmt.Println("Please enter your configuration. Press Enter to use default values.")

	disk := prompt(reader, "Disk (e.g., /dev/sda)", defaults.Disk)
	bootPartition := prompt(reader, "Boot Partition (e.g., /dev/sda1)", defaults.BootPartition)
	rootPartition := prompt(reader, "Root Partition (e.g., /dev/sda2)", defaults.RootPartition)
	hostname := prompt(reader, "Hostname", defaults.Hostname)
	timezone := prompt(reader, "Timezone", defaults.Timezone)
	locale := prompt(reader, "Locale", defaults.Locale)
	localeLang := prompt(reader, "Locale Language (e.g., en_US.UTF-8)", defaults.LocaleLang)
	keymap := prompt(reader, "Keyboard Layout (e.g., de, us)", defaults.Keymap)
	dotfilePath := prompt(reader, "Path to your dotfile repository (e.g., /home/dotfiles, /mnt/dotfiles)", defaults.Dotfilepath)
	grubThemePath := prompt(reader, "Path to grub theme (e.g., /usr/share/grub/themes/themename/theme.txt)", defaults.GrubThemePath)
	sddmConfigPath := prompt(reader, "Path to sddm config (e.g., /etc/sddm.conf)", defaults.SddmCfgPath)
	sddmTheme := prompt(reader, "Name for sddm theme (e.g., theme-git)", defaults.SddmTheme)
	rootPassword := prompt(reader, "Root Password", defaults.RootPassword)
	username := prompt(reader, "Username", defaults.Username)
	userPassword := prompt(reader, "User Password", defaults.UserPassword)

	color.Yellow("\nDefault package list:")
	for _, pkg := range defaults.Packages {
		fmt.Printf("  - %s\n", pkg)
	}
	fmt.Println("You can install more packages later if you need them.")

	return SystemConfig{
		Disk:          disk,
		BootPartition: bootPartition,
		RootPartition: rootPartition,
		Hostname:      hostname,
		Timezone:      timezone,
		Locale:        locale,
		LocaleLang:    localeLang,
		Keymap:				 keymap,
		Dotfilepath:	 dotfilePath,
		GrubThemePath: grubThemePath,
		SddmCfgPath:	 sddmConfigPath,
		SddmTheme:     sddmTheme,
		RootPassword:  rootPassword,
		Username:      username,
		UserPassword:  userPassword,
		Packages:      defaults.Packages,
	}
}

func main() {
	PrintSpecs()

	config := getUserInput()

	// Execute Arch Installtion steps
	if err := InstallArchBase(config); err != nil {
		log.Fatalf("Installation failed: %v", err)
	}

	// Install Desktop Environment
	InstallDesktopEnvironment()

	// Configure System
	InstallDotFiles()

	// Configure Rice
	InstallRice()

	// Reboot System
	color.Green("\nInstallation complete! Rebooting the system. Please remove your boot stick...")
	Reboot()
}
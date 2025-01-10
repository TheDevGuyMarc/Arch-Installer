package specs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ANSI Color codes
const (
	COLOR_RESET  = "\033[0m"
	COLOR_RED    = "\033[31m"
	COLOR_GREEN  = "\033[32m"
	COLOR_YELLOW = "\033[33m"
	COLOR_BLUE   = "\033[34m"
	COLOR_PURPLE = "\033[35m"
	COLOR_CYAN   = "\033[36m"
	COLOR_WHITE  = "\033[97m"
)

// RunCommand executes a shell command and returns its output.
func RunCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error running command '%s': %v", command, err)
	}
	return strings.TrimSpace(out.String())
}

// GetCPUInfo retrieves CPU information.
func GetCPUInfo() string {
	return RunCommand("lscpu")
}

// GetGPUInfo retrieves GPU information.
func GetGPUInfo() string {
	return RunCommand("lspci", "|", "grep", "-A2", "VGA")
}

// GetRAMInfo retrieves memory information.
func GetRAMInfo() string {
	return RunCommand("free", "-h")
}

// GetBootMode determines if the system uses BIOS or UEFI.
func GetBootMode() string {
	efi := RunCommand("efibootmgr")
	if efi != "" {
		return "UEFI"
	}
	return "BIOS"
}

// GetBlockDevices retrieves information about partitions and block devices.
func GetBlockDevices() string {
	return RunCommand("lsblk", "-o", "NAME,FSTYPE,SIZE,TYPE,MOUNTPOINT")
}

// PrintHeader prints a formatted header.
func PrintHeader(title string) {
	fmt.Printf("%s\n%s%s%s\n\n", ColorCyan, strings.Repeat("=", 40), ColorReset, title)
}

// PrintSection prints a section of specs with a colored header.
func PrintSection(title, data string, color string) {
	fmt.Printf("%s[%s]%s\n", color, title, ColorReset)
	fmt.Println(data)
	fmt.Println(strings.Repeat("-", 40))
}

// PrintSpecs gathers all specs and prints them in an organized way.
func PrintSpecs() {
	fmt.Printf("%s=== Arch Linux Installation Specs ===%s\n\n", ColorPurple, ColorReset)

	PrintHeader("System Information")
	PrintSection("CPU Info", GetCPUInfo(), ColorGreen)
	PrintSection("GPU Info", GetGPUInfo(), ColorYellow)
	PrintSection("RAM Info", GetRAMInfo(), ColorBlue)

	PrintHeader("Boot Information")
	PrintSection("Boot Mode", GetBootMode(), ColorCyan)

	PrintHeader("Disk Information")
	PrintSection("Block Devices", GetBlockDevices(), ColorRed)

	fmt.Printf("%s%s\n", strings.Repeat("=", 40), ColorReset)
}
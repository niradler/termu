package tools

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Tool struct {
	Name        string
	BinaryName  string
	Description string
	CargoName   string
	BrewName    string
	ScoopName   string
	AptName     string
}

var ModernTools = []Tool{
	{
		Name:        "fd",
		BinaryName:  "fd",
		Description: "Fast file finder",
		CargoName:   "fd-find",
		BrewName:    "fd",
		ScoopName:   "fd",
		AptName:     "fd-find",
	},
	{
		Name:        "ripgrep",
		BinaryName:  "rg",
		Description: "Fast content search",
		CargoName:   "ripgrep",
		BrewName:    "ripgrep",
		ScoopName:   "ripgrep",
		AptName:     "ripgrep",
	},
	{
		Name:        "bat",
		BinaryName:  "bat",
		Description: "File preview with syntax highlighting",
		CargoName:   "bat",
		BrewName:    "bat",
		ScoopName:   "bat",
		AptName:     "bat",
	},
	{
		Name:        "sd",
		BinaryName:  "sd",
		Description: "Search and replace",
		CargoName:   "sd",
		BrewName:    "sd",
		ScoopName:   "sd",
		AptName:     "",
	},
	{
		Name:        "xsv",
		BinaryName:  "xsv",
		Description: "CSV toolkit",
		CargoName:   "xsv",
		BrewName:    "xsv",
		ScoopName:   "",
		AptName:     "",
	},
	{
		Name:        "jaq",
		BinaryName:  "jaq",
		Description: "JSON processor",
		CargoName:   "jaq",
		BrewName:    "jaq",
		ScoopName:   "",
		AptName:     "",
	},
	{
		Name:        "yq",
		BinaryName:  "yq",
		Description: "YAML/JSON processor",
		CargoName:   "",
		BrewName:    "yq",
		ScoopName:   "yq",
		AptName:     "",
	},
	{
		Name:        "dua",
		BinaryName:  "dua",
		Description: "Disk usage analyzer",
		CargoName:   "dua-cli",
		BrewName:    "dua-cli",
		ScoopName:   "dua",
		AptName:     "",
	},
	{
		Name:        "eza",
		BinaryName:  "eza",
		Description: "Modern directory listing",
		CargoName:   "eza",
		BrewName:    "eza",
		ScoopName:   "eza",
		AptName:     "eza",
	},
}

type PackageManager struct {
	Name       string
	CheckCmd   string
	InstallCmd func(string) []string
	UpdateCmd  []string
	NeedsSudo  bool
	Available  bool
}

type Installer struct {
	OS             string
	PackageManager *PackageManager
	HasCargo       bool
}

func NewInstaller() *Installer {
	inst := &Installer{
		OS:       runtime.GOOS,
		HasCargo: commandExists("cargo"),
	}
	inst.detectPackageManager()
	return inst
}

func (i *Installer) detectPackageManager() {
	var managers []*PackageManager

	switch i.OS {
	case "linux":
		managers = []*PackageManager{
			{
				Name:       "apt",
				CheckCmd:   "apt",
				InstallCmd: func(pkg string) []string { return []string{"apt", "install", "-y", pkg} },
				UpdateCmd:  []string{"apt", "update", "-y"},
				NeedsSudo:  true,
			},
			{
				Name:       "dnf",
				CheckCmd:   "dnf",
				InstallCmd: func(pkg string) []string { return []string{"dnf", "install", "-y", pkg} },
				UpdateCmd:  []string{"dnf", "check-update"},
				NeedsSudo:  true,
			},
			{
				Name:       "pacman",
				CheckCmd:   "pacman",
				InstallCmd: func(pkg string) []string { return []string{"pacman", "-Syu", "--noconfirm", pkg} },
				UpdateCmd:  []string{"pacman", "-Sy"},
				NeedsSudo:  true,
			},
			{
				Name:       "apk",
				CheckCmd:   "apk",
				InstallCmd: func(pkg string) []string { return []string{"apk", "add", pkg} },
				UpdateCmd:  []string{"apk", "update"},
				NeedsSudo:  true,
			},
		}
	case "darwin":
		managers = []*PackageManager{
			{
				Name:       "brew",
				CheckCmd:   "brew",
				InstallCmd: func(pkg string) []string { return []string{"brew", "install", pkg} },
				UpdateCmd:  []string{"brew", "update"},
				NeedsSudo:  false,
			},
		}
	case "windows":
		managers = []*PackageManager{
			{
				Name:       "scoop",
				CheckCmd:   "scoop",
				InstallCmd: func(pkg string) []string { return []string{"scoop", "install", pkg} },
				UpdateCmd:  []string{"scoop", "update"},
				NeedsSudo:  false,
			},
			{
				Name:       "choco",
				CheckCmd:   "choco",
				InstallCmd: func(pkg string) []string { return []string{"choco", "install", "-y", pkg} },
				UpdateCmd:  []string{"choco", "upgrade", "all"},
				NeedsSudo:  false,
			},
		}
	}

	for _, mgr := range managers {
		if commandExists(mgr.CheckCmd) {
			mgr.Available = true
			i.PackageManager = mgr
			return
		}
	}
}

func (i *Installer) CheckInstalledTools() (installed, missing []Tool) {
	for _, tool := range ModernTools {
		if commandExists(tool.BinaryName) {
			installed = append(installed, tool)
		} else {
			missing = append(missing, tool)
		}
	}
	return
}

func (i *Installer) InstallTool(tool Tool) error {
	fmt.Printf("➡️  Installing %s (%s)...\n", tool.Name, tool.Description)

	var pkgName string
	useCargo := false

	if i.PackageManager != nil && i.PackageManager.Available {
		switch i.PackageManager.Name {
		case "brew":
			pkgName = tool.BrewName
		case "scoop", "choco":
			pkgName = tool.ScoopName
		case "apt":
			pkgName = tool.AptName
		default:
			pkgName = tool.Name
		}

		if pkgName == "" {
			useCargo = true
		}
	} else {
		useCargo = true
	}

	if useCargo {
		if !i.HasCargo {
			return fmt.Errorf("no package manager found and cargo is not available")
		}
		return i.installWithCargo(tool)
	}

	return i.installWithPackageManager(pkgName)
}

func (i *Installer) installWithPackageManager(pkgName string) error {
	cmdParts := i.PackageManager.InstallCmd(pkgName)

	if i.PackageManager.NeedsSudo {
		cmdParts = append([]string{"sudo"}, cmdParts...)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	return nil
}

func (i *Installer) installWithCargo(tool Tool) error {
	if tool.CargoName == "" {
		return fmt.Errorf("tool not available via cargo")
	}

	fmt.Printf("   Using cargo to install %s\n", tool.CargoName)

	cmd := exec.Command("cargo", "install", tool.CargoName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cargo installation failed: %w", err)
	}

	return nil
}

func (i *Installer) UpdatePackageManager() error {
	if i.PackageManager == nil || !i.PackageManager.Available {
		return nil
	}

	if len(i.PackageManager.UpdateCmd) == 0 {
		return nil
	}

	fmt.Printf("🔄 Updating %s package list...\n", i.PackageManager.Name)

	cmdParts := i.PackageManager.UpdateCmd
	if i.PackageManager.NeedsSudo {
		cmdParts = append([]string{"sudo"}, cmdParts...)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (i *Installer) PrintStatus() {
	fmt.Println("🔍 Detecting environment...")
	fmt.Printf("🧠 Detected OS: %s\n", i.OS)

	if i.PackageManager != nil && i.PackageManager.Available {
		fmt.Printf("📦 Package Manager: %s\n", i.PackageManager.Name)
	} else if i.HasCargo {
		fmt.Println("📦 Package Manager: cargo (fallback)")
	} else {
		fmt.Println("⚠️  No package manager detected")
	}
}

func (i *Installer) VerifyInstallation() {
	fmt.Println("\n✅ Installation complete! Verifying...")
	fmt.Println()

	installed, missing := i.CheckInstalledTools()

	for _, tool := range installed {
		path, _ := exec.LookPath(tool.BinaryName)
		fmt.Printf("✔️  %s: %s\n", tool.Name, path)
	}

	if len(missing) > 0 {
		fmt.Println()
		for _, tool := range missing {
			fmt.Printf("⚠️  %s failed to install. Try installing manually.\n", tool.Name)
		}
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (i *Installer) InstallAll() error {
	i.PrintStatus()
	fmt.Println()

	installed, missing := i.CheckInstalledTools()

	fmt.Println("🔧 Checking installed tools...")
	for _, tool := range installed {
		fmt.Printf("✅ %s already installed\n", tool.Name)
	}

	if len(missing) == 0 {
		fmt.Println("\n🎉 All tools are already installed!")
		return nil
	}

	fmt.Println()
	var missingNames []string
	for _, tool := range missing {
		fmt.Printf("❌ %s not found\n", tool.Name)
		missingNames = append(missingNames, tool.Name)
	}

	fmt.Printf("\n📦 Installing missing tools: %s\n", strings.Join(missingNames, ", "))
	fmt.Println()

	if i.PackageManager != nil && i.PackageManager.Available {
		if err := i.UpdatePackageManager(); err != nil {
			fmt.Printf("⚠️  Failed to update package manager: %v\n", err)
		}
	}

	for _, tool := range missing {
		if err := i.InstallTool(tool); err != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", tool.Name, err)
		}
	}

	i.VerifyInstallation()

	fmt.Println()
	fmt.Println("✨ Setup finished!")

	return nil
}

package boot

import (
	"fmt"
	"strings"

	"PXE-Manager/internal/storage"
)

func GenerateBootScript(serverBase string, profile storage.Profile) string {
	serverBase = strings.TrimRight(serverBase, "/")

	switch profile.BootType {
	case "kernel_initrd":
		return fmt.Sprintf(`#!ipxe
kernel %s%s %s
initrd %s%s
boot
`, serverBase, profile.Kernel, profile.Cmdline, serverBase, profile.Initrd)

	case "iso":
		return fmt.Sprintf(`#!ipxe
dhcp
sanboot %s%s
`, serverBase, profile.ImagePath)

	case "sanboot":
		return fmt.Sprintf(`#!ipxe
dhcp
sanboot %s%s
`, serverBase, profile.ImagePath)

	default:
		return "#!ipxe\necho Unsupported boot type\nshell\n"
	}
}

func GenerateMenuScript(serverBase string, profiles []storage.Profile) string {
	serverBase = strings.TrimRight(serverBase, "/")

	var b strings.Builder
	b.WriteString("#!ipxe\n\n")
	b.WriteString("menu PXE Boot Menu\n")

	for _, p := range profiles {
		if !p.Enabled {
			continue
		}
		key := menuKey(p.Name)
		b.WriteString(fmt.Sprintf("item %s %s\n", key, p.Name))
	}

	b.WriteString("choose target && goto ${target}\n\n")

	for _, p := range profiles {
		if !p.Enabled {
			continue
		}
		key := menuKey(p.Name)
		b.WriteString(fmt.Sprintf(":%s\n", key))
		b.WriteString(generateProfileEntry(serverBase, p))
		b.WriteString("\n")
	}

	return b.String()
}

func generateProfileEntry(serverBase string, profile storage.Profile) string {
	switch profile.BootType {
	case "kernel_initrd":
		return fmt.Sprintf("kernel %s%s %s\ninitrd %s%s\nboot\n",
			serverBase, profile.Kernel, profile.Cmdline,
			serverBase, profile.Initrd,
		)

	case "iso":
		return fmt.Sprintf("dhcp\nsanboot %s%s\n",
			serverBase, profile.ImagePath,
		)

	case "sanboot":
		return fmt.Sprintf("dhcp\nsanboot %s%s\n",
			serverBase, profile.ImagePath,
		)

	default:
		return "echo Unsupported boot type\nshell\n"
	}
}

func menuKey(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

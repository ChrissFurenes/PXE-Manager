package boot

import "fmt"

func DefaultIPXEMenu(serverIP string, httpPort string) string {
	baseURL := fmt.Sprintf("http://%s:%s", serverIP, httpPort)

	return fmt.Sprintf(`#!ipxe

menu PXE Boot Menu
item ubuntu Ubuntu Installer
item rescue Rescue System
item shell iPXE Shell
choose target && goto ${target}

:ubuntu
kernel %s/ubuntu/vmlinuz ip=dhcp
initrd %s/ubuntu/initrd
boot

:rescue
kernel %s/rescue/vmlinuz ip=dhcp
initrd %s/rescue/initramfs.img
boot

:shell
shell
`, baseURL, baseURL, baseURL, baseURL)
}

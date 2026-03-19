# Pingu PX 

## ⚠️**Under development!!!**

## Testing
when testing sett default boot file in dhcp server to ``undionly.kpxe``
remember to change ip in ``config.yaml``

## some links
* [Github ipxe](https://github.com/ipxe/ipxe)
* [IPXE commands](https://ipxe.org/cmd)
* [Talos PXE](https://docs.siderolabs.com/talos/v1.12/platform-specific-installations/bare-metal-platforms/pxe)
* [Ubuntu netboot](https://ubuntu.com/server/docs/how-to/installation/netboot-the-server-installer-via-uefi-pxe-on-arm-aarch64-arm64-and-x86-64-amd64/)
* [Debian](https://wiki.debian.org/DebianInstaller/NetbootAssistant)
* [Proxmox](https://forum.proxmox.com/threads/automated-installation-pxe-boot.169009/)
* [Proxmox1](https://forum.proxmox.com/threads/install-proxmox-host-via-pxe.108281/)
* [ESXI](https://techdocs.broadcom.com/us/en/vmware-cis/vsphere/vsphere/8-0/esx-upgrade/upgrading-esxi-hosts-upgrade/how-to-boot-an-esxi-host-from-a-network-device-upgrade/boot-the-esxi-installer-by-using-pxe-and-tftp-upgrade.html)
* [UNIFI OS](https://help.ui.com/hc/en-us/articles/34210126298775-Self-Hosting-UniFi)



## Status
some stuff work.
you ich boot from selected os.

Web UI runs on port ``:8080``

## How to run it
Clone repo and run with go

``git clone https://github.com/ChrissFurenes/PXE-Manager.git``

``cd PXE-Manager``

``sudo go run cmd/pxe-server/main.go``
## Talos info
to boot talos with config add this to ``Kernel cmdline``:

```talos.platform=metal slab_nomerge pti=on ip=eth0:dhcp talos.config=http://10.230.0.212:8080/files/talos/worker.yaml```

## ideas and planing
- [X] Does have a web UI
- [X] Can select image on boot
- [ ] Making so iso boot works 100%
- [ ] rules to boot from image
- [ ] online status? 
- [ ] client groups
- [ ] image groups
- [ ] logs
- [ ] logs per client
- [ ] logs per assets
- [ ] logs per profile
- [ ] audit log
- [ ] Use jinja2 for web?????
- [ ] Talos/cloud-init profile designer
- [ ] Making web UI fancy
- [ ] DARK MODE!!
- [ ] API
- [ ] Webhook
- [ ] Visual boot flow builder
- [ ] Approval queue unknown clients
- [ ] Web terminal for boot logs
- [ ] Embedded mini inventory
- [ ] Image marketplace
- [ ] support for VLAN
- [ ] QR support????
- [ ] mobile UX
- [ ] boot recipes
- [ ] plugins
- [ ] scan network for MAC and hostname
- [ ] Adding support for env
- [ ] auth on web UI
- [ ] roles access
- [ ] metrics for eks grafana
- [ ] drag and drop in web UI
- [ ] Withe list
- [ ] Block list
- [ ] Password auth to boot from image
- [ ] Support for talos config files (worker/controlplane)
- [ ] support for cloud-init

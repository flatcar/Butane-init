<div style="text-align: center">

[![Flatcar OS](https://img.shields.io/badge/Flatcar-Website-blue?logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhLS0gR2VuZXJhdG9yOiBBZG9iZSBJbGx1c3RyYXRvciAyNi4wLjMsIFNWRyBFeHBvcnQgUGx1Zy1JbiAuIFNWRyBWZXJzaW9uOiA2LjAwIEJ1aWxkIDApICAtLT4NCjxzdmcgdmVyc2lvbj0iMS4wIiBpZD0ia2F0bWFuXzEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4Ig0KCSB2aWV3Qm94PSIwIDAgODAwIDYwMCIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgODAwIDYwMDsiIHhtbDpzcGFjZT0icHJlc2VydmUiPg0KPHN0eWxlIHR5cGU9InRleHQvY3NzIj4NCgkuc3Qwe2ZpbGw6IzA5QkFDODt9DQo8L3N0eWxlPg0KPHBhdGggY2xhc3M9InN0MCIgZD0iTTQ0MCwxODIuOGgtMTUuOXYxNS45SDQ0MFYxODIuOHoiLz4NCjxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik00MDAuNSwzMTcuOWgtMzEuOXYxNS45aDMxLjlWMzE3Ljl6Ii8+DQo8cGF0aCBjbGFzcz0ic3QwIiBkPSJNNTQzLjgsMzE3LjlINTEydjE1LjloMzEuOVYzMTcuOXoiLz4NCjxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik02NTUuMiw0MjAuOXYtOTUuNGgtMTUuOXY5NS40aC0xNS45VjI2MmgtMzEuOVYxMzQuOEgyMDkuNFYyNjJoLTMxLjl2MTU5aC0xNS45di05NS40aC0xNnY5NS40aC0xNS45djMxLjINCgloMzEuOXYxNS44aDQ3Ljh2LTE1LjhoMTUuOXYxNS44SDI3M3YtMTUuOGgyNTQuOHYxNS44aDQ3Ljh2LTE1LjhoMTUuOXYxNS44aDQ3Ljh2LTE1LjhoMzEuOXYtMzEuMkg2NTUuMnogTTQ4Ny44LDE1MWg3OS42djMxLjgNCgloLTIzLjZ2NjMuNkg1MTJ2LTYzLjZoLTI0LjJMNDg3LjgsMTUxTDQ4Ny44LDE1MXogTTIzMywyMTQuNlYxNTFoNjMuN3YyMy41aC0zMS45djE1LjhoMzEuOXYyNC4yaC0zMS45djMxLjhIMjMzVjIxNC42eiBNMzA1LDMxNy45DQoJdjE1LjhoLTQ3Ljh2MzEuOEgzMDV2NDcuN2gtOTUuNVYyODYuMUgzMDVMMzA1LDMxNy45eiBNMzEyLjYsMjQ2LjRWMTUxaDMxLjl2NjMuNmgzMS45djMxLjhMMzEyLjYsMjQ2LjRMMzEyLjYsMjQ2LjRMMzEyLjYsMjQ2LjR6DQoJIE00NDguMywzMTcuOXY5NS40aC00Ny44di00Ny43aC0zMS45djQ3LjdoLTQ3LjhWMzAyaDE1Ljl2LTE1LjhoOTUuNVYzMDJoMTUuOUw0NDguMywzMTcuOXogTTQ0MCwyNDYuNHYtMzEuOGgtMTUuOXYzMS44aC0zMS45DQoJdi03OS41aDE1Ljl2LTE1LjhoNDcuOHYxNS44aDE1Ljl2NzkuNUg0NDB6IE01OTEuNiwzMTcuOXY0Ny43aC0xNS45djE1LjhoMTUuOXYzMS44aC00Ny44di0zMS43SDUyOHYtMTUuOGgtMTUuOXY0Ny43aC00Ny44VjI4Ni4xDQoJaDEyNy4zVjMxNy45eiIvPg0KPC9zdmc+DQo=)](https://www.flatcar.org/)
[![Discord](https://img.shields.io/badge/Discord-Chat%20with%20us!-5865F2?logo=discord)](https://discord.gg/PMYjFUsJyq)
[![Slack](https://img.shields.io/badge/Slack-Chat%20with%20us!-4A154B?logo=slack)](https://kubernetes.slack.com/archives/C03GQ8B5XNJ)
[![Twitter Follow](https://img.shields.io/twitter/follow/flatcar?style=social)](https://x.com/flatcar)
[![Mastodon Follow](https://img.shields.io/badge/Mastodon-Follow-6364FF?logo=mastodon)](https://hachyderm.io/@flatcar)
[![Bluesky](https://img.shields.io/badge/Bluesky-Follow-0285FF?logo=bluesky)](https://bsky.app/profile/flatcar.org)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10926/badge)](https://www.bestpractices.dev/projects/10926)



</div>



# Butane-Init: a Cloud-Init to Butane Config Transpiler

### Background

Flatcar uses Butane and Ignition for provisioning. Much of the wider ecosystem still produces cloud-config (cloud-init) YAML, which makes it harder to adopt Flatcar in environments where cloud-config is already the norm.

### The Project

A transpiler in Go that converts cloud-config YAML into Butane YAML.
The transpiler implements the minimum functionality needed to support ClusterAPI worker node provisioning, and to replace [coreos-cloudinit](https://github.com/flatcar/coreos-cloudinit)
In practice that limits the transpiler to basic Butane features such as users, groups, certificates, files and systemd units.

#### Goals and Non-Goals
- A working Go transpiler covering a documented subset of cloud-config.
- Test coverage for each supported feature, running in CI.
- End to end integration with Flatcar and ClusterAPI.
- Documentation covering supported fields and known limitations.

### Usage

Build `bt`, then pass a cloud-config file or pipe one on standard input:

```sh
go build -o bt ./cmd/bt
bt cloud-config.yaml > butane.yaml
# or: cat cloud-config.yaml | bt -o butane.yaml
```

`bt` currently converts only the Cluster API worker `users` subset. A config
is either converted completely or rejected; unsupported fields are never
silently discarded.

| cloud-config field      | Butane field          |
| ----------------------- | --------------------- |
| `name`                  | `name`                |
| `passwd`                | `password_hash`       |
| `gecos`                 | `gecos`               |
| `homedir`               | `home_dir`            |
| `shell`                 | `shell`               |
| `ssh_authorized_keys`   | `ssh_authorized_keys` |

Password hashes remain locked, matching cloud-init's default
`lock_passwd: true`. Supplementary and primary groups, sudo configuration,
password unlocking, and every non-`users` stanza are not yet supported.

### A Flatcar Container Linux project

Flatcar Container Linux is a fully open source, minimal-footprint, secure by default and always up-to-date Linux distribution for running containers at scale.

This repository is maintained by the Flatcar community and contributors. Development and design decisions are guided by the principles of transparency, modularity, and collaboration.

Please find information on:

- [Contribution Guide](https://github.com/flatcar/Flatcar/blob/main/CONTRIBUTING.md) – how to get involved and submit patches.
- [Repository Codeowners](./CODEOWNERS) – who's assigned for contribution reviews in this repository.
- [Repository Maintainers](./MAINTAINERS.md) – who’s responsible for what in this repository.
- [Flatcar Project Maintainers](https://github.com/flatcar/Flatcar/blob/main/MAINTAINERS.md) – who’s responsible for what at Flatcar.
- [Governance](https://github.com/flatcar/Flatcar/blob/main/governance.md) – how decisions are made and who the stakeholders are.
- [Code of Conduct](./code-of-conduct.md) – our commitment to a welcoming and inclusive community.
- [LICENSE](./LICENSE) - License under which the software is released.

---

## Community & Project Documentation

- [Contributing Guidelines](CONTRIBUTING.md) — How to contribute, find issues, and submit pull requests
- [Code of Conduct](CODE_OF_CONDUCT.md) — Standards for respectful and inclusive community participation
- [Security Policy](SECURITY.md) — How to report vulnerabilities and security-related information
- [Maintainers](MAINTAINERS.md) — Current project maintainers and their responsibilities
- [Governance](GOVERNANCE.md) — Project governance model, decision-making process, and roles

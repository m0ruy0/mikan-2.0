# Mikan 2.0

**Mikan 2.0 is an independent development fork of [Mikan](https://github.com/Miroshka000/mikan), a VPN panel built on the mihomo core.**

Our goal is to keep the core panel practical and open while building additional capabilities for community operators and developers.

> **Development status:** experimental fork. Do not use this branch as a production installer until builds, migration paths, release signing and upgrades have been validated.

## Project priorities

- **Pterodactyl nodes:** investigate a least-privilege Egg for regular game hosting, without requiring root access or a Docker socket inside the server container.
- **A real addon marketplace:** multiple Telegram bots, payment integrations and third-party extensions behind a documented permission model.
- **Developer API:** scoped API keys, signed webhooks, stable versioning and examples for building custom bots.
- **Free core features:** panel migration, external subscription sources and a separate Telegram admin bot.

The detailed roadmap and technical constraints are in [Mikan 2.0 roadmap](docs/mikan20-roadmap.md).

## Upstream and license

- Upstream project: [Miroshka000/mikan](https://github.com/Miroshka000/mikan)
- License: [GNU GPL v3](LICENSE)

This fork retains upstream attribution and applicable license notices. Mikan 2.0 is a modified project and is not an official release of the upstream maintainers.

## Current status

The roadmap describes planned work, not completed functionality. The Pterodactyl mode in particular depends on publishing and verifying standalone node binaries, then testing network allocations against a real Pterodactyl/Wings installation.

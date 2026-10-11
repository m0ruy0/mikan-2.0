<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset=".github/assets/banner-en-dark.svg">
  <img alt="Mikan 2.0 — open fork of a mihomo VPN panel" src=".github/assets/banner-en-light.svg" width="100%">
</picture>

<br>

**Mikan 2.0 is an open Mikan fork with a free roadmap and its own release process.**

[![Repository](https://img.shields.io/badge/GitHub-Mikan%202.0-ff853d?style=flat-square&logo=github)](https://github.com/m0ruy0/mikan-2.0)
[![License](https://img.shields.io/badge/license-GPL--3.0-2f9e6b?style=flat-square)](LICENSE)
[![Core](https://img.shields.io/badge/core-mihomo-9e8eff?style=flat-square)](https://github.com/MetaCubeX/mihomo)

[Русский](README.md) · **English**

</div>

---

## Quick start

> **Important:** the first signed Mikan 2.0 release has not been published yet. The install command below will become functional after the release secret is configured and the first release succeeds. Do not run it on a production server yet.

After the first release is published, installation on a clean Ubuntu 22.04+ or Debian 12+ VPS (amd64/arm64) will be a single command:

~~~bash
curl -fsSL https://github.com/m0ruy0/mikan-2.0/releases/latest/download/install.sh | sudo bash
~~~

The installer verifies the release manifest signature and the downloaded binary's SHA-256 checksum before starting the setup wizard. Do not use the original Mikan installer or Docker image for Mikan 2.0.

**Documentation:** [installation](https://m0ruy0.github.io/mikan-2.0/en/start/install/) · [panel overview](https://m0ruy0.github.io/mikan-2.0/en/start/overview/) · [protocols](https://m0ruy0.github.io/mikan-2.0/en/protocols/presets/)

## What the inherited core includes

- **VPS/VDS nodes** — separate panel nodes, connection management and statistics.
- **mihomo and multiple protocols** — configuration presets and traffic accounting for supported protocols.
- **App-specific subscriptions** — formats tailored for compatible VPN clients.
- **Telegram bot and Mini App** — subscription management and supported payment adapters.
- **Plans, devices and promo codes** — core customer-management features.
- **Updates and backups** — inherited infrastructure being adapted for this fork.

## What Mikan 2.0 plans to add for free

| Feature | Plan |
|---|---|
| Pterodactyl nodes | An Egg for game hosting without root or a Docker socket inside the container; compatibility depends on the provider's TCP/UDP allocations. |
| Addon marketplace | Multiple ready-made Telegram bots, payment adapters and third-party integrations. |
| API for custom bots | Documented API, scoped keys and webhooks. |
| External subscriptions as nodes | Safe import and refresh with last-known-good configuration. |
| Panel migration | Backup, restore validation and rollback. |
| Telegram admin bot | Roles, administrator allowlist and action audit. |

This is a **roadmap**, not a claim that these features are already implemented. See [docs/mikan20-roadmap.md](docs/mikan20-roadmap.md).

## Development and releases

The project is derived from [Mikan](https://github.com/Miroshka000/mikan). Mikan 2.0 is a modified fork under GNU GPL v3, not an official upstream release. Original attribution and license notices are retained.

- [Mikan 2.0 source](https://github.com/m0ruy0/mikan-2.0)
- [Releases](https://github.com/m0ruy0/mikan-2.0/releases)
- [Roadmap](docs/mikan20-roadmap.md)
- [Signed release setup](docs/releasing.md)

Go tests, web/site builds and Rust tests must pass before a release. Signed releases must not be published until RELEASE_SIGNING_KEY is configured and all CI checks succeed.

## License

[GNU GPL v3](LICENSE). Commercial use and redistribution are permitted under the license terms.

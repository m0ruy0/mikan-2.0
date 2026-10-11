# Mikan 2.0 — roadmap

Mikan 2.0 is an independent fork of [Mikan](https://github.com/Miroshka000/mikan). It keeps the inherited core while adding capabilities that remain part of the free core release.

> Status: roadmap and prototype planning. Items below are not implemented merely because they appear here.

## Design principles

- Keep changes in the open and preserve upstream attribution and GPL-3.0 notices.
- Keep the core feature set free; fundraising is not a gate for the features below.
- Preserve existing VPS node compatibility where possible.
- Never give an addon the host's Docker socket, panel database, private keys or broad admin token by default.
- Require tests, documentation, error handling, security review and an upgrade/rollback plan before a feature is called ready.

## Pterodactyl nodes

Goal: run a node on regular Pterodactyl game hosting when the customer only has a server allocation and cannot use SSH/root or a Docker socket.

Current starting point: `cmd/mikan-node` is a separate Go binary, but the existing VPS installer depends on Docker, host networking and host-level setup. Those assumptions cannot be made inside a game-server container.

MVP:
1. Publish a standalone Linux amd64/arm64 node binary with a SHA-256 checksum and release-manifest signature.
2. Provide a Pterodactyl Egg that runs as the container user and stores state in persistent storage.
3. Configure the node API port and client TCP/UDP ports from explicit allocations.
4. Add a panel workflow that validates the API endpoint and allocations before marking the node healthy.
5. Test start/restart, saved counters, reconnection, logs and graceful shutdown against a real Pterodactyl/Wings host.

Limitations: host providers may filter UDP or not expose enough ports. TUN, WARP, cascades and other host-network-dependent features must not be advertised unless independently verified in that environment. No root, privileged-container, host-network or Docker-socket access is required for the MVP.

Acceptance: a test node joins a test panel, accepts client traffic only on assigned ports, survives a container restart and reports healthy status and counters.

## Addon marketplace

Extend the existing payment-adapter marketplace into a common catalog while keeping payment adapters as a category.

Planned categories:
- several separate ready-made Telegram bots (sales/renewal, customer support, notifications and administration);
- payment adapters;
- analytics and monitoring;
- referrals/promotions;
- CRM and other external integrations.

Each addon manifest should declare its ID, publisher, version, compatibility, category, signed artifact digest, documentation, permissions and configuration schema. Installs/updates must be verifiable and reversible. Addons run isolated and receive explicit, narrow tokens; they do not receive unrestricted panel secrets or Docker access.

## API for third-party bots

Build on the existing API-key mechanism:
- documented endpoints and stable versioning;
- scopes for clients, subscriptions, plans, devices, nodes and read-only reports;
- expiry/revocation, rate limits, usage tracking and audit events;
- signed webhooks with timestamps and replay protection;
- sample Python bot/SDK.

Payment secrets, refunds and dangerous admin actions remain session-only unless explicitly delegated through a stronger reviewed permission flow.

## Free core features

### Panel migration
Move the database, users, plans, settings, nodes and configuration with a backup, dry-run, restore verification and rollback. DNS/IP changes remain a separate operator task.

### External subscriptions as node sources
Fetch and refresh an external subscription, normalize only supported formats, show last refresh status and preserve the last-known-good configuration when an upstream is unavailable. Add SSRF safeguards for DNS results, redirects, private/loopback/link-local addresses, cloud metadata, ports and response size.

### Telegram admin bot
A separate bot for administrators, with explicit enablement, an admin-ID allowlist, role checks, confirmations for destructive actions and an audit trail. Never send bot tokens or API keys into chat or logs.

## Work order

1. Adapt branding, release metadata and CI to the Mikan 2.0 fork.
2. Verify the signed release/installer path.
3. Build and test the Pterodactyl Egg proof of concept.
4. Define addon manifest and API scopes/webhooks; preserve payment-adapter compatibility.
5. Implement migration, external sources and the admin bot in small, tested changes.

A feature is complete only when implementation, tests, documentation, failure handling, security review and safe upgrade/rollback are present.

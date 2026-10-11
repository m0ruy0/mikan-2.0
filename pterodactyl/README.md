# Mikan 2.0 on Pterodactyl

This Egg runs the standalone Go `mikan-node` process in a regular Pterodactyl container. It does **not** run Docker inside the game-server container and does not require root privileges.

## Release prerequisite

The Egg verifies the Ed25519 signature on the release manifest and the SHA-256 checksum of the selected node binary. It will refuse to start until the `latest` Mikan 2.0 release contains a signed `node` section for both `x86_64` and `aarch64`. This repository change adds those assets to future release builds; it does not retroactively add binaries to already published releases.

## Setup

1. In Pterodactyl, import [egg-mikan-node.json](egg-mikan-node.json).
2. Create a server using the **Python 3.12** image. The container must provide `curl`, `openssl`, `python3`, `base64` and `sha256sum`.
3. Allocate a primary port for the node's management API, plus enough additional ports for the client inbounds you intend to offer. The provider must allow inbound TCP and/or UDP as required by those protocols.
4. In Mikan 2.0, open **Nodes → Add node**. In **API port**, enter exactly the same port as the Pterodactyl server's primary allocation. Use the public address that clients and the panel can reach.
5. Copy the one-time join key shown by Mikan. In the Pterodactyl server's Startup variables, paste only the key into `MIKAN_NODE_JOIN`, not the entire generated command.
6. Start the server. The bootstrap downloads the signed manifest, verifies its Ed25519 signature, selects the current architecture, downloads the node binary and verifies SHA-256. It also checks that the join key's API port matches `SERVER_PORT` before launching.
7. In Mikan, configure the node's **client inbounds** to use only TCP/UDP ports allocated to that Pterodactyl server. The API allocation and client listening ports are separate allocations. Do not expose unallocated ports in subscription profiles.

## Limitations

- Compatibility depends on the hosting provider: UDP may be filtered and port allocations may be limited.
- The Egg is for user-space TCP/UDP listeners. Do not expect host-network features such as TUN, WARP or routing changes that require kernel capabilities unavailable inside a game-server container.
- A passing CI build proves the binary compiles; it does not replace an end-to-end test against a real Pterodactyl/Wings host.
- Keep `MIKAN_NODE_JOIN` secret. Re-issue the node key in Mikan if it is exposed.

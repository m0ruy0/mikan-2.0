# Publishing Mikan 2.0 releases

Mikan 2.0 uses its own release namespace and Ed25519 signing key. **Never use the original Mikan signing key or copy its signed manifest.**

## One-time repository setup

1. In repository **Settings → Secrets and variables → Actions**, create a repository secret named `RELEASE_SIGNING_KEY`.
2. Its value must be the full PKCS#8 PEM from the private key generated specifically for Mikan 2.0. Do not commit this key or paste it into an issue/PR.
3. Enable GitHub Actions for this fork.
4. Run the CI workflow and fix all failures before creating a release tag.

The public half of the key is embedded in `internal/release/release.go`, `installer/src/release.rs` and `installer/install.sh`; those values must remain synchronized.

## First release

After CI is green and the repository secret is configured, create a version tag on `main` using the repository's release policy (for example `v2.0.0`). The release workflow builds the web bundle, Go binaries and installer, publishes the multi-architecture image to GitHub Container Registry and creates signed manifest/index assets.

The one-command installer will work only after the first successful signed release is published:

```bash
curl -fsSL https://github.com/m0ruy0/mikan-2.0/releases/latest/download/install.sh | sudo bash
```

Before that, the command intentionally has no release asset to download. Do not install from original Mikan release URLs.

## Security

Treat the private signing key as a production secret. If it is exposed, create a new keypair, update the public key in all three source locations and revalidate the update/rotation plan before publishing new releases.

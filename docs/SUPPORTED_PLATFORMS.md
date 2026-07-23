# Supported Platforms

## Verified Baseline

SHantilly's required build-and-test platform is the repository's
`src/dev.Dockerfile`: Debian Trixie on x86-64 with the distribution Qt6
packages. The GitHub Actions job **Debian Trixie / Qt6** must compile the
production executable and pass the complete CTest suite.

Use the same baseline locally:

```bash
docker build --pull -t shantilly-dev -f src/dev.Dockerfile .
docker run --rm -v "$PWD:/workspace" -w /workspace shantilly-dev \
  bash -lc 'cmake -S . -B /tmp/build && cmake --build /tmp/build --parallel 2 && ctest --test-dir /tmp/build --output-on-failure'
```

## Other Platforms

Ubuntu, Fedora, Arch Linux, macOS, and Windows packaging exists or is planned,
but those environments are not yet compatibility baselines. Treat them as
experimental until a required CI job builds production code and passes the
same CTest suite. Platform promotion requires updating this document and the
roadmap with CI evidence.

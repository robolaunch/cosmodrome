# <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/logos/svg/rocket.svg" width="40" height="40" align="top"> **cosmodrome** - Image Pipeline

<div align="center">
  <p align="center">
    <a href="https://github.com/robolaunch/robot-image-pipeline/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/robolaunch/robot-image-pipeline" alt="license">
    </a>
    <a href="https://github.com/robolaunch/robot-image-pipeline/issues">
      <img src="https://img.shields.io/github/issues/robolaunch/robot-image-pipeline" alt="issues">
    </a>
  </p>
</div>

**cosmodrome** is an image pipeline that produces images.
## Quick Start

### Installation

Install with Go:

```bash
go install github.com/robolaunch/cosmodrome@latest
```

Install binary:

```bash
wget https://github.com/robolaunch/cosmodrome/releases/download/v0.1.0-alpha.1/cosmodrome-amd64
chmod +x cosmodrome-amd64
mv cosmodrome-amd64 /usr/local/bin/cosmodrome
cosmodrome --help
```

### Usage

For example, building Freecad, run:

```bash
cosmodrome launch --config pipelines/freecad.yaml
```

To enable pushes, add your personal access token to your environment first:

```bash
export REGISTRY_PAT=<YOUR-PAT>
```

## Contributing

Please see [this guide](./CONTRIBUTING) if you want to contribute.

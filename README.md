# ink‑stream – Kindle e‑book Converter

`ink‑stream` is a command‑line toolkit that extracts, processes, and converts image‑heavy documents (e.g. comic strips,
scanned books, PDFs) into Kindle‑friendly formats (EPUB, MOBI, AZW3).  
It provides a **pipeline** of image‑processing steps (grayscale, rescale, autocrop, margin wrap, etc.) that can be tuned
with flags or by editing the Go source.

---

## Table of Contents

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Command‑line Options](#command-line-options)
6. [Contributing](#contributing)
7. [License](#license)

---

## Features

| Feature                         | Description                                                                                                            |
|---------------------------------|------------------------------------------------------------------------------------------------------------------------|
| **Multi‑step image processing** | Grayscale conversion, auto‑crop, margin wrap, auto‑contrast, gamma correction, rescaling, Gaussian blur, etc.          |
| **Device profiles**             | Pre‑defined Kindle device profiles (`profile`) that set optimal resolution, margins, orientation, and colour handling. |
| **Multiple output formats**     | EPUB, MOBI, AZW3 (via `-format`).                                                                                      |
| **Batch processing**            | Process whole directories (`-src`/`-out`) or individual files.                                                         |
| **CLI flags**                   | Toggle each step, set crop level, rotate, stretch, etc.                                                                |
| **Docker devcontainer**         | Ready‑to‑run development environment.                                                                                  |
| **Test suite**                  | Unit tests for image pipelines and palette handling.                                                                   |

---

## Prerequisites

- **Go** 1.25.1 or newer
- **make** (optional, for convenience scripts)

---

## Installation

```shell script
go install github.com/jictyvoo/ink_stream/cmd/inkonverter@latest
```

---

## Usage

The tool ships with three executables:

| Sub‑command       | Purpose                                                                          |
|-------------------|----------------------------------------------------------------------------------|
| `comiconverter`   | Convert comic‑style documents (e.g. .zip, .cbz, .pdf) to Kindle‑friendly format. |
| `extractor`       | Pull out images from an archive or document for further processing.              |
| `kindleconverter` | Take already‑processed images and package them into an EPUB/MOBI/AZW3.           |

---

## Command‑line Options

| Flag              | Type   | Default     | Description                                                      |
|-------------------|--------|-------------|------------------------------------------------------------------|
| `-src`            | string | `""`        | Path to the source folder containing images or archives.         |
| `-out`            | string | `""`        | Output folder where converted files will be written.             |
| `-rotate`         | bool   | `false`     | Rotate images 90° clockwise before processing.                   |
| `-colored`        | bool   | `false`     | Keep pages in colour; otherwise convert to grayscale.            |
| `-margins`        | bool   | `false`     | Add margin around images (margin color defaults to white).       |
| `-stretch`        | bool   | `false`     | Stretch images to fit target resolution.                         |
| `-crop-level`     | uint   | `CropBasic` | Level of auto‑cropping (basic, aggressive, etc.).                |
| `-profile`        | string | `""`        | Name of a pre‑defined device profile (e.g. `kindle_paperwhite`). |
| `-format`         | string | `epub`      | Output format: `epub`, `mobi`, or `azw3`.                        |
| `-read-direction` | string | `""`        | Reading direction (`ltr`, `rtl`, `vertical`).                    |

> See the source of `cmd/kindleconverter/args_parser.go` for the full enumeration of available options and default
> values.

---

## Sub‑commands

### 1. `extractor`

```shell script
inkextract -src ./input -out ./images
```

- Pulls out all images from the source folder without further processing. Handy for debugging or manual re‑processing.

### 2. `kindleconverter`

```shell script
inkonverter -src ./images -out ./epub -format epub
```

- Packages a folder of already‑processed images into an EPUB (or MOBI/azw3 if specified).

---

## Contributing

Pull requests are welcome!  
Please run tests before submitting:

```shell script
go test ./...
```

Follow the Go style guidelines (`gofmt -w .`) and keep commit messages concise.

---

## License

MIT – see [LICENSE](LICENSE) for details.

---

> *If you have any questions or suggestions, feel free to open an issue or reach out via the project's issue tracker.*

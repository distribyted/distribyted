[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL3 License][license-shield]][license-url]
[![Coveralls][coveralls-shield]][coveralls-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/distribyted/distribyted">
    <img src="docs/images/distribyted_icon.png" alt="Logo" width="100">
  </a>

  <h3 align="center">distribyted</h3>

  <p align="center">
    Torrent client with on-demand file downloading as a filesystem.
    <br />
    <br />
    <a href="https://github.com/distribyted/distribyted/issues">Report a Bug</a>
    ·
    <a href="https://github.com/distribyted/distribyted/issues">Request Feature</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->
## Table of Contents

- [Table of Contents](#table-of-contents)
- [About The Project](#about-the-project)
  - [Use Cases](#use-cases)
  - [Supported _Expandable_ File Formats](#supported-_expandable_-file-formats)
    - [Supported](#supported)
    - [To Be Supported](#to-be-supported)
    - [Not Supported](#not-supported)
- [Getting Started](#getting-started)
  - [Prerequisites on windows](#prerequisites-on-windows)
- [Usage](#usage)
  - [Configuration File](#configuration-file)
    - [root](#root)
    - [mountpoints](#mountpoints)
- [Contributing](#contributing)
- [License](#license)

<!-- ABOUT THE PROJECT -->
## About The Project

![Distribyted Screen Shot][product-screenshot]

Distribyted tries to make easier integrations with other applications among torrent files, presenting them as a standard filesystem. 

We aim to use some compressed file characteristics to avoid download it entirely, just the parts that we'll need.

Also, if the file format is not supported, distribyted can stream and seek through the file if needed.

**Note that distribyted is in alpha version, it is a proof of concept with a lot of bugs.**

### Use Cases

- Play **multimedia files** on your favorite video or audio player. These files will be downloaded on demand and only the needed parts.
- Explore TBs of data from public **datasets** only downloading the parts you need. Use **Jupyter Notebooks** directly to process or analyze this data.
- Play your **ROM backups** directly from the torrent file. You can have virtually GBs in games and only downloaded the needed ones.

### Supported _Expandable_ File Formats
Distribyted can show some kind of files directly as folders, making it possible for applications read only the parts that they need. Here is a list of supported, to be supported and not supported formats.

#### Supported
- zip: Able to uncompress just one file. The file is decompressed to a temporal file sequentially to make possible seek over it. The decompression stops if no one is reading it.

#### To Be Supported
- tar: Seek to any file and inside that files using a [modified standard library](https://github.com/ajnavarro/go-tar). Not useful on `.tar.gz` files.
- 7zip: Similar to Zip. Need for a library similar to [zip](https://github.com/saracen/go7z).
- xz: Only worth it when the file is created using blocks. Possible library [here](https://github.com/ulikunitz/xz) and [here](https://github.com/frrad/bxzf).

#### Not Supported
- gzip: As far as I know, it doesn't support random access.

## Getting Started

Get the latest release from [releases][releases-url] page or download the source code and execute `make build`.

Run the program: `./distribyted-[VERSION]-[OS]-[ARCH]`

Defaults are good enough for starters, but you can change them. Here is the output of `./distribyted -help`:

```text
NAME:
   distribyted - Torrent client with on-demand file downloading as a filesystem.

USAGE:
   distribyted [global options] [arguments...]

GLOBAL OPTIONS:
   --config value      YAML file containing distribyted configuration. (default: "./distribyted-data/config.yaml") [$DISTRIBYTED_CONFIG]
   --http-port value   HTTP port for web interface (default: 4444) [$DISTRIBYTED_HTTP_PORT]
   --fuse-allow-other  Allow other users to acces to all fuse mountpoints. You need to add user_allow_other flag to /etc/fuse.conf file. (default: false) [$DISTRIBYTED_FUSE_ALLOW_OTHER]
   --help, -h          show help (default: false)
```

### Prerequisites on windows

Download and install [WinFsp](http://www.secfs.net/winfsp/).

## Usage

After executing and load all torrent or magnet files, a web interface will be available here: `http://localhost:4444`
It contains information about the mounted routes and torrent files like download/upload speed, leechers, seeders...

You can also modify the configuration file and reload the server from here: `http://localhost:4444/config` .

### Configuration File

#### root

|Config key|Description|
|-|-|
|max-cache-size| Size in MB for the cache. This is the maximum space used by distribyted to store torrent data. Less used torrent data will be discarded if this value is reached.|
|metadata-folder-name| Folder where distribyted metadata will be stored.|
|mountPoints|List of folders where torrents will be mounted as a filesystem. Possible configuration keys described [here](#mountpoints).|

#### mountpoints

|Config key|Description|
|-|-|
|path|Path where a new fuse mount will be initialized. On Windows you can use a drive letter (`X:` per example) or a folder path that **does not exist**.|
|torrents|List of `magnetUri`s or/and `torrentPath`s to be loaded on this fuse mount.|

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the GPL3 license. See `LICENSE` for more information.

[contributors-shield]: https://img.shields.io/github/contributors/distribyted/distribyted.svg?style=flat-square
[contributors-url]: https://github.com/distribyted/distribyted/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/distribyted/distribyted.svg?style=flat-square
[forks-url]: https://github.com/distribyted/distribyted/network/members
[stars-shield]: https://img.shields.io/github/stars/distribyted/distribyted.svg?style=flat-square
[stars-url]: https://github.com/distribyted/distribyted/stargazers
[issues-shield]: https://img.shields.io/github/issues/distribyted/distribyted.svg?style=flat-square
[issues-url]: https://github.com/distribyted/distribyted/issues
[releases-url]: https://github.com/distribyted/distribyted/releases
[license-shield]: https://img.shields.io/github/license/distribyted/distribyted.svg?style=flat-square
[license-url]: https://github.com/distribyted/distribyted/blob/master/LICENSE
[product-screenshot]: docs/images/distribyted_demo.gif
[example-config]: https://github.com/distribyted/distribyted/blob/master/examples/conf_example.yaml
[coveralls-shield]: https://img.shields.io/coveralls/github/distribyted/distribyted?style=flat-square
[coveralls-url]: https://coveralls.io/github/distribyted/distribyted

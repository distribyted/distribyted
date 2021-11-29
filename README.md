[![Releases][releases-shield]][releases-url]
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL3 License][license-shield]][license-url]
[![Coveralls][coveralls-shield]][coveralls-url]
[![Docker Image][docker-pulls-shield]][docker-pulls-url]
<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/distribyted/distribyted">
    <img src="mkdocs/docs/images/distribyted_icon.png" alt="Logo" width="100">
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

## About The Project

![Distribyted Screen Shot][product-screenshot]

Distribyted is an alternative torrent client. 
It can expose torrent files as a standard FUSE, webDAV or HTTP endpoint and download them on demand, allowing random reads using a fixed amount of disk space. 

Distribyted tries to make easier integrations with other applications using torrent files, presenting them as a standard filesystem. 

**Note that distribyted is in beta version, it is a proof of concept with a lot of bugs.**

## Use Cases

- Play **multimedia files** on your favorite video or audio player. These files will be downloaded on demand and only the needed parts.
- Explore TBs of data from public **datasets** only downloading the parts you need. Use **Jupyter Notebooks** directly to process or analyze this data.
- Share your latest dataset creation just sharing a magnet link. People will start access your data in seconds.
- Play your **ROM backups** directly from the torrent file. You can have virtually GBs in games and only downloaded the needed ones.

## Documentation

Check [here][main-url] for further documentation.

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

Some areas need more care than others:
- Windows and macOS tests and compatibility. I don't have any easy way to test distribyted on these operating systems.
- Web interface. Web development is not my _forte_.
- Tutorials. Share with the community your use case!

## License

Distributed under the GPL3 license. See `LICENSE` for more information.

[main-url]: https://distribyted.com
[releases-shield]: https://img.shields.io/github/v/release/distribyted/distribyted.svg?style=flat-square
[releases-url]: https://github.com/distribyted/distribyted/releases
[docker-pulls-shield]:https://img.shields.io/docker/pulls/distribyted/distribyted.svg?style=flat-square
[docker-pulls-url]:https://hub.docker.com/r/distribyted/distribyted
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
[product-screenshot]: mkdocs/docs/images/distribyted.gif
[example-config]: https://github.com/distribyted/distribyted/blob/master/examples/conf_example.yaml
[coveralls-shield]: https://img.shields.io/coveralls/github/distribyted/distribyted?style=flat-square
[coveralls-url]: https://coveralls.io/github/distribyted/distribyted

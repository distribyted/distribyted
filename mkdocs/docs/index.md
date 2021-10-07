# Main

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

### Docker

Docker run example:

```shell
docker run \
  --rm -p 4444:4444 -p 36911:36911 \
  --cap-add SYS_ADMIN \
  --device /dev/fuse \
  --security-opt apparmor:unconfined \
  -v /tmp/mount:/distribyted-data/mount:shared \
  -v /tmp/metadata:/distribyted-data/metadata \
  -v /tmp/config:/distribyted-data/config \
  distribyted/distribyted:latest
```

Docker compose example:

```yaml
distribyted:
    container_name: distribyted
    image: distribyted/distribyted:latest
    restart: always
    ports:
      - "4444:4444/tcp"
      - "36911:36911/tcp"
    volumes:
      - /home/user/mount:/distribyted-data/mount:shared
      - /home/user/metadata:/distribyted-data/metadata
      - /home/user/config:/distribyted-data/config
    security_opt:
      - apparmor:unconfined
    devices:
      - /dev/fuse
    cap_add:
      - SYS_ADMIN
```

### Configuration File

You can see the default configuration file with some explanation comments [here](https://github.com/distribyted/distribyted/blob/master/templates/config_template.yaml).

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the GPL3 license. See `LICENSE` for more information.


[product-screenshot]: images/distribyted.gif

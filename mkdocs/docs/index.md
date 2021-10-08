
Distribyted is an alternative torrent client. 
It can expose torrent files as a standard FUSE mount or webDAV endpoint and download them on demand, allowing random reads using a fixed amount of disk space. 

![Distribyted Screen Shot][product-screenshot]

[product-screenshot]: images/distribyted.gif


## Features

### User Interfaces

Distribyted supports several ways to expose the files to the user or external applications:

#### Supported

- FUSE: Other applications can access to torrent files directly as a filesystem.
- WebDAV: Applications that supports WebDAV can access torrent files using this protocol. It is recommended when distribyted is running in a remote machine or using docker.

#### To be supported
- HTTP: distribyted will support direct HTTP access to files.

### _Expandable_ File Formats
Distribyted can show some kind of files directly as folders, making it possible for applications read only the parts that they need. Here is a list of supported, to be supported and not supported formats.

#### Supported
- zip: Able to uncompress just one file. The file is decompressed to a temporal file sequentially to make possible seek over it. The decompression stops if no one is reading it.

#### To Be Supported
- tar: Seek to any file and inside that files using a [modified standard library](https://github.com/ajnavarro/go-tar). Not useful on `.tar.gz` files.
- 7zip: Similar to Zip. Need for a library similar to [zip](https://github.com/saracen/go7z).
- xz: Only worth it when the file is created using blocks. Possible library [here](https://github.com/ulikunitz/xz) and [here](https://github.com/frrad/bxzf).

#### Not Supported
- gzip: As far as I know, it doesn't support random access.

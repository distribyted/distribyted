
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
- HTTP: A simple HTTP interface for all the available routes. You can acces it from `http://[HOST]:[PORT]/fs`

### _Expandable_ File Formats
Distribyted can show some kind of files directly as folders, making it possible for applications read only the parts that they need. Here is a list of supported, to be supported and not supported formats.

#### Supported
- zip: Able to uncompress just one file. The file is decompressed to a temporal file sequentially to make possible seek over it. The decompression stops if no one is reading it.
- rar: Thanks to [rardecode][rardecode-url] experimental branch library, it is possible to seek through rar files.
- 7zip: Thanks to [sevenzip][sevenzip-url] library, it is possible to read `7z` files in a similar way that is done using the `zip` implementation.

#### To Be Supported
- xz: Only worth it when the file is created using blocks. Possible library [here][xz-url] and [here][bxzf-url].

#### Not Supported
- gzip: As far as I know, it doesn't support random access.

{% include "../_links.md" %}
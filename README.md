<p align="center" width="100%">
    <img width="30%" src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/4df4e58fef369005a7fe65ac04b74b3a667b296b.png"> 
</p>

![build](https://github.com/edbighead/rc/workflows/build/badge.svg)  [![Go Report Card](https://goreportcard.com/badge/github.com/edbighead/rc)](https://goreportcard.com/report/github.com/edbighead/rc)

[Overview](#overview) | [Installation Guide ](#install) | [Usage](#usage) | [Documentation](#docs) | [Demo](#demo)

## <a name="overview"></a>Overview
`rc` is used to delete images from a private docker registry of version 2.3 or later. Images are ordered by creation date in descending order. When calling `rc`, specify image name and number of images you want to keep. Optionally specify `--dry-run` flag to only output the images.

## <a name="usage"></a>Usage
### Authentication
Before running `rc`, you should first configure your credentials details. There are several ways of providing credentials, listed from highest to lowest priority:
1. Config file specified by `--config /path/to/file.yaml` flag
2. Environment variables
```bash
export RC_URL=https://registry.private
export RC_USERNAME=user
export RC_PASSWORD=password
```
3. Default config file `$HOME/.rc.yaml`

Example of config file:
```yaml
---
url: "https://registry.private"
username: "user"
password: "password"
```
### Example
The following command will remove all `alpine` images, keeping only last `3` tags. Images are ordered by creation date in descending order.

`rc cleanup -i alpine -k 3`

### Output
```
2020/11/03 21:16:22 registry.ping url=https://registry.private/v2/
2020/11/03 21:16:22 registry.repositories url=https://registry.private/v2/_catalog
2020/11/03 21:16:22 registry.tags url=https://registry.private/v2/alpine/tags/list repository=alpine
2020/11/03 21:16:22 Some images will be deleted
+---+--------+---------------------+-----+--------+
| № | IMAGE  |       CREATED       | TAG | DELETE |
+---+--------+---------------------+-----+--------+
| 1 | alpine | 3 Nov 2020 19:14:24 | v5  | no     |
| 2 | alpine | 3 Nov 2020 19:03:34 | v4  | no     |
| 3 | alpine | 3 Nov 2020 12:48:33 | v3  | no     |
| 4 | alpine | 3 Nov 2020 11:39:43 | v2  | yes    |
| 5 | alpine | 3 Nov 2020 11:39:10 | v1  | yes    |
+---+--------+---------------------+-----+--------+
2020/11/03 21:16:22 alpine:v2 successfully deleted!
2020/11/03 21:16:23 alpine:v1 successfully deleted!
```

### --dry-run
You can use `--dry-run` flag to list the images, without removing them.

## <a name="install"></a>Installation Guide
####  From the Binary Releases
1. Download your [desired version](https://github.com/edbighead/rc/releases)
2. Unpack it `(tar -zxvf rc.tar.gz)`
3. Find the `rc` binary in the unpacked directory, and move it to its desired destination `(mv rc /usr/local/bin/rc)`


## <a name="docs"></a>Documentation
* [rc](docs/rc.md)
* [rc cleanup](docs/rc_cleanup.md)

## <a name="demo"></a>Demo
[![asciicast](https://asciinema.org/a/370225.svg)](https://asciinema.org/a/370225)

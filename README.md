# OpenSlide Go (unofficial)

OpenSlide Go is a Golang interface to the [OpenSlide] library.
The library derives code from [OpenSlide Python] and another Golang library [gophenslide]. 
Compared to the latter, additional functionality has been implemented.

[OpenSlide]: https://openslide.org/
[OpenSlide Python]: https://github.com/openslide/openslide-python
[gophenslide]: https://github.com/jammy-dodgers/gophenslide


## Requirements

* Go 1.19
* golang.org/x/image v0.1.0
* OpenSlide &ge; 3.4.0


## Installation

OpenSlide Go requires [OpenSlide]. Run `make` in the main directory. This requires that you
have downloaded the test image in `testdata`. To do so run

```shell
mkdir openslide/testdata && cd openslide/testdata
wget https://openslide.cs.cmu.edu/download/openslide-testdata/Generic-TIFF/CMU-1.tiff && cd ../..
make
```

Using brew on OS X it might be required to set the environment flags `CGO_FLAGS` and `CGO_LDFLAGS` to the
appropriate values. For instance
```shell
export CGO_CFLAGS="-I/opt/homebrew/Cellar/openslide/3.4.1_7/include/openslide/ -g -Wall"
export CGO_LDFLAGS="-L/opt/homebrew/Cellar/openslide/3.4.1_7/lib -lopenslide"
```

## More Information

- [Website][OpenSlide]
- [GitHub](https://github.com/NKI-AI/openslide-go)
- [Sample data](https://openslide.cs.cmu.edu/download/openslide-testdata/)


## License

OpenSlide Go is released under the terms of the [GNU Lesser General
Public License, version 2.1](https://openslide.org/license/).

OpenSlide Go is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public
License for more details.
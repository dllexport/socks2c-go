# Socks2c-go

golang implementation of socks2c

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Build

You need only 2 dependencies to build socks2c-go

1. On Linux || MacOS

```
1. go get -v github.com/rakyll/statik
2. compile and install libsodium static lib manually
3. go build (it will link /usr/local/lib/libsodium.a)
```

2. On Windows

```
1. go get -v github.com/rakyll/statik
2. compile libsodium static lib manually
3. crate a directory /lib at the path of the source code
4. copy libsodium.lib into /lib
5. go build
```

## Built With

- [libsodium](https://github.com/jedisct1/libsodium) - A modern, portable, easy to use crypto library
- [statik](https://github.com/rakyll/statik) - embed static files into your Go binary

## Authors

- **Mario Lau** - [Blog](https://dllexport.com)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

- No support for udp over utcp
- No traffic obfuscation data is append for small packet
- For more usage details, check [Release Page](https://code.dllexport.com/mario/socks2c-go-release)

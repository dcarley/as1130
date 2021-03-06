[![Documentation](https://godoc.org/github.com/dcarley/as1130?status.svg)][godoc]
[![CircleCI](https://circleci.com/gh/dcarley/as1130.svg?style=svg)](https://circleci.com/gh/dcarley/as1130)

[godoc]: http://godoc.org/github.com/dcarley/as1130

# AS1130

[Go][] library for the [AS1130][] LED driver, as used by [The Matrix][] from
[Boldport][].

[Go]: https://golang.org/
[AS1130]: http://ams.com/eng/Products/Power-Management/LED-Drivers/AS1130
[The Matrix]: https://www.boldport.com/products/the-matrix/
[Boldport]: http://www.boldport.club/

It can be used from hardware that supports [I²C][] such as a [Raspberry
Pi][].

[I²C]: https://en.wikipedia.org/wiki/I%C2%B2C
[Raspberry Pi]: https://pinout.xyz/pinout/i2c

![as1130 text scrolling](https://user-images.githubusercontent.com/260438/30202937-1f3b391c-9478-11e7-9482-1fd4fc82b096.gif)

## Library

To fetch the library:

```
go get -u github.com/dcarley/as1130
```

Refer to [godoc][] for documentation and examples.

## CLI

The `as1130ctl` command line utility lets you:

- turn on individual LEDs
- scroll text
- read registers
- and more

To fetch and install it:

```
go get -u github.com/dcarley/as1130/as1130ctl
```

To see detailed usage information:

```
as1130ctl --help
```

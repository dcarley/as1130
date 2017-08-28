# embedfont

This utility compresses and embeds a TTF font file into Go code so that it
can be used from a compiled binary.

## Updating

If you modify the font file then you'll need to run `go generate` from the
`as1130` package and commit the modified `.ttf` and `.go` files.

## License

The font is:

- https://fontstruct.com/fontstructions/show/1424354/led-5px-narrow

Which is derived from a CC0 licensed font by vyznev:

- https://fontstruct.com/fontstructions/show/1404171/cg-pixel-4x5

[CC0]: https://creativecommons.org/publicdomain/zero/1.0/

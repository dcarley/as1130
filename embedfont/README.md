# embedfont

This utility compresses and embeds a TTF font file into Go code so that it
can be used from a compiled binary.

## Updating

If you modify the font file then you'll need to run `go generate` from the
`as1130` package and commit the modified `.ttf` and `.go` files.

## License

The font that we're currently using is by vyznev and licensed under Creative
Commons CC0.

- https://fontstruct.com/fontstructions/show/1404171/cg-pixel-4x5

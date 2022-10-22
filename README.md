# Structure
Structure is a library for Dragonfly implementing support for reading and writing Minecraft Bedrock Edition structures.

## Installation
Structure requires at least Go 1.18. The library may be installed using:
```shell
go get github.com/df-mc/structure
```

## Usage
Structures may be read (from a file) using the `structure.Read` and `structure.ReadFile` functions. These structures may
be edited and written afterwards using the `structure.Write` and `structure.WriteFile` functions. Alternatively, a new
structure can be created using `structure.New`.

An example of reading and building a structure in a world:
```go
package main

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/structure"
)

func main() {
	s, err := structure.ReadFile("example.mcstructure")
	if err != nil {
		panic(err)
    }
	
	var w *world.World
	w.BuildStructure(world.BlockPos{}, s)
}
```

## Documentation
[![Go Reference](https://pkg.go.dev/badge/github.com/df-mc/structure.svg)](https://pkg.go.dev/github.com/df-mc/structure)

## Contact
[![Discord Banner 2](https://discordapp.com/api/guilds/623638955262345216/widget.png?style=banner2)](https://discord.gg/U4kFWHhTNR)
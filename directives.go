package structure

import (
	"github.com/df-mc/dragonfly/dragonfly/world"
	_ "unsafe"
)

//go:linkname world_blockByNameAndProperties github.com/df-mc/dragonfly/dragonfly/world.blockByNameAndProperties
// noinspection ALL
func world_blockByNameAndProperties(name string, properties map[string]interface{}) (b world.Block, found bool)

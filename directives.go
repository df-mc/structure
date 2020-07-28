package structure

import (
	"github.com/brentp/intintmap"
	"github.com/df-mc/dragonfly/dragonfly/world"
	_ "unsafe"
)

//go:linkname world_blockByNameAndProperties github.com/df-mc/dragonfly/dragonfly/world.blockByNameAndProperties
//noinspection ALL
func world_blockByNameAndProperties(name string, properties map[string]interface{}) (b world.Block, found bool)

//go:linkname world_runtimeIDsHashes github.com/df-mc/dragonfly/dragonfly/world.runtimeIDsHashes
//noinspection ALL
var world_runtimeIDsHashes *intintmap.Map

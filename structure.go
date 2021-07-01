package structure

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"io"
	"os"
)

// Structure holds the data of an .mcstructure file. Structure implements the world.Structure interface. It
// may be built in a Dragonfly world by using (world.World).BuildStructure.
type Structure struct {
	*structure
}

// Read attempts to read a Structure from the io.Reader passed. If successful, the Structure returned is
// valid and the error is nil.
// Read uses a palette name of 'default' by default. UsePalette may be used to change the name of the
// palette to use.
func Read(r io.Reader) (Structure, error) {
	s := &structure{}
	if err := nbt.NewDecoderWithEncoding(r, nbt.LittleEndian).Decode(s); err != nil {
		return Structure{}, fmt.Errorf("decode structure: %w", err)
	}
	if err := s.check(); err != nil {
		return Structure{}, fmt.Errorf("verify structure: %w", err)
	}
	str := Structure{structure: s}
	str.UsePalette("default")
	return str, nil
}

// ReadFile attempts to read a Structure from a file at the path passed. If successful, the error returned is
// nil.
// ReadFile, like Read, uses a palette name of 'default' by default. UsePalette may be used to change
// the name of the palette to use.
func ReadFile(file string) (Structure, error) {
	f, err := os.Open(file)
	if err != nil {
		return Structure{}, fmt.Errorf("open file: %w", err)
	}
	return Read(f)
}

// Write writes a Structure to the io.Writer passed. If successful, the error returned is nil.
func Write(w io.Writer, s Structure) error {
	s.Structure.Palettes[s.paletteName] = *s.palette

	if err := nbt.NewEncoderWithEncoding(w, nbt.LittleEndian).Encode(s.structure); err != nil {
		return fmt.Errorf("encode structure: %w", err)
	}
	return nil
}

// WriteFile writes a Structure to the file passed. If successful, the error returned is nil. WriteFile
// creates a file if it doesn't yet exist and truncates it if one does exist.
func WriteFile(file string, s Structure) error {
	f, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	return Write(f, s)
}

// New creates a new Structure and initialises it with air blocks. The Structure returned may be written to
// using Structure.Set and Structure.SetAdditionalLiquid and the palette may be changed by using UsePalette.
func New(dimensions [3]int) Structure {
	front := make([]int32, dimensions[0]*dimensions[1]*dimensions[2])
	liquids := make([]int32, dimensions[0]*dimensions[1]*dimensions[2])
	for i := range liquids {
		liquids[i] = -1
	}

	s := Structure{structure: &structure{
		FormatVersion: version,
		Size:          []int32{int32(dimensions[0]), int32(dimensions[1]), int32(dimensions[2])},
		Origin:        []int32{0, 0, 0},
		Structure: structureData{
			BlockIndices: [][]int32{front, liquids},
			Palettes:     map[string]palette{},
		},
	}}
	s.UsePalette("default")
	s.palette.BlockPalette = append(s.palette.BlockPalette, block{
		Name:    "minecraft:air",
		States:  map[string]interface{}{},
		Version: chunk.CurrentBlockVersion,
	})
	return s
}

// UsePalette changes the palette name to use for the Structure. When reading a Structure, this will change
// the palette used to read blocks from. When writing a Structure, the palette will be written with this name,
// so that subsequent readers of the Structure must first call UsePalette with this name to get the right
// palette.
func (s Structure) UsePalette(name string) {
	if current := s.palette; current != nil {
		s.Structure.Palettes[s.paletteName] = *s.palette
	}

	p, _ := s.Structure.Palettes[name]
	if p.BlockPositionData == nil {
		p.BlockPositionData = map[string]blockPositionData{}
	}
	s.palette = &p
	s.paletteName = name
}

// RotateLeft returns a new structure with the same contents but rotated 90 degrees anti-clockwise.
func (s Structure) RotateLeft() Structure {
	return rotate(-1)
}

// RotateRight returns a new structure with the same contents but rotated 90 degrees clockwise.
func (s Structure) RotateRight() Structure {
	return rotate(1)
}

// rotate returns a new structure with the same contents but rotated 90 degrees clockwise the specified
// amount of times.
func (s Structure) rotate(iterations int) Structure {
	iterations %= 4
	if iterations < 0 {
		iterations += 4
	}
	sizeX, sizeY, sizeZ := int(s.Size[0]), int(s.Size[1]), int(s.Size[2])

	var newStructure Structure
	switch iterations {
	case 2:
		newStructure = New([3]int{sizeX, sizeY, sizeZ})
	case 1, 3:
		newStructure = New([3]int{sizeZ, sizeY, sizeX})
	default:
		return s
	}

	maxX, maxZ := sizeX-1, sizeZ-1
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			for z := 0; z < sizeZ; z++ {
				newX, newZ := x, z
				switch iterations {
				case 1:
					newX = -z + maxZ
					newZ = x
				case 2:
					newX = -x + maxX
					newZ = -z + maxZ
				case 3:
					newX = z
					newZ = -x + maxX
				}

				b, l := s.At(x, y, z, nil)
				newStructure.Set(newX, y, newZ, b, l)
			}
		}
	}
	return newStructure
}

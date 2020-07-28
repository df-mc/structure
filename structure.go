package structure

import (
	"bufio"
	"fmt"
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
	structure := Structure{structure: s}
	structure.UsePalette("default")
	return structure, nil
}

// ReadFile attempts to read a Structure from a file at the path passed. If successful, the error returned is
// nil. ReadFile wraps the file in a bufio.Reader and feeds it into the Read function.
// ReadFile, like Read, uses a palette name of 'default' by default. UsePalette may be used to change
// the name of the palette to use.
func ReadFile(file string) (Structure, error) {
	f, err := os.Open(file)
	if err != nil {
		return Structure{}, fmt.Errorf("open file: %w", err)
	}
	return Read(bufio.NewReader(f))
}

// Write writes a Structure to the io.Writer passed. If successful, the error returned is nil.
func Write(w io.Writer, s Structure) error {
	if err := nbt.NewEncoderWithEncoding(w, nbt.LittleEndian).Encode(s.structure); err != nil {
		return fmt.Errorf("encode structure: %w", err)
	}
	return nil
}

// WriteFile writes a Structure to the file passed. If successful, the error returned is nil. WriteFile
// creates a file if it doesn't yet exist and truncates it if one does exist.
// The os.File, when opened, is wrapped in a bufio.Writer.
func WriteFile(file string, s Structure) error {
	f, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	return Write(bufio.NewWriter(f), s)
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

package config

type Config struct {
	Inputs []string
	Output string
	Cut    Cut
	Merge  Merge
	Swap   Swap
}

type Cut struct {
	Page int
	Skip int
}

type Merge struct {
	ByBit   bool
	ByByte  bool
	ByWord  bool
	ByDword bool
}

type Swap struct {
	Bits   bool
	Halfs  bool
	Words  bool
	Dwords bool
}

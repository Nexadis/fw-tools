package config

type Config struct {
	Inputs []string
	Cut    Cut
	Merge  Merge
	Swap   Swap
}

type Cut struct {
	PageSize int
	SkipSize int
}

type Merge struct {
	Output  string
	ByBit   bool
	ByByte  bool
	ByWord  bool
	ByDword bool
}

type Swap struct {
	Bits   bool
	Halfs  bool
	Bytes  bool
	Words  bool
	Dwords bool
}

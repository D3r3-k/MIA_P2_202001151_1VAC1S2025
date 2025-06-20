package Structs

type ParamDef struct {
	Required bool
	Default  string
	NotValue bool
}

type LoginSession struct {
	GID         int32
	UID         int32
	User        string
	Password    string
	PartitionID [4]byte
}

package protos

type CommState = int64

const PageLimit int64 = 1000

const (
	CommStateNil CommState = iota
	CommStateEnabled
	CommStateDisabled
)

package config

// Confs Confs
type Confs []*Value

// NowVer NowVer
var NowVer = &version{
	ver:   1024,
	diffs: []int64{988, 543},
}

// Vers Vers
var Vers = &Confs{
	{CID: 1, Name: "linac.micro.users.timeout", Config: "10"},
	{CID: 2, Name: "linac.micro.users.traffic", Config: "300"},
	{CID: 3, Name: "linac.micro.users.clc", Config: "50"},
}

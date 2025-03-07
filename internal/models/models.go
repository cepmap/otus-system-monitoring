package models

type LoadAverage struct {
	Load1Min  float64
	Load5Min  float64
	Load15Min float64
}
type CPUStat struct {
	User   float64
	System float64
	Idle   float64
}
type DiskLoad struct {
	Tps float64
	Kps float64
}

type NetStats struct {
	Sockets   float64
	ConnCount float64
}

type DiskStats struct {
	DiskStats []DiskStat
}

type DiskStat struct {
	FileSystem string
	Usage      DiskUsage
	Inodes     InodeUsage
}
type DiskUsage struct {
	Used  uint64
	Usage string
}
type InodeUsage struct {
	Used  uint64
	Usage string
}

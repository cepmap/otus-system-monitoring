package models

type LoadAverage struct {
	Load1Min  float64 `protobuf:"fixed64,1,opt,name=load1min,proto3" json:"load1min"`
	Load5Min  float64 `protobuf:"fixed64,2,opt,name=load5min,proto3" json:"load5min"`
	Load15Min float64 `protobuf:"fixed64,3,opt,name=load15min,proto3" json:"load15min"`
}

type CPUStat struct {
	User   float64 `protobuf:"fixed64,1,opt,name=user,proto3" json:"user"`
	System float64 `protobuf:"fixed64,2,opt,name=system,proto3" json:"system"`
	Idle   float64 `protobuf:"fixed64,3,opt,name=idle,proto3" json:"idle"`
}

type DisksLoad struct {
	DisksLoad []DiskLoad `protobuf:"bytes,1,rep,name=disks_load,proto3" json:"disksLoad"`
}

type DiskLoad struct {
	FSName string  `protobuf:"bytes,1,opt,name=fs_name,proto3" json:"fs_name"`
	Tps    float64 `protobuf:"fixed64,1,opt,name=tps,proto3" json:"tps"`
	Kps    float64 `protobuf:"fixed64,2,opt,name=kps,proto3" json:"kps"`
}

type DiskStats struct {
	DiskStats []DiskStat `protobuf:"bytes,1,rep,name=disk_stats,proto3" json:"diskStats"`
}

type DiskStat struct {
	FileSystem string     `protobuf:"bytes,1,opt,name=filesystem,proto3" json:"filesystem"`
	Usage      DiskUsage  `protobuf:"bytes,2,opt,name=usage,proto3" json:"usage"`
	Inodes     InodeUsage `protobuf:"bytes,3,opt,name=inodes,proto3" json:"inodes"`
}

type DiskUsage struct {
	Used  uint64 `protobuf:"varint,1,opt,name=used,proto3" json:"used"`
	Usage string `protobuf:"bytes,2,opt,name=usage,proto3" json:"usage"`
}

type InodeUsage struct {
	Used  uint64 `protobuf:"varint,1,opt,name=used,proto3" json:"used"`
	Usage string `protobuf:"bytes,2,opt,name=usage,proto3" json:"usage"`
}

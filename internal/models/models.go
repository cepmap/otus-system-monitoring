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
type DiskInfo struct {
	Tps float64
	Kps float64
}
type DiskLoad struct {
	MUsed float64
	IUsed float64
}
type NetStats struct {
	TUSockets float64
	TCount    float64
}
type TopTalkers struct {
	TUSockets    string
	TConnections string
}

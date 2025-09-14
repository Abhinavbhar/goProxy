package helperfunc

import "sync"

var (
	IpBandwidth = make(map[string]int64)
	IpMutex     = sync.RWMutex{}
)

package main

func CheckIp(ip string) bool {
	_, ok := ipBandwidth[ip]
	if ok {
		return true
	}
	//call backend using grpc to see if the ip exist in the database
	return false
}

package dnsdb

func checkRDATAFormat(format string) bool {
	switch format {
	case "name", "ip", "raw":
		return true
	default:
		return false
	}
}

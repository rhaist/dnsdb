package dnsdb

// checkRDATAFormat is used for validation of the format string parameter to
// only contain one of the three valid strings.
func checkRDATAFormat(format string) bool {
	switch format {
	case "name", "ip", "raw":
		return true
	default:
		return false
	}
}

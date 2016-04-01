package dnsdb

import (
	"fmt"
	"strconv"
	"time"
)

// timestamp is used to marshal and unmarshal unix
// timestamps returned by the API
type timestamp time.Time

func (t *timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = timestamp(time.Unix(int64(ts), 0))

	return nil
}

func (t *timestamp) String() string {
	if t != nil {
		return time.Time(*t).UTC().String()
	}
	return ""
}

// RRSET contains the result set of a rrset API query
// Endpoint: /lookup/rrset
type RRSET struct {
	Count         int        `json:"count"`
	TimeFirst     *timestamp `json:"time_first,omitempty"`
	TimeLast      *timestamp `json:"time_last,omitempty"`
	ZoneTimeFirst *timestamp `json:"zone_time_first,omitempty"`
	ZoneTimeLast  *timestamp `json:"zone_time_last,omitempty"`
	Rrname        string     `json:"rrname"`
	Rrtype        string     `json:"rrtype"`
	Rdata         []string   `json:"rdata"`
	Bailiwick     string     `json:"bailiwick"`
}

// RDATA contains the result set of a rdata API query
// Endpoint: /lookup/rdata
type RDATA struct {
	Count         int        `json:"count"`
	TimeFirst     *timestamp `json:"time_first,omitempty"`
	TimeLast      *timestamp `json:"time_last,omitempty"`
	ZoneTimeFirst *timestamp `json:"zone_time_first,omitempty"`
	ZoneTimeLast  *timestamp `json:"zone_time_last,omitempty"`
	Rrname        string     `json:"rrname"`
	Rrtype        string     `json:"rrtype"`
	Rdata         string     `json:"rdata"`
}

// RateLimit contains the current rate limit information for the
// used API key.
// Endpoint: /lookup/rate_limit
type RateLimit struct {
	Rate struct {
		Reset     *timestamp `json:"reset"`
		Limit     int        `json:"limit"`
		Remaining int        `json:"remaining"`
	} `json:"rate"`
}

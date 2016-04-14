package dnsdb

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const apiFormat = "application/json"

var (
	APIKEY        string // APIKEY is required to access the service
	SERVER        string // SERVER must be set to the DNSDB API server
	RateLimit     int
	RateRemaining int
	client        *http.Client
)

// Initialize and configure our http.Client to only use TLS 1.2 crypto for
// securely encrypted transport.
func init() {
	t := &tls.Config{}
	t.PreferServerCipherSuites = true
	t.MinVersion = tls.VersionTLS12
	t.MaxVersion = tls.VersionTLS12

	tr := &http.Transport{
		TLSClientConfig: t,
	}

	client = &http.Client{Transport: tr}
}

// baseAPICall is used by the query functions to perform an API request and get
// the response data. This func handles the authentication and format of the
// DNSDB API.
func baseAPICall(url string) (*http.Response, error) {
	// check if we hit the rate limit yet. If RateLimit is 0
	// we asume this is the first request and let it pass
	if !checkRateLimit() && RateLimit != 0 {
		return nil, errors.New("DNSDB API quota limit reached")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", apiFormat)
	req.Header.Add("X-API-Key", APIKEY)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = updateRateLimit(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// updateRateLimit is a convenient function to parse the HTTP header of every
// DNSDB API Response for the status of rate limits for the current API Key. We
// use this to throw errors when the limit is reached. This passive approach is
// more elegant than doing a RateLimitQuery after each API query.
func updateRateLimit(resp *http.Response) error {
	var err error
	rl := resp.Header.Get("X-RateLimit-Limit")
	rr := resp.Header.Get("X-RateLimit-Remaining")

	if rl == "unlimited" {
		RateLimit = -1
		RateRemaining = -1
	} else {
		RateLimit, err = strconv.Atoi(rl)
		if err != nil {
			return err
		}
		RateRemaining, err = strconv.Atoi(rr)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkRateLimit is a convenient function to check if the API allows us at
// least one more query. If we don't have information about the rate limit we
// assume this is the first query and use RateLimitQuery once.
func checkRateLimit() bool {
	if RateRemaining > 0 {
		return true
	}
	return false
}

// RateLimitQuery returns the rate limits for the currently used DNSDB APIKEY.
func RateLimitQuery() (RateLimitInfo, error) {
	var limits RateLimitInfo
	url := SERVER + "/lookup/rate_limit/"

	resp, err := baseAPICall(url)
	if err != nil {
		return limits, err
	}

	err = json.NewDecoder(resp.Body).Decode(&limits)
	if err != nil {
		return limits, err
	}
	return limits, nil
}

// RRSETQuery takes a query string to search for rrset records in the DNSDB and
// returns a RRSETArray struct.
func RRSETQuery(query string) ([]RRSET, error) {
	var rrset RRSET
	var rrsetArray []RRSET
	url := SERVER + "/lookup/rrset/name/" + query

	resp, err := baseAPICall(url)
	if err != nil {
		return rrsetArray, err
	}

	// The API does not return a valid JSON array at
	// the moment so we need to marshal line by line
	// and put the objects into our rrsetArray slice.
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		err = json.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&rrset)
		if err != nil {
			return rrsetArray, err
		}
		rrsetArray = append(rrsetArray, rrset)
	}

	if err := scanner.Err(); err != nil {
		return rrsetArray, err
	}
	return rrsetArray, nil
}

// RDATAQuery takes a query string to search for rdata records in the DNSDB and
// returns a RDATA struct. Allowed format strings are "ip", "name" and "raw".
func RDATAQuery(query string, format string) ([]RDATA, error) {
	if !checkRDATAFormat(format) {
		return nil, errors.New("Wrong rdata format - allowed are (name|ip|raw)")
	}

	var rdata RDATA
	var rdataArray []RDATA
	url := SERVER + "/lookup/rdata/" + format + "/" + query

	fmt.Println(url)

	resp, err := baseAPICall(url)
	if err != nil {
		return rdataArray, err
	}

	// The API does not return a valid JSON array at
	// the moment so we need to marshal line by line
	// and put the objects into our rdataArray slice.
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		err = json.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&rdata)
		if err != nil {
			return rdataArray, err
		}
		rdataArray = append(rdataArray, rdata)
	}

	if err := scanner.Err(); err != nil {
		return rdataArray, err
	}
	return rdataArray, nil
}

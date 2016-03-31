package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/SleuthKid/dnsdb"
	"github.com/voxelbrain/goptions"
)

func main() {
	options := struct {
		Server    string        `goptions:"-s, --server, description='DNSDB API servers to connect to'"`
		Config    string        `goptions:"-c, --config, description='Path to config file'"`
		RateLimit bool          `goptions:"-r, --ratelimit, description='Print current rate limit data'"`
		Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
		RRSET struct {
			Query string `goptions:"-q, --query, obligatory, description='Query string'"`
		} `goptions:"rrset"`
		RDATA struct {
			Query  string `goptions:"-q, --query, obligatory, description='Query string'"`
			Format string `goptions:"-f, --format, description='Specify rdata format (name|ip|raw)'"`
		} `goptions:"rdata"`
	}{ // Default values
		Server: "https://api.dnsdb.info",
		Config: getDefaultConfPath(),
	}
	goptions.ParseAndFail(&options)

	err := loadConfig(options.Config)
	if err != nil {
		panic(err)
	}

	// Config the dnsdb package
	dnsdb.APIKEY = CONF.APIKEY
	dnsdb.SERVER = options.Server

	if options.RateLimit {
		rl, err := dnsdb.RateLimitQuery()
		if err != nil {
			panic(err)
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "Limit:\t"+strconv.Itoa(rl.Rate.Limit))
		fmt.Fprintln(w, "Remaining:\t"+strconv.Itoa(rl.Rate.Remaining))
		fmt.Fprintln(w, "Reset:\t"+rl.Rate.Reset.String())
		w.Flush()
	}

	if options.Verbs == "rrset" {
		rs, err := dnsdb.RRSETQuery(options.RRSET.Query)
		if err != nil {
			fmt.Println("[Error] rrset query exited with the following error:", err)
			os.Exit(1)

		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		for _, v := range rs {
			fmt.Fprintln(w, strconv.Itoa(v.Count)+"\t"+
				v.TimeFirst.String()+"\t"+
				v.TimeLast.String()+"\t"+
				v.ZoneTimeFirst.String()+"\t"+
				v.ZoneTimeLast.String()+"\t"+
				v.Bailiwick+"\t"+
				v.Rrtype+"\t"+
				v.Rrname+"\t"+
				fmt.Sprintf("%v", v.Rdata))
		}
		w.Flush()
	}

	if options.Verbs == "rdata" {
		rs, err := dnsdb.RDATAQuery(strings.TrimSpace(options.RDATA.Query), options.RDATA.Format)
		if err != nil {
			fmt.Println("[Error] rdata query exited with the following error:", err)
			os.Exit(1)
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		for _, v := range rs {
			fmt.Fprintln(w, strconv.Itoa(v.Count)+"\t"+
				v.TimeFirst.String()+"\t"+
				v.TimeLast.String()+"\t"+
				v.ZoneTimeFirst.String()+"\t"+
				v.ZoneTimeLast.String()+"\t"+
				v.Rrtype+"\t"+
				v.Rrname+"\t"+
				v.Rdata)
		}
		w.Flush()
	}
}

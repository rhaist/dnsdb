# dnsdb

[![GoDoc](https://godoc.org/github.com/rhaist/dnsdb?status.svg)](https://godoc.org/github.com/rhaist/dnsdb)


dnsdb is a go API wrapper for the [DNSDB API](https://api.dnsdb.info/) provided by [Farsight Security, Inc.](https://www.farsightsecurity.com/)

Some of the functionality is based on information from the [official manual](https://api.dnsdb.info/) and the example clients available [here](https://github.com/dnsdb/dnsdb-query)

**Please note:** access to the DNSDB is not freely available. You will need to get a valid service subscription and API key first from Farsight. Information about the subscription is
available [here](https://www.farsightsecurity.com/OrderServices/).

## Install

to install the library with the example command line client do:

    go get -u github.com/rhaist/dnsdb/...

to only install the library **without** the command line client do:

    go get -u github.com/rhaist/dnsdb

## Usage

Please refer to the godocs available [here](https://godoc.org/github.com/rhaist/dnsdb)

The command line client allows the following Parameters:
```
Usage: dnsdb [global options] <verb> [verb options]

Global options:
        -s, --server    DNSDB API servers to connect to (default: https://api.dnsdb.info)
        -c, --config    Path to config file (default: /home/user/.dnsdb-query.conf)
        -r, --ratelimit Print current rate limit data
        -h, --help      Show this help

Verbs:
    rdata:
        -q, --query     Query string (*)
        -f, --format    Specify rdata format (name|ip|raw)
    rrset:
        -q, --query     Query string (*)
```

An example query with the command line client for all rrsets of google.com would look like this:

    dnsdb rrset -q google.com

## Tests

Providing test cases / code for this package would allways require having a valid
API key and consumption of a considerable amount of API quota. This is why the
author has written the example command line client to allow more granular testing
of this packages functions.

Testing and feedback from other users of this package is much appreciated.

## License
This software is distributed under the MIT license.
Please have a look at the LICENSE file in the source distribution.

## Author

Robert Haist / rhaist at mailbox dot org

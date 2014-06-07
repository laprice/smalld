NOTE: certain forward looking statements in this README do not yet hold true
This is a WIP.

This is a micro-service. It accepts GET and POST requests that contain
a geolocation and return a data structure expressing labels for the
the location queried.

Additionally it provides a feed of queries and responses that can be displayed.

The two functions are presented as http endpoints:

    `location/` accepts GET or POST with variables (lat, lon, acc, label)

    `feed/` GET and returns a feed of Server Sent Events

`smalld` requires a postgres database with the postgis extensions installed.

It aquires its configuration from the following environment variables:

    `SMALLD_DB_CONNECTION` :  a database connection string as specified in
    http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
    
    `SMALLD_URL_BASE` : a string containing the prefix you wish to use
    for generating urls in responses

    `SMALLD_OPTIONS` : a string containing options for the daemon
    see options below.


options:
	`--no-feed` turn off the feed portion, do not respond to requests on
	the `feed/` endpoint, do not push updates. 

	`--fore` foreground and log to stderr 

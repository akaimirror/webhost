webhost command runs simple http server mapping local directories to different
domain names and subpaths.

	Usage of webhost:
	  -addr string
		address to listen at (default "localhost:8080")
	  -map string
		file with host/backend mapping (default "mapping.yml")
	  -rto duration
		maximum duration before timing out read of the request (default 10s)
	  -wto duration
		maximum duration before timing out write of the response (default 5m0s)

Format of mapping file:

	hostname1/path/: /local/filesystem/path1
	hostname2: /local/filesystem/path2
	/path/: /local/filesystem/path3

Note how path has trailing slash — if it's omitted, directory would only be
available under exact `/path`, but not `/path/`. If mapping is configured with
`/path/`, then both `/path` and `/path/` would work, with former doing redirect
to latter.

# fmp-json
A simple Filemaker XML to JSON proxy

	Usage of ./fmp-json:
	  -debug
			Debug requests
	  -http string
	    	HTTP port and interface the server will use, format interface:port
	  -redirect-to string
	    	When using TLS, redirect all request using http to this port
	  -server string
	    	The filemaker server to use as host
	  -ssl-cert string
	    	Path to the ssl cert to use, if empty it will use http
	  -ssl-key string
	    	Path to the ssl key to use, if empty it will use http
	  -tls string
	    	TLS port and interface the server will use, format interface:port
	  -usesyslog
	    	Use syslog

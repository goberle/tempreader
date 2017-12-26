tempreader
==========

RESTful API to read data from DS18B20 temperature sensors.

[![Build Status](https://travis-ci.org/jtyr/tempreader.svg?branch=master)](https://travis-ci.org/jtyr/tempreader)


Usage
-----

Run it locally:

```
$ go run tempreader.go
```

Configuration:

```
$ ./tempreader -help
Usage of ./tempreader:
  -addr string
        port or ip:port for the HTTP server (default "0.0.0.0:8000")
  -debug
        show more verbose output
  -help
        show this help message and exit
  -root string
        root directory with sensors (default "/sys/bus/w1/devices")
  -version
        show version and exit
```

Configurable via ENV variables:

```
TEMPREADER_ADDR  - port or ip:port for the HTTP server (default: 0.0.0.0:8000)
TEMPREADER_DEBUG - show more verbose output (default: false)
TEMPREADER_ROOT  - root directory with sensors (default: /sys/bus/w1/devices)
```

Testing:

```
# Create testing data
$ mkdir -p /tmp/tempreader/28-{000005e2fdc3,03168bf4edff}
cat > /tmp/tempreader/28-000005e2fdc3/w1_slave <<END
8c 01 4b 46 7f ff 0c 10 58 : crc=58 YES
8c 01 4b 46 7f ff 0c 10 58 t=24750
END
$ cat > /tmp/tempreader/28-03168bf4edff/w1_slave <<END
8c 01 4b 46 7f ff 0c 10 58 : crc=58 YES
8c 01 4b 46 7f ff 0c 10 58 t=25335
END

# Point the script to the testing data
$ TEMPREADER_ROOT=/tmp/tempreader go run tempreader.go

# Query the sernsors
$ curl http://localhost:8000/tempreader/api/v1.0/sensors
{"sensors":["28-000005e2fdc3","28-03168bf4edff"]}
$ curl http://localhost:8000/tempreader/api/v1.0/sensors/28-03168bf4edff
{"crc":true,"temp":25.335}
```

Croscompile for Raspberry Pi:

```
# Dynamically:
$ GOOS=linux GOARCH=arm go build
# Statically:
$ GOOS=linux GOARCH=arm CGO_ENABLED=0 go build
```


License
-------

MIT


Author
------

Jiri Tyr

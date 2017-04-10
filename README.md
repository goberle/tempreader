tempreader
==========

RESTful API to read data from DS18B20 temperature sensors.


Usage
-----

Run it locally:

```
cd tempreader
go run main.go
```

Configurable via ENV variables:

```
TEMPREADER_ROOT  - root directory with sensors (default: /sys/bus/w1/devices)
TEMPREADER_DEBUG - show more verbose output (default: false)
TEMPREADER_ADDR  - port or ip:port for the HTTP server (default: 0.0.0.0:8000)
```

Testing:

```
# Create testing data
mkdir -p /tmp/tempreader/28-{000005e2fdc3,03168bf4edff}
cat > /tmp/tempreader/28-000005e2fdc3/w1_slave <<END
8c 01 4b 46 7f ff 0c 10 58 : crc=58 YES
8c 01 4b 46 7f ff 0c 10 58 t=24750
END
cat > /tmp/tempreader/28-03168bf4edff/w1_slave <<END
8c 01 4b 46 7f ff 0c 10 58 : crc=58 YES
8c 01 4b 46 7f ff 0c 10 58 t=25335
END

# Point the script to the testing data
cd tempreader
TEMPREADER_ROOT=/tmp/tempreader go run main.go

# Query the sernsors
curl http://localhost:8000/tempreader/api/v1.0/sensors
curl http://localhost:8000/tempreader/api/v1.0/sensors/28-03168bf4edff
```

Croscompile for Raspberry Pi:

```
# Dynamically:
GOOS=linux GOARCH=arm go build
# Statically:
GOOS=linux GOARCH=arm CGO_ENABLED=0 go build
```


License
-------

MIT


Author
------

Jiri Tyr

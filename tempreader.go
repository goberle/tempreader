package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/jtyr/tempreader/ds18b20"
    "github.com/jtyr/tempreader/utils"
)

// SensorHandler inherits the ds18b20.Sensor.
type SensorHandler struct {
    ds18b20.Sensor
    debug bool
}

const version = "1.0.0"

// SensorInfoHandler provides HTTP handler for information about a sensor.
func (sh *SensorHandler) SensorInfoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    if sh.debug {
        log.Println("I: SensorInfo (" + vars["id"] + ") from " + r.RemoteAddr)
    }

    w.Header().Set("Content-Type", "application/json")

    // Read the sensor info
    ri, err := sh.GetSensorInfo(vars["id"])
    if err != nil {
        log.Printf("E: %s\n", err)
        w.Write([]byte("{}"))
        return
    }

    j, _ := json.Marshal(ri)

    w.Write([]byte(j))
}

// SensorsListHandler provides HTTP handler for the list of available sensors.
func (sh *SensorHandler) SensorsListHandler(w http.ResponseWriter, r *http.Request) {
    if sh.debug {
        log.Println("I: SensorList from " + r.RemoteAddr)
    }

    w.Header().Set("Content-Type", "application/json")

    // Read the sensor IDs
    rl, err := sh.GetSensorList()
    if err != nil {
        log.Printf("E: %s\n", err)
        w.Write([]byte("{\"sensors\": []}"))
        return
    }

    j, _ := json.Marshal(rl)

    w.Write([]byte(j))
}

// For the test
var run = flag.PrintDefaults

func main() {
    var argsDebug, argsHelp, argsVersion bool
    var argsAddr, argsRoot string

    flag.StringVar(
        &argsRoot, "root",
        utils.GetEnv("TEMPREADER_ROOT", "/sys/bus/w1/devices"),
        "root directory with sensors")
    flag.StringVar(
        &argsAddr, "addr",
        utils.GetEnv("TEMPREADER_ADDR", "0.0.0.0:8000"),
        "port or ip:port for the HTTP server")
    flag.BoolVar(
        &argsDebug, "debug",
        utils.GetEnvBool("TEMPREADER_DEBUG", false),
        "show more verbose output")
    flag.BoolVar(
        &argsHelp, "help",
        false,
        "show this help message and exit")
    flag.BoolVar(
        &argsVersion, "version",
        false,
        "show version and exit")
    flag.Parse()

    if argsHelp {
        fmt.Printf("Usage of %s:\n", os.Args[0])
        run()
        os.Exit(0)
    } else if argsVersion {
        fmt.Printf("tempreader v%s\n", version)
        os.Exit(0)
    }

    sh := SensorHandler{}
    sh.Root = argsRoot
    sh.debug = argsDebug

    // Internal variable
    prefix := "/tempreader/api/v1.0"

    r := mux.NewRouter()
    s := r.PathPrefix(prefix).Subrouter()

    // Routes
    s.HandleFunc("/{sensors:sensors(?:\\/)?}", sh.SensorsListHandler).Methods("GET")
    s.HandleFunc("/sensors/{id:[0-9a-f-]+(?:\\/)?}", sh.SensorInfoHandler).Methods("GET")

    if argsDebug {
        log.Println("I: Running server on http://" + argsAddr)
    }

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(argsAddr, r))
}

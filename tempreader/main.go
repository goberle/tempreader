package main

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gorilla/mux"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
)


type ResponseList struct {
    Sensors []string `json:"sensors"`
}

type ResponseInfo struct {
    Crc bool     `json:"crc"`
    Temp float64 `json:"temp"`
}

var root string
var debug bool


func ReadSensorInfo(id string) (ResponseInfo, error) {
    if debug {
        log.Println("I: Reading Sensor data...")
    }

    f, err := os.Open(fmt.Sprintf("%s/%s/w1_slave", root, id))
    if err != nil {
        return ResponseInfo{}, err
    }
    defer f.Close()

    crc := false
    temp := 0.0

    i := 0
    r := bufio.NewReader(f)
    line, err := r.ReadString(10)

    for err != io.EOF {
        if i == 0 {
            if line[36:len(line)-1] == "YES" {
                crc = true
            }
        } else if i == 1 {
            t, err := strconv.ParseFloat(line[29:len(line)-1], 64)
            if err != nil {
                return ResponseInfo{}, err
            }

            temp = t/1000
        }

        i += 1
        line, err = r.ReadString(10)
    }

    return ResponseInfo{crc, temp}, nil
}

func ReadSensorIDs() (ResponseList, error) {
    if debug {
        log.Println("I: Reading Sensor IDs...")
    }

    deviceDirInfo, err := os.Stat(root)
    if err != nil {
        return ResponseList{}, err
    }
    if ! deviceDirInfo.IsDir() {
        return ResponseList{}, errors.New("Not directory.")
    }

    files, err := ioutil.ReadDir(root)
    if err != nil {
        return ResponseList{}, err
    }

    var lst []string

    for _, file := range files {
        fileInfo, err := os.Stat(root + "/" + file.Name())
        if err != nil {
            return ResponseList{}, err
        }

        if fileInfo.IsDir() {
            // Support only sensors starting with 28
            if file.Name()[:2] == "28" {
                lst = append(lst, file.Name())
            }
        }
    }

    return ResponseList{lst}, nil
}


func SensorInfoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    // Read the sensor info
    ri, err := ReadSensorInfo(vars["id"])
    if err != nil {
        log.Printf("E: %s\n", err)
        w.Write([]byte("{}"))
        return
    }

    j, _ := json.Marshal(ri)

    w.Write([]byte(j))
}

func SensorsListHandler(w http.ResponseWriter, r *http.Request) {
    // Read the sensor IDs
    rl, err := ReadSensorIDs()
    if err != nil {
        log.Printf("E: %s\n", err)
        w.Write([]byte("{\"sensors\": []}"))
        return
    }

    j, _ := json.Marshal(rl)

    w.Write([]byte(j))
}


func GetEnvVar(key string, d string) (r string) {
    val, ok := os.LookupEnv(key)
    if !ok {
        r = d
    } else {
        r = val
    }

    return r
}


func main() {
    // Configurable variables
    root = GetEnvVar("TEMPREADER_ROOT", "/sys/bus/w1/devices")
    address := GetEnvVar("TEMPREADER_ADDR", "0.0.0.0:8000")
    d, err := strconv.ParseBool(GetEnvVar("TEMPREADER_DEBUG", "false"))
    if err != nil {
        log.Fatal("Cannot convert TEMPREADER_DEBUG to bool.")
    }
    debug = d

    // Internal variable
    prefix := "/tempreader/api/v1.0"

    r := mux.NewRouter()
    s := r.PathPrefix(prefix).Subrouter()

    // Routes
    s.HandleFunc("/sensors", SensorsListHandler).Methods("GET")
    s.HandleFunc("/sensors/{id:[0-9a-f-]+}", SensorInfoHandler).Methods("GET")

    if debug {
        log.Println("I: Running server on http://" + address)
    }

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(address, r))
}

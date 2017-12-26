package ds18b20

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "strconv"
)

// SensorList holds information about the list of sensors.
type SensorList struct {
    Sensors []string `json:"sensors"`
}

// SensorInfo holds information about one particular sensor.
type SensorInfo struct {
    Crc  bool    `json:"crc"`
    Temp float64 `json:"temp"`
}

// Sensor is a type defining parameters required to obtain sensor information.
type Sensor struct {
    Root string
}

// GetSensorList returns a list of available sensors.
func (s *Sensor) GetSensorList() (SensorList, error) {
    deviceDirInfo, err := os.Stat(s.Root)
    if err != nil {
        return SensorList{}, err
    }
    if ! deviceDirInfo.IsDir() {
        return SensorList{}, errors.New("not a directory")
    }

    files, err := ioutil.ReadDir(s.Root)
    if err != nil {
        return SensorList{}, err
    }

    var lst []string

    for _, file := range files {
        fileInfo, err := os.Stat(s.Root + "/" + file.Name())
        if err != nil {
            return SensorList{}, err
        }

        if fileInfo.IsDir() {
            // Support only sensors starting with 28
            if file.Name()[:2] == "28" {
                lst = append(lst, file.Name())
            }
        }
    }

    return SensorList{lst}, nil
}

// GetSensorInfo reads the sensor file.
func (s *Sensor) GetSensorInfo(id string) (SensorInfo, error) {
    f, err := os.Open(fmt.Sprintf("%s/%s/w1_slave", s.Root, id))
    if err != nil {
        return SensorInfo{}, err
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
                return SensorInfo{}, err
            }

            temp = t/1000
        }

        i++
        line, err = r.ReadString(10)
    }

    return SensorInfo{crc, temp}, nil
}

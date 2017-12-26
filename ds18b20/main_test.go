package ds18b20

import (
    "io/ioutil"
    "os"
    "testing"

    "github.com/Flaque/filet"
)

func TestGetSensorList(t *testing.T) {
    defer filet.CleanUp(t)

    tmp := filet.TmpDir(t, "")

    var lst []string
    lst = append(lst, "28-000005e2fdc3")
    lst = append(lst, "28-03168bf4edff")

    for i := range lst {
        os.Mkdir(tmp + "/" + lst[i], 0755)
    }

    s := Sensor{}
    s.Root = tmp

    var sl SensorList
    sl, _ = s.GetSensorList()

    for i, item := range sl.Sensors {
        if item != lst[i] {
            t.Errorf("Expected %s, found %s.", lst[i], item)
        }
    }
}

func TestGetSensorInfo(t *testing.T) {
    defer filet.CleanUp(t)

    tmp := filet.TmpDir(t, "")
    id  := "28-000005e2fdc3"
    dir := tmp + "/" + id

    os.Mkdir(dir, 0755)

    ioutil.WriteFile(
        dir + "/w1_slave",
        []byte("8c 01 4b 46 7f ff 0c 10 58 : crc=58 YES\n8c 01 4b 46 7f ff 0c 10 58 t=24750\n"),
        0644)

    s := Sensor{}
    s.Root = tmp

    var si SensorInfo
    si, _ = s.GetSensorInfo(id)

    if si.Crc != true || si.Temp != 24.750 {
        t.Errorf("Expected (%t, %f), found (%t, %f).", true, 24.750, si.Crc, si.Temp)
    }
}

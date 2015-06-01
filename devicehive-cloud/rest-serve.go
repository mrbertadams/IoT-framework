package main

import (
	"encoding/json"

	"github.com/devicehive/IoT-framework/devicehive-cloud/conf"
	"github.com/devicehive/IoT-framework/devicehive-cloud/rest"
	"github.com/devicehive/IoT-framework/devicehive-cloud/say"
	"github.com/godbus/dbus"
)

const (
	restObjectPath  = "/com/devicehive/cloud"
	restCommandName = "com.devicehive.cloud.CommandReceived"
)

type DbusRestWrapper struct {
	URL        string
	AccessKey  string
	DeviceID   string
	DeviceName string
}

func NewDbusRestWrapper(c conf.Conf) *DbusRestWrapper {
	w := new(DbusRestWrapper)
	w.URL = c.URL
	w.AccessKey = c.AccessKey
	w.DeviceID = c.DeviceID
	w.DeviceName = c.DeviceName
	return w
}

func (w *DbusRestWrapper) SendNotification(name, parameters string, priority uint64) *dbus.Error {
	say.Verbosef("SendNotification(name='%s',params='%s',priority=%d)\n", name, parameters, priority)
	dat, err := parseJSON(parameters)

	if err != nil {
		return newDHError(err.Error())
	}

	rest.DeviceNotificationInsert(w.URL, w.DeviceID, w.AccessKey, name, dat)

	return nil
}

func (w *DbusRestWrapper) UpdateCommand(id uint32, status string, result string) *dbus.Error {

	dat, err := parseJSON(result)

	if err != nil {
		return newDHError(err.Error())
	}

	m := map[string]interface{}{
		"action":    "command/update",
		"commandId": id,
		"command": map[string]interface{}{
			"status": status,
			"result": dat,
		},
	}

	//TODO: update to real method
	rest.DeviceCmdInsert(w.URL, w.DeviceID, w.AccessKey, "SendCommand", m)

	return nil
}

func restImplementation(bus *dbus.Conn, config conf.Conf) {

	// var info rest.ApiInfo

	// for {
	// 	var err error
	// 	info, err = rest.GetApiInfo(config.URL)
	// 	if err == nil {
	// 		say.Verbosef("API info: %+v", info)
	// 		break
	// 	}
	// 	say.Infof("API info error: %s", err.Error())
	// 	time.Sleep(5 * time.Second)
	// }

	go func() {
		control := rest.NewPollAsync()
		//out := make(chan rest.DeviceCmdResource, 16)
		out := make(chan rest.DeviceNotificationResource, 16)

		go rest.DeviceNotificationPollAsync(config.URL, config.DeviceID, config.AccessKey, out, control)

		for {
			select {
			case n := <-out:
				parameters := ""
				if n.Parameters != nil {
					b, err := json.Marshal(n.Parameters)
					if err != nil {
						say.Infof("Could not generate JSON from parameters of command %+v\nWith error %s", n, err.Error())
						continue
					}

					parameters = string(b)
				}
				say.Verbosef("NOTIFICATION %s -> %s(%v)", config.URL, n.Notification, parameters)
				bus.Emit(restObjectPath, restCommandName, uint32(n.Id), n.Notification, parameters)
			}
		}
	}()

	w := NewDbusRestWrapper(config)
	bus.Export(w, "/com/devicehive/cloud", DBusConnName)

	select {}
}

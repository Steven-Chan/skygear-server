// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"testing"

	"github.com/skygeario/skygear-server/pkg/server/handler/handlertest"
	. "github.com/skygeario/skygear-server/pkg/server/skytest"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/skygeario/skygear-server/pkg/server/push"
	"github.com/skygeario/skygear-server/pkg/server/router"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
)

func TestPushToDevice(t *testing.T) {
	Convey("push to device", t, func() {
		testdevice := skydb.Device{
			ID:         "device",
			Type:       "ios",
			Token:      "token",
			AuthInfoID: "userid",
		}
		conn := simpleDeviceConn{
			devices: []skydb.Device{testdevice},
		}

		r := handlertest.NewSingleRouteRouter(&PushToDeviceHandler{}, func(p *router.Payload) {
			p.DBConn = &conn
		})

		originalSendFunc := sendPushNotification
		defer func() {
			sendPushNotification = originalSendFunc
		}()

		Convey("push to single device", func(c C) {
			called := false
			sendPushNotification = func(sender push.Sender, device skydb.Device, m push.Mapper) {
				c.So(device, ShouldResemble, testdevice)
				c.So(m.Map(), ShouldResemble, map[string]interface{}{
					"aps": map[string]interface{}{
						"alert": "This is a message.",
						"sound": "sosumi.mp3",
					},
					"acme": "interesting",
				})
				called = true
			}

			resp := r.POST(`{
					"device_ids": ["device"],
					"notification": {
						"aps": {
							"alert": "This is a message.",
							"sound": "sosumi.mp3"
						},
						"acme": "interesting"
					}
				}`)
			So(resp.Code, ShouldEqual, 200)
			So(resp.Body.Bytes(), ShouldEqualJSON, `{
	"result": [{
		"_id": "device"
	}]
}`)
			So(called, ShouldBeTrue)
		})

		Convey("push to non-existent device", func() {
			called := false
			sendPushNotification = func(sender push.Sender, device skydb.Device, m push.Mapper) {
				called = true
			}
			resp := r.POST(`{
						"device_ids": ["nonexistent"],
						"notification": {
							"aps": {
								"alert": "This is a message.",
								"sound": "sosumi.mp3"
							},
							"acme": "interesting"
						}
					}`)
			So(resp.Code, ShouldEqual, 200)
			So(resp.Body.Bytes(), ShouldEqualJSON, `{
	"result": [{
		"_id": "nonexistent",
		"_type": "error",
		"message": "cannot find device \"nonexistent\"",
		"name": "ResourceNotFound",
		"code": 110,
		"info": {"id": "nonexistent"}
	}]
}`)
			So(called, ShouldBeFalse)
		})
	})

}

func TestPushToUser(t *testing.T) {
	Convey("push to user", t, func() {
		testdevice1 := skydb.Device{
			ID:         "device1",
			Type:       "ios",
			Token:      "token1",
			AuthInfoID: "johndoe",
		}
		testdevice2 := skydb.Device{
			ID:         "device2",
			Type:       "ios",
			Token:      "token2",
			AuthInfoID: "johndoe",
		}
		testdevice3 := skydb.Device{
			ID:         "device2",
			Type:       "ios",
			Token:      "token3",
			AuthInfoID: "janedoe",
		}
		conn := simpleDeviceConn{
			devices: []skydb.Device{testdevice1, testdevice2, testdevice3},
		}

		r := handlertest.NewSingleRouteRouter(&PushToUserHandler{}, func(p *router.Payload) {
			p.DBConn = &conn
		})

		originalSendFunc := sendPushNotification
		defer func() {
			sendPushNotification = originalSendFunc
		}()

		Convey("push to single user", func(c C) {
			sentDevices := []skydb.Device{}
			sendPushNotification = func(sender push.Sender, device skydb.Device, m push.Mapper) {
				c.So(m.Map(), ShouldResemble, map[string]interface{}{
					"aps": map[string]interface{}{
						"alert": "This is a message.",
						"sound": "sosumi.mp3",
					},
					"acme": "interesting",
				})
				sentDevices = append(sentDevices, device)
			}

			resp := r.POST(`{
					"user_ids": ["johndoe"],
					"notification": {
						"aps": {
							"alert": "This is a message.",
							"sound": "sosumi.mp3"
						},
						"acme": "interesting"
					}
				}`)
			So(resp.Code, ShouldEqual, 200)
			So(resp.Body.Bytes(), ShouldEqualJSON, `{
	"result": [{"_id":"johndoe"}]
}`)

			So(len(sentDevices), ShouldEqual, 2)
			So(sentDevices[0], ShouldResemble, testdevice1)
			So(sentDevices[1], ShouldResemble, testdevice2)
		})

		Convey("push to non-existent user", func() {
			called := false
			sendPushNotification = func(sender push.Sender, device skydb.Device, m push.Mapper) {
				called = true
			}
			resp := r.POST(`{
					"user_ids": ["nonexistent"],
					"notification": {
						"aps": {
							"alert": "This is a message.",
							"sound": "sosumi.mp3"
						},
						"acme": "interesting"
					}
				}`)
			So(resp.Code, ShouldEqual, 200)
			So(resp.Body.Bytes(), ShouldEqualJSON, `{
	"result": [{
		"_id": "nonexistent",
		"_type": "error",
		"message": "cannot find user \"nonexistent\"",
		"name": "ResourceNotFound",
		"code": 110,
		"info": {"id": "nonexistent"}
	}]
}`)
			So(called, ShouldBeFalse)
		})
	})

}

type simpleDeviceConn struct {
	devices []skydb.Device
	skydb.Conn
}

func (conn *simpleDeviceConn) GetDevice(id string, device *skydb.Device) error {
	for _, prospectiveDevice := range conn.devices {
		if prospectiveDevice.ID == id {
			*device = prospectiveDevice
			return nil
		}
	}
	return skydb.ErrDeviceNotFound
}

func (conn *simpleDeviceConn) QueryDevicesByUser(user string) ([]skydb.Device, error) {
	var result []skydb.Device
	for _, prospectiveDevice := range conn.devices {
		if prospectiveDevice.AuthInfoID == user {
			result = append(result, prospectiveDevice)
		}
	}
	if len(result) == 0 {
		return nil, skydb.ErrUserNotFound
	}
	return result, nil
}

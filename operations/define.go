package operations

import "github.com/thingio/edge-device-std/models"

type DriverInitialization struct {
	Products []*models.Product `json:"products"`
	Devices  []*models.Device  `json:"devices"`
}

type ProductMutation = models.Product
type DeviceMutation = models.Device
type DeviceStatus = models.DeviceStatus
type DriverStatus = models.DriverStatus

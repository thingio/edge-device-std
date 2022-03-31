package models

import (
	"context"
	"github.com/thingio/edge-device-std/logger"
)

// DeviceTwin indicates a connection with a real device.
type DeviceTwin interface {
	// Initialize will try to initialize a device connector to
	// create the connection with device which needs to activate.
	// It must always return nil if the device needn't be initialized.
	Initialize(lg *logger.Logger) error

	// Start will to try to open a connection with the real device.
	// It must always return nil if the device needn't be initialized.
	Start(ctx context.Context) error
	// Stop will to try to close a connection with the real device.
	// It must always return nil if the device needn't be initialized.
	Stop(force bool) error
	// HealthCheck is used to check the connectivity with the real device.
	HealthCheck() (*DeviceStatus, error)

	// Read indicates hard read, it will read the specified property from the real device.
	Read(propertyID ProductPropertyID) (map[ProductPropertyID]*DeviceData, error)
	// Write will write the specified property to the real device.
	Write(propertyID ProductPropertyID, values map[ProductPropertyID]*DeviceData) error
	// Subscribe will subscribe the specified event,
	// and you should put the event into the bus.
	Subscribe(eventID ProductEventID, bus chan<- *DeviceDataWrapper) error
	// Call is used to call the specified method defined in product,
	// then waiting for a while to receive its response.
	// If the call is timeout, it will return a timeout errors.
	Call(methodID ProductMethodID, ins map[ProductPropertyID]*DeviceData) (outs map[ProductPropertyID]*DeviceData, err error)
}

// DeviceTwinBuilder is used to create a new device twin using the specified product and device.
type DeviceTwinBuilder func(product *Product, device *Device) (DeviceTwin, error)

package operations

import (
	"github.com/thingio/edge-device-std/logger"
	"github.com/thingio/edge-device-std/models"
	bus "github.com/thingio/edge-device-std/msgbus"
)

func NewDriverClient(mb bus.MessageBus, lg *logger.Logger) (DriverClient, error) {
	mdc, err := newMetaDriverClient(mb, lg)
	if err != nil {
		return nil, err
	}
	ddc, err := newDataDriverClient(mb, lg)
	if err != nil {
		return nil, err
	}
	return &driverClient{
		mdc,
		ddc,
	}, nil
}

type (
	DriverClient interface {
		MetaDriverClient
		DataDriverClient
	}
	driverClient struct {
		MetaDriverClient
		DataDriverClient
	}
)

type (
	MetaDriverClient interface {
		PublishDriverStatus(status *models.DriverStatus) error
	}
	metaDriverClient struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newMetaDriverClient(mb bus.MessageBus, lg *logger.Logger) (MetaDriverClient, error) {
	return &metaDriverClient{mb: mb, lg: lg}, nil
}

func (m *metaDriverClient) PublishDriverStatus(status *models.DriverStatus) error {
	o := NewMetaOperation(OperationModeUp, status.Protocol.ID,
		MetaOperationTypeDriverHealthCheck, EmptyReqID())
	o.SetValue(status)
	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return m.mb.Publish(msg)
}

type (
	DataDriverClient interface {
		PublishDeviceStatus(protocolID, productID, deviceID string, status *models.DeviceStatus) error
		PublishDeviceProps(protocolID, productID, deviceID string, propertyID models.ProductPropertyID,
			props map[models.ProductPropertyID]*models.DeviceData) error
		PublishDeviceEvent(protocolID, productID, deviceID string, eventID models.ProductEventID,
			props map[models.ProductPropertyID]*models.DeviceData) error
	}
	dataDriverClient struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newDataDriverClient(mb bus.MessageBus, lg *logger.Logger) (DataDriverClient, error) {
	return &dataDriverClient{mb: mb, lg: lg}, nil
}

func (d *dataDriverClient) PublishDeviceStatus(protocolID, productID, deviceID string, status *models.DeviceStatus) error {
	o := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, "-",
		DataOperationTypeHealthCheck, EmptyReqID())
	o.SetValue(status)
	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return d.mb.Publish(msg)
}

func (d *dataDriverClient) PublishDeviceProps(protocolID, productID, deviceID string, propertyID models.ProductPropertyID,
	props map[models.ProductPropertyID]*models.DeviceData) error {
	o := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, propertyID,
		DataOperationTypeWatch, EmptyReqID())
	o.SetValue(props)
	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return d.mb.Publish(msg)
}

func (d *dataDriverClient) PublishDeviceEvent(protocolID, productID, deviceID string, eventID models.ProductEventID,
	props map[models.ProductPropertyID]*models.DeviceData) error {
	o := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, eventID,
		DataOperationTypeEvent, EmptyReqID())
	o.SetValue(props)
	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return d.mb.Publish(msg)
}

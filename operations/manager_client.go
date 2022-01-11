package operations

import (
	"fmt"
	"github.com/thingio/edge-device-std/errors"
	"github.com/thingio/edge-device-std/logger"
	"github.com/thingio/edge-device-std/models"
	bus "github.com/thingio/edge-device-std/msgbus"
)

func NewManagerClient(mb bus.MessageBus, lg *logger.Logger) (ManagerClient, error) {
	mmc, err := newMetaManagerClient(mb, lg)
	if err != nil {
		return nil, err
	}
	dmc, err := newDataManagerClient(mb, lg)
	if err != nil {
		return nil, err
	}
	return &managerClient{
		mmc,
		dmc,
	}, nil
}

type (
	ManagerClient interface {
		MetaManagerClient
		DataManagerClient
	}
	managerClient struct {
		MetaManagerClient
		DataManagerClient
	}
)

type (
	MetaManagerClient interface {
		InitDriver(protocolID string, products []*models.Product, devices []*models.Device) error

		UpdateProduct(protocolID string, product *models.Product) error
		DeleteProduct(protocolID string, productID string) error

		UpdateDevice(protocolID string, device *models.Device) error
		DeleteDevice(protocolID string, deviceID string) error
	}
	metaManagerClient struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newMetaManagerClient(mb bus.MessageBus, lg *logger.Logger) (MetaManagerClient, error) {
	return &metaManagerClient{mb: mb, lg: lg}, nil
}

func (m *metaManagerClient) InitDriver(protocolID string, products []*models.Product, devices []*models.Device) error {
	o := NewMetaOperation(OperationModeDown, protocolID,
		MetaOperationTypeDriverInit, EmptyReqID())
	o.SetValue(&DriverInitialization{
		Products: products,
		Devices:  devices,
	})

	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return m.mb.Publish(msg)
}

func (m *metaManagerClient) UpdateProduct(protocolID string, product *models.Product) error {
	o := NewMetaOperation(OperationModeDown, protocolID,
		MetaOperationTypeProductMutation, product.ID)
	o.SetValue(product)

	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return m.mb.Publish(msg)
}

func (m *metaManagerClient) DeleteProduct(protocolID string, productID string) error {
	o := NewMetaOperation(OperationModeDown, protocolID,
		MetaOperationTypeProductMutation, productID)

	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return m.mb.Publish(msg)
}

func (m *metaManagerClient) UpdateDevice(protocolID string, device *models.Device) error {
	o := NewMetaOperation(OperationModeDown, protocolID,
		MetaOperationTypeDeviceMutation, device.ID)
	o.SetValue(device)

	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return m.mb.Publish(msg)
}

func (m metaManagerClient) DeleteDevice(protocolID string, deviceID string) error {
	o := NewMetaOperation(OperationModeDown, protocolID,
		MetaOperationTypeDeviceMutation, deviceID)

	msg, err := o.ToMessage()
	if err != nil {
		return err
	}
	return m.mb.Publish(msg)
}

type (
	DataManagerClient interface {
		Read(protocolID, productID, deviceID string,
			propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error)
		HardRead(protocolID, productID, deviceID string,
			propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error)
		Write(protocolID, productID, deviceID string,
			propertyID models.ProductPropertyID, props map[models.ProductPropertyID]*models.DeviceData) error
		Call(protocolID, productID, deviceID string, methodID models.ProductMethodID,
			ins map[string]*models.DeviceData) (outs map[string]*models.DeviceData, err error)
	}
	dataManagerClient struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newDataManagerClient(mb bus.MessageBus, lg *logger.Logger) (DataManagerClient, error) {
	return &dataManagerClient{mb: mb, lg: lg}, nil
}

func (d *dataManagerClient) Read(protocolID, productID, deviceID string,
	propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error) {
	reqID := NewReqID()
	request := NewDataOperation(OperationModeDown, protocolID, productID, deviceID, propertyID,
		DataOperationTypeRead, reqID)
	reqMsg, err := request.ToMessage()
	if err != nil {
		return nil, err
	}
	rspTpc := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, propertyID, DataOperationTypeRead, reqID).Topic().String()
	errTpc := NewDataOperation(OperationModeUpErr, protocolID, productID, deviceID, propertyID, DataOperationTypeRead, reqID).Topic().String()
	rspMsg, err := d.mb.Call(reqMsg, rspTpc, errTpc)
	if err != nil {
		return nil, errors.NewCommonEdgeErrorWrapper(err)
	}
	props = make(map[models.ProductPropertyID]*models.DeviceData)
	if err = rspMsg.Unmarshal(&props); err != nil {
		return nil, errors.NewCommonEdgeError(errors.Internal,
			fmt.Sprintf("fail to unmarshal the payload of the response"), err)
	}
	return props, nil
}

func (d *dataManagerClient) HardRead(protocolID, productID, deviceID string,
	propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error) {
	reqID := NewReqID()
	request := NewDataOperation(OperationModeDown, protocolID, productID, deviceID, propertyID,
		DataOperationTypeHardRead, reqID)
	reqMsg, err := request.ToMessage()
	if err != nil {
		return nil, err
	}
	rspTpc := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, propertyID, DataOperationTypeHardRead, reqID).Topic().String()
	errTpc := NewDataOperation(OperationModeUpErr, protocolID, productID, deviceID, propertyID, DataOperationTypeHardRead, reqID).Topic().String()
	rspMsg, err := d.mb.Call(reqMsg, rspTpc, errTpc)
	if err != nil {
		return nil, errors.NewCommonEdgeErrorWrapper(err)
	}
	props = make(map[models.ProductPropertyID]*models.DeviceData)
	if err = rspMsg.Unmarshal(&props); err != nil {
		return nil, errors.NewCommonEdgeError(errors.Internal,
			fmt.Sprintf("fail to unmarshal the payload of the response"), err)
	}
	return props, nil
}

func (d *dataManagerClient) Write(protocolID, productID, deviceID string,
	propertyID models.ProductPropertyID, props map[models.ProductPropertyID]*models.DeviceData) error {
	reqID := NewReqID()
	request := NewDataOperation(OperationModeDown, protocolID, productID, deviceID, propertyID,
		DataOperationTypeWrite, reqID)
	request.SetValue(props)
	reqMsg, err := request.ToMessage()
	if err != nil {
		return err
	}
	rspTpc := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, propertyID, DataOperationTypeWrite, reqID).Topic().String()
	errTpc := NewDataOperation(OperationModeUpErr, protocolID, productID, deviceID, propertyID, DataOperationTypeWrite, reqID).Topic().String()
	if _, err = d.mb.Call(reqMsg, rspTpc, errTpc); err != nil {
		return errors.NewCommonEdgeErrorWrapper(err)
	}
	return nil
}

func (d *dataManagerClient) Call(protocolID, productID, deviceID string, methodID models.ProductMethodID,
	ins map[string]*models.DeviceData) (outs map[string]*models.DeviceData, err error) {
	reqID := NewReqID()
	request := NewDataOperation(OperationModeDown, protocolID, productID, deviceID, methodID,
		DataOperationTypeCall, reqID)
	request.SetValue(ins)
	reqMsg, err := request.ToMessage()
	if err != nil {
		return nil, err
	}
	rspTpc := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, methodID, DataOperationTypeCall, reqID).Topic().String()
	errTpc := NewDataOperation(OperationModeUpErr, protocolID, productID, deviceID, methodID, DataOperationTypeCall, reqID).Topic().String()
	rspMsg, err := d.mb.Call(reqMsg, rspTpc, errTpc)
	if err != nil {
		return nil, errors.NewCommonEdgeErrorWrapper(err)
	}
	outs = make(map[models.ProductPropertyID]*models.DeviceData)
	if err = rspMsg.Unmarshal(&outs); err != nil {
		return nil, errors.NewCommonEdgeError(errors.Internal,
			fmt.Sprintf("fail to unmarshal the payload of the response"), err)
	}
	return outs, nil
}

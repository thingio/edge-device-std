package operations

import (
	"github.com/thingio/edge-device-std/errors"
	"github.com/thingio/edge-device-std/logger"
	"github.com/thingio/edge-device-std/models"
	bus "github.com/thingio/edge-device-std/msgbus"
	"github.com/thingio/edge-device-std/msgbus/message"
)

func NewDriverService(mb bus.MessageBus, lg *logger.Logger) (DriverService, error) {
	mds, err := newMetaDriverService(mb, lg)
	if err != nil {
		return nil, err
	}
	dds, err := newDataDriverService(mb, lg)
	if err != nil {
		return nil, err
	}
	return &driverService{
		mds,
		dds,
	}, nil
}

type (
	DriverService interface {
		MetaDriverService
		DataDriverService
	}
	driverService struct {
		MetaDriverService
		DataDriverService
	}
)

type (
	MetaDriverService interface {
		InitializeDriverHandler(protocolID string, handler func(products []*models.Product, devices []*models.Device) error) error

		UpdateProductHandler(protocolID string, handler func(product *models.Product) error) error
		DeleteProductHandler(protocolID string, handler func(productID string) error) error

		UpdateDeviceHandler(protocolID string, handler func(device *models.Device) error) error
		DeleteDeviceHandler(protocolID string, handler func(deviceID string) error) error
	}
	metaDriverService struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newMetaDriverService(mb bus.MessageBus, lg *logger.Logger) (MetaDriverService, error) {
	return &metaDriverService{mb: mb, lg: lg}, nil
}

func (m *metaDriverService) InitializeDriverHandler(protocolID string,
	handler func(products []*models.Product, devices []*models.Device) error) error {
	return m.metaHandler(protocolID, MetaOperationTypeDriverInit, func(o *MetaOperation) error {
		v := new(DriverInitialization)
		if err := o.Unmarshal(v); err != nil {
			return err
		}
		return handler(v.Products, v.Devices)
	})
}

func (m *metaDriverService) UpdateProductHandler(protocolID string, handler func(product *models.Product) error) error {
	return m.metaHandler(protocolID, MetaOperationTypeProductMutation, func(o *MetaOperation) error {
		v := new(ProductMutation)
		if err := o.Unmarshal(v); err != nil {
			return err
		}
		return handler(v)
	})
}

func (m *metaDriverService) DeleteProductHandler(protocolID string, handler func(productID string) error) error {
	return m.metaHandler(protocolID, MetaOperationTypeProductMutation, func(o *MetaOperation) error {
		return handler(o.reqID)
	})
}

func (m *metaDriverService) UpdateDeviceHandler(protocolID string, handler func(device *models.Device) error) error {
	return m.metaHandler(protocolID, MetaOperationTypeDeviceMutation, func(o *MetaOperation) error {
		v := new(DeviceMutation)
		if err := o.Unmarshal(v); err != nil {
			return err
		}
		return handler(v)
	})
}

func (m *metaDriverService) DeleteDeviceHandler(protocolID string, handler func(deviceID string) error) error {
	return m.metaHandler(protocolID, MetaOperationTypeDeviceMutation, func(o *MetaOperation) error {
		return handler(o.reqID)
	})
}

func (m *metaDriverService) metaHandler(protocolID string, optType MetaOperationType,
	handler func(o *MetaOperation) error) error {
	schema := NewMetaOperation(OperationModeDown, protocolID, optType, TopicSingleLevelWildcard)
	topic := schema.Topic().String()
	if err := m.mb.Subscribe(func(msg *message.Message) {
		o, err := ParseMetaOperation(msg)
		if err != nil {
			m.lg.WithError(err).Errorf("fail to parse the meta operation: %s", topic)
			return
		}
		if err = handler(o); err != nil {
			m.lg.Errorf(err.Error())
			return
		}
	}, topic); err != nil {
		return err
	}
	return nil
}

type (
	DataDriverService interface {
		ReadHandler(protocolID string, handler func(productID, deviceID string,
			propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error)) error
		HardReadHandler(protocolID string, handler func(productID, deviceID string,
			propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error)) error
		WriteHandler(protocolID string, handler func(productID, deviceID string,
			propertyID models.ProductPropertyID, props map[models.ProductPropertyID]*models.DeviceData) error) error
		CallHandler(protocolID string, handler func(productID, deviceID string, methodID models.ProductMethodID,
			ins map[string]*models.DeviceData) (outs map[string]*models.DeviceData, err error)) error
	}
	dataDriverService struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newDataDriverService(mb bus.MessageBus, lg *logger.Logger) (DataDriverService, error) {
	return &dataDriverService{mb: mb, lg: lg}, nil
}

func (d *dataDriverService) ReadHandler(protocolID string, handler func(productID string, deviceID string,
	propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error)) error {
	return d.dataHandler(protocolID, DataOperationTypeRead,
		func(o *DataOperation) (outs map[string]*models.DeviceData, err error) {
			productID, deviceID, propertyID := o.productID, o.deviceID, o.funcID
			return handler(productID, deviceID, propertyID)
		},
	)
}

func (d *dataDriverService) HardReadHandler(protocolID string, handler func(productID string, deviceID string,
	propertyID models.ProductPropertyID) (props map[models.ProductPropertyID]*models.DeviceData, err error)) error {
	return d.dataHandler(protocolID, DataOperationTypeHardRead,
		func(o *DataOperation) (outs map[string]*models.DeviceData, err error) {
			productID, deviceID, propertyID := o.productID, o.deviceID, o.funcID
			return handler(productID, deviceID, propertyID)
		},
	)
}

func (d *dataDriverService) WriteHandler(protocolID string, handler func(productID string, deviceID string,
	propertyID models.ProductPropertyID, props map[models.ProductPropertyID]*models.DeviceData) error) error {
	return d.dataHandler(protocolID, DataOperationTypeWrite,
		func(o *DataOperation) (outs map[string]*models.DeviceData, err error) {
			productID, deviceID, propertyID := o.productID, o.deviceID, o.funcID
			props := make(map[models.ProductPropertyID]*models.DeviceData)
			if err = o.Unmarshal(&props); err != nil {
				d.lg.WithError(err).Errorf("fail to unmarshal the property[%s] "+
					"from the device[%s]", propertyID, deviceID)
				return
			}
			return map[models.ProductPropertyID]*models.DeviceData{}, handler(productID, deviceID, propertyID, props)
		},
	)
}

func (d *dataDriverService) CallHandler(protocolID string, handler func(productID string, deviceID string, methodID models.ProductMethodID,
	ins map[string]*models.DeviceData) (outs map[string]*models.DeviceData, err error)) error {
	return d.dataHandler(protocolID, DataOperationTypeCall,
		func(o *DataOperation) (outs map[string]*models.DeviceData, err error) {
			productID, deviceID, methodID := o.productID, o.deviceID, o.funcID
			ins := make(map[string]*models.DeviceData)
			if err = o.Unmarshal(&ins); err != nil {
				d.lg.WithError(err).Errorf("fail to unmarshal the ins of the method[%s] "+
					"of the device[%s]", methodID, deviceID)
				return
			}
			return handler(productID, deviceID, methodID, ins)
		},
	)
}

func (d *dataDriverService) dataHandler(protocolID string, optType DataOperationType,
	handler func(o *DataOperation) (outs map[string]*models.DeviceData, err error)) error {
	schema := NewDataOperation(OperationModeDown, protocolID,
		TopicSingleLevelWildcard, TopicSingleLevelWildcard, TopicSingleLevelWildcard,
		optType, TopicSingleLevelWildcard)
	topic := schema.Topic().String()
	if err := d.mb.Subscribe(func(msg *message.Message) {
		var err error
		var request, response *DataOperation
		var outs map[string]*models.DeviceData
		defer func() {
			response = NewDataOperation(OperationModeUp, request.protocolID, request.productID, request.deviceID,
				request.funcID, optType, request.reqID)
			if err != nil {
				response.optMode = OperationModeUpErr
				response.SetValue(errors.NewCommonEdgeErrorWrapper(err))
			} else {
				response.SetValue(outs)
			}
			rspMsg, err := response.ToMessage()
			if err != nil {
				d.lg.WithError(err).Errorf("fail to parse the message of the response")
				return
			}
			_ = d.mb.Publish(rspMsg)
		}()

		if request, err = ParseDataOperation(msg); err != nil {
			d.lg.WithError(err).Errorf("fail to parse the data operation")
			return
		}
		if outs, err = handler(request); err != nil {
			return
		}
	}, topic); err != nil {
		return err
	}
	return nil
}

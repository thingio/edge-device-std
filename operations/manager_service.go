package operations

import (
	"github.com/thingio/edge-device-std/logger"
	"github.com/thingio/edge-device-std/models"
	bus "github.com/thingio/edge-device-std/msgbus"
	"github.com/thingio/edge-device-std/msgbus/message"
)

func NewManagerService(mb bus.MessageBus, lg *logger.Logger) (ManagerService, error) {
	mms, err := newMetaManagerService(mb, lg)
	if err != nil {
		return nil, err
	}
	dms, err := newDataManagerService(mb, lg)
	if err != nil {
		return nil, err
	}
	return &managerService{
		mms,
		dms,
	}, nil
}

type (
	ManagerService interface {
		MetaManagerService
		DataManagerService
	}
	managerService struct {
		MetaManagerService
		DataManagerService
	}
)

type (
	MetaManagerService interface {
		SubscribeDriverStatus() (bus <-chan interface{}, stop func(), err error)
	}
	metaManagerService struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newMetaManagerService(mb bus.MessageBus, lg *logger.Logger) (MetaManagerService, error) {
	return &metaManagerService{mb: mb, lg: lg}, nil
}

func (m *metaManagerService) SubscribeDriverStatus() (<-chan interface{}, func(), error) {
	return m.subscribe(MetaOperationTypeDriverHealthCheck, func(o *MetaOperation) (interface{}, error) {
		status := new(DriverStatus)
		if err := o.Unmarshal(status); err != nil {
			return nil, err
		}
		return status, nil
	})
}

func (m *metaManagerService) subscribe(optType MetaOperationType,
	parser func(o *MetaOperation) (interface{}, error)) (<-chan interface{}, func(), error) {
	schema := NewMetaOperation(OperationModeUp, TopicSingleLevelWildcard, optType, TopicSingleLevelWildcard)
	return subscribe(m.mb, m.lg, schema.Topic().String(), func(msg *message.Message) (interface{}, error) {
		o, err := ParseMetaOperation(msg)
		if err != nil {
			return nil, err
		}
		return parser(o)
	})
}

type (
	DataManagerService interface {
		SubscribeDeviceStatus(protocolID string) (<-chan interface{}, func(), error)
		SubscribeDeviceProps(protocolID, productID, deviceID string, propertyID models.ProductPropertyID) (<-chan interface{}, func(), error)
		SubscribeDeviceEvent(protocolID, productID, deviceID string, eventID models.ProductEventID) (<-chan interface{}, func(), error)
	}
	dataManagerService struct {
		mb bus.MessageBus
		lg *logger.Logger
	}
)

func newDataManagerService(mb bus.MessageBus, lg *logger.Logger) (DataManagerService, error) {
	return &dataManagerService{mb: mb, lg: lg}, nil
}

func (d *dataManagerService) SubscribeDeviceStatus(protocolID string) (<-chan interface{}, func(), error) {
	return d.subscribe(protocolID, TopicSingleLevelWildcard, TopicSingleLevelWildcard, TopicSingleLevelWildcard, DataOperationTypeHealthCheck,
		func(o *DataOperation) (interface{}, error) {
			status := new(DeviceStatus)
			if err := o.Unmarshal(status); err != nil {
				return nil, err
			}
			return status, nil
		},
	)
}

func (d *dataManagerService) SubscribeDeviceProps(protocolID, productID, deviceID string, propertyID models.ProductPropertyID) (<-chan interface{}, func(), error) {
	return d.subscribe(protocolID, productID, deviceID, propertyID, DataOperationTypeWatch,
		func(o *DataOperation) (interface{}, error) {
			props := make(map[models.ProductPropertyID]*models.DeviceData)
			if err := o.Unmarshal(&props); err != nil {
				return nil, err
			}
			return props, nil
		},
	)
}

func (d *dataManagerService) SubscribeDeviceEvent(protocolID, productID, deviceID string, eventID models.ProductEventID) (<-chan interface{}, func(), error) {
	return d.subscribe(protocolID, productID, deviceID, eventID, DataOperationTypeEvent,
		func(o *DataOperation) (interface{}, error) {
			props := make(map[models.ProductPropertyID]*models.DeviceData)
			if err := o.Unmarshal(&props); err != nil {
				return nil, err
			}
			return props, nil
		},
	)
}

func (d *dataManagerService) subscribe(protocolID, productID, deviceID string, funcID models.ProductEventID,
	optType DataOperationType, parser func(o *DataOperation) (interface{}, error)) (<-chan interface{}, func(), error) {
	schema := NewDataOperation(OperationModeUp, protocolID, productID, deviceID, funcID, optType, TopicSingleLevelWildcard)
	return subscribe(d.mb, d.lg, schema.Topic().String(),
		func(msg *message.Message) (interface{}, error) {
			o, err := ParseDataOperation(msg)
			if err != nil {
				return nil, err
			}
			return parser(o)
		},
	)
}

func subscribe(mb bus.MessageBus, lg *logger.Logger, topic string,
	parser func(msg *message.Message) (interface{}, error)) (<-chan interface{}, func(), error) {
	buffer := make(chan interface{}, 1000)
	if err := mb.Subscribe(func(msg *message.Message) {
		v, err := parser(msg)
		if err != nil {
			lg.WithError(err).Errorf("fail to parse the payload of the data operation")
		}
		buffer <- v
	}, topic); err != nil {
		return nil, nil, err
	}
	return buffer, func() {
		_ = mb.Unsubscribe(topic)

		close(buffer)
	}, nil
}

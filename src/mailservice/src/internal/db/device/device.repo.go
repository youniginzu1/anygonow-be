package device

import "context"

type DeviceRepo interface {
	InsertDevice(context.Context, *Device) (*Device, error)
	SelectDevice(context.Context, *Search) (*Device, error)
	ListDevices(context.Context, *Search) ([]*Device, error)
	UpdateDevice(context.Context, *Search, *Device) error
	DeleteDevice(context.Context, *Search, *Device) error
}

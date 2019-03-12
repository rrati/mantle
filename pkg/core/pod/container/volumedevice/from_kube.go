package volumedevice

import (
	"fmt"
	"reflect"

	"k8s.io/api/core/v1"
)

// NewVolumeDeviceFromKubeVolumeDevice will create a new
// VolumeDevice object with the data from a provided kubernetes
// VolumeDevice object
func NewVolumeDeviceFromKubeVolumeDevice(obj interface{}) (*VolumeDevice, error) {
	switch reflect.TypeOf(obj) {
	case reflect.TypeOf(v1.VolumeDevice{}):
		return fromKubeVolumeDeviceV1(obj.(v1.VolumeDevice))
	case reflect.TypeOf(&v1.VolumeDevice{}):
		o := obj.(*v1.VolumeDevice)
		return fromKubeVolumeDeviceV1(*o)
	default:
		return nil, fmt.Errorf("unknown VolumeDevice version: %s", reflect.TypeOf(obj))
	}
}

func fromKubeVolumeDeviceV1(device v1.VolumeDevice) (*VolumeDevice, error) {
	return &VolumeDevice{
		Name: device.Name,
		Path: device.DevicePath,
	}, nil
}

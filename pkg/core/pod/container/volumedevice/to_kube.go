package volumedevice

import (
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
)

// ToKube will return a kubernetes VolumeDevice object of the api version provided
func (vd *VolumeDevice) ToKube(version string) (interface{}, error) {
	switch strings.ToLower(version) {
	case "v1":
		return vd.toKubeV1()
	case "":
		return vd.toKubeV1()
	default:
		return nil, fmt.Errorf("unsupported api version for VolumeDevice: %s", version)
	}
}

func (vd *VolumeDevice) toKubeV1() (*v1.VolumeDevice, error) {
	return &v1.VolumeDevice{
		Name:       vd.Name,
		DevicePath: vd.Path,
	}, nil
}

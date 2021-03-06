package projected

import (
	"fmt"
	"reflect"

	"mantle/internal/converterutils"
	"mantle/internal/pkg/core/pod/volume/keyandmode"

	"k8s.io/api/core/v1"
)

// NewSecretProjectionFromKubeSecretProjection will create a new
// SecretProjection object with the data from a provided kubernetes
// SecretProjection object
func NewSecretProjectionFromKubeSecretProjection(obj interface{}) (*SecretProjection, error) {
	switch reflect.TypeOf(obj) {
	case reflect.TypeOf(v1.SecretProjection{}):
		o := obj.(v1.SecretProjection)
		return fromKubeSecretProjectionV1(&o)
	case reflect.TypeOf(&v1.SecretProjection{}):
		return fromKubeSecretProjectionV1(obj.(*v1.SecretProjection))
	default:
		return nil, fmt.Errorf("unknown SecretProjection version: %s", reflect.TypeOf(obj))
	}
}

func fromKubeSecretProjectionV1(vol *v1.SecretProjection) (*SecretProjection, error) {
	return &SecretProjection{
		Name:  converterutils.FromKubeLocalObjectReferenceV1(&vol.LocalObjectReference),
		Items: keyandmode.NewKeyToPathFromKubeKeyToPathV1(vol.Items),
	}, nil
}

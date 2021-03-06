package pod

import (
	"fmt"
	"reflect"
	"strings"

	. "mantle/internal/marshal"
	. "mantle/internal/pkg/core/pod/volume/aws"
	. "mantle/internal/pkg/core/pod/volume/azure"
	. "mantle/internal/pkg/core/pod/volume/ceph"
	. "mantle/internal/pkg/core/pod/volume/cinder"
	. "mantle/internal/pkg/core/pod/volume/configmap"
	. "mantle/internal/pkg/core/pod/volume/downwardapi"
	. "mantle/internal/pkg/core/pod/volume/emptydir"
	. "mantle/internal/pkg/core/pod/volume/fiberchannel"
	. "mantle/internal/pkg/core/pod/volume/flex"
	. "mantle/internal/pkg/core/pod/volume/flocker"
	. "mantle/internal/pkg/core/pod/volume/gcepd"
	. "mantle/internal/pkg/core/pod/volume/git"
	. "mantle/internal/pkg/core/pod/volume/gluster"
	. "mantle/internal/pkg/core/pod/volume/hostpath"
	. "mantle/internal/pkg/core/pod/volume/iscsi"
	. "mantle/internal/pkg/core/pod/volume/nfs"
	. "mantle/internal/pkg/core/pod/volume/photon"
	. "mantle/internal/pkg/core/pod/volume/portworx"
	. "mantle/internal/pkg/core/pod/volume/projected"
	. "mantle/internal/pkg/core/pod/volume/pvc"
	. "mantle/internal/pkg/core/pod/volume/quobyte"
	. "mantle/internal/pkg/core/pod/volume/rbd"
	. "mantle/internal/pkg/core/pod/volume/scaleio"
	. "mantle/internal/pkg/core/pod/volume/secret"
	. "mantle/internal/pkg/core/pod/volume/storageos"
	. "mantle/internal/pkg/core/pod/volume/vsphere"

	"github.com/koki/json"
	"github.com/koki/json/jsonutil"
	serrors "github.com/koki/structurederrors"
)

type Volume struct {
	HostPath     *HostPathVolume
	EmptyDir     *EmptyDirVolume
	GcePD        *GcePDVolume
	AwsEBS       *AwsEBSVolume
	AzureDisk    *AzureDiskVolume
	AzureFile    *AzureFileVolume
	CephFS       *CephFSVolume
	Cinder       *CinderVolume
	FibreChannel *FibreChannelVolume
	Flex         *FlexVolume
	Flocker      *FlockerVolume
	Glusterfs    *GlusterfsVolume
	ISCSI        *ISCSIVolume
	NFS          *NFSVolume
	PhotonPD     *PhotonPDVolume
	Portworx     *PortworxVolume
	PVC          *PVCVolume
	Quobyte      *QuobyteVolume
	ScaleIO      *ScaleIOVolume
	Vsphere      *VsphereVolume
	ConfigMap    *ConfigMapVolume
	Secret       *SecretVolume
	DownwardAPI  *DownwardAPIVolume
	Projected    *ProjectedVolume
	Git          *GitVolume
	RBD          *RBDVolume
	StorageOS    *StorageOSVolume
}

func (v *Volume) UnmarshalJSON(data []byte) error {
	var err error
	str := ""
	err = json.Unmarshal(data, &str)
	if err == nil {
		segments := strings.Split(str, ":")
		return v.Unmarshal(nil, segments[0], segments[1:])
	}

	obj := map[string]interface{}{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return serrors.InvalidValueErrorf(string(data), "expected either string or dictionary")
	}

	selector := []string{}
	if val, ok := obj["vol_id"]; ok {
		if volName, ok := val.(string); ok {
			selector = append(selector, volName)
		} else {
			return serrors.InvalidValueErrorf(string(data), "expected string for key \"vol_id\"")
		}
	}

	volType, err := jsonutil.GetStringEntry(obj, "vol_type")
	if err != nil {
		return err
	}

	return v.Unmarshal(obj, volType, selector)
}

func (v *Volume) Unmarshal(obj map[string]interface{}, volType string, selector []string) error {
	switch volType {
	case VolumeTypeHostPath:
		v.HostPath = &HostPathVolume{}
		return v.HostPath.Unmarshal(selector)
	case VolumeTypeEmptyDir:
		v.EmptyDir = &EmptyDirVolume{}
		return v.EmptyDir.Unmarshal(obj, selector)
	case VolumeTypeGcePD:
		v.GcePD = &GcePDVolume{}
		return v.GcePD.Unmarshal(obj, selector)
	case VolumeTypeAwsEBS:
		v.AwsEBS = &AwsEBSVolume{}
		return v.AwsEBS.Unmarshal(obj, selector)
	case VolumeTypeAzureDisk:
		v.AzureDisk = &AzureDiskVolume{}
		return v.AzureDisk.Unmarshal(obj, selector)
	case VolumeTypeAzureFile:
		v.AzureFile = &AzureFileVolume{}
		return v.AzureFile.Unmarshal(selector)
	case VolumeTypeCephFS:
		v.CephFS = &CephFSVolume{}
		return v.CephFS.Unmarshal(obj, selector)
	case VolumeTypeCinder:
		v.Cinder = &CinderVolume{}
		return v.Cinder.Unmarshal(obj, selector)
	case VolumeTypeFibreChannel:
		v.FibreChannel = &FibreChannelVolume{}
		return v.FibreChannel.Unmarshal(obj, selector)
	case VolumeTypeFlex:
		v.Flex = &FlexVolume{}
		return v.Flex.Unmarshal(obj, selector)
	case VolumeTypeFlocker:
		v.Flocker = &FlockerVolume{}
		return v.Flocker.Unmarshal(selector)
	case VolumeTypeGlusterfs:
		v.Glusterfs = &GlusterfsVolume{}
		return v.Glusterfs.Unmarshal(obj, selector)
	case VolumeTypeISCSI:
		v.ISCSI = &ISCSIVolume{}
		return v.ISCSI.Unmarshal(obj, selector)
	case VolumeTypeNFS:
		v.NFS = &NFSVolume{}
		return v.NFS.Unmarshal(selector)
	case VolumeTypePhotonPD:
		v.PhotonPD = &PhotonPDVolume{}
		return v.PhotonPD.Unmarshal(selector)
	case VolumeTypePortworx:
		v.Portworx = &PortworxVolume{}
		return v.Portworx.Unmarshal(obj, selector)
	case VolumeTypePVC:
		v.PVC = &PVCVolume{}
		return v.PVC.Unmarshal(selector)
	case VolumeTypeQuobyte:
		v.Quobyte = &QuobyteVolume{}
		return v.Quobyte.Unmarshal(obj, selector)
	case VolumeTypeScaleIO:
		v.ScaleIO = &ScaleIOVolume{}
		return v.ScaleIO.Unmarshal(obj, selector)
	case VolumeTypeVsphere:
		v.Vsphere = &VsphereVolume{}
		return v.Vsphere.Unmarshal(obj, selector)
	case VolumeTypeConfigMap:
		v.ConfigMap = &ConfigMapVolume{}
		return v.ConfigMap.Unmarshal(obj, selector)
	case VolumeTypeSecret:
		v.Secret = &SecretVolume{}
		return v.Secret.Unmarshal(obj, selector)
	case VolumeTypeDownwardAPI:
		v.DownwardAPI = &DownwardAPIVolume{}
		return v.DownwardAPI.Unmarshal(obj, selector)
	case VolumeTypeProjected:
		v.Projected = &ProjectedVolume{}
		return v.Projected.Unmarshal(obj, selector)
	case VolumeTypeGit:
		v.Git = &GitVolume{}
		return v.Git.Unmarshal(obj, selector)
	case VolumeTypeRBD:
		v.RBD = &RBDVolume{}
		return v.RBD.Unmarshal(obj, selector)
	case VolumeTypeStorageOS:
		v.StorageOS = &StorageOSVolume{}
		return v.StorageOS.Unmarshal(obj, selector)
	default:
		return serrors.InvalidValueErrorf(volType, "unsupported volume type (%s)", volType)
	}
}

func (v Volume) MarshalJSON() ([]byte, error) {
	var marshalledVolume *MarshalledVolume
	var err error
	if v.HostPath != nil {
		marshalledVolume, err = v.HostPath.Marshal()
	}
	if v.EmptyDir != nil {
		marshalledVolume, err = v.EmptyDir.Marshal()
	}
	if v.GcePD != nil {
		marshalledVolume, err = v.GcePD.Marshal()
	}
	if v.AwsEBS != nil {
		marshalledVolume, err = v.AwsEBS.Marshal()
	}
	if v.AzureDisk != nil {
		marshalledVolume, err = v.AzureDisk.Marshal()
	}
	if v.AzureFile != nil {
		marshalledVolume, err = v.AzureFile.Marshal()
	}
	if v.CephFS != nil {
		marshalledVolume, err = v.CephFS.Marshal()
	}
	if v.Cinder != nil {
		marshalledVolume, err = v.Cinder.Marshal()
	}
	if v.FibreChannel != nil {
		marshalledVolume, err = v.FibreChannel.Marshal()
	}
	if v.Flex != nil {
		marshalledVolume, err = v.Flex.Marshal()
	}
	if v.Flocker != nil {
		marshalledVolume, err = v.Flocker.Marshal()
	}
	if v.Glusterfs != nil {
		marshalledVolume, err = v.Glusterfs.Marshal()
	}
	if v.ISCSI != nil {
		marshalledVolume, err = v.ISCSI.Marshal()
	}
	if v.NFS != nil {
		marshalledVolume, err = v.NFS.Marshal()
	}
	if v.PhotonPD != nil {
		marshalledVolume, err = v.PhotonPD.Marshal()
	}
	if v.Portworx != nil {
		marshalledVolume, err = v.Portworx.Marshal()
	}
	if v.PVC != nil {
		marshalledVolume, err = v.PVC.Marshal()
	}
	if v.Quobyte != nil {
		marshalledVolume, err = v.Quobyte.Marshal()
	}
	if v.ScaleIO != nil {
		marshalledVolume, err = v.ScaleIO.Marshal()
	}
	if v.Vsphere != nil {
		marshalledVolume, err = v.Vsphere.Marshal()
	}
	if v.ConfigMap != nil {
		marshalledVolume, err = v.ConfigMap.Marshal()
	}
	if v.Secret != nil {
		marshalledVolume, err = v.Secret.Marshal()
	}
	if v.DownwardAPI != nil {
		marshalledVolume, err = v.DownwardAPI.Marshal()
	}
	if v.Projected != nil {
		marshalledVolume, err = v.Projected.Marshal()
	}
	if v.Git != nil {
		marshalledVolume, err = v.Git.Marshal()
	}
	if v.RBD != nil {
		marshalledVolume, err = v.RBD.Marshal()
	}
	if v.StorageOS != nil {
		marshalledVolume, err = v.StorageOS.Marshal()
	}

	if err != nil {
		return nil, err
	}

	if marshalledVolume == nil {
		return nil, serrors.InvalidInstanceErrorf(v, "empty volume definition")
	}

	if len(marshalledVolume.ExtraFields) == 0 {
		segments := []string{marshalledVolume.Type}
		segments = append(segments, marshalledVolume.Selector...)
		return json.Marshal(strings.Join(segments, ":"))
	}

	obj := marshalledVolume.ExtraFields
	obj["vol_type"] = marshalledVolume.Type
	if len(marshalledVolume.Selector) > 0 {
		obj["vol_id"] = strings.Join(marshalledVolume.Selector, ":")
	}

	return json.Marshal(obj)
}

func (v *Volume) ToKube(version string) (interface{}, interface{}) {
	fields := reflect.ValueOf(v).Elem()
	for n := 0; n < fields.NumField(); n++ {
		field := fields.Field(n)
		if field.IsValid() && !field.IsNil() {
			convFunc := field.MethodByName("ToKube")
			resp := convFunc.Call([]reflect.Value{reflect.ValueOf(version)})
			return resp[0].Interface(), resp[1].Interface()
		}
	}

	return nil, fmt.Errorf("no volume type set")
}

func NewVolumeFromKubeVolume(obj interface{}) (*Volume, error) {
	return nil, nil
}

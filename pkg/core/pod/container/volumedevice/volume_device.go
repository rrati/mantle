package volumedevice

// VolumeDevice defines a volume device for a container
type VolumeDevice struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

package container

var FullTestMantleContainer = container.Container{
	Command: []string{"cmd1", "cmd2"},
	Args: []floatstr.FloatOrString{
		{
			Type:     floatstr.Float,
			FloatVal: 3.1459267,
		},
		{
			Type:      floatstr.String,
			StringVal: "strVal",
		},
	},
	Env: []env.Env{
		{
			Type: env.EnvValEnvType,
			Val: &env.EnvVal{
				Key: "key1",
				Val: "val1",
			},
		},
		{
			Type: env.EnvFromEnvType,
			From: &env.EnvFrom{
				From:                  env.EnvFromTypeSecret,
				VarNameOrPrefix:       "EnvName",
				ConfigMapOrSecretName: "SecretName",
				ConfigMapOrSecretKey:  "SecretKey",
				Required:              &boolEntry,
			},
		},
	},
	Image: "registry.io/path/to/image",
	Pull:  container.PullAlways,
	OnStart: &action.Action{
		ActionType: action.ActionTypeTCP,
		Host:       "actionHost",
		Port:       "22",
	},
	PreStop: &action.Action{
		ActionType: action.ActionTypeHTTPS,
		Path:       "/action/path",
		Host:       "actionHost",
		Port:       "443",
		Headers:    []string{"Key:Value", "HEADER:VALUE"},
	},
	CPU: &resources.CPU{
		Min: "1",
		Max: "3",
	},
	Mem: &resources.Mem{
		Min: "10Gi",
		Max: "10Ti",
	},
	Name:            "container-name",
	AddCapabilities: []string{"CAP1", "CAP2"},
	DelCapabilities: []string{"CAP3"},
	Privileged:      &boolEntry,
	AllowEscalation: &boolEntry,
	RO:              &boolEntry,
	ForceNonRoot:    &boolEntry,
	UID:             &int64Val,
	GID:             &int64Val,
	SELinux: &selinux.SELinux{
		Level: "selinuxlevel",
		Role:  "selinuxrole",
		Type:  "selinuxtype",
		User:  "selinuxuser",
	},
	LivenessProbe: &probe.Probe{
		Action: action.Action{
			ActionType: action.ActionTypeCommand,
			Command:    []string{"cmd1", "cmd2"},
		},
		Delay:           int32Val,
		Interval:        int32Val,
		MinCountSuccess: int32Val,
		MinCountFailure: int32Val,
		Timeout:         int32Val,
	},
	ReadinessProbe: &probe.Probe{
		Action: action.Action{
			ActionType: action.ActionTypeCommand,
			Command:    []string{"cmd3", "cmd4"},
		},
		Delay:           int32Val,
		Interval:        int32Val,
		MinCountSuccess: int32Val,
		MinCountFailure: int32Val,
		Timeout:         int32Val,
	},
	Expose: []port.Port{
		{
			Name:          "port1",
			Protocol:      protocol.ProtocolTCP,
			IP:            "1.1.1.1",
			HostPort:      "1000",
			ContainerPort: "1000",
		},
	},
	Stdin:                boolEntry,
	StdinOnce:            boolEntry,
	TTY:                  boolEntry,
	ProcMount:            &mountType,
	WorkingDir:           "/path/to/dir",
	TerminationMsgPath:   "/msg/path",
	TerminationMsgPolicy: container.TerminationMessageReadFile,
	VolumeMounts: []volumemount.VolumeMount{
		{
			MountPath:   "/path/to/mount",
			Propagation: &mountPropagation,
			Store:       "storeName",
			ReadOnly:    false,
		},
	},
	VolumeDevices: []volumedevice.VolumeDevice{
		{
			Name: "deviceName",
			Path: "/path/to/device",
		},
	},
}

var EmptyTestMantleContainer = container.Container{
	Command:         []string{},
	Args:            []floatstr.FloatOrString{},
	Env:             []env.Env{},
	Image:           "",
	Pull:            container.PullDefault,
	OnStart:         &action.Action{},
	PreStop:         &action.Action{},
	CPU:             &resources.CPU{},
	Mem:             &resources.Mem{},
	Name:            "",
	AddCapabilities: []string{},
	DelCapabilities: []string{},
	Privileged:      nil,
	AllowEscalation: nil,
	RO:              nil,
	ForceNonRoot:    nil,
	UID:             nil,
	GID:             nil,
	SELinux:         &selinux.SELinux{},
	LivenessProbe: &probe.Probe{
		Action: action.Action{},
	},
	ReadinessProbe:       &probe.Probe{},
	Expose:               []port.Port{},
	Stdin:                false,
	StdinOnce:            false,
	TTY:                  false,
	ProcMount:            nil,
	WorkingDir:           "",
	TerminationMsgPath:   "",
	TerminationMsgPolicy: container.TerminationMessageDefault,
	VolumeMounts:         []volumemount.VolumeMount{},
	VolumeDevices:        []volumedevice.VolumeDevice{},
}

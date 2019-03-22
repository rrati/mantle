package pod

import (
	"reflect"
	"testing"

	"mantle/pkg/core/action"
	"mantle/pkg/core/pod/affinity"
	"mantle/pkg/core/pod/container"
	"mantle/pkg/core/pod/container/env"
	"mantle/pkg/core/pod/container/port"
	"mantle/pkg/core/pod/container/probe"
	"mantle/pkg/core/pod/container/resources"
	"mantle/pkg/core/pod/container/volumedevice"
	"mantle/pkg/core/pod/container/volumemount"
	"mantle/pkg/core/pod/hostalias"
	"mantle/pkg/core/pod/podtemplate"
	"mantle/pkg/core/pod/toleration"
	"mantle/pkg/core/pod/volume"
	"mantle/pkg/core/protocol"
	"mantle/pkg/core/selinux"
	"mantle/pkg/util/floatstr"

	"k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	"github.com/google/go-cmp/cmp"
)

var int64Val = int64(5)
var boolEntry = true
var inverseBoolEntry = !boolEntry
var int32Val = int32(1)
var strEntry = "testStr"
var strArray = []string{strEntry}
var propagation = volumemount.MountPropagationNone
var kubePropagation = v1.MountPropagationNone
var procMount = container.MountTypeDefault
var kubeProcMount = v1.DefaultProcMount
var cpuMin = "1"
var cpuMax = "3"
var memMin = "10Gi"
var memMax = "10Ti"
var cpuLimit, _ = resource.ParseQuantity(cpuMax)
var cpuRequest, _ = resource.ParseQuantity(cpuMin)
var memLimit, _ = resource.ParseQuantity(memMax)
var memRequest, _ = resource.ParseQuantity(memMin)
var fullMantleContainer = container.Container{
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
				VarNameOrPrefix:       "SecretEnvName",
				ConfigMapOrSecretName: "SecretName",
				ConfigMapOrSecretKey:  "SecretKey",
				Required:              &boolEntry,
			},
		},
		{
			Type: env.EnvFromEnvType,
			From: &env.EnvFrom{
				From:                  env.EnvFromTypeConfig,
				VarNameOrPrefix:       "ConfigEnvName",
				ConfigMapOrSecretName: "ConfigName",
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
		Min: cpuMin,
		Max: cpuMax,
	},
	Mem: &resources.Mem{
		Min: memMin,
		Max: memMax,
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
	ProcMount:            &procMount,
	WorkingDir:           "/path/to/dir",
	TerminationMsgPath:   "/msg/path",
	TerminationMsgPolicy: container.TerminationMessageReadFile,
	VolumeMounts: []volumemount.VolumeMount{
		{
			MountPath:   "/path/to/mount",
			Propagation: &propagation,
			Store:       "volumeName:/path/to/volume",
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

var fullKubeContainer = v1.Container{
	Name:       "container-name",
	Image:      "registry.io/path/to/image",
	Command:    []string{"cmd1", "cmd2"},
	Args:       []string{"3.1459267", "strVal"},
	WorkingDir: "/path/to/dir",
	Ports: []v1.ContainerPort{
		{
			Name:          "port1",
			Protocol:      v1.ProtocolTCP,
			HostIP:        "1.1.1.1",
			HostPort:      int32(1000),
			ContainerPort: int32(1000),
		},
	},
	EnvFrom: []v1.EnvFromSource{
		{
			Prefix: "ConfigEnvName",
			ConfigMapRef: &v1.ConfigMapEnvSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "ConfigName",
				},
				Optional: &inverseBoolEntry,
			},
		},
	},
	Env: []v1.EnvVar{
		{
			Name:  "key1",
			Value: "val1",
		},
		{
			Name: "SecretEnvName",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "SecretName",
					},
					Key:      "SecretKey",
					Optional: &inverseBoolEntry,
				},
			},
		},
	},
	Resources: v1.ResourceRequirements{
		Limits: map[v1.ResourceName]resource.Quantity{
			v1.ResourceCPU:    cpuLimit,
			v1.ResourceMemory: memLimit,
		},
		Requests: map[v1.ResourceName]resource.Quantity{
			v1.ResourceCPU:    cpuRequest,
			v1.ResourceMemory: memRequest,
		},
	},
	VolumeMounts: []v1.VolumeMount{
		{
			MountPath:        "/path/to/mount",
			MountPropagation: &kubePropagation,
			Name:             "volumeName",
			SubPath:          "/path/to/volume",
			ReadOnly:         false,
		},
	},
	VolumeDevices: []v1.VolumeDevice{
		{
			Name:       "deviceName",
			DevicePath: "/path/to/device",
		},
	},
	LivenessProbe: &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{"cmd1", "cmd2"},
			},
		},
		InitialDelaySeconds: int32Val,
		TimeoutSeconds:      int32Val,
		PeriodSeconds:       int32Val,
		SuccessThreshold:    int32Val,
		FailureThreshold:    int32Val,
	},
	ReadinessProbe: &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{"cmd3", "cmd4"},
			},
		},
		InitialDelaySeconds: int32Val,
		TimeoutSeconds:      int32Val,
		PeriodSeconds:       int32Val,
		SuccessThreshold:    int32Val,
		FailureThreshold:    int32Val,
	},
	Lifecycle: &v1.Lifecycle{
		PostStart: &v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Host: "actionHost",
				Port: intstr.FromString("22"),
			},
		},
		PreStop: &v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path: "/action/path",
				Host: "actionHost",
				Port: intstr.FromString("443"),
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "Key",
						Value: "Value",
					},
					{
						Name:  "HEADER",
						Value: "VALUE",
					},
				},
				Scheme: v1.URISchemeHTTPS,
			},
		},
	},
	TerminationMessagePath:   "/msg/path",
	TerminationMessagePolicy: v1.TerminationMessageReadFile,
	ImagePullPolicy:          v1.PullAlways,
	SecurityContext: &v1.SecurityContext{
		Capabilities: &v1.Capabilities{
			Add:  []v1.Capability{"CAP1", "CAP2"},
			Drop: []v1.Capability{"CAP3"},
		},
		Privileged: &boolEntry,
		SELinuxOptions: &v1.SELinuxOptions{
			Level: "selinuxlevel",
			Role:  "selinuxrole",
			Type:  "selinuxtype",
			User:  "selinuxuser",
		},
		RunAsUser:                &int64Val,
		RunAsGroup:               &int64Val,
		RunAsNonRoot:             &boolEntry,
		ReadOnlyRootFilesystem:   &boolEntry,
		AllowPrivilegeEscalation: &boolEntry,
		ProcMount:                &kubeProcMount,
	},
	Stdin:     boolEntry,
	StdinOnce: boolEntry,
	TTY:       boolEntry,
}

var emptyMantleContainer = container.Container{
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

var emptyKubeStatus = v1.PodStatus{
	Message: "",
	Reason:  "",
	HostIP:  "",
	PodIP:   "",
}

type kubePodConfig struct {
	basePod        v1.Pod
	volumes        []v1.Volume
	initContainers []v1.Container
	containers     []v1.Container
	context        *v1.PodSecurityContext
	secrets        []v1.LocalObjectReference
	affinity       *v1.Affinity
	tolerations    []v1.Toleration
	aliases        []v1.HostAlias
	dnsConfig      *v1.PodDNSConfig
	gates          []v1.PodReadinessGate
}

func (k *kubePodConfig) Generate() v1.Pod {
	pod := k.basePod

	if k.volumes != nil {
		pod.Spec.Volumes = k.volumes
	}

	if k.initContainers != nil {
		pod.Spec.InitContainers = k.initContainers
	}

	if k.containers != nil {
		pod.Spec.Containers = k.containers
	}

	if k.secrets != nil {
		pod.Spec.ImagePullSecrets = k.secrets
	}

	if k.context != nil {
		pod.Spec.SecurityContext = k.context
	}
	if k.affinity != nil {
		pod.Spec.Affinity = k.affinity
	}

	if k.tolerations != nil {
		pod.Spec.Tolerations = k.tolerations
	}

	if k.aliases != nil {
		pod.Spec.HostAliases = k.aliases
	}

	if k.dnsConfig != nil {
		pod.Spec.DNSConfig = k.dnsConfig
	}

	if k.gates != nil {
		pod.Spec.ReadinessGates = k.gates
	}

	return pod
}

type mantlePodConfig struct {
	basePod         Pod
	volumes         map[string]volume.Volume
	initContainers  []container.Container
	containers      []container.Container
	fsgid           *int64
	gids            []int64
	registries      []string
	affinity        *affinity.Affinity
	tolerations     []toleration.Toleration
	hostAliases     []hostalias.HostAlias
	nameservers     []string
	searchDomains   []string
	resolverOptions []podtemplate.ResolverOptions
	gates           []podtemplate.PodConditionType
}

func (m *mantlePodConfig) Generate() Pod {
	pod := m.basePod

	if m.volumes != nil {
		pod.Volumes = m.volumes
	}

	if m.initContainers != nil {
		pod.InitContainers = m.initContainers
	}

	if m.containers != nil {
		pod.Containers = m.containers
	}

	if m.fsgid != nil {
		pod.FSGID = m.fsgid
	}

	if m.gids != nil {
		pod.GIDs = m.gids
	}

	if m.registries != nil {
		pod.Registries = m.registries
	}

	if m.affinity != nil {
		pod.Affinity = m.affinity
	}

	if m.tolerations != nil {
		pod.Tolerations = m.tolerations
	}

	if m.hostAliases != nil {
		pod.HostAliases = m.hostAliases
	}

	if m.nameservers != nil {
		pod.Nameservers = m.nameservers
	}

	if m.searchDomains != nil {
		pod.SearchDomains = m.searchDomains
	}

	if m.resolverOptions != nil {
		pod.ResolverOptions = m.resolverOptions
	}

	if m.gates != nil {
		pod.Gates = m.gates
	}

	return pod
}

var kubePod = v1.Pod{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Pod",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "testPod",
		Namespace:   "testNS",
		ClusterName: "testCluster",
		Labels:      map[string]string{"label1": "value1"},
		Annotations: map[string]string{"annotation1": "value2"},
	},
	Spec: v1.PodSpec{
		Volumes:                       []v1.Volume{},
		RestartPolicy:                 v1.RestartPolicyAlways,
		TerminationGracePeriodSeconds: &int64Val,
		ActiveDeadlineSeconds:         &int64Val,
		DNSPolicy:                     v1.DNSDefault,
		NodeSelector:                  map[string]string{"label": "value"},
		ServiceAccountName:            "svcAccount",
		AutomountServiceAccountToken:  &boolEntry,
		NodeName:                      "testNode",
		HostNetwork:                   true,
		HostPID:                       true,
		HostIPC:                       true,
		ShareProcessNamespace:         &boolEntry,
		SecurityContext:               nil,
		Hostname:                      "testHostname",
		Subdomain:                     "testSubdomain.com",
		Affinity:                      &v1.Affinity{},
		//Affinity:          nil,
		SchedulerName:     "schedName",
		Tolerations:       []v1.Toleration{},
		HostAliases:       []v1.HostAlias{},
		PriorityClassName: "className",
		Priority:          &int32Val,
		DNSConfig:         &v1.PodDNSConfig{},
		//DNSConfig:          nil,
		ReadinessGates:     []v1.PodReadinessGate{},
		RuntimeClassName:   &strEntry,
		EnableServiceLinks: &boolEntry,
	},
}

var mantlePod = Pod{
	Version: "v1",
	Phase:   PodPhaseNone,
	QOS:     PodQOSClassNone,
	PodTemplateMeta: PodTemplateMeta{
		Name:        "testPod",
		Namespace:   "testNS",
		Cluster:     "testCluster",
		Labels:      map[string]string{"label1": "value1"},
		Annotations: map[string]string{"annotation1": "value2"},
	},
	PodTemplate: podtemplate.PodTemplate{
		Volumes:                map[string]volume.Volume{},
		RestartPolicy:          podtemplate.RestartPolicyAlways,
		TerminationGracePeriod: &int64Val,
		ActiveDeadline:         &int64Val,
		DNSPolicy:              podtemplate.DNSDefault,
		NodeSelector:           map[string]string{"label": "value"},
		Account:                "svcAccount",
		AutomountAccountToken:  &boolEntry,
		Node:                   "testNode",
		HostMode: []podtemplate.HostMode{
			podtemplate.HostModeNet,
			podtemplate.HostModePID,
			podtemplate.HostModeIPC,
		},
		ShareNamespace: &boolEntry,
		Hostname:       "testHostname.testSubdomain.com",
		Affinity:       &affinity.Affinity{},
		//Affinity:        nil,
		SchedulerName:   "schedName",
		Tolerations:     []toleration.Toleration{},
		HostAliases:     []hostalias.HostAlias{},
		PriorityClass:   "className",
		Priority:        &int32Val,
		ResolverOptions: []podtemplate.ResolverOptions{},
		Gates:           []podtemplate.PodConditionType{},
		RuntimeClass:    &strEntry,
		ServiceLinks:    &boolEntry,
	},
}

func TestPodFromKube(t *testing.T) {
	testCases := []struct {
		name     string
		original kubePodConfig
		expected mantlePodConfig
		fail     bool
	}{
		{
			name:     "empty kube v1 pod",
			original: kubePodConfig{basePod: v1.Pod{TypeMeta: metav1.TypeMeta{APIVersion: "v1"}}},
			expected: mantlePodConfig{basePod: Pod{Version: "v1", Phase: PodPhaseNone, QOS: PodQOSClassNone, PodTemplate: podtemplate.PodTemplate{DNSPolicy: podtemplate.DNSUnset}}},
			fail:     false,
		},
		{
			name: "all fields defined v1 pod",
			original: kubePodConfig{
				basePod:        kubePod,
				initContainers: []v1.Container{fullKubeContainer},
				containers:     []v1.Container{fullKubeContainer},
				secrets:        []v1.LocalObjectReference{{Name: strEntry}},
				context: &v1.PodSecurityContext{
					FSGroup:            &int64Val,
					SupplementalGroups: []int64{int64Val},
				},
				dnsConfig: &v1.PodDNSConfig{
					Nameservers: []string{strEntry},
					Searches:    []string{strEntry},
					Options: []v1.PodDNSConfigOption{
						{
							Name:  strEntry,
							Value: &strEntry,
						},
					},
				},
			},
			expected: mantlePodConfig{
				basePod:        mantlePod,
				containers:     []container.Container{fullMantleContainer},
				initContainers: []container.Container{fullMantleContainer},
				fsgid:          &int64Val,
				gids:           []int64{int64Val},
				registries:     []string{strEntry},
				nameservers:    []string{strEntry},
				searchDomains:  []string{strEntry},
				resolverOptions: []podtemplate.ResolverOptions{
					{
						Name:  strEntry,
						Value: &strEntry,
					},
				},
			},
			fail: false,
		},
		{
			name:     "invalid pod version",
			original: kubePodConfig{basePod: v1.Pod{TypeMeta: metav1.TypeMeta{APIVersion: "invalid"}}},
			expected: mantlePodConfig{},
			fail:     true,
		},
	}

	for _, tc := range testCases {
		pod, err := NewPodFromKubePod(tc.original.Generate())
		if err != nil && !tc.fail {
			t.Errorf("%s: error converting pod: %+v", tc.name, err)
		}

		if pod != nil {
			expectedPod := tc.expected.Generate()
			if !reflect.DeepEqual(*pod, expectedPod) {
				t.Errorf("%s: bad pod conversion from kube: %s", tc.name, cmp.Diff(*pod, expectedPod))
			}
		}
	}
}

func TestPodToKube(t *testing.T) {
	testCases := []struct {
		name     string
		original mantlePodConfig
		expected kubePodConfig
		fail     bool
	}{
		{
			name:     "empty pod with version v1",
			original: mantlePodConfig{basePod: Pod{Version: "v1", Phase: PodPhaseNone, QOS: PodQOSClassNone, PodTemplate: podtemplate.PodTemplate{DNSPolicy: podtemplate.DNSUnset}}},
			expected: kubePodConfig{basePod: v1.Pod{TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"}, Status: emptyKubeStatus}},
			fail:     false,
		},
		{
			name: "all fields defined pod",
			original: mantlePodConfig{
				basePod:        mantlePod,
				containers:     []container.Container{fullMantleContainer},
				initContainers: []container.Container{fullMantleContainer},
				fsgid:          &int64Val,
				gids:           []int64{int64Val},
				registries:     []string{strEntry},
				nameservers:    []string{strEntry},
				searchDomains:  []string{strEntry},
			},
			expected: kubePodConfig{basePod: kubePod},
			fail:     false,
		},
		{
			name:     "invalid pod version",
			original: mantlePodConfig{basePod: Pod{Version: "invalid"}},
			expected: kubePodConfig{},
			fail:     true,
		},
	}

	for _, tc := range testCases {
		p := tc.original.Generate()
		pod, err := p.ToKube()
		if err != nil && !tc.fail {
			t.Errorf("%s: error converting pod: %+v", tc.name, err)
		}

		if pod != nil {
			expectedPod := tc.expected.Generate()
			p := pod.(*v1.Pod)
			if !cmp.Equal(p, expectedPod) {
				t.Errorf("%s: bad pod conversion to kube: %s", tc.name, cmp.Diff(p, expectedPod))
			}
		}
	}
}

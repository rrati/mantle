package pod

import (
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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/go-cmp/cmp"
)

var int64Val = int64(5)
var boolEntry = true
var int32Val = int32(1)
var strEntry = "testStr"
var mountType = container.MountTypeDefault
var mountPropagation = volumemount.MountPropagationNone
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
	pod.Spec.Volumes = k.volumes
	pod.Spec.InitContainers = k.initContainers
	pod.Spec.Containers = k.containers
	pod.Spec.SecurityContext = k.context
	pod.Spec.Affinity = k.affinity
	pod.Spec.Tolerations = k.tolerations
	pod.Spec.HostAliases = k.aliases
	pod.Spec.DNSConfig = k.dnsConfig
	pod.Spec.ReadinessGates = k.gates

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
	pod.Volumes = m.volumes
	pod.InitContainers = m.initContainers
	pod.Containers = m.containers
	pod.FSGID = m.fsgid
	pod.GIDs = m.gids
	pod.Registries = m.registries
	pod.Affinity = m.affinity
	pod.Tolerations = m.tolerations
	pod.HostAliases = m.hostAliases
	pod.Nameservers = m.nameservers
	pod.SearchDomains = m.searchDomains
	pod.ResolverOptions = m.resolverOptions
	pod.Gates = m.gates

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
		Volumes:        []v1.Volume{},
		InitContainers: []v1.Container{},
		//		Containers:                    []v1.Container{},
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
		//SecurityContext:               &v1.PodSecurityContext{},
		SecurityContext:  nil,
		ImagePullSecrets: []v1.LocalObjectReference{},
		Hostname:         "testHostname",
		Subdomain:        "testSubdomain.com",
		Affinity:         &v1.Affinity{},
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
		Volumes:        map[string]volume.Volume{},
		InitContainers: []container.Container{},
		//		Containers:             []container.Container{},
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
		FSGID:          nil,
		GIDs:           nil,
		Registries:     []string{},
		Hostname:       "testHostname.testSubdomain.com",
		Affinity:       &affinity.Affinity{},
		//Affinity:        nil,
		SchedulerName:   "schedName",
		Tolerations:     []toleration.Toleration{},
		HostAliases:     []hostalias.HostAlias{},
		PriorityClass:   "className",
		Priority:        &int32Val,
		Nameservers:     []string{},
		SearchDomains:   []string{},
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
			name:     "all fields defined v1 pod",
			original: kubePodConfig{basePod: kubePod},
			expected: mantlePodConfig{basePod: mantlePod, containers: []container.Container{fullMantleContainer}},
			fail:     false,
		},
		{
			name:     "invalid pod version",
			original: kubePodConfig{basePod: v1.Pod{TypeMeta: metav1.TypeMeta{APIVersion: "invalid"}}},
			expected: nil,
			fail:     true,
		},
	}

	for _, tc := range testCases {
		pod, err := NewPodFromKubePod(tc.original.Generate())
		if err != nil && !tc.fail {
			t.Errorf("%s: error converting pod: %+v", tc.name, err)
		}

		if pod != nil {
			if !cmp.Equal(pod, tc.expected) {
				t.Errorf("%s: bad pod conversion from kube: %s", tc.name, cmp.Diff(*pod, tc.expected))
			}
		}
	}
}

func TestPodToKube(t *testing.T) {
	testCases := []struct {
		name     string
		original Pod
		expected *v1.Pod
		fail     bool
	}{
		{
			name:     "empty pod with version v1",
			original: mantlePodConfig{basePod: Pod{Version: "v1", Phase: PodPhaseNone, QOS: PodQOSClassNone, PodTemplate: podtemplate.PodTemplate{DNSPolicy: podtemplate.DNSUnset}}},
			expected: kubePodConfig{basePod: v1.Pod{TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"}}},
			fail:     false,
		},
		{
			name:     "all fields defined pod",
			original: mantlPodConfig{basePod: mantlePod, []container.Container{fullMantleContainer}},
			expected: kubePodConfig{basePod: kubePod},
			fail:     false,
		},
		{
			name:     "invalid pod version",
			original: Pod{Version: "invalid"},
			expected: nil,
			fail:     true,
		},
	}

	for _, tc := range testCases {
		pod, err := tc.original.Generate().ToKube()
		if err != nil && !tc.fail {
			t.Errorf("%s: error converting pod: %+v", tc.name, err)
		}

		if pod != nil {
			if !cmp.Equal(pod, tc.expected) {
				t.Errorf("%s: bad pod conversion to kube: %s", tc.name, cmp.Diff(pod, tc.expected))
			}
		}
	}
}

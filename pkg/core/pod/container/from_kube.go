package container

import (
	"fmt"
	"reflect"

	"mantle/pkg/core/action"
	"mantle/pkg/core/pod/container/env"
	"mantle/pkg/core/pod/container/port"
	"mantle/pkg/core/pod/container/probe"
	"mantle/pkg/core/pod/container/resources"
	"mantle/pkg/core/pod/container/volumedevice"
	"mantle/pkg/core/pod/container/volumemount"
	"mantle/pkg/core/selinux"
	"mantle/pkg/util/floatstr"

	"k8s.io/api/core/v1"

	serrors "github.com/koki/structurederrors"
)

// NewContainerFromKubeContainer will create a new Container object with
// the data from a provided kubernetes container object
func NewContainerFromKubeContainer(container interface{}) (*Container, error) {
	switch reflect.TypeOf(container) {
	case reflect.TypeOf(v1.Container{}):
		obj := container.(v1.Container)
		return fromKubeContainerV1(&obj)
	case reflect.TypeOf(&v1.Container{}):
		return fromKubeContainerV1(container.(*v1.Container))
	default:
		return nil, fmt.Errorf("unknown Container version: %s", reflect.TypeOf(container))
	}
}

func fromKubeContainerV1(container *v1.Container) (*Container, error) {
	mantleContainer := &Container{}

	mantleContainer.Name = container.Name
	mantleContainer.Command = container.Command
	mantleContainer.Image = container.Image
	mantleContainer.Args = fromKubeArgsV1(container.Args)
	mantleContainer.WorkingDir = container.WorkingDir
	mantleContainer.Pull = fromKubePullPolicyV1(container.ImagePullPolicy)

	onStart, preStop, err := fromKubeLifeCycleV1(container.Lifecycle)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "life cycle")
	}
	mantleContainer.OnStart = onStart
	mantleContainer.PreStop = preStop

	cpu, err := resources.NewCPUFromKubeResourceRequirements(container.Resources)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "cpu resource requirements")
	}
	mantleContainer.CPU = cpu

	mem, err := resources.NewMemFromKubeResourceRequirements(container.Resources)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "mem resource requirements")
	}
	mantleContainer.Mem = mem

	if container.SecurityContext != nil {
		mantleContainer.Privileged = container.SecurityContext.Privileged
		mantleContainer.AllowEscalation = container.SecurityContext.AllowPrivilegeEscalation
		mantleContainer.RO = container.SecurityContext.ReadOnlyRootFilesystem
		mantleContainer.ForceNonRoot = container.SecurityContext.RunAsNonRoot
		mantleContainer.UID = container.SecurityContext.RunAsUser
		mantleContainer.GID = container.SecurityContext.RunAsGroup

		sel, err := selinux.NewSELinuxFromKubeSELinuxOptions(container.SecurityContext.SELinuxOptions)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "selinux options")
		}
		mantleContainer.SELinux = sel

		mantleContainer.AddCapabilities, mantleContainer.DelCapabilities = fromKubeCapabilitiesV1(container.SecurityContext.Capabilities)

		if container.SecurityContext.ProcMount != nil {
			procMount, err := fromKubeProcMountV1(container.SecurityContext.ProcMount)
			if err != nil {
				return nil, serrors.ContextualizeErrorf(err, "proc mount")
			}
			mantleContainer.ProcMount = &procMount
		}
	}

	livenessProbe, err := probe.NewProbeFromKubeProbe(container.LivenessProbe)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "liveness probe")
	}
	mantleContainer.LivenessProbe = livenessProbe

	readinessProbe, err := probe.NewProbeFromKubeProbe(container.ReadinessProbe)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "readiness probe")
	}
	mantleContainer.ReadinessProbe = readinessProbe

	ports, err := fromKubeContainerPortsV1(container.Ports)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "container ports")
	}
	mantleContainer.Expose = ports

	mantleContainer.Stdin = container.Stdin
	mantleContainer.StdinOnce = container.StdinOnce
	mantleContainer.TTY = container.TTY
	mantleContainer.TerminationMsgPath = container.TerminationMessagePath
	mantleContainer.TerminationMsgPolicy = fromKubeTerminationMessagePolicyV1(container.TerminationMessagePolicy)

	envs, err := fromKubeEnvVarsV1(container.Env)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "env var")
	}

	envFroms, err := fromKubeEnvFromSourceV1(container.EnvFrom)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "env from source")
	}

	envs = append(envs, envFroms...)
	mantleContainer.Env = envs

	volumeMounts, err := fromKubeVolumeMountsV1(container.VolumeMounts)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "volume mount")
	}
	mantleContainer.VolumeMounts = volumeMounts

	volumeDevices, err := fromKubeVolumeDevicesV1(container.VolumeDevices)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "volume device")
	}
	mantleContainer.VolumeDevices = volumeDevices

	return mantleContainer, nil
}

func fromKubeArgsV1(kubeArgs []string) []floatstr.FloatOrString {
	if kubeArgs == nil {
		return nil
	}

	mantleArgs := make([]floatstr.FloatOrString, len(kubeArgs))
	for i, kubeArg := range kubeArgs {
		mantleArgs[i] = *floatstr.Parse(kubeArg)
	}

	return mantleArgs
}

func fromKubePullPolicyV1(pullPolicy v1.PullPolicy) PullPolicy {
	switch pullPolicy {
	case v1.PullAlways:
		return PullAlways

	case v1.PullNever:
		return PullNever

	case v1.PullIfNotPresent:
		return PullIfNotPresent

	default:
		return PullDefault
	}
}

func fromKubeLifeCycleV1(lifecycle *v1.Lifecycle) (*action.Action, *action.Action, error) {
	if lifecycle == nil {
		return nil, nil, nil
	}

	actionOnStart, err := action.NewActionFromKubeHandler(lifecycle.PostStart)
	if err != nil {
		return nil, nil, err
	}

	actionPreStop, err := action.NewActionFromKubeHandler(lifecycle.PreStop)
	if err != nil {
		return nil, nil, err
	}

	return actionOnStart, actionPreStop, nil
}

func fromKubeCapabilitiesV1(caps *v1.Capabilities) ([]string, []string) {
	if caps == nil {
		return nil, nil
	}

	var addCaps []string
	var delCaps []string

	if caps.Add != nil {
		for _, add := range caps.Add {
			addCaps = append(addCaps, string(add))
		}
	}

	if caps.Drop != nil {
		for _, del := range caps.Drop {
			delCaps = append(delCaps, string(del))
		}
	}

	return addCaps, delCaps
}

func fromKubeContainerPortsV1(ports []v1.ContainerPort) ([]port.Port, error) {
	if ports == nil {
		return nil, nil
	}

	var portList []port.Port
	for _, p := range ports {
		mantlePort, err := port.NewPortFromKubeContainerPort(p)
		if err != nil {
			return nil, err
		}
		portList = append(portList, *mantlePort)
	}

	return portList, nil
}

func fromKubeTerminationMessagePolicyV1(p v1.TerminationMessagePolicy) TerminationMessagePolicy {
	if p == v1.TerminationMessageReadFile {
		return TerminationMessageReadFile
	}

	if p == v1.TerminationMessageFallbackToLogsOnError {
		return TerminationMessageFallbackToLogsOnError
	}

	return TerminationMessageDefault
}

func fromKubeEnvVarsV1(kubeEnvs []v1.EnvVar) ([]env.Env, error) {
	var envs []env.Env

	for _, e := range kubeEnvs {
		new, err := env.NewEnvFromKubeEnvVar(e)
		if err != nil {
			return nil, err
		}
		envs = append(envs, *new)
	}

	return envs, nil
}

func fromKubeEnvFromSourceV1(kubeEnvFromSource []v1.EnvFromSource) ([]env.Env, error) {
	var envs []env.Env

	for _, e := range kubeEnvFromSource {
		new, err := env.NewEnvFromKubeEnvFromSource(e)
		if err != nil {
			return nil, err
		}
		envs = append(envs, *new)
	}

	return envs, nil
}

func fromKubeVolumeMountsV1(kubeMounts []v1.VolumeMount) ([]volumemount.VolumeMount, error) {
	mounts := []volumemount.VolumeMount{}

	for _, mount := range kubeMounts {
		m, err := volumemount.NewVolumeMountFromKubeVolumeMount(mount)
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, *m)
	}

	return mounts, nil
}

func fromKubeVolumeDevicesV1(kubeDevices []v1.VolumeDevice) ([]volumedevice.VolumeDevice, error) {
	devices := []volumedevice.VolumeDevice{}

	for _, device := range kubeDevices {
		d, err := volumedevice.NewVolumeDeviceFromKubeVolumeDevice(device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, *d)
	}

	return devices, nil
}

func fromKubeProcMountV1(procMount *v1.ProcMountType) (MountType, error) {
	switch *procMount {
	case v1.DefaultProcMount:
		return MountTypeDefault, nil
	case v1.UnmaskedProcMount:
		return MountTypeUnmasked, nil
	default:
		return MountTypeInvalid, fmt.Errorf("unknown ProcMountType: %s", *procMount)
	}
}

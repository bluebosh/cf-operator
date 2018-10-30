package deployment_manifest

type InstanceGroup struct {
	Name      string
	Instances int
}
type Manifest struct {
	InstanceGroups []InstanceGroup `yaml:"instance-groups"`
}

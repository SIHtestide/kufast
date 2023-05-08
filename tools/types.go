package tools

// Target represents a deployment target and contains its name and the type of access (either group or node)
type Target struct {
	Name       string
	AccessType string
}

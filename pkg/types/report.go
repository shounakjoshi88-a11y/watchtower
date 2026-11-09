package types

// Report contains reports for all the containers processed during a session
type Report interface {
	Scanned() []ContainerReport
	Updated() []ContainerReport
	Failed() []ContainerReport
	Skipped() []ContainerReport
	Stale() []ContainerReport
	Fresh() []ContainerReport
	All() []ContainerReport
	// AllUpdatedNoNotify reports whether every container that was updated in
	// this session carried the no-notify label while nothing failed or was
	// skipped, meaning the session notification can be suppressed entirely
	AllUpdatedNoNotify() bool
}

// ContainerReport represents a container that was included in watchtower session
type ContainerReport interface {
	ID() ContainerID
	Name() string
	CurrentImageID() ImageID
	LatestImageID() ImageID
	ImageName() string
	Error() string
	State() string
}

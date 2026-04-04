package dto

type (
	// GroupingContainer - comment struct.
	GroupingContainer struct {
		ID   uint64
		Code string
		// Tags   []string
		Images []string
	}

	// UpdateGroupContainer - comment struct.
	UpdateGroupContainer struct {
		ID     uint64
		Tags   []string
		Images []string
	}
)

package skizze

import (
	pb "github.com/skizzehq/goskizze/datamodel"
)

// Defaults are the default settings for newly created Sketches.
type Defaults struct {
	Rank     int64
	Capacity int64
}

func newDefaultsFromRaw(d *pb.Defaults) *Defaults {
	return &Defaults{
		Rank:     d.GetRank(),
		Capacity: d.GetCapacity(),
	}
}

func getRawDefaultsFromDefaults(d *Defaults) *pb.Defaults {
	return &pb.Defaults{
		Rank:     &d.Rank,
		Capacity: &d.Capacity,
	}
}

package skizze

import (
	pb "github.com/skizzehq/goskizze/datamodel"
)

// Domain describes details of a Skizze domain
type Domain struct {
	Name     string
	Sketches []*Sketch
}

// DomainProperties details configuration for Sketches in a domain
type DomainProperties struct {
	MembershipProperties Properties
	FrequencyProperties  Properties
	RankingsProperties   Properties
}

func newDomainFromRaw(d *pb.Domain) *Domain {
	ret := &Domain{}
	ret.Name = d.GetName()

	for _, s := range d.GetSketches() {
		ret.Sketches = append(ret.Sketches, newSketchFromRaw(s))
	}

	return ret
}

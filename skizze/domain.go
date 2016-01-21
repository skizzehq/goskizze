package skizze

import (
	pb "github.com/skizzehq/goskizze/datamodel"
)

// Domain describes details of a Skizze domain
type Domain struct {
	Name     string
	Defaults *Defaults
	Sketches []*Sketch
}

func newDomainFromRaw(d *pb.Domain) *Domain {
	ret := &Domain{}
	ret.Name = d.GetName()
	ret.Defaults = newDefaultsFromRaw(d.GetDefaults())

	for _, s := range d.GetSketches() {
		ret.Sketches = append(ret.Sketches, newSketchFromRaw(s))
	}

	return ret
}

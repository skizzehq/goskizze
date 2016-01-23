package skizze

import (
	pb "github.com/skizzehq/goskizze/datamodel"
)

var (
	defaultMembUnique    int64   = 1000000
	defaultMembErrorRate float32 = 0.01
	defaultFreqUnique    int64   = 100000
	defaultFreqErrorRate float32 = 0.01
	defaultRankSize      int64   = 101
)

// Properties are configuration settings for a Sketch.
// Size is used by Rankings Sketches to determine the number of rankings this Sketch should track e.g. top 10, top 100, top 1000
type Properties struct {
	MaxUniqueItems int64
	ErrorRate      float32
	Size           int64
}

func newPropertiesFromRaw(r *pb.SketchProperties) *Properties {
	if r == nil {
		return nil
	}
	return &Properties{
		MaxUniqueItems: r.GetMaxUniqueItems(),
		ErrorRate:      r.GetErrorRate(),
		Size:           r.GetSize(),
	}
}

func newRawPropertiesFromProperties(p *Properties) *pb.SketchProperties {
	if p == nil {
		return nil
	}
	return &pb.SketchProperties{
		MaxUniqueItems: &p.MaxUniqueItems,
		ErrorRate:      &p.ErrorRate,
		Size:           &p.Size,
	}
}

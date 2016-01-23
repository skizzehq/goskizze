package skizze

import (
	"log"

	pb "github.com/skizzehq/goskizze/datamodel"
)

// SketchType is the kind of sketch
type SketchType int

const (
	// Membership sketches can evaluate if a data set contains an element.
	Membership SketchType = iota
	// Frequency sketches can evaluate the frequency of an element in a data set.
	Frequency
	// Ranking sketches return the top ranked elements of the data set.
	Ranking
	// Cardinality sketches evaluate how many distinct elements are in the data set.
	Cardinality
)

// Sketch describes the details of a sketch
type Sketch struct {
	Name       string
	Type       SketchType
	Properties *Properties
}

func newSketchFromRaw(s *pb.Sketch) *Sketch {
	if s == nil {
		return nil
	}
	return &Sketch{
		Name:       s.GetName(),
		Type:       getSketchTypeForRawType(s.GetType()),
		Properties: newPropertiesFromRaw(s.GetProperties()),
	}
}

func getRawSketchFromSketch(s Sketch) *pb.Sketch {
	t := getRawSketchForSketchType(s.Type)
	return &pb.Sketch{
		Name:       &s.Name,
		Type:       &t,
		Properties: newRawPropertiesFromProperties(s.Properties),
	}
}

func getSketchTypeForRawType(t pb.SketchType) SketchType {
	switch t {
	case pb.SketchType_MEMB:
		return Membership
	case pb.SketchType_FREQ:
		return Frequency
	case pb.SketchType_RANK:
		return Ranking
	case pb.SketchType_CARD:
		return Cardinality
	default:
		log.Panicf("SketchType %v unknown", t)
		return 0
	}
}

func getRawSketchForSketchType(t SketchType) pb.SketchType {
	switch t {
	case Membership:
		return pb.SketchType_MEMB
	case Frequency:
		return pb.SketchType_FREQ
	case Ranking:
		return pb.SketchType_RANK
	case Cardinality:
		return pb.SketchType_CARD
	default:
		log.Panicf("SketchType %v unknown", t)
		return 0
	}
}

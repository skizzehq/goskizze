package skizze

// MembershipResult indicates the result of a membership query for a value.
type MembershipResult struct {
	Value    string
	IsMember bool
}

// FrequencyResult indicates the result of a frequency query for a value.
type FrequencyResult struct {
	Value string
	Count int64
}

// RankingsResult indicates the result of a ranking query for a value.
type RankingsResult struct {
	Value string
	Count int64
}

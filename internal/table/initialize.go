package table

// AllInitials stores all the possible initials
var AllInitials []string = []string{}

// AllFinals stores all the possible finals
var AllFinals []string = []string{}

// initializes the InitialsAndFinals
func init() {
	rawInitials := []string{}
	for initial := range InitialToKana {
		rawInitials = append(
			rawInitials,
			initial,
		)
	}
	AllInitials = SortStringSlice(rawInitials)
	rawFinals := []string{}
	for final := range FinalToKana {
		rawFinals = append(
			rawFinals,
			final,
		)
	}
	AllFinals = SortStringSlice(rawFinals)
}

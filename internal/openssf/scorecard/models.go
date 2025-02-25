package scorecard

import "time"

type Scorecard struct {
	// Date is the time at which the scorecard report was generated.
	// Could be a date string 2006-01-02, or a RFC 3339-formatted time string.
	Date string `json:"date"`
	// Score is a value 0-10.
	Score float32 `json:"score"`
}

// Time returns the time at which the scorecard report was generated.
// Returns an error if the time is invalid or unsupported.
func (s *Scorecard) Time() (time.Time, error) {
	if len(s.Date) == 10 {
		return time.Parse("2006-01-02", s.Date)
	}

	return time.Parse(time.RFC3339, s.Date)

}

package scorecard

import "time"

type Scorecard struct {
	Date  string  `json:"date"`
	Score float32 `json:"score"`
}

func (s *Scorecard) Time() (time.Time, error) {
	if len(s.Date) == 10 {
		return time.Parse("2006-01-02", s.Date)
	}

	return time.Parse(time.RFC3339, s.Date)

}

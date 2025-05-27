package factories

//
// import (
// 	"github.com/hengadev/leviosa/internal/domain/vote"
// )
//
// func NewBasicVote(overrides map[string]any) *vote.Vote {
// 	userID := NewBasicUser(nil).ID
// 	vote := &vote.Vote{
// 		UserID: userID,
// 		Day:    23,
// 		Month:  4,
// 		Year:   2025,
// 	}
// 	for key, value := range overrides {
// 		switch {
// 		case key == "UserID":
// 			vote.UserID = value.(string)
// 		case key == "Day":
// 			vote.Day = value.(int)
// 		case key == "Month":
// 			vote.Month = value.(int)
// 		case key == "Year":
// 			vote.Year = value.(int)
// 		}
// 	}
// 	return vote
// }
// func NewAvailableVotesList() []*vote.AvailableVote {
// 	const YEAR = 2025
// 	return []*vote.AvailableVote{
// 		{Day: 23, Month: 4, Year: YEAR},
// 		{Day: 12, Month: 4, Year: YEAR},
// 		{Day: 3, Month: 4, Year: YEAR},
// 		{Day: 9, Month: 4, Year: YEAR},
// 		{Day: 17, Month: 4, Year: YEAR},
// 		{Day: 12, Month: 5, Year: YEAR},
// 		{Day: 18, Month: 5, Year: YEAR},
// 		{Day: 5, Month: 5, Year: YEAR},
// 		{Day: 7, Month: 6, Year: YEAR},
// 		{Day: 21, Month: 6, Year: YEAR},
// 		{Day: 18, Month: 7, Year: YEAR},
// 		{Day: 2, Month: 7, Year: YEAR},
// 		{Day: 30, Month: 7, Year: YEAR},
// 	}
// }

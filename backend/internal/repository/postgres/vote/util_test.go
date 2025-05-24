package voteRepository_test

import (
	"github.com/hengadev/leviosa/internal/domain/vote"
)

var (
	baseVote = &vote.Vote{
		UserID: "1",
		Day:    23,
		Month:  4,
		Year:   2025,
	}
	baseVote2 = &vote.Vote{
		UserID: "1",
		Day:    12,
		Month:  4,
		Year:   2025,
	}
	baseVote3 = &vote.Vote{
		UserID: "1",
		Day:    6,
		Month:  4,
		Year:   2025,
	}
)

const YEAR = 2025

var availableVotesArr = []*vote.AvailableVote{
	{Day: 23, Month: 4, Year: YEAR},
	{Day: 12, Month: 4, Year: YEAR},
	{Day: 3, Month: 4, Year: YEAR},
	{Day: 9, Month: 4, Year: YEAR},
	{Day: 17, Month: 4, Year: YEAR},
	{Day: 12, Month: 5, Year: YEAR},
	{Day: 18, Month: 5, Year: YEAR},
	{Day: 5, Month: 5, Year: YEAR},
	{Day: 7, Month: 6, Year: YEAR},
	{Day: 21, Month: 6, Year: YEAR},
	{Day: 18, Month: 7, Year: YEAR},
	{Day: 2, Month: 7, Year: YEAR},
	{Day: 30, Month: 7, Year: YEAR},
}

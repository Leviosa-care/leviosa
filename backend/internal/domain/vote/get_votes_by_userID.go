package vote

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hengadev/errsx"
)

// Function that returns the votes (order is important) for a specific user
func (s *service) GetVotesByUserID(ctx context.Context, monthStr, yearStr string, userID string) ([]*Vote, error) {
	monthInt, err := strconv.Atoi(monthStr)
	if err != nil {
		return nil, fmt.Errorf("fail to convert string month to int")
	}
	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, fmt.Errorf("fail to convert string year to int")
	}
	votesStr, err := s.Repo.FindVotesByUserID(ctx, monthStr, yearInt, userID)
	if err != nil {
		return nil, fmt.Errorf("get votes by userID: %w", err)
	}
	votes, err := parseVotes(ctx, votesStr, monthInt, yearInt)
	if err != nil {
		return nil, fmt.Errorf("parse votes by userID: %w", err)
	}
	return votes, nil
}

// vote du mois, the table is going to be vote_january_2024
// userID - someformatted vote thing

// two tables are made
// votes [month-year-availabledates]
// votes_april_2024

// parseVotes parses string stored in database into votes.
func parseVotes(ctx context.Context, daysStr string, month, year int) ([]*Vote, error) {
	var errs errsx.Map
	if daysStr == "" {
		return nil, nil
	}
	days := strings.Split(daysStr, VoteSeparator)
	var votes = make([]*Vote, len(days))
	for i, day := range days {
		day, err := strconv.Atoi(day)
		if err != nil {
			errs.Set("convert string day to int", err)
		}
		vote := &Vote{Day: day, Month: month, Year: year}
		if err := vote.Valid(ctx); err != nil {
			errs.Set("vote validation", err)
		}
		votes[i] = vote
	}
	return votes, errs.AsError()
}

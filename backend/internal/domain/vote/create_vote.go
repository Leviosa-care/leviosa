package vote

import (
	"context"
	"fmt"
	"strings"

	"github.com/hengadev/leviosa/internal/domain"
)

func (s *service) CreateVote(ctx context.Context, votes []*Vote) error {
	for _, vote := range votes {
		if err := vote.Valid(ctx); err != nil {
			return domain.NewInvalidValueErr(err.Error())
		}
	}

	month := votes[0].Month
	year := votes[0].Year
	userID := votes[0].UserID

	days := formatVote(votes)
	// check if has vote
	err := s.Repo.HasVote(ctx, month, year, userID)
	if err != nil {
		// remove previous vote
		if err := s.Repo.RemoveVote(ctx, userID, month, year); err != nil {
			return fmt.Errorf("remove existing vote: %w", err)
		}
	}
	// create new vote with the new information
	if err := s.Repo.CreateVote(ctx, userID, days, month, year); err != nil {
		return fmt.Errorf("add vote: %w", err)
	}
	return nil
}

// formatVote takes an array of Vote type to return formatted days respecting the order of Votes (indicating user's preference) to write in database. The separator string is a package constant.
func formatVote(votes []*Vote) string {
	var daysArr = make([]string, len(votes))
	for i, vote := range votes {
		daysArr[i] = fmt.Sprintf("%d", vote.Day)
	}
	days := strings.Join(daysArr, VoteSeparator)
	return days
}

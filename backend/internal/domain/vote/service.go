package vote

import "context"

type Service interface {
	CreateVote(ctx context.Context, votes []*Vote) error
	GetVotesByUserID(ctx context.Context, monthStr, yearStr string, userID string) ([]*Vote, error)
}

type service struct {
	Repo ReadWriter
}

func New(repo ReadWriter) Service {
	return &service{Repo: repo}
}

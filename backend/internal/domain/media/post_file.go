package mediaService

import (
	"context"
	"fmt"
	"mime/multipart"
)

// NOTE: make that thing better brother

// TODO: make a domain function to get the name of the folder to place the thing

func (s *service) PostFile(ctx context.Context, file multipart.File, filename, folder string) (string, error) {
	// NOTE: Here the best thing to do is to change the value of that thing brother
	// key := fmt.Sprintf("%s/%s", eventID, filename)
	key := fmt.Sprintf("events_%s/%s", folder, filename)
	_ = key
	// url, err := s.Repo.AddFile(ctx, file, key)
	// if err != nil {
	// 	return "", fmt.Errorf("add file: %w", err)
	// }
	return "", nil
}

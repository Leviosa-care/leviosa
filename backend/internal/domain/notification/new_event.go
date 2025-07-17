package notification

import (
	"context"
	"sync"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/errsx"
)

// Send an email to all users notifying new event creation.
func (s *mailService) NewEvent(ctx context.Context, users []*models.User, eventTime string) error {
	var errs errsx.Map
	var wg sync.WaitGroup
	var errMutex sync.Mutex

	for _, user := range users {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			// this is just to test of both clients
			emails := []string{user.Email, "henry.gary@hotmail.com"}
			// here I just test with oulook since it does not work
			// emails := []string{"henry.gary@hotmail.com"}
			templData := struct {
				Firstname string
				Heure     string
			}{
				Firstname: user.FirstName,
				Heure:     eventTime,
			}
			for _, email := range emails {
				if err := s.sendMail(
					ctx,
					email,
					"Nouvel Évènement disponible",
					"event_notification",
					templData,
					nil,
					nil,
				); err != nil {
					errMutex.Lock()
					errs.Set("send mail", err)
					errMutex.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	return errs.AsError()
}

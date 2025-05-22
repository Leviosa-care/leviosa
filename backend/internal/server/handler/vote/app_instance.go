package vote

import (
	"github.com/hengadev/leviosa/internal/server/app"
)

type AppInstance struct {
	*app.App
}

func NewHandler(appCtx *app.App) *AppInstance {
	return &AppInstance{appCtx}
}

// old api

// // TODO: change that function, the vote should be submitted to the Vote[Month][Year] table.
// func OldMakeVote() http.Handler {
// 	// use the sessionID from the application
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx, cancel := context.WithCancel(r.Context())
// 		defer cancel()
// 		// get userID from context
// 		userID := ctx.Value(mw.SessionIDKey)
// 		// TODO: change that type brother
// 		// get the votes from the request body
// 		reqBody := []*types.VoteSent{}
// 		votes, err := jsonio.Decode[*types.VoteSent](r)
// 		if err != nil {
// 			slog.ErrorContext(ctx, "failed to decode the votes from request body", "error", err)
// 			jsonio.WriteResponse(w, handler.NewInternalErr(err), http.StatusBadRequest)
// 			return
// 		}
//
// 		if s.Store.HasVote(reqBody[0].Month, reqBody[0].Year, userID) { // get the month and the year in this
// 			err := s.Store.RemoveVoteByUserID(reqBody[0].Month, reqBody[0].Year, userId)
// 			if err != nil {
// 				fmt.Println("failed to remove the vote ! - ", err)
// 				WriteResponse(w, fmt.Sprintf("Failed to remove previous vote from the database- %s", err), http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 		// TODO: That thing is not working, why ? (before I use the writeResponse at the end.)
// 		err = s.Store.CreateVoteByUserId(reqBody, userId)
// 		if err != nil {
// 			WriteResponse(w, fmt.Sprintf("Failed to create the vote in the database - %s", err), http.StatusInternalServerError)
// 			return
// 		}
// 		WriteResponse(w, fmt.Sprintf("Vote created successfully"), http.StatusCreated)
// 	})
// }

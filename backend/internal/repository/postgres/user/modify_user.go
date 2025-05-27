package userRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/pkg/sqliteutil"
)

// ModifyAccount modifies the user account information based on the provided fields.
// It constructs an SQL update query, executes it, and handles the result. The function allows certain fields to be excluded from modification through the `prohibitedFields` argument.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the database query.
//   - user: The user struct containing the updated information to be saved in the database.
//   - whereMap: A map of conditions used to identify which user record to update.
//   - prohibitedFields: A list of fields that should not be modified during the update.
//
// Returns:
//   - error: Returns an error if the query fails, the user record cannot be updated, or if any issues arise during execution.
func (u *repository) ModifyAccount(
	ctx context.Context,
	user *models.User,
	whereMap map[string]any,
) error {
	// TODO: in the domain, change the table to use the QualifiedTable function
	query, values, writeUpdateErr := sqliteutil.WriteUpdateQuery(*user, whereMap)
	fmt.Println("the query is:", query)
	result, err := u.DB.ExecContext(ctx, query, values...)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotUpdatedErr(errors.New("no rows affected"), fmt.Sprintf("user with ID %q", user.ID))
	}
	if writeUpdateErr != nil {
		return rp.NewValidationErr(err, "'models.User' fields for modification")
	}
	return nil
}

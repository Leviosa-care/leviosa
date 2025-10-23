package partnerRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreatePartner(t *testing.T) {
	ctx := context.Background()
	t.Run("should successfully create a new partner", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		partner := NewTestPartnerEncx()
		err := repo.CreatePartner(ctx, partner)
		assert.NoError(t, err)
		// TODO: then check with a select if the parnter is created
	})
}

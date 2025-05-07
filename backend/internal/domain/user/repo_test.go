package userService_test

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

type MockRepo struct {
	FindAccountByIDFunc          func(ctx context.Context, id string) (*models.User, error)
	GetHashedPasswordByEmailFunc func(ctx context.Context, email string) (string, error)
	GetUserFromEmailHashFunc     func(ctx context.Context, emailHash string) (*models.User, error)
	GetUnverifiedUserFunc        func(ctx context.Context, emailHash string) (*models.User, error)
	GetPendingUserFunc           func(ctx context.Context, emailHash string, provider models.ProviderType) (*models.User, error)
	GetPendingUsersFunc          func(ctx context.Context) ([]*models.User, error)
	GetOAuthUserFunc             func(ctx context.Context, email, provider string) (*models.User, error)
	GetUserSessionDataFunc       func(ctx context.Context, email string) (string, models.Role, error)
	AddUserFunc                  func(ctx context.Context, user *models.User, provider models.ProviderType) error
	AddPendingUserFunc           func(ctx context.Context, user *models.User, provider models.ProviderType) error
	AddUnverifiedUserFunc        func(ctx context.Context, user *models.User) error
	ModifyAccountFunc            func(ctx context.Context, user *models.User, whereMap map[string]any) error
	DeleteUserFunc               func(ctx context.Context, id string) error
	HasUserFunc                  func(ctx context.Context, emailHash string) error
	HasOAuthUserFunc             func(ctx context.Context, emailHash string, provider models.ProviderType) error
}

func (m *MockRepo) FindAccountByID(ctx context.Context, id string) (*models.User, error) {
	if m.FindAccountByIDFunc != nil {
		return m.FindAccountByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepo) GetHashedPasswordByEmail(ctx context.Context, email string) (string, error) {
	if m.GetHashedPasswordByEmailFunc != nil {
		return m.GetHashedPasswordByEmailFunc(ctx, email)
	}
	return "", nil
}

func (m *MockRepo) GetOAuthUser(ctx context.Context, email, provider string) (*models.User, error) {
	if m.GetOAuthUserFunc != nil {
		return m.GetOAuthUserFunc(ctx, email, provider)
	}
	return nil, nil
}

func (m *MockRepo) GetUserSessionData(ctx context.Context, email string) (string, models.Role, error) {
	if m.GetUserSessionDataFunc != nil {
		return m.GetUserSessionDataFunc(ctx, email)
	}
	return "", models.VISITOR, nil
}
func (m *MockRepo) AddUser(ctx context.Context, user *models.User, provider models.ProviderType) error {
	if m.AddUserFunc != nil {
		return m.AddUserFunc(ctx, user, provider)
	}
	return nil
}
func (m *MockRepo) ModifyAccount(ctx context.Context, user *models.User, whereMap map[string]any) error {
	if m.ModifyAccountFunc != nil {
		return m.ModifyAccountFunc(ctx, user, whereMap)
	}
	return nil
}
func (m *MockRepo) DeleteUser(ctx context.Context, id string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *MockRepo) GetUserFromEmailHash(ctx context.Context, emailHash string) (*models.User, error) {
	if m.GetUserFromEmailHashFunc != nil {
		return m.GetUserFromEmailHashFunc(ctx, emailHash)
	}
	return nil, nil
}

func (m *MockRepo) GetPendingUsers(ctx context.Context) ([]*models.User, error) {
	if m.GetPendingUsersFunc != nil {
		return m.GetPendingUsersFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepo) AddPendingUser(ctx context.Context, user *models.User, provider models.ProviderType) error {
	if m.AddPendingUserFunc != nil {
		return m.AddPendingUserFunc(ctx, user, provider)
	}
	return nil
}

func (m *MockRepo) AddUnverifiedUser(ctx context.Context, user *models.User) error {
	if m.AddUnverifiedUserFunc != nil {
		return m.AddUnverifiedUserFunc(ctx, user)
	}
	return nil
}

func (m *MockRepo) GetPendingUser(ctx context.Context, emailHash string, provider models.ProviderType) (*models.User, error) {
	if m.GetPendingUserFunc != nil {
		return m.GetPendingUserFunc(ctx, emailHash, provider)
	}
	return nil, nil
}

func (m *MockRepo) GetUnverifiedUser(ctx context.Context, emailHash string) (*models.User, error) {
	if m.GetUnverifiedUserFunc != nil {
		return m.GetUnverifiedUserFunc(ctx, emailHash)
	}
	return nil, nil
}

func (m *MockRepo) HasUser(ctx context.Context, emailHash string) error {
	if m.HasUserFunc != nil {
		return m.HasUserFunc(ctx, emailHash)
	}
	return nil
}

func (m *MockRepo) HasOAuthUser(ctx context.Context, emailHash string, provider models.ProviderType) error {
	if m.HasOAuthUserFunc != nil {
		return m.HasOAuthUserFunc(ctx, emailHash, provider)
	}
	return nil
}

package factories

import (
	"context"
	"io"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/mock"
)

// MockCryptoService is a mock type for the encx.CryptoService interface
type MockCryptoService struct {
	mock.Mock
}

// GetPepper provides a mocked function with given fields:
func (_m *MockCryptoService) GetPepper() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// GetArgon2Params provides a mocked function with given fields:
func (_m *MockCryptoService) GetArgon2Params() *encx.Argon2Params {
	ret := _m.Called()

	var r0 *encx.Argon2Params
	if rf, ok := ret.Get(0).(func() *encx.Argon2Params); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*encx.Argon2Params)
		}
	}

	return r0
}

// GetAlias provides a mocked function with given fields:
func (_m *MockCryptoService) GetAlias() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GenerateDEK provides a mocked function with given fields:
func (_m *MockCryptoService) GenerateDEK() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		r0, r1 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptData provides a mocked function with given fields: ctx, plaintext, dek
func (_m *MockCryptoService) EncryptData(ctx context.Context, plaintext []byte, dek []byte) ([]byte, error) {
	ret := _m.Called(ctx, plaintext, dek)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte) ([]byte, error)); ok {
		r0, r1 = rf(ctx, plaintext, dek)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecryptData provides a mocked function with given fields: ctx, ciphertext, dek
func (_m *MockCryptoService) DecryptData(ctx context.Context, ciphertext []byte, dek []byte) ([]byte, error) {
	ret := _m.Called(ctx, ciphertext, dek)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte) ([]byte, error)); ok {
		r0, r1 = rf(ctx, ciphertext, dek)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessStruct provides a mocked function with given fields: ctx, object
func (_m *MockCryptoService) ProcessStruct(ctx context.Context, object any) error {
	ret := _m.Called(ctx, object)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, any) error); ok {
		r0 = rf(ctx, object)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecryptStruct provides a mocked function with given fields: ctx, object
func (_m *MockCryptoService) DecryptStruct(ctx context.Context, object any) error {
	ret := _m.Called(ctx, object)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, any) error); ok {
		r0 = rf(ctx, object)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EncryptDEK provides a mocked function with given fields: plaintextDEK
func (_m *MockCryptoService) EncryptDEK(ctx context.Context, plaintextDEK []byte) ([]byte, error) {
	ret := _m.Called(ctx, plaintextDEK)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) ([]byte, error)); ok {
		r0, r1 = rf(ctx, plaintextDEK)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecryptDEKWithVersion provides a mocked function with given fields: ctx, ciphertextDEK, kekVersion
func (_m *MockCryptoService) DecryptDEKWithVersion(ctx context.Context, ciphertextDEK []byte, kekVersion int) ([]byte, error) {
	ret := _m.Called(ctx, ciphertextDEK, kekVersion)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, int) ([]byte, error)); ok {
		r0, r1 = rf(ctx, ciphertextDEK, kekVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RotateKEK provides a mocked function with given fields: ctx
func (_m *MockCryptoService) RotateKEK(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HashBasic provides a mocked function with given fields: ctx, value
func (_m *MockCryptoService) HashBasic(ctx context.Context, value []byte) string {
	ret := _m.Called(ctx, value)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, []byte) string); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// HashSecure provides a mocked function with given fields: ctx, value
func (_m *MockCryptoService) HashSecure(ctx context.Context, value []byte) (string, error) {
	ret := _m.Called(ctx, value)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) (string, error)); ok {
		r0, r1 = rf(ctx, value)
	} else {
		r0 = ret.Get(0).(string)
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CompareSecureHashAndValue provides a mocked function with given fields: ctx, value, hashValue
func (_m *MockCryptoService) CompareSecureHashAndValue(ctx context.Context, value any, hashValue string) (bool, error) {
	ret := _m.Called(ctx, value, hashValue)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, any, string) (bool, error)); ok {
		r0, r1 = rf(ctx, value, hashValue)
	} else {
		r0 = ret.Get(0).(bool)
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CompareBasicHashAndValue provides a mocked function with given fields: ctx, value, hashValue
func (_m *MockCryptoService) CompareBasicHashAndValue(ctx context.Context, value any, hashValue string) (bool, error) {
	ret := _m.Called(ctx, value, hashValue)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, any, string) (bool, error)); ok {
		r0, _ = rf(ctx, value, hashValue)
	} else {
		r0 = ret.Get(0).(bool)
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptStream provides a mocked function with given fields: ctx, reader, writer, dek
func (_m *MockCryptoService) EncryptStream(ctx context.Context, reader io.Reader, writer io.Writer, dek []byte) error {
	ret := _m.Called(ctx, reader, writer, dek)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader, io.Writer, []byte) error); ok {
		r0 = rf(ctx, reader, writer, dek)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecryptStream provides a mocked function with given fields: ctx, reader, writer, dek
func (_m *MockCryptoService) DecryptStream(ctx context.Context, reader io.Reader, writer io.Writer, dek []byte) error {
	ret := _m.Called(ctx, reader, writer, dek)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader, io.Writer, []byte) error); ok {
		r0 = rf(ctx, reader, writer, dek)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AssertExpectations is a helper method that asserts that all the expectations set on the mock were met.
func (_m *MockCryptoService) AssertExpectations(t mock.TestingT) {
	_m.Mock.AssertExpectations(t)
}

// AssertNumberOfCalls is a helper method that asserts that the method with the specified name was called the specified number of times.
func (_m *MockCryptoService) AssertNumberOfCalls(t mock.TestingT, name string, calls int) {
	_m.Mock.AssertNumberOfCalls(t, name, calls)
}

// AssertCalled is a helper method that asserts that the method with the specified name was called with the specified arguments.
func (_m *MockCryptoService) AssertCalled(t mock.TestingT, name string, arguments ...interface{}) {
	_m.Mock.AssertCalled(t, name, arguments...)
}

// AssertNotCalled is a helper method that asserts that the method with the specified name was not called.
func (_m *MockCryptoService) AssertNotCalled(t mock.TestingT, name string) {
	_m.Mock.AssertNotCalled(t, name)
}

# Backend API Documentation

This document provides comprehensive documentation for all HTTP endpoints available in the `authuser`, `settings`, and `catalog` microservices.

## Table of Contents

- [Authentication & Authorization](#authentication--authorization)
- [Status Codes Reference](#status-codes-reference)
- [AuthUser Service](#authuser-service)
  - [Authentication Endpoints](#authentication-endpoints)
  - [User Management Endpoints](#user-management-endpoints)
  - [Partner Management Endpoints](#partner-management-endpoints)
- [Settings Service](#settings-service)
  - [Company Settings Endpoints](#company-settings-endpoints)
  - [OTP Settings Endpoints](#otp-settings-endpoints)
  - [Token Settings Endpoints](#token-settings-endpoints)
  - [Bulk Settings Endpoints](#bulk-settings-endpoints)
- [Catalog Service](#catalog-service)
  - [Product Endpoints](#product-endpoints)
  - [Category Endpoints](#category-endpoints)
  - [Price Endpoints](#price-endpoints)
  - [Image Endpoints](#image-endpoints)
  - [Coupon Endpoints](#coupon-endpoints)
  - [Promotion Code Endpoints](#promotion-code-endpoints)

---

## Authentication & Authorization

### Authentication Mechanisms

The API uses cookie-based authentication with dual tokens:
- **Access Token**: Short-lived token for API requests
- **Refresh Token**: Long-lived token for obtaining new access tokens

### Role Hierarchy

From least to most privileged:
1. **Visitor** - Newly registered users (email verified but profile incomplete)
2. **Standard** - Regular users with complete profiles
3. **Partner** - Service providers
4. **Administrator** - System administrators

Each role has access to all endpoints available to lower-privileged roles.

### Protected Endpoint Notation

In this documentation, endpoints are marked with their minimum required role:
- 🔓 **Public** - No authentication required
- 🎫 **Visitor** - Requires Visitor role or higher
- 👤 **Standard** - Requires Standard role or higher
- 🤝 **Partner** - Requires Partner role or higher
- 👑 **Administrator** - Requires Administrator role only

---

## Status Codes Reference

### Success Codes (2xx)
- **200 OK** - Request succeeded
- **201 Created** - Resource created successfully (e.g., user logged in)
- **207 Multi-Status** - Partial success (used in bulk operations)

### Client Error Codes (4xx)
- **400 Bad Request** - Invalid input data or validation failure
- **401 Unauthorized** - Authentication required or failed
- **403 Forbidden** - Authenticated but not authorized for this resource
- **404 Not Found** - Resource does not exist
- **408 Request Timeout** - Client cancelled request or took too long
- **409 Conflict** - Resource already exists or state conflict
- **415 Unsupported Media Type** - Missing or wrong Content-Type header
- **423 Locked** - Account is locked
- **429 Too Many Requests** - Rate limit exceeded

### Server Error Codes (5xx)
- **500 Internal Server Error** - Generic server error
- **502 Bad Gateway** - External service returned invalid response
- **503 Service Unavailable** - Temporary failure, client should retry (database connections, deadlocks, etc.)
- **504 Gateway Timeout** - Upstream service timed out

---

## AuthUser Service

The AuthUser service handles user authentication, registration, profile management, and partner operations.

---

## Authentication Endpoints

### Check Email & Send OTP

🔓 **Public**

Checks if an email is available for registration and sends a verification OTP.

- **Method**: `POST`
- **Path**: `/auth/email`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "email": "user@example.com"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Verification email sent successfully",
  "status": "sent"
}
```

#### Status Codes
- **200 OK** - OTP sent successfully
- **400 Bad Request** - Invalid email format or validation error
- **409 Conflict** - Email already registered
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **429 Too Many Requests** - Rate limit exceeded for OTP requests
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database or email service temporarily unavailable

---

### Validate OTP & Create Pending User

🔓 **Public**

Validates the OTP code and creates a pending user entry.

- **Method**: `POST`
- **Path**: `/auth/otp`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

#### Success Response (200 OK)
```json
{
  "message": "OTP validated successfully",
  "status": "validated"
}
```

#### Status Codes
- **200 OK** - OTP validated, pending user created
- **400 Bad Request** - Invalid email or OTP format
- **401 Unauthorized** - Invalid or expired OTP
- **404 Not Found** - No OTP found for this email
- **409 Conflict** - User already exists
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **423 Locked** - Too many failed OTP attempts, account locked
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Complete User Registration

🎫 **Visitor**

Completes user registration by setting password and profile information.

- **Method**: `POST`
- **Path**: `/auth/complete`
- **Content-Type**: `application/json`
- **Authentication**: Requires Visitor role (obtained after OTP validation)

#### Request Body
```json
{
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe",
  "birth_date": "1990-01-15T00:00:00Z",
  "gender": "male",
  "telephone": "+33612345678",
  "postal_code": "75001",
  "city": "Paris",
  "address1": "123 Rue de Rivoli",
  "address2": "Apartment 4B"
}
```

**Field Validations**:
- `password`: Minimum 8 characters, must contain uppercase, lowercase, number, and special character
- `gender`: One of: "male", "female", "other", "prefer_not_to_say"
- `telephone`: Valid phone number format (E.164 recommended)
- `postal_code`: Valid French postal code format
- `address2`: Optional

#### Success Response (200 OK)
```json
{
  "message": "User registration completed successfully",
  "status": "completed"
}
```

#### Status Codes
- **200 OK** - User profile completed successfully
- **400 Bad Request** - Validation error (invalid password, missing required fields, etc.)
- **401 Unauthorized** - Not authenticated or invalid token
- **403 Forbidden** - User already completed registration
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database or Stripe service temporarily unavailable

---

### Complete Partner Registration

🎫 **Visitor**

Completes partner registration with user profile and partner-specific information.

- **Method**: `POST`
- **Path**: `/auth/complete/partner`
- **Content-Type**: `application/json`
- **Authentication**: Requires Visitor role

#### Request Body
```json
{
  "password": "SecurePassword123!",
  "first_name": "Jane",
  "last_name": "Smith",
  "birth_date": "1985-05-20T00:00:00Z",
  "gender": "female",
  "telephone": "+33698765432",
  "postal_code": "69001",
  "city": "Lyon",
  "address1": "45 Rue de la République",
  "address2": "",
  "bio": "Experienced massage therapist with 10 years of practice",
  "experience": "Certified in Swedish massage, deep tissue, and reflexology",
  "category_ids": ["550e8400-e29b-41d4-a716-446655440000"],
  "product_ids": ["660e8400-e29b-41d4-a716-446655440000"]
}
```

**Partner-Specific Field Validations**:
- `bio`: Optional, max 1000 characters
- `experience`: Optional, max 2000 characters
- `category_ids`: Optional, array of valid category UUIDs
- `product_ids`: Optional, array of valid product UUIDs

#### Success Response (200 OK)
```json
{
  "message": "Partner registration completed successfully",
  "status": "completed"
}
```

#### Status Codes
- **200 OK** - Partner profile completed successfully
- **400 Bad Request** - Validation error (invalid UUIDs, text too long, etc.)
- **401 Unauthorized** - Not authenticated or invalid token
- **403 Forbidden** - User already completed registration
- **404 Not Found** - Invalid category or product IDs
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database or Stripe service temporarily unavailable

---

### Sign In

🔓 **Public**

Authenticates a user with email and password.

- **Method**: `POST`
- **Path**: `/auth/login`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "email": "user@example.com",
  "password": "UserPassword123!"
}
```

#### Success Response (201 Created)

Sets `access_token` and `refresh_token` HTTP-only cookies.

```json
{
  "message": "user logged in successfully",
  "status": "created"
}
```

#### Status Codes
- **201 Created** - Login successful, tokens set in cookies
- **400 Bad Request** - Invalid email or password format
- **401 Unauthorized** - Invalid credentials
- **403 Forbidden** - Account not approved or inactive
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **423 Locked** - Account locked due to too many failed attempts
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Sign Out

👤 **Standard**

Logs out the current user by invalidating their refresh token.

- **Method**: `POST`
- **Path**: `/auth/logout`
- **Content-Type**: `application/json`
- **Authentication**: Requires Standard role

#### Request Body
```json
{
  "token": "refresh_token_value"
}
```

#### Success Response (200 OK)
```json
{
  "message": "user logged out successfully",
  "status": "logged_out"
}
```

#### Status Codes
- **200 OK** - Logout successful
- **400 Bad Request** - Invalid token format
- **401 Unauthorized** - Not authenticated
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database or Redis temporarily unavailable

---

### Refresh Session

🔐 **Refresh Token Required**

Refreshes the user session and issues new access and refresh tokens.

- **Method**: `POST`
- **Path**: `/auth/refresh`
- **Authentication**: Requires valid refresh token in cookies

#### Request Body
None (uses refresh token from cookies)

#### Success Response (200 OK)

Sets new `access_token` and `refresh_token` HTTP-only cookies.

```json
{
  "message": "session refreshed successfully",
  "status": "refreshed"
}
```

#### Status Codes
- **200 OK** - Session refreshed, new tokens issued
- **401 Unauthorized** - Invalid or expired refresh token
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Request Password Reset

🔓 **Public**

Initiates password reset flow by sending a reset OTP to the user's email.

- **Method**: `POST`
- **Path**: `/auth/password/reset/request`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "email": "user@example.com"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Password reset email sent successfully",
  "status": "sent"
}
```

#### Status Codes
- **200 OK** - Reset OTP sent successfully
- **400 Bad Request** - Invalid email format
- **404 Not Found** - No user found with this email
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **429 Too Many Requests** - Rate limit exceeded
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Email service temporarily unavailable

---

### Validate Password Reset OTP

🔓 **Public**

Validates the password reset OTP and issues a password reset token.

- **Method**: `POST`
- **Path**: `/auth/password/reset/validate`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

#### Success Response (200 OK)
```json
{
  "token": "password_reset_token_here",
  "expires_at": "2025-12-29T15:30:00Z"
}
```

#### Status Codes
- **200 OK** - OTP validated, reset token issued
- **400 Bad Request** - Invalid email or OTP format
- **401 Unauthorized** - Invalid or expired OTP
- **404 Not Found** - No reset request found for this email
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **423 Locked** - Too many failed attempts
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Confirm Password Reset

🔓 **Public**

Confirms password reset with a valid token and updates the password.

- **Method**: `POST`
- **Path**: `/auth/password/reset/confirm`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "token": "password_reset_token_here",
  "new_password": "NewSecurePassword123!"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Password reset successfully",
  "status": "reset"
}
```

#### Status Codes
- **200 OK** - Password reset successfully
- **400 Bad Request** - Invalid password format or validation error
- **401 Unauthorized** - Invalid or expired reset token
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### OAuth Start

🔓 **Public**

Initiates OAuth flow by redirecting to the provider's authorization screen.

- **Method**: `GET`
- **Path**: `/auth/oauth/{provider}`
- **Path Parameters**:
  - `provider`: OAuth provider name (e.g., "google", "apple")

#### Success Response (302 Redirect)

Redirects to OAuth provider's consent screen.

#### Status Codes
- **302 Found** - Redirect to OAuth provider
- **400 Bad Request** - Invalid or unsupported provider
- **500 Internal Server Error** - OAuth configuration error

---

### OAuth Callback

🔓 **Public**

Handles OAuth provider callback, exchanges code for tokens, and creates/logs in the user.

- **Method**: `GET`
- **Path**: `/auth/oauth/{provider}/callback`
- **Path Parameters**:
  - `provider`: OAuth provider name
- **Query Parameters**:
  - `code`: Authorization code from provider
  - `state`: CSRF protection token

#### Success Response (302 Redirect)

Redirects to application with authentication cookies set.

#### Status Codes
- **302 Found** - Redirect to application after successful authentication
- **400 Bad Request** - Invalid callback parameters or CSRF token mismatch
- **401 Unauthorized** - OAuth provider rejected authentication
- **500 Internal Server Error** - Token exchange or user creation failed
- **502 Bad Gateway** - OAuth provider returned invalid response
- **503 Service Unavailable** - Database temporarily unavailable

---

### Delete Own Account

👤 **Standard**

Deletes the currently authenticated user's account.

- **Method**: `DELETE`
- **Path**: `/auth/me`
- **Authentication**: Requires Standard role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "message": "Account deleted successfully",
  "status": "deleted"
}
```

#### Status Codes
- **200 OK** - Account deleted successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - User doesn't have permission to delete own account
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Delete User by Admin

👑 **Administrator**

Deletes any user account (admin only).

- **Method**: `DELETE`
- **Path**: `/admin/auth/users/{id}`
- **Path Parameters**:
  - `id`: User UUID to delete
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "message": "User deleted successfully",
  "status": "deleted"
}
```

#### Status Codes
- **200 OK** - User deleted successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - User not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## User Management Endpoints

### Get Current User

👤 **Standard**

Retrieves the profile of the currently authenticated user.

- **Method**: `GET`
- **Path**: `/users/me`
- **Authentication**: Requires Standard role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "state": "active",
  "email": "user@example.com",
  "picture": "https://example.com/avatar.jpg",
  "created_at": "2025-01-15T10:30:00Z",
  "logged_in_at": "2025-12-29T14:20:00Z",
  "role": "standard",
  "birthdate": "1990-01-15T00:00:00Z",
  "last_name": "Doe",
  "first_name": "John",
  "gender": "male",
  "telephone": "+33612345678",
  "postal_code": "75001",
  "city": "Paris",
  "address1": "123 Rue de Rivoli",
  "address2": "Apartment 4B",
  "google_id": "",
  "apple_id": ""
}
```

#### Status Codes
- **200 OK** - User profile retrieved successfully
- **401 Unauthorized** - Not authenticated
- **404 Not Found** - User not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Update User Profile

👤 **Standard**

Updates the current user's profile information.

- **Method**: `PATCH`
- **Path**: `/users/me`
- **Content-Type**: `application/json`
- **Authentication**: Requires Standard role

#### Request Body

All fields are optional. Only include fields you want to update.

```json
{
  "picture": "https://example.com/new-avatar.jpg",
  "first_name": "Jane",
  "last_name": "Smith",
  "birthdate": "1990-06-20T00:00:00Z",
  "gender": "female",
  "email": "newemail@example.com",
  "telephone": "+33698765432",
  "postal_code": "69001",
  "city": "Lyon",
  "address1": "45 Rue de la République",
  "address2": ""
}
```

#### Success Response (200 OK)
```json
{
  "message": "User profile updated successfully",
  "status": "updated"
}
```

#### Status Codes
- **200 OK** - Profile updated successfully
- **400 Bad Request** - Validation error (invalid email, phone, etc.)
- **401 Unauthorized** - Not authenticated
- **409 Conflict** - Email already in use by another user
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Change Password

👤 **Standard**

Changes the password for the authenticated user.

- **Method**: `PATCH`
- **Path**: `/users/me/password`
- **Content-Type**: `application/json`
- **Authentication**: Requires Standard role

#### Request Body
```json
{
  "old_password": "CurrentPassword123!",
  "new_password": "NewSecurePassword456!"
}
```

**Validations**:
- `old_password`: Must match current password
- `new_password`: Must be different from old password, meet password requirements

#### Success Response (200 OK)
```json
{
  "message": "Password changed successfully",
  "status": "changed"
}
```

#### Status Codes
- **200 OK** - Password changed successfully
- **400 Bad Request** - New password doesn't meet requirements or is same as old password
- **401 Unauthorized** - Not authenticated or old password incorrect
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get All Users

👑 **Administrator**

Retrieves all registered users (admin only).

- **Method**: `GET`
- **Path**: `/admin/users`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "state": "active",
    "email": "user1@example.com",
    "picture": "",
    "created_at": "2025-01-15T10:30:00Z",
    "logged_in_at": "2025-12-29T14:20:00Z",
    "role": "standard",
    "birthdate": "1990-01-15T00:00:00Z",
    "last_name": "Doe",
    "first_name": "John",
    "gender": "male",
    "telephone": "+33612345678",
    "postal_code": "75001",
    "city": "Paris",
    "address1": "123 Rue de Rivoli",
    "address2": ""
  }
]
```

#### Status Codes
- **200 OK** - Users retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Pending Users

👑 **Administrator**

Retrieves all users in pending state awaiting approval (admin only).

- **Method**: `GET`
- **Path**: `/admin/auth/admin/users/pending`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "state": "pending",
    "email": "pending@example.com",
    "created_at": "2025-12-28T10:30:00Z",
    "first_name": "Pending",
    "last_name": "User"
  }
]
```

#### Status Codes
- **200 OK** - Pending users retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get User by ID

👑 **Administrator**

Retrieves details of a specific user by their ID (admin only).

- **Method**: `GET`
- **Path**: `/admin/users/{id}`
- **Path Parameters**:
  - `id`: User UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "state": "active",
  "email": "user@example.com",
  "picture": "",
  "created_at": "2025-01-15T10:30:00Z",
  "logged_in_at": "2025-12-29T14:20:00Z",
  "role": "standard",
  "birthdate": "1990-01-15T00:00:00Z",
  "last_name": "Doe",
  "first_name": "John",
  "gender": "male",
  "telephone": "+33612345678",
  "postal_code": "75001",
  "city": "Paris",
  "address1": "123 Rue de Rivoli",
  "address2": ""
}
```

#### Status Codes
- **200 OK** - User retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - User not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Approve User

👑 **Administrator**

Approves a pending user by setting their role and activating their account (admin only).

- **Method**: `PATCH`
- **Path**: `/admin/users/approve`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "user_id": "660e8400-e29b-41d4-a716-446655440000",
  "role": "standard"
}
```

**Valid Roles**: "visitor", "standard", "partner", "administrator"

#### Success Response (200 OK)
```json
{
  "message": "User approved successfully",
  "status": "approved"
}
```

#### Status Codes
- **200 OK** - User approved successfully
- **400 Bad Request** - Invalid role value
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - User not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Update User Role

👑 **Administrator**

Updates the role of a specific user (admin only).

- **Method**: `PATCH`
- **Path**: `/admin/users/{id}/role`
- **Path Parameters**:
  - `id`: User UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "role": "partner"
}
```

**Valid Roles**: "visitor", "standard", "partner", "administrator"

#### Success Response (200 OK)
```json
{
  "message": "User role updated successfully",
  "status": "updated"
}
```

#### Status Codes
- **200 OK** - Role updated successfully
- **400 Bad Request** - Invalid role value
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - User not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Partner Management Endpoints

### Get Partner by ID

🔓 **Public**

Retrieves partner details by their ID.

- **Method**: `GET`
- **Path**: `/partners/{id}`
- **Path Parameters**:
  - `id`: Partner UUID

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440000",
  "bio": "Experienced massage therapist",
  "experience": "10 years in Swedish and deep tissue massage",
  "category_ids": ["880e8400-e29b-41d4-a716-446655440000"],
  "product_ids": ["990e8400-e29b-41d4-a716-446655440000"],
  "created_at": "2025-01-10T08:00:00Z",
  "updated_at": "2025-12-20T16:30:00Z"
}
```

#### Status Codes
- **200 OK** - Partner retrieved successfully
- **404 Not Found** - Partner not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Authenticated Partner Profile

🤝 **Partner**

Retrieves the authenticated partner's own profile.

- **Method**: `GET`
- **Path**: `/partners/me`
- **Authentication**: Requires Partner role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440000",
  "bio": "Experienced massage therapist",
  "experience": "10 years in Swedish and deep tissue massage",
  "category_ids": ["880e8400-e29b-41d4-a716-446655440000"],
  "product_ids": ["990e8400-e29b-41d4-a716-446655440000"],
  "created_at": "2025-01-10T08:00:00Z",
  "updated_at": "2025-12-20T16:30:00Z"
}
```

#### Status Codes
- **200 OK** - Partner profile retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - User is not a partner
- **404 Not Found** - Partner profile not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get All Partners

🔓 **Public** (Admin restricted for full details)

Retrieves all partners.

- **Method**: `GET`
- **Path**: `/admin/partners`
- **Authentication**: None required for listing, Administrator for full details

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "bio": "Experienced massage therapist",
    "experience": "10 years in Swedish and deep tissue massage",
    "category_ids": ["880e8400-e29b-41d4-a716-446655440000"],
    "product_ids": ["990e8400-e29b-41d4-a716-446655440000"],
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-12-20T16:30:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Partners retrieved successfully
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Partners by Category

🔓 **Public**

Retrieves all partners offering services in a specific category.

- **Method**: `GET`
- **Path**: `/partners/categories/{id}`
- **Path Parameters**:
  - `id`: Category UUID

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "bio": "Swedish massage specialist",
    "experience": "8 years experience",
    "category_ids": ["880e8400-e29b-41d4-a716-446655440000"],
    "product_ids": ["990e8400-e29b-41d4-a716-446655440000"],
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-12-20T16:30:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Partners retrieved successfully
- **404 Not Found** - Category not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Partners by Categories

🔓 **Public**

Retrieves all partners offering services in multiple categories.

- **Method**: `GET`
- **Path**: `/partners/categories?ids={uuid1},{uuid2}`
- **Query Parameters**:
  - `ids`: Comma-separated list of category UUIDs

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "bio": "Multi-specialty therapist",
    "experience": "15 years experience",
    "category_ids": ["880e8400-e29b-41d4-a716-446655440000", "881e8400-e29b-41d4-a716-446655440000"],
    "product_ids": ["990e8400-e29b-41d4-a716-446655440000"],
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-12-20T16:30:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Partners retrieved successfully
- **400 Bad Request** - Invalid UUID format in query parameter
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Partners by Product

🔓 **Public**

Retrieves all partners offering a specific product/service.

- **Method**: `GET`
- **Path**: `/partners/products/{id}`
- **Path Parameters**:
  - `id`: Product UUID

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "bio": "Deep tissue massage expert",
    "experience": "12 years experience",
    "category_ids": ["880e8400-e29b-41d4-a716-446655440000"],
    "product_ids": ["990e8400-e29b-41d4-a716-446655440000"],
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-12-20T16:30:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Partners retrieved successfully
- **404 Not Found** - Product not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Partners by Products

🔓 **Public**

Retrieves all partners offering multiple products/services.

- **Method**: `GET`
- **Path**: `/partners/products?ids={uuid1},{uuid2}`
- **Query Parameters**:
  - `ids`: Comma-separated list of product UUIDs

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "bio": "Full-service massage therapist",
    "experience": "20 years experience",
    "category_ids": ["880e8400-e29b-41d4-a716-446655440000"],
    "product_ids": ["990e8400-e29b-41d4-a716-446655440000", "991e8400-e29b-41d4-a716-446655440000"],
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-12-20T16:30:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Partners retrieved successfully
- **400 Bad Request** - Invalid UUID format in query parameter
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Update Partner Profile

🤝 **Partner**

Updates the authenticated partner's profile.

- **Method**: `PUT`
- **Path**: `/partners/{id}`
- **Path Parameters**:
  - `id`: Partner UUID (must match authenticated partner or be admin)
- **Content-Type**: `application/json`
- **Authentication**: Requires Partner role (own profile) or Administrator (any profile)

#### Request Body

All fields are optional. Only include fields you want to update.

```json
{
  "bio": "Updated bio with new certifications",
  "experience": "Now 11 years of experience in massage therapy"
}
```

**Field Validations**:
- `bio`: Optional, max 1000 characters
- `experience`: Optional, max 2000 characters

#### Success Response (200 OK)
```json
{
  "message": "Partner profile updated successfully",
  "status": "updated"
}
```

#### Status Codes
- **200 OK** - Profile updated successfully
- **400 Bad Request** - Validation error (text too long, etc.)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not authorized to update this partner profile
- **404 Not Found** - Partner not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Delete Partner

👑 **Administrator**

Deletes a partner profile (admin only).

- **Method**: `DELETE`
- **Path**: `/admin/partners/{id}`
- **Path Parameters**:
  - `id`: Partner UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "message": "Partner deleted successfully",
  "status": "deleted"
}
```

#### Status Codes
- **200 OK** - Partner deleted successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Partner not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Verify Partner

👑 **Administrator**

Verifies partner credentials and activates their account (admin only).

- **Method**: `POST`
- **Path**: `/admin/partners/{id}/verify`
- **Path Parameters**:
  - `id`: Partner UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "partner_id": "770e8400-e29b-41d4-a716-446655440000"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Partner verified successfully",
  "status": "verified"
}
```

#### Status Codes
- **200 OK** - Partner verified successfully
- **400 Bad Request** - Invalid partner ID
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Partner not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Settings Service

The Settings service manages system configuration including company information, OTP settings, and token durations.

---

## Company Settings Endpoints

### Get Company Name

🔓 **Public**

Retrieves the company name.

- **Method**: `GET`
- **Path**: `/settings/name`

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "name": "Leviosa Spa & Wellness"
}
```

#### Status Codes
- **200 OK** - Company name retrieved successfully
- **404 Not Found** - Company name not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Company Name

👑 **Administrator**

Sets or updates the company name.

- **Method**: `POST`
- **Path**: `/admin/settings/name`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "name": "Leviosa Spa & Wellness Center"
}
```

**Validations**:
- `name`: Required, min 1 character, max 255 characters, cannot be only whitespace

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Company name updated successfully"
}
```

#### Status Codes
- **200 OK** - Company name updated successfully
- **400 Bad Request** - Validation error (empty name, too long, etc.)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Company Email

🔓 **Public**

Retrieves the company contact email.

- **Method**: `GET`
- **Path**: `/settings/email`

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "email": "contact@leviosa-spa.com"
}
```

#### Status Codes
- **200 OK** - Company email retrieved successfully
- **404 Not Found** - Company email not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Company Email

👑 **Administrator**

Sets or updates the company contact email.

- **Method**: `POST`
- **Path**: `/admin/settings/email`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "email": "info@leviosa-spa.com"
}
```

**Validations**:
- `email`: Required, valid email format, max 255 characters

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Company email updated successfully"
}
```

#### Status Codes
- **200 OK** - Company email updated successfully
- **400 Bad Request** - Invalid email format
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Company Phone

👑 **Administrator**

Retrieves the company phone number (admin only, sensitive data).

- **Method**: `GET`
- **Path**: `/admin/settings/phone`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "telephone": "+33123456789"
}
```

#### Status Codes
- **200 OK** - Company phone retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Company phone not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Company Phone

👑 **Administrator**

Sets or updates the company phone number.

- **Method**: `POST`
- **Path**: `/admin/settings/phone`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "telephone": "+33123456789"
}
```

**Validations**:
- `telephone`: Required, valid phone format, min 10 characters, max 20 characters

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Company phone updated successfully"
}
```

#### Status Codes
- **200 OK** - Company phone updated successfully
- **400 Bad Request** - Invalid phone format
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Company Address

🔓 **Public**

Retrieves the company legal address.

- **Method**: `GET`
- **Path**: `/settings/address`

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "address": "123 Wellness Boulevard, 75001 Paris, France"
}
```

#### Status Codes
- **200 OK** - Company address retrieved successfully
- **404 Not Found** - Company address not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Company Address

👑 **Administrator**

Sets or updates the company legal address.

- **Method**: `POST`
- **Path**: `/admin/settings/address`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "address": "456 Spa Avenue, 69001 Lyon, France"
}
```

**Validations**:
- `address`: Required, min 1 character, max 500 characters, cannot be only whitespace

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Company address updated successfully"
}
```

#### Status Codes
- **200 OK** - Company address updated successfully
- **400 Bad Request** - Validation error (empty, too long, etc.)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Company Instagram

🔓 **Public**

Retrieves the company Instagram profile URL.

- **Method**: `GET`
- **Path**: `/settings/instagram`

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "instagram": "https://instagram.com/leviosa_spa"
}
```

#### Status Codes
- **200 OK** - Instagram URL retrieved successfully
- **404 Not Found** - Instagram URL not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Company Instagram

👑 **Administrator**

Sets or updates the company Instagram profile URL.

- **Method**: `POST`
- **Path**: `/admin/settings/instagram`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "instagram": "https://instagram.com/leviosa_wellness"
}
```

**Validations**:
- `instagram`: Required, valid URL format (http/https), max 255 characters

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Instagram URL updated successfully"
}
```

#### Status Codes
- **200 OK** - Instagram URL updated successfully
- **400 Bad Request** - Invalid URL format
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Company Logo

🔓 **Public**

Retrieves the company logo URL and metadata.

- **Method**: `GET`
- **Path**: `/settings/logo`

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "logo_url": "https://s3.amazonaws.com/leviosa-assets/logo.png",
  "content_type": "image/png"
}
```

#### Status Codes
- **200 OK** - Logo retrieved successfully
- **404 Not Found** - Logo not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database or S3 temporarily unavailable

---

### Set Company Logo

👑 **Administrator**

Uploads and sets the company logo.

- **Method**: `POST`
- **Path**: `/admin/settings/logo`
- **Content-Type**: `multipart/form-data`
- **Authentication**: Requires Administrator role

#### Request Body (Multipart Form)
- `file`: Logo image file
- `content_type`: MIME type (image/jpeg, image/png, or image/gif)
- `file_size`: File size in bytes

**Validations**:
- `content_type`: Must be image/jpeg, image/png, or image/gif
- `file_size`: Min 1 byte, max 5MB (5,242,880 bytes)

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Company logo uploaded successfully"
}
```

#### Status Codes
- **200 OK** - Logo uploaded successfully
- **400 Bad Request** - Invalid file type or size
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server or S3 error
- **503 Service Unavailable** - S3 temporarily unavailable

---

## OTP Settings Endpoints

All OTP settings endpoints require Administrator role.

### Get OTP Duration

👑 **Administrator**

Retrieves the OTP validity duration in seconds.

- **Method**: `GET`
- **Path**: `/admin/settings/otp/duration`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "duration": 300
}
```

#### Status Codes
- **200 OK** - OTP duration retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - OTP duration not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set OTP Duration

👑 **Administrator**

Sets the OTP validity duration in seconds.

- **Method**: `POST`
- **Path**: `/admin/settings/otp/duration`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "duration": 600
}
```

**Validations**:
- `duration`: Required, min 60 seconds, max 3600 seconds (1 hour)

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "OTP duration updated successfully"
}
```

#### Status Codes
- **200 OK** - OTP duration updated successfully
- **400 Bad Request** - Duration out of valid range (60-3600)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get OTP Length

👑 **Administrator**

Retrieves the OTP code length (number of digits).

- **Method**: `GET`
- **Path**: `/admin/settings/otp/length`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "length": 6
}
```

#### Status Codes
- **200 OK** - OTP length retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - OTP length not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set OTP Length

👑 **Administrator**

Sets the OTP code length (number of digits).

- **Method**: `POST`
- **Path**: `/admin/settings/otp/length`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "length": 8
}
```

**Validations**:
- `length`: Required, min 4 digits, max 10 digits

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "OTP length updated successfully"
}
```

#### Status Codes
- **200 OK** - OTP length updated successfully
- **400 Bad Request** - Length out of valid range (4-10)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get OTP Max Attempts

👑 **Administrator**

Retrieves the maximum number of OTP validation attempts before account lockout.

- **Method**: `GET`
- **Path**: `/admin/settings/otp/max-attempts`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "max_attempts": 5
}
```

#### Status Codes
- **200 OK** - Max attempts retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Max attempts not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set OTP Max Attempts

👑 **Administrator**

Sets the maximum number of OTP validation attempts.

- **Method**: `POST`
- **Path**: `/admin/settings/otp/max-attempts`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "max_attempts": 3
}
```

**Validations**:
- `max_attempts`: Required, min 1, max 10

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "OTP max attempts updated successfully"
}
```

#### Status Codes
- **200 OK** - Max attempts updated successfully
- **400 Bad Request** - Value out of valid range (1-10)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Token Settings Endpoints

All token settings endpoints require Administrator role.

### Get Access Token Duration

👑 **Administrator**

Retrieves the access token validity duration in minutes.

- **Method**: `GET`
- **Path**: `/admin/settings/tokens/access-duration`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "duration": 15
}
```

#### Status Codes
- **200 OK** - Access token duration retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Access token duration not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Access Token Duration

👑 **Administrator**

Sets the access token validity duration in minutes.

- **Method**: `POST`
- **Path**: `/admin/settings/tokens/access-duration`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "duration": 30
}
```

**Validations**:
- `duration`: Required, min 1 minute, max 240 minutes (4 hours)

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Access token duration updated successfully"
}
```

#### Status Codes
- **200 OK** - Access token duration updated successfully
- **400 Bad Request** - Duration out of valid range (1-240)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Refresh Token Duration

👑 **Administrator**

Retrieves the refresh token validity duration in hours.

- **Method**: `GET`
- **Path**: `/admin/settings/tokens/refresh-duration`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "duration": 168
}
```

#### Status Codes
- **200 OK** - Refresh token duration retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Refresh token duration not configured
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Refresh Token Duration

👑 **Administrator**

Sets the refresh token validity duration in hours.

- **Method**: `POST`
- **Path**: `/admin/settings/tokens/refresh-duration`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "duration": 336
}
```

**Validations**:
- `duration`: Required, min 1 hour, max 720 hours (30 days)

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "Refresh token duration updated successfully"
}
```

#### Status Codes
- **200 OK** - Refresh token duration updated successfully
- **400 Bad Request** - Duration out of valid range (1-720)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Bulk Settings Endpoints

### Get Bulk Settings

👑 **Administrator**

Retrieves multiple settings in a single request.

- **Method**: `GET`
- **Path**: `/admin/settings/bulk?keys={key1},{key2},{key3}`
- **Query Parameters**:
  - `keys`: Comma-separated list of setting keys
- **Authentication**: Requires Administrator role

#### Valid Setting Keys
- `company_name`
- `company_email`
- `company_phone`
- `company_address`
- `company_instagram`
- `company_logo`
- `otp_duration`
- `otp_length`
- `otp_max_attempts`

#### Request Example
```
GET /admin/settings/bulk?keys=company_name,company_email,otp_duration
```

#### Success Response (200 OK)
```json
[
  {
    "key": "company_name",
    "value": "Leviosa Spa & Wellness"
  },
  {
    "key": "company_email",
    "value": "contact@leviosa-spa.com"
  },
  {
    "key": "otp_duration",
    "value": "300"
  }
]
```

#### Partial Success Response (207 Multi-Status)

When some settings succeed and others fail:

```json
{
  "data": [
    {
      "key": "company_name",
      "value": "Leviosa Spa & Wellness"
    }
  ],
  "errors": {
    "company_logo": "setting not found",
    "invalid_key": "invalid key: invalid_key"
  }
}
```

#### Status Codes
- **200 OK** - All settings retrieved successfully
- **207 Multi-Status** - Some settings succeeded, some failed (see response for details)
- **400 Bad Request** - Missing `keys` query parameter
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Internal Service-to-Service Endpoints

The Settings service also provides internal endpoints for service-to-service communication. These endpoints are protected by service authentication (not user authentication) and follow the same patterns as the external endpoints but with different base paths:

- **Base Path**: `/internal/settings`
- **Authentication**: Service authentication token required

These endpoints mirror the public GET endpoints but are intended for microservice-to-microservice communication within the backend infrastructure.

---

## Catalog Service

The Catalog service manages products, categories, pricing, images, coupons, and promotion codes. It integrates with Stripe for payment processing and pricing management.

---

## Product Endpoints

### Get All Published Products

🔓 **Public**

Retrieves all published products with their associated pricing and category information.

- **Method**: `GET`
- **Path**: `/products`

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Swedish Massage - 60min",
    "description": "Classic Swedish massage for relaxation",
    "category": "880e8400-e29b-41d4-a716-446655440000",
    "duration": 60,
    "published_status": "published",
    "availability": "both",
    "buffer_time": 15,
    "cancellation_hours": 24,
    "metadata": {},
    "prices": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "stripe_price_id": "price_1234567890",
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "amount": 8000,
        "currency": "eur",
        "interval": "one_time",
        "active": true,
        "nickname": "Standard Rate",
        "metadata": {},
        "created_at": "2025-01-15T10:30:00Z"
      }
    ],
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
]
```

#### Field Descriptions
- `availability`: Enum values - "in_salon", "home_visit", "both"
- `published_status`: Enum values - "draft", "published", "archived"
- `duration`: Service duration in minutes
- `buffer_time`: Buffer time after service in minutes
- `cancellation_hours`: Minimum hours before appointment for cancellation
- `interval`: Price billing interval - "one_time", "month", "year"
- `amount`: Price in cents (e.g., 8000 = €80.00)

#### Status Codes
- **200 OK** - Products retrieved successfully
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Product by ID

🔓 **Public**

Retrieves a single published product by its UUID.

- **Method**: `GET`
- **Path**: `/products/{id}`
- **Path Parameters**:
  - `id`: Product UUID

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Swedish Massage - 60min",
  "description": "Classic Swedish massage for relaxation",
  "category": "880e8400-e29b-41d4-a716-446655440000",
  "duration": 60,
  "published_status": "published",
  "availability": "both",
  "buffer_time": 15,
  "cancellation_hours": 24,
  "metadata": {},
  "prices": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "stripe_price_id": "price_1234567890",
      "product_id": "550e8400-e29b-41d4-a716-446655440000",
      "amount": 8000,
      "currency": "eur",
      "interval": "one_time",
      "active": true,
      "nickname": "Standard Rate",
      "metadata": {},
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

#### Status Codes
- **200 OK** - Product retrieved successfully
- **404 Not Found** - Product not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get All Products (Admin)

👑 **Administrator**

Retrieves all products including draft and archived ones (admin only).

- **Method**: `GET`
- **Path**: `/admin/products`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as public product list, but includes products with all published_status values.

#### Status Codes
- **200 OK** - Products retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Create Product with Price

👑 **Administrator**

Creates a new product with an initial price. Automatically creates corresponding Stripe product and price.

- **Method**: `POST`
- **Path**: `/admin/products`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "product": {
    "name": "Deep Tissue Massage - 90min",
    "description": "Intensive deep tissue massage for muscle relief",
    "category": "880e8400-e29b-41d4-a716-446655440000",
    "duration": 90,
    "availability": "both",
    "buffer_time": 20,
    "cancellation_hours": 48,
    "metadata": {
      "intensity": "high"
    }
  },
  "price": {
    "amount": 12000,
    "currency": "eur",
    "interval": "one_time",
    "nickname": "Standard Rate",
    "metadata": {}
  }
}
```

**Product Field Validations**:
- `name`: Required, max 255 characters
- `description`: Optional, max 1000 characters
- `category`: Required, valid category UUID
- `duration`: Required, positive integer (minutes)
- `availability`: Required, one of: "in_salon", "home_visit", "both"
- `buffer_time`: Required, non-negative integer (minutes)
- `cancellation_hours`: Required, non-negative integer
- `metadata`: Optional JSON object

**Price Field Validations**:
- `amount`: Required, positive integer (cents)
- `currency`: Required, 3-character ISO currency code (e.g., "eur", "usd")
- `interval`: Required, one of: "one_time", "month", "year"
- `nickname`: Optional, max 255 characters
- `metadata`: Optional JSON object

#### Success Response (201 Created)
```json
{
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "price_id": "660e8400-e29b-41d4-a716-446655440000",
  "message": "Product created successfully!"
}
```

#### Status Codes
- **201 Created** - Product and price created successfully
- **400 Bad Request** - Validation error or invalid input
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Category not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Modify Product

👑 **Administrator**

Updates an existing product. All fields are optional - only include fields you want to change.

- **Method**: `PATCH`
- **Path**: `/admin/products/{id}`
- **Path Parameters**:
  - `id`: Product UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "name": "Swedish Massage - 75min",
  "description": "Extended Swedish massage session",
  "category": "880e8400-e29b-41d4-a716-446655440001",
  "duration": 75,
  "published_status": "published",
  "availability": "in_salon",
  "buffer_time": 20,
  "cancellation_hours": 48,
  "metadata": {
    "featured": true
  }
}
```

**Field Validations**:
- `name`: Optional, max 255 characters
- `description`: Optional, max 1000 characters
- `category`: Optional, valid category UUID
- `duration`: Optional, positive integer
- `published_status`: Optional, one of: "draft", "published", "archived"
- `availability`: Optional, one of: "in_salon", "home_visit", "both"
- `buffer_time`: Optional, non-negative integer
- `cancellation_hours`: Optional, non-negative integer
- `metadata`: Optional JSON object
- At least one field must be provided

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Product updated successfully
- **400 Bad Request** - No fields provided or validation error
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Product or category not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Remove Product

👑 **Administrator**

Deletes a product and all associated prices. Also deletes from Stripe.

- **Method**: `DELETE`
- **Path**: `/admin/products/{id}`
- **Path Parameters**:
  - `id`: Product UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Product deleted successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Product not found
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

## Category Endpoints

### Get All Published Categories

🔓 **Public**

Retrieves all published categories with their active images.

- **Method**: `GET`
- **Path**: `/categories`

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "name": "Massage Therapy",
    "description": "Professional massage and bodywork services",
    "status": "published",
    "metadata": {},
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-01-15T10:30:00Z",
    "images": [
      {
        "id": "990e8400-e29b-41d4-a716-446655440000",
        "parent_id": "880e8400-e29b-41d4-a716-446655440000",
        "parent_type": "category",
        "url": "https://s3.amazonaws.com/images/category-massage.jpg",
        "title": "Massage Room",
        "is_active": true,
        "created_at": "2025-01-10T08:15:00Z"
      }
    ]
  }
]
```

#### Field Descriptions
- `status`: Enum values - "draft", "published", "archived"
- `parent_type`: Always "category" for category images
- `is_active`: Only one image per category can be active

#### Status Codes
- **200 OK** - Categories retrieved successfully
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Category by ID

🔓 **Public**

Retrieves a single published category by its UUID.

- **Method**: `GET`
- **Path**: `/categories/{id}`
- **Path Parameters**:
  - `id`: Category UUID

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "name": "Massage Therapy",
  "description": "Professional massage and bodywork services",
  "status": "published",
  "metadata": {},
  "created_at": "2025-01-10T08:00:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "images": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440000",
      "parent_id": "880e8400-e29b-41d4-a716-446655440000",
      "parent_type": "category",
      "url": "https://s3.amazonaws.com/images/category-massage.jpg",
      "title": "Massage Room",
      "is_active": true,
      "created_at": "2025-01-10T08:15:00Z"
    }
  ]
}
```

#### Status Codes
- **200 OK** - Category retrieved successfully
- **404 Not Found** - Category not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get All Categories (Admin)

👑 **Administrator**

Retrieves all categories including draft and archived ones (admin only).

- **Method**: `GET`
- **Path**: `/admin/categories`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as public category list, but includes categories with all status values.

#### Status Codes
- **200 OK** - Categories retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Create Category

👑 **Administrator**

Creates a new category.

- **Method**: `POST`
- **Path**: `/admin/categories`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "name": "Aromatherapy",
  "description": "Essential oil-based therapeutic treatments",
  "metadata": {
    "color": "purple"
  }
}
```

**Field Validations**:
- `name`: Required, min 1 character, max 255 characters
- `description`: Required, min 1 character, max 1000 characters
- `metadata`: Optional JSON object

#### Success Response (201 Created)
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440001",
  "message": "Category created successfully!"
}
```

#### Status Codes
- **201 Created** - Category created successfully
- **400 Bad Request** - Validation error (missing required fields, too long, etc.)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **409 Conflict** - Category name already exists
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Modify Category

👑 **Administrator**

Updates an existing category. All fields are optional.

- **Method**: `PATCH`
- **Path**: `/admin/categories/{id}`
- **Path Parameters**:
  - `id`: Category UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "name": "Massage & Bodywork",
  "description": "Updated description",
  "status": "published",
  "metadata": {
    "featured": true
  }
}
```

**Field Validations**:
- `name`: Optional, min 1 character, max 255 characters
- `description`: Optional, min 1 character, max 1000 characters
- `status`: Optional, one of: "draft", "published", "archived"
- `metadata`: Optional JSON object
- At least one field must be provided

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Category updated successfully
- **400 Bad Request** - No fields provided or validation error
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Category not found
- **409 Conflict** - Category name already exists
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Remove Category

👑 **Administrator**

Deletes a category. Cannot delete if products are associated with it.

- **Method**: `DELETE`
- **Path**: `/admin/categories/{id}`
- **Path Parameters**:
  - `id`: Category UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Category deleted successfully
- **400 Bad Request** - Category has associated products
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Category not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Price Endpoints

### Get Price by ID

👑 **Administrator**

Retrieves a single price by its UUID.

- **Method**: `GET`
- **Path**: `/admin/prices/{id}`
- **Path Parameters**:
  - `id`: Price UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "stripe_price_id": "price_1234567890",
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 8000,
  "currency": "eur",
  "interval": "one_time",
  "active": true,
  "nickname": "Standard Rate",
  "metadata": {},
  "created_at": "2025-01-15T10:30:00Z"
}
```

#### Status Codes
- **200 OK** - Price retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Price not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Prices by Product ID

👑 **Administrator**

Retrieves all prices for a specific product.

- **Method**: `GET`
- **Path**: `/admin/products/{id}/prices`
- **Path Parameters**:
  - `id`: Product UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "stripe_price_id": "price_1234567890",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 8000,
    "currency": "eur",
    "interval": "one_time",
    "active": true,
    "nickname": "Standard Rate",
    "metadata": {},
    "created_at": "2025-01-15T10:30:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Prices retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Product not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Create Price

👑 **Administrator**

Creates a new price for an existing product. Also creates corresponding Stripe price.

- **Method**: `POST`
- **Path**: `/admin/products/{id}/prices`
- **Path Parameters**:
  - `id`: Product UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "amount": 10000,
  "currency": "eur",
  "interval": "one_time",
  "nickname": "Premium Rate",
  "metadata": {
    "tier": "premium"
  }
}
```

**Field Validations**:
- `amount`: Required, positive integer (cents)
- `currency`: Required, 3-character ISO currency code
- `interval`: Required, one of: "one_time", "month", "year"
- `nickname`: Optional, max 255 characters
- `metadata`: Optional JSON object

#### Success Response (201 Created)
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "stripe_price_id": "price_0987654321",
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 10000,
  "currency": "eur",
  "interval": "one_time",
  "active": true,
  "nickname": "Premium Rate",
  "metadata": {
    "tier": "premium"
  },
  "created_at": "2025-01-15T14:00:00Z"
}
```

#### Status Codes
- **201 Created** - Price created successfully
- **400 Bad Request** - Validation error or invalid input
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Product not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Update Price

👑 **Administrator**

Updates price metadata and active status. Note: Amount and currency cannot be changed (Stripe limitation).

- **Method**: `PATCH`
- **Path**: `/admin/prices/{id}`
- **Path Parameters**:
  - `id`: Price UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "active": false,
  "nickname": "Legacy Rate",
  "metadata": {
    "deprecated": true
  }
}
```

**Field Validations**:
- `active`: Optional, boolean
- `nickname`: Optional, max 255 characters
- `metadata`: Optional JSON object
- At least one field must be provided

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Price updated successfully
- **400 Bad Request** - No fields provided or validation error
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Price not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

## Image Endpoints

### Upload Image

👑 **Administrator**

Associates an image with a category or product.

- **Method**: `POST`
- **Path**: `/admin/images`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "parent_id": "880e8400-e29b-41d4-a716-446655440000",
  "parent_type": "category",
  "title": "Main Category Image",
  "is_active": true
}
```

**Field Validations**:
- `parent_id`: Required, valid UUID
- `parent_type`: Required, one of: "category", "product"
- `title`: Required, max 255 characters
- `is_active`: Optional, boolean (default: false)

**Important**: If `is_active` is true, automatically deactivates other images for the same parent.

#### Success Response (201 Created)
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440001",
  "parent_id": "880e8400-e29b-41d4-a716-446655440000",
  "parent_type": "category",
  "url": "https://s3.amazonaws.com/images/category-image.jpg",
  "title": "Main Category Image",
  "is_active": true,
  "created_at": "2025-01-15T15:00:00Z"
}
```

#### Status Codes
- **201 Created** - Image uploaded successfully
- **400 Bad Request** - Validation error or invalid parent type
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Parent category or product not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database or S3 temporarily unavailable

---

### Remove Image

👑 **Administrator**

Deletes an image association.

- **Method**: `DELETE`
- **Path**: `/admin/images`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "image_id": "990e8400-e29b-41d4-a716-446655440001",
  "parent_id": "880e8400-e29b-41d4-a716-446655440000",
  "parent_type": "category"
}
```

**Field Validations**:
- `image_id`: Required, valid UUID
- `parent_id`: Required, valid UUID
- `parent_type`: Required, one of: "category", "product"

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Image deleted successfully
- **400 Bad Request** - Validation error or invalid parent type
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Image not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Set Active Image

👑 **Administrator**

Sets an image as the active image for a category or product. Automatically deactivates other images.

- **Method**: `POST`
- **Path**: `/admin/images/set-active`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "image_id": "990e8400-e29b-41d4-a716-446655440001",
  "parent_id": "880e8400-e29b-41d4-a716-446655440000",
  "parent_type": "category"
}
```

**Field Validations**:
- `image_id`: Required, valid UUID
- `parent_id`: Required, valid UUID
- `parent_type`: Required, one of: "category", "product"

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Active image updated successfully
- **400 Bad Request** - Validation error or invalid parent type
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Image or parent not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

## Coupon Endpoints

### Validate Coupon

🔓 **Public**

Validates a coupon by its Stripe coupon ID.

- **Method**: `POST`
- **Path**: `/coupons/validate`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "stripeCouponId": "SUMMER2025"
}
```

**Field Validations**:
- `stripeCouponId`: Required, non-empty string

#### Success Response (200 OK)

**Valid Coupon**:
```json
{
  "valid": true,
  "coupon": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "stripe_coupon_id": "SUMMER2025",
    "name": "Summer Sale",
    "percent_off": 20.0,
    "amount_off": null,
    "currency": null,
    "duration": "once",
    "duration_in_months": null,
    "max_redemptions": 100,
    "times_redeemed": 45,
    "redeem_by": "2025-09-01T00:00:00Z",
    "valid": true,
    "metadata": {},
    "created_at": "2025-06-01T00:00:00Z"
  }
}
```

**Invalid Coupon**:
```json
{
  "valid": false,
  "reason": "Coupon has expired"
}
```

#### Field Descriptions
- `percent_off`: Percentage discount (0.1-100)
- `amount_off`: Fixed discount in cents (mutually exclusive with percent_off)
- `currency`: Required if amount_off is set
- `duration`: Enum values - "once", "repeating", "forever"
- `duration_in_months`: Required if duration is "repeating"
- `times_redeemed`: Current redemption count

#### Status Codes
- **200 OK** - Validation completed (check `valid` field for result)
- **400 Bad Request** - Missing or invalid stripeCouponId
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Get Valid Coupons

🔓 **Public**

Retrieves all currently valid and active coupons.

- **Method**: `GET`
- **Path**: `/coupons/valid`

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "stripe_coupon_id": "SUMMER2025",
    "name": "Summer Sale",
    "percent_off": 20.0,
    "amount_off": null,
    "currency": null,
    "duration": "once",
    "duration_in_months": null,
    "max_redemptions": 100,
    "times_redeemed": 45,
    "redeem_by": "2025-09-01T00:00:00Z",
    "valid": true,
    "metadata": {},
    "created_at": "2025-06-01T00:00:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Valid coupons retrieved successfully
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get All Coupons (Admin)

👑 **Administrator**

Retrieves all coupons including invalid and expired ones.

- **Method**: `GET`
- **Path**: `/admin/coupons`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as valid coupons list, but includes all coupons regardless of validity.

#### Status Codes
- **200 OK** - Coupons retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Coupon by ID (Admin)

👑 **Administrator**

Retrieves a single coupon by its UUID.

- **Method**: `GET`
- **Path**: `/admin/coupons/{id}`
- **Path Parameters**:
  - `id`: Coupon UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440000",
  "stripe_coupon_id": "SUMMER2025",
  "name": "Summer Sale",
  "percent_off": 20.0,
  "amount_off": null,
  "currency": null,
  "duration": "once",
  "duration_in_months": null,
  "max_redemptions": 100,
  "times_redeemed": 45,
  "redeem_by": "2025-09-01T00:00:00Z",
  "valid": true,
  "metadata": {},
  "created_at": "2025-06-01T00:00:00Z"
}
```

#### Status Codes
- **200 OK** - Coupon retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Coupon not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Coupon by Stripe ID (Admin)

👑 **Administrator**

Retrieves a coupon by its Stripe coupon ID.

- **Method**: `GET`
- **Path**: `/admin/coupons/stripe/{stripeId}`
- **Path Parameters**:
  - `stripeId`: Stripe coupon ID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as Get Coupon by ID.

#### Status Codes
- **200 OK** - Coupon retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Coupon not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Create Coupon (Admin)

👑 **Administrator**

Creates a new coupon in both the database and Stripe.

- **Method**: `POST`
- **Path**: `/admin/coupons`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "name": "Holiday Special",
  "percentOff": 25.0,
  "duration": "once",
  "maxRedemptions": 50,
  "redeemBy": "2025-12-31T23:59:59Z",
  "metadata": {
    "campaign": "holiday2025"
  }
}
```

**Field Validations**:
- `name`: Required, max 255 characters
- `percentOff`: Optional, 0.1-100 (mutually exclusive with amountOff)
- `amountOff`: Optional, positive integer in cents (mutually exclusive with percentOff)
- `currency`: Required if amountOff is set, 3-character ISO code
- `duration`: Required, one of: "once", "repeating", "forever"
- `durationInMonths`: Required if duration is "repeating"
- `maxRedemptions`: Optional, positive integer
- `redeemBy`: Optional, ISO 8601 timestamp
- `metadata`: Optional JSON object

**Important**: Must provide either `percentOff` OR `amountOff`, not both.

#### Success Response (201 Created)
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440001",
  "stripe_coupon_id": "HOLIDAY2025_ABC123",
  "name": "Holiday Special",
  "percent_off": 25.0,
  "amount_off": null,
  "currency": null,
  "duration": "once",
  "duration_in_months": null,
  "max_redemptions": 50,
  "times_redeemed": 0,
  "redeem_by": "2025-12-31T23:59:59Z",
  "valid": true,
  "metadata": {
    "campaign": "holiday2025"
  },
  "created_at": "2025-01-15T16:00:00Z"
}
```

#### Status Codes
- **201 Created** - Coupon created successfully
- **400 Bad Request** - Validation error (both percentOff and amountOff provided, missing currency, etc.)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Update Coupon (Admin)

👑 **Administrator**

Updates coupon name and metadata. Most fields cannot be changed after creation (Stripe limitation).

- **Method**: `PATCH`
- **Path**: `/admin/coupons/{id}`
- **Path Parameters**:
  - `id`: Coupon UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "name": "Updated Holiday Special",
  "metadata": {
    "campaign": "holiday2025",
    "featured": true
  }
}
```

**Field Validations**:
- `name`: Optional, max 255 characters
- `metadata`: Optional JSON object
- At least one field must be provided

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Coupon updated successfully
- **400 Bad Request** - No fields provided or validation error
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Coupon not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Deactivate Coupon (Admin)

👑 **Administrator**

Deactivates a coupon, preventing new redemptions. Does not affect existing subscriptions.

- **Method**: `POST`
- **Path**: `/admin/coupons/{id}/deactivate`
- **Path Parameters**:
  - `id`: Coupon UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Coupon deactivated successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Coupon not found
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Delete Coupon (Admin)

👑 **Administrator**

Deletes a coupon from both database and Stripe.

- **Method**: `DELETE`
- **Path**: `/admin/coupons/{id}`
- **Path Parameters**:
  - `id`: Coupon UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Coupon deleted successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Coupon not found
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

## Promotion Code Endpoints

### Validate Promotion Code

🔓 **Public**

Validates a promotion code and checks if it can be applied to an order.

- **Method**: `POST`
- **Path**: `/promotion-codes/validate`
- **Content-Type**: `application/json`

#### Request Body
```json
{
  "code": "SUMMER20",
  "orderAmount": 10000,
  "orderCurrency": "eur",
  "customerId": "cus_123456789"
}
```

**Field Validations**:
- `code`: Required, non-empty string
- `orderAmount`: Optional, positive integer in cents
- `orderCurrency`: Optional, 3-character ISO code
- `customerId`: Optional, Stripe customer ID

#### Success Response (200 OK)

**Valid Promotion Code**:
```json
{
  "valid": true,
  "promotionCode": {
    "id": "aa0e8400-e29b-41d4-a716-446655440000",
    "stripe_promotion_code_id": "promo_123456",
    "coupon_id": "770e8400-e29b-41d4-a716-446655440000",
    "code": "SUMMER20",
    "active": true,
    "max_redemptions": 100,
    "times_redeemed": 45,
    "expires_at": "2025-09-01T00:00:00Z",
    "first_time_transaction": false,
    "minimum_amount": 5000,
    "minimum_amount_currency": "eur",
    "restrictions": {
      "currency_options": ["eur", "usd"]
    },
    "metadata": {},
    "created_at": "2025-06-01T00:00:00Z"
  }
}
```

**Invalid Promotion Code**:
```json
{
  "valid": false,
  "reason": "Order amount below minimum required (€50.00)"
}
```

#### Validation Rules
- Code must be active
- Must not be expired (expires_at check)
- Must not exceed max_redemptions
- Order amount must meet minimum_amount requirement
- Order currency must match restrictions.currency_options
- If first_time_transaction is true, customer must be new

#### Status Codes
- **200 OK** - Validation completed (check `valid` field for result)
- **400 Bad Request** - Missing or invalid code
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Get Promotion Code with Coupon (Public)

🔓 **Public**

Retrieves promotion code details with associated coupon information.

- **Method**: `GET`
- **Path**: `/promotion-codes/code/{code}`
- **Path Parameters**:
  - `code`: Promotion code string

#### Request Body
None

#### Success Response (200 OK)
```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440000",
  "stripe_promotion_code_id": "promo_123456",
  "coupon_id": "770e8400-e29b-41d4-a716-446655440000",
  "code": "SUMMER20",
  "active": true,
  "max_redemptions": 100,
  "times_redeemed": 45,
  "expires_at": "2025-09-01T00:00:00Z",
  "first_time_transaction": false,
  "minimum_amount": 5000,
  "minimum_amount_currency": "eur",
  "restrictions": {
    "currency_options": ["eur", "usd"]
  },
  "metadata": {},
  "created_at": "2025-06-01T00:00:00Z",
  "coupon": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "stripe_coupon_id": "SUMMER2025",
    "name": "Summer Sale",
    "percent_off": 20.0,
    "amount_off": null,
    "currency": null,
    "duration": "once",
    "duration_in_months": null,
    "max_redemptions": 100,
    "times_redeemed": 45,
    "redeem_by": "2025-09-01T00:00:00Z",
    "valid": true,
    "metadata": {},
    "created_at": "2025-06-01T00:00:00Z"
  }
}
```

#### Status Codes
- **200 OK** - Promotion code retrieved successfully
- **404 Not Found** - Promotion code not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get All Promotion Codes (Admin)

👑 **Administrator**

Retrieves all promotion codes.

- **Method**: `GET`
- **Path**: `/admin/promotion-codes`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
```json
[
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440000",
    "stripe_promotion_code_id": "promo_123456",
    "coupon_id": "770e8400-e29b-41d4-a716-446655440000",
    "code": "SUMMER20",
    "active": true,
    "max_redemptions": 100,
    "times_redeemed": 45,
    "expires_at": "2025-09-01T00:00:00Z",
    "first_time_transaction": false,
    "minimum_amount": 5000,
    "minimum_amount_currency": "eur",
    "restrictions": {
      "currency_options": ["eur", "usd"]
    },
    "metadata": {},
    "created_at": "2025-06-01T00:00:00Z"
  }
]
```

#### Status Codes
- **200 OK** - Promotion codes retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Active Promotion Codes (Admin)

👑 **Administrator**

Retrieves only active promotion codes.

- **Method**: `GET`
- **Path**: `/admin/promotion-codes/active`
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as Get All Promotion Codes, but filtered to active codes only.

#### Status Codes
- **200 OK** - Active promotion codes retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Promotion Code by ID (Admin)

👑 **Administrator**

Retrieves a single promotion code by its UUID.

- **Method**: `GET`
- **Path**: `/admin/promotion-codes/{id}`
- **Path Parameters**:
  - `id`: Promotion code UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as promotion code object in Get All Promotion Codes.

#### Status Codes
- **200 OK** - Promotion code retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Promotion code not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Get Promotion Code by Code String (Admin)

👑 **Administrator**

Retrieves a promotion code by its code string.

- **Method**: `GET`
- **Path**: `/admin/promotion-codes/code/{code}`
- **Path Parameters**:
  - `code`: Promotion code string
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (200 OK)
Same structure as promotion code object in Get All Promotion Codes.

#### Status Codes
- **200 OK** - Promotion code retrieved successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Promotion code not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Database temporarily unavailable

---

### Create Promotion Code (Admin)

👑 **Administrator**

Creates a new promotion code in both the database and Stripe.

- **Method**: `POST`
- **Path**: `/admin/promotion-codes`
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "couponId": "770e8400-e29b-41d4-a716-446655440000",
  "code": "WELCOME-2025",
  "maxRedemptions": 500,
  "expiresAt": "2025-12-31T23:59:59Z",
  "firstTimeTransaction": true,
  "minimumAmount": 3000,
  "minimumAmountCurrency": "eur",
  "restrictions": {
    "currencyOptions": ["eur"]
  },
  "metadata": {
    "campaign": "welcome"
  }
}
```

**Field Validations**:
- `couponId`: Required, valid coupon UUID
- `code`: Required, 3-50 characters, uppercase alphanumeric with `-` and `_` allowed
- `maxRedemptions`: Optional, positive integer
- `expiresAt`: Optional, ISO 8601 timestamp
- `firstTimeTransaction`: Required, boolean
- `minimumAmount`: Optional, positive integer in cents
- `minimumAmountCurrency`: Optional, 3-character ISO code (required if minimumAmount set)
- `restrictions.currencyOptions`: Optional, array of 3-character ISO codes
- `metadata`: Optional JSON object

#### Success Response (201 Created)
```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440001",
  "stripe_promotion_code_id": "promo_987654",
  "coupon_id": "770e8400-e29b-41d4-a716-446655440000",
  "code": "WELCOME-2025",
  "active": true,
  "max_redemptions": 500,
  "times_redeemed": 0,
  "expires_at": "2025-12-31T23:59:59Z",
  "first_time_transaction": true,
  "minimum_amount": 3000,
  "minimum_amount_currency": "eur",
  "restrictions": {
    "currency_options": ["eur"]
  },
  "metadata": {
    "campaign": "welcome"
  },
  "created_at": "2025-01-15T17:00:00Z"
}
```

#### Status Codes
- **201 Created** - Promotion code created successfully
- **400 Bad Request** - Validation error (invalid code format, missing currency, etc.)
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Coupon not found
- **409 Conflict** - Promotion code already exists
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Update Promotion Code (Admin)

👑 **Administrator**

Updates promotion code active status and metadata.

- **Method**: `PATCH`
- **Path**: `/admin/promotion-codes/{id}`
- **Path Parameters**:
  - `id`: Promotion code UUID
- **Content-Type**: `application/json`
- **Authentication**: Requires Administrator role

#### Request Body
```json
{
  "active": false,
  "metadata": {
    "campaign": "welcome",
    "deprecated": true
  }
}
```

**Field Validations**:
- `active`: Optional, boolean
- `metadata`: Optional JSON object
- At least one field must be provided

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Promotion code updated successfully
- **400 Bad Request** - No fields provided or validation error
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Promotion code not found
- **415 Unsupported Media Type** - Missing `application/json` Content-Type
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Deactivate Promotion Code (Admin)

👑 **Administrator**

Deactivates a promotion code, preventing new redemptions.

- **Method**: `POST`
- **Path**: `/admin/promotion-codes/{id}/deactivate`
- **Path Parameters**:
  - `id`: Promotion code UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Promotion code deactivated successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Promotion code not found
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

### Delete Promotion Code (Admin)

👑 **Administrator**

Deletes a promotion code from both database and Stripe.

- **Method**: `DELETE`
- **Path**: `/admin/promotion-codes/{id}`
- **Path Parameters**:
  - `id`: Promotion code UUID
- **Authentication**: Requires Administrator role

#### Request Body
None

#### Success Response (204 No Content)
No response body

#### Status Codes
- **204 No Content** - Promotion code deleted successfully
- **401 Unauthorized** - Not authenticated
- **403 Forbidden** - Not an administrator
- **404 Not Found** - Promotion code not found
- **500 Internal Server Error** - Server error
- **502 Bad Gateway** - Stripe API error
- **503 Service Unavailable** - Database or Stripe temporarily unavailable

---

## Error Response Format

All error responses follow a consistent format:

```json
{
  "error": "Human-readable error message",
  "details": {
    "field_name": "Specific validation error"
  },
  "code": "ERROR_CODE"
}
```

For validation errors (400 Bad Request), the `details` field contains field-specific error messages.

---

## Notes

### Content-Type Requirements

All `POST`, `PUT`, and `PATCH` endpoints require `Content-Type: application/json` header (except file uploads which use `multipart/form-data`). Missing or incorrect Content-Type will result in **415 Unsupported Media Type**.

### Authentication Cookies

After successful authentication, the API sets two HTTP-only cookies:
- `access_token`: Short-lived token for API requests
- `refresh_token`: Long-lived token for refreshing sessions

The frontend should include these cookies automatically in subsequent requests.

### Rate Limiting

Some endpoints (particularly OTP-related) implement rate limiting to prevent abuse. Excessive requests will result in **429 Too Many Requests**.

### Retryable Errors

**503 Service Unavailable** errors are generally retryable - they indicate temporary infrastructure issues (database connection pools, deadlocks, etc.). Clients should implement exponential backoff retry logic.

### CORS

All endpoints support Cross-Origin Resource Sharing (CORS) for browser-based clients.

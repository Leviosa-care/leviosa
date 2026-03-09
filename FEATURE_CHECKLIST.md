# Leviosa Feature Checklist

> Complete feature checklist for Leviosa - A wellness booking platform with comprehensive room scheduling, product catalog, and partner management.

---

## Authentication & User Management

### User Registration Flow
#### Frontend
- [ ] `POST /auth/email` - Email entry and validation
  - [ ] Email format validation
  - [ ] Check if email already exists
- [ ] `POST /auth/otp` - OTP verification
  - [ ] OTP input form (6-digit code)
  - [ ] Resend OTP functionality
  - [ ] OTP expiry timer
  - [ ] Store verified OTP in Redis (`otp_verified:{email}`)
- [ ] `POST /auth/general` - General information
  - [ ] First name, last name input
  - [ ] Age, gender selection
  - [ ] Form validation
- [ ] `POST /auth/address` - Address information
  - [ ] Street address, city, postal code
  - [ ] Country/region selection
  - [ ] Address validation
- [ ] `POST /auth/password` - Password creation
  - [ ] Password strength validation
  - [ ] Confirm password matching
  - [ ] Password visibility toggle
- [ ] `/auth/pending` - Pending user page
  - [ ] Redirect after successful registration
  - [ ] Limited access for pending users

#### Backend (authuser service)
- [ ] `POST /auth/send-otp` - Send OTP to email
  - [ ] Generate random 6-digit OTP
  - [ ] Store in Redis with TTL
  - [ ] Email sending via notification service
- [ ] `POST /auth/verify-otp` - Verify OTP code
  - [ ] Validate OTP from Redis
  - [ ] Create pending user on first verification
  - [ ] Return user session
- [ ] `POST /auth/register` - Complete registration
  - [ ] Validate all required fields
  - [ ] Encrypt sensitive data (GDPR compliance)
  - [ ] Create user in PostgreSQL
  - [ ] Update user status from pending to active
- [ ] `POST /auth/forgot-password` - Password reset flow
  - [ ] Email validation
  - [ ] Send reset OTP
  - [ ] Verify reset OTP
  - [ ] Update password

### User Sessions & Security
#### Frontend
- [ ] Session management via cookies
  - [ ] Access token cookie
  - [ ] Refresh token cookie
  - [ ] Automatic token refresh
- [ ] Protected route handling
  - [ ] Redirect to auth if not authenticated
  - [ ] Role-based route access
- [ ] Logout functionality
  - [ ] Clear cookies
  - [ ] Redirect to home

#### Backend (authuser service)
- [ ] `POST /auth/login` - User login
  - [ ] Email/password validation
  - [ ] Generate JWT access token
  - [ ] Generate refresh token
  - [ ] Set httpOnly cookies
- [ ] `POST /auth/logout` - User logout
  - [ ] Invalidate refresh token
  - [ ] Clear cookies
- [ ] `POST /auth/refresh` - Token refresh
  - [ ] Validate refresh token
  - [ ] Generate new access token
- [ ] `GET /users/me` - Get current user
  - [ ] Return user profile
  - [ ] Include role and status
- [ ] `PUT /users/me` - Update profile
  - [ ] Update personal information
  - [ ] Encrypt sensitive fields

### User Roles & Permissions
#### Backend
- [ ] Role definitions
  - [ ] Standard - Regular users
  - [ ] Partner - Service providers
  - [ ] Admin - Platform administrators
- [ ] Permission middleware
  - [ ] `require_minimum_role` - Minimum role check
  - [ ] `require_admin` - Admin-only access
  - [ ] `require_access_token` - JWT validation
  - [ ] `require_refresh_token` - Refresh token validation
- [ ] Partner management
  - [ ] `POST /partners` - Create partner
  - [ ] `GET /partners/{id}` - Get partner details
  - [ ] `PUT /partners/{id}` - Update partner
  - [ ] `DELETE /partners/{id}` - Delete partner
  - [ ] `GET /partners` - List all partners (admin)

---

## Product Catalog

### Categories Management
#### Frontend (ops/catalog)
- [ ] Categories page (`/ops/catalog`)
  - [ ] List all categories
  - [ ] Search and filter categories
  - [ ] Category cards with image
  - [ ] Create/Edit category modal

#### Backend (catalog service)
- [ ] `POST /categories` - Create category
  - [ ] Name validation
  - [ ] Description support
  - [ ] Image upload to S3
  - [ ] Encrypted sensitive fields
- [ ] `GET /categories` - List categories
  - [ ] Pagination support
  - [ ] Filter by active status
  - [ ] Search by name
- [ ] `GET /categories/{id}` - Get category details
- [ ] `PUT /categories/{id}` - Update category
- [ ] `DELETE /categories/{id}` - Delete category
  - [ ] Check for product dependencies
  - [ ] Soft delete option

### Products Management
#### Frontend (admin/products)
- [ ] Products list page (`/admin/products`)
  - [ ] Product cards with images
  - [ ] Category filter dropdown
  - [ ] Search functionality
  - [ ] Status indicators (active/inactive)
- [ ] Product creation/edit modal
  - [ ] Name input
  - [ ] Category selection
  - [ ] Description text area
  - [ ] Duration (minutes)
  - [ ] Buffer time (minutes)
  - [ ] Image upload
  - [ ] Status toggle
- [ ] Product detail view
  - [ ] Full description
  - [ ] Associated prices
  - [ ] Promotion codes
  - [ ] Image gallery

#### Backend (catalog service)
- [ ] `POST /products` - Create product
  - [ ] Validate required fields
  - [ ] Associate with category
  - [ ] Upload images to S3
  - [ ] Encrypt sensitive data
- [ ] `GET /products` - List products
  - [ ] Pagination
  - [ ] Filter by category
  - [ ] Filter by status
  - [ ] Search by name
  - [ ] Include prices
- [ ] `GET /products/{id}` - Get product details
  - [ ] Include category
  - [ ] Include images
  - [ ] Include prices
- [ ] `PUT /products/{id}` - Update product
- [ ] `DELETE /products/{id}` - Delete product
  - [ ] Check for booking dependencies
  - [ ] Delete associated images

### Pricing Management
#### Frontend (ops/catalog)
- [ ] Prices page (`/ops/catalog/Prices`)
  - [ ] List prices by product
  - [ ] Create price form
  - [ ] Amount input
  - [ ] Currency selection
  - [ ] Active date range

#### Backend (catalog service)
- [ ] `POST /prices` - Create price
  - [ ] Associate with product
  - [ ] Amount validation
  - [ ] Currency code (ISO 4217)
  - [ ] Start/end date for temporal pricing
- [ ] `GET /prices/product/{productId}` - Get prices by product
  - [ ] Return active prices
  - [ ] Filter by date range
- [ ] `GET /prices/{id}` - Get price details
- [ ] `PUT /prices/{id}` - Update price
- [ ] `DELETE /prices/{id}` - Delete price
  - [ ] Check if used in active bookings

### Coupons & Promotion Codes
#### Frontend (ops/catalog)
- [ ] Coupons page (`/ops/catalog/Coupons`)
  - [ ] List all coupons
  - [ ] Create coupon form
  - [ ] Discount type (percentage/fixed)
  - [ ] Discount amount
  - [ ] Usage limit
  - [ ] Expiry date
- [ ] Promotion Codes page (`/ops/catalog/PromotionCodes`)
  - [ ] List promotion codes
  - [ ] Generate codes for coupon
  - [ ] Code prefix
  - [ ] Batch generation

#### Backend (catalog service)
- [ ] `POST /coupons` - Create coupon
  - [ ] Discount type validation
  - [ ] Amount validation
  - [ ] Usage limits
  - [ ] Expiry date
  - [ ] Max uses per customer
- [ ] `GET /coupons` - List coupons
  - [ ] Filter by status
  - [ ] Pagination
- [ ] `GET /coupons/{id}` - Get coupon details
- [ ] `PUT /coupons/{id}` - Update coupon
- [ ] `DELETE /coupons/{id}` - Delete coupon
- [ ] `POST /promotion-codes` - Create promotion code
  - [ ] Associate with coupon
  - [ ] Generate unique code
  - [ ] Stripe integration
- [ ] `GET /promotion-codes` - List promotion codes
  - [ ] Filter by coupon
  - [ ] Filter by status
- [ ] `GET /promotion-codes/{id}` - Get promotion code
- [ ] `DELETE /promotion-codes/{id}` - Delete promotion code
  - [ ] Also deletes from Stripe

### Media Management
#### Backend (catalog service)
- [ ] `POST /images` - Upload image
  - [ ] Upload to S3
  - [ ] Associate with product or category
  - [ ] Alt text support
  - [ ] Caption support
  - [ ] Sort order
- [ ] `GET /images/{parentId}` - Get images by parent
  - [ ] Return sorted list
  - [ ] Include URLs
- [ ] `DELETE /images/{id}` - Delete image
  - [ ] Delete from S3
  - [ ] Delete database record

---

## Booking & Availability Management

### Building Management
#### Backend (booking service)
- [ ] `POST /buildings` - Create building
  - [ ] Name, address, contact info
  - [ ] Encrypt personal data (GDPR)
  - [ ] Active status
- [ ] `GET /buildings` - List buildings
  - [ ] Filter by active status
  - [ ] Pagination
- [ ] `GET /buildings/{id}` - Get building details
- [ ] `PUT /buildings/{id}` - Update building
- [ ] `DELETE /buildings/{id}` - Delete building
  - [ ] Check for room dependencies

### Room Management
#### Backend (booking service)
- [ ] `POST /buildings/{buildingId}/rooms` - Create room
  - [ ] Name, description
  - [ ] Room number
  - [ ] Capacity
  - [ ] Equipment specifications
  - [ ] Encrypt identification data
- [ ] `GET /buildings/{buildingId}/rooms` - List rooms in building
  - [ ] Filter by active status
- [ ] `GET /rooms/{id}` - Get room details
  - [ ] Include building info
- [ ] `PUT /rooms/{id}` - Update room
- [ ] `DELETE /rooms/{id}` - Delete room
  - [ ] Check for allocation/availability dependencies

### Room Allocations
#### Backend (booking service)
- [ ] `POST /allocations` - Create allocation
  - [ ] Dedicated vs Shared type
  - [ ] Partner assignment
  - [ ] Room assignment
  - [ ] Time period (for dedicated)
  - [ ] Active status
- [ ] `GET /partners/{partnerId}/allocations` - Get partner allocations
  - [ ] Include room details
  - [ ] Include time periods
- [ ] `GET /allocations/{id}` - Get allocation details
- [ ] `PUT /allocations/{id}` - Update allocation
- [ ] `DELETE /allocations/{id}` - Delete allocation
  - [ ] Check for active bookings

### Availability Slots
#### Frontend
- [ ] Partner availability calendar
  - [ ] Calendar view for room
  - [ ] Create availability modal
  - [ ] Date/time pickers
  - [ ] Recurring pattern selection
  - [ ] Service type selection
  - [ ] Capacity input
- [ ] Availability suggestions
  - [ ] Block duration suggestions
  - [ ] Based on partner products
  - [ ] Single/multi-session options

#### Backend (booking service)
- [ ] `POST /rooms/{roomId}/availability` - Create availability
  - [ ] Start/end time validation (min 15min, max 12hr)
  - [ ] Capacity validation (max 50)
  - [ ] Recurring pattern support
  - [ ] Check for overlaps
  - [ ] Validate partner has allocation
- [ ] `GET /rooms/{roomId}/availability` - List room availability
  - [ ] Date range filter
  - [ ] Status filter
  - [ ] Pagination
- [ ] `GET /availability/{id}` - Get availability details
- [ ] `PUT /availability/{id}` - Update availability
- [ ] `DELETE /availability/{id}` - Delete availability
  - [ ] Check for active bookings
- [ ] `GET /partners/{partnerId}/rooms/{roomId}/suggest-blocks` - Get block suggestions
  - [ ] Analyze partner products
  - [ ] Suggest standard durations
  - [ ] Suggest multi-session blocks
  - [ ] Priority ranking

### Bookings
#### Frontend (premium/bookings, admin/bookings)
- [ ] User booking page (`/premium/bookings`)
  - [ ] Available slots listing
  - [ ] Slot details (time, duration, price)
  - [ ] Booking creation form
  - [ ] Payment integration
- [ ] Admin bookings page (`/admin/bookings/events`, `/admin/bookings/consultations`)
  - [ ] List all bookings
  - [ ] Filter by date/status
  - [ ] Filter by room
  - [ ] View booking details
  - [ ] Cancel booking
  - [ ] Complete booking
- [ ] Booking detail view
  - [ ] Client information
  - [ ] Service details
  - [ ] Payment status
  - [ ] Notes from client/partner

#### Backend (booking service)
- [ ] `POST /bookings` - Create booking
  - [ ] Validate availability exists
  - [ ] Check availability is bookable
  - [ ] Check capacity limits
  - [ ] Calculate pricing
  - [ ] Create Stripe payment intent
  - [ ] Mark availability as booked
  - [ ] Atomic transaction
- [ ] `GET /bookings/{id}` - Get booking details
  - [ ] Include availability info
  - [ ] Include room/building info
  - [ ] Include client/partner info
- [ ] `GET /clients/{clientId}/bookings` - Get client bookings
  - [ ] Filter by status
  - [ ] Pagination
- [ ] `GET /partners/{partnerId}/bookings` - Get partner bookings
  - [ ] Filter by status
  - [ ] Pagination
- [ ] `PUT /bookings/{id}/notes` - Update booking notes
  - [ ] Client notes
  - [ ] Partner notes
  - [ ] Encrypt notes
- [ ] `POST /bookings/{id}/cancel` - Cancel booking
  - [ ] Validate booking can be cancelled
  - [ ] Update status to Cancelled
  - [ ] Refund via Stripe if applicable
  - [ ] Release availability
- [ ] `POST /bookings/{id}/complete` - Complete booking
  - [ ] Update status to Completed
  - [ ] Finalize payment
- [ ] `POST /bookings/{id}/payment` - Process payment
  - [ ] Create/confirm Stripe payment
  - [ ] Update payment status
  - [ ] Handle payment failures

### Analytics & Metrics
#### Backend (booking service)
- [ ] `GET /availabilities/rooms/{roomId}/gaps` - Find scheduling gaps
  - [ ] Analyze room schedule for date
  - [ ] Find gaps between bookings
  - [ ] Suggest products for each gap
  - [ ] Sort by duration
- [ ] `GET /rooms/{roomId}/metrics` - Get room metrics
  - [ ] Date range filter
  - [ ] Utilization percentage
  - [ ] Fragmentation count
  - [ ] Idle minutes
  - [ ] Efficiency score
- [ ] `GET /partners/{partnerId}/metrics` - Get partner metrics
  - [ ] Aggregate across rooms
  - [ ] Date range filter
  - [ ] Utilization trends
  - [ ] Efficiency tracking

---

## Settings & Configuration

### System Settings
#### Frontend (ops/settings)
- [ ] Settings page (`/ops/settings`)
  - [ ] Company information form
  - [ ] Logo upload
  - [ ] Contact details
  - [ ] Social media links
- [ ] OTP Settings
  - [ ] OTP length configuration
  - [ ] Expiry time
  - [ ] Max attempts
- [ ] Notification Settings
  - [ ] Email configuration
  - [ ] SMS configuration
  - [ ] Notification templates

#### Backend (settings service)
- [ ] `POST /settings` - Set setting value
  - [ ] String values
  - [ ] Encrypted values
  - [ ] Integer values
  - [ ] JSON values
- [ ] `GET /settings` - Get all settings
  - [ ] Include encrypted values
  - [ ] Bulk endpoint
- [ ] `GET /settings/{key}` - Get single setting
- [ ] `DELETE /settings/{key}` - Delete setting
- [ ] RabbitMQ integration for updates

---

## Notifications

### Email Notifications
#### Backend (notification service)
- [ ] `POST /notifications/send-email` - Send email
  - [ ] Recipient validation
  - [ ] Subject and body
  - [ ] HTML support
  - [ ] Template support
- [ ] `POST /notifications/send-bulk-email` - Bulk emails
  - [ ] Multiple recipients
  - [ ] Queue processing
- [ ] Email templates
  - [ ] OTP verification email
  - [ ] Booking confirmation
  - [ ] Booking cancellation
  - [ ] Payment receipt
  - [ ] Password reset

### SMS Notifications
#### Backend (notification service)
- [ ] `POST /notifications/send-sms` - Send SMS
  - [ ] Phone number validation
  - [ ] Message content
  - [ ] Character limit
- [ ] SMS templates
  - [ ] OTP verification
  - [ ] Booking reminders
  - [ ] Cancellation notices

---

## Marketing Pages

### Public Pages
#### Frontend (marketing routes)
- [ ] Homepage (`/`)
  - [ ] Hero section
  - [ ] Services overview
  - [ ] Why choose us section
  - [ ] Testimonials
  - [ ] CTA to book
- [ ] Services page (`/services`)
  - [ ] Service cards
  - [ ] Category breakdown
  - [ ] Pricing information
- [ ] About page (`/about`)
  - [ ] Company story
  - [ ] Team profiles
  - [ ] Mission/values
- [ ] Team page (`/team`)
  - [ ] Team member cards
  - [ ] Roles and bios
  - [ ] Social links
- [ ] Book page (`/book`)
  - [ ] Service selection
  - [ ] Availability calendar
  - [ ] Booking flow
- [ ] Footer component
  - [ ] Navigation links
  - [ ] Social media links
  - [ ] Contact information

### Legal Pages
#### Frontend (legal routes)
- [ ] Terms of Service (`/legal/terms`)
  - [ ] User agreement
  - [ ] Service terms
  - [ ] Cancellation policy
- [ ] Privacy Policy (`/legal/privacy`)
  - [ ] Data collection
  - [ ] GDPR compliance
  - [ ] Cookie policy

---

## Operations Dashboard

### Dashboard Overview
#### Frontend (ops routes)
- [ ] Dashboard home (`/ops`)
  - [ ] Activity cards
  - [ ] Volume metrics
  - [ ] Agenda view
  - [ ] Quick actions

### User Management
- [ ] Users page (`/ops/users`)
  - [ ] List all users
  - [ ] Filter by role
  - [ ] Filter by status
  - [ ] User detail view
  - [ ] Edit user
  - [ ] Delete user
  - [ ] Change role

---

## Infrastructure & Security

### GDPR Compliance
- [ ] Data encryption at rest
  - [ ] Personal data encrypted with encx
  - [ ] Encryption key management via Vault
- [ ] Right to erasure
  - [ ] User data deletion
  - [ ] Booking data anonymization
- [ ] Data export
  - [ ] User data export endpoint
  - [ ] GDPR-compliant format

### Security Features
- [ ] Authentication
  - [ ] JWT token validation
  - [ ] Refresh token rotation
  - [ ] Session management
- [ ] Authorization
  - [ ] Role-based access control
  - [ ] Resource ownership validation
  - [ ] Admin-only endpoints
- [ ] Data Protection
  - [ ] SQL injection prevention
  - [ ] XSS protection
  - [ ] CSRF protection
  - [ ] Rate limiting

### Monitoring & Logging
- [ ] Application logging
  - [ ] Structured logging
  - [ ] Request/response logging
  - [ ] Error tracking
- [ ] Metrics
  - [ ] Room utilization
  - [ ] Booking trends
  - [ ] Partner performance
- [ ] Health checks
  - [ ] Database connectivity
  - [ ] Redis connectivity
  - [ ] RabbitMQ connectivity
  - [ ] External service status

---

## Technical Implementation Notes

### Database Schema
- **auth schema**: users, partners, sessions
- **catalog schema**: categories, products, prices, coupons, promotion_codes, images
- **booking schema**: buildings, rooms, room_allocations, availabilities, bookings, payments, room_daily_metrics
- **settings schema**: application_settings (encrypted and plain)

### External Integrations
- **Stripe**: Payment processing, promotion codes
- **AWS S3**: Image and file storage
- **Twilio**: SMS notifications
- **Gmail SMTP**: Email notifications
- **HashiCorp Vault**: Encryption key management

### Message Queue (RabbitMQ)
- Settings update notifications
- Async notification processing
- Booking state changes

### Caching (Redis)
- Session storage
- OTP storage
- Frequently accessed data

---

## Development Status Legend

- [ ] Not started
- [~] In progress
- [x] Complete

---

## Notes

This checklist represents the full scope of the Leviosa platform:

1. **Frontend Routes**:
   - `(app)/admin/*` - Admin management interface
   - `(app)/premium/*` - Premium user booking interface
   - `(ops)/*` - Operations dashboard
   - `(marketing)/*` - Public marketing pages
   - `(legal)/*` - Legal documentation
   - `auth/*` - Authentication flow

2. **Backend Services**:
   - **authuser** - Authentication, user/partner management
   - **catalog** - Products, categories, prices, promotions, images
   - **booking** - Buildings, rooms, allocations, availability, bookings
   - **settings** - System configuration
   - **notification** - Email/SMS notifications

3. **Key Business Flows**:
   - Multi-step registration with OTP verification
   - Partner onboarding with room allocation
   - Availability creation with recurring patterns
   - Booking creation with Stripe payment
   - Utilization analytics and gap detection

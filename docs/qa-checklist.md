# QA Test Checklist

## Authentication & Onboarding

- [ ] Sign in with valid credentials
- [ ] Sign in with invalid credentials (error state)
- [ ] Sign up flow — complete all steps (general info, address, password, email verification)
- [ ] Complete partner onboarding (`/auth/complete-partner`)
- [ ] Forgotten password flow
- [ ] Email verification (OTP code)
- [ ] OAuth link / unlink a provider
- [ ] Logout

---

## Public Pages

- [ ] Home page loads with hero, services, team, stats, FAQ sections
- [ ] `/services` lists all available services
- [ ] `/services/[id]` shows a specific service detail
- [ ] `/team` lists staff members
- [ ] `/team/[id]` shows a staff member profile
- [ ] `/book` — booking flow: select service, date/time, confirm
- [ ] `/book/confirmation` — confirmation page loads after booking; guest claim card visible if applicable
- [ ] `/bookings` — public booking lookup works
- [ ] Legal pages (privacy, terms) render

---

## Client Area (`/client`)

- [ ] Dashboard loads with correct data
- [ ] `/client/bookings` lists all bookings for the logged-in client
- [ ] `/client/bookings/[id]` shows booking detail
- [ ] Cancel a booking from the client side
- [ ] `/client/messages` — thread list loads, can send a message
- [ ] `/client/profile` — view and update profile

---

## Staff / Partner Area (`/staff`)

- [ ] Dashboard loads
- [ ] `/staff/agenda/disponibilites` — view availabilities
- [ ] Create a single availability slot
- [ ] Create recurring availability slots
- [ ] Edit an availability
- [ ] Delete an availability
- [ ] Cancel an availability
- [ ] `/staff/agenda/reservations` — view bookings
- [ ] Complete a booking (mark as done)
- [ ] Mark a booking as no-show
- [ ] Add notes to a booking
- [ ] `/staff/catalog` — view categories, products, prices, exercises, coupons, promo codes
- [ ] `/staff/messages` — send and receive messages
- [ ] `/staff/profile` — update profile
- [ ] `/staff/settings` — save settings
- [ ] `/staff/statistics/analytics` loads
- [ ] `/staff/statistics/finances` loads
- [ ] Stripe onboarding link generates and redirects

---

## Admin — Catalog (`/admin/catalog`)

- [ ] Product list loads with search/filter
- [ ] Create a new product (`/admin/catalog/products/new`)
- [ ] View product detail (`/admin/catalog/products/[id]`)
- [ ] Edit an existing product
- [ ] Delete a product (confirm dialog)
- [ ] Create a category
- [ ] Edit / delete a category
- [ ] Add a price to a product
- [ ] View prices for a product

---

## Admin — Users & Misc (`/admin`)

- [ ] Dashboard loads with summary cards
- [ ] `/admin/users` — user list loads
- [ ] `/admin/bookings/consultations` — consultation bookings list
- [ ] `/admin/buildings` — buildings list loads
- [ ] `/admin/planning` — planning view loads
- [ ] `/admin/messages` — messages overview
- [ ] `/admin/compta` — accounting view loads
- [ ] `/admin/analytics` — analytics view loads

---

## Cross-Cutting / Edge Cases

- [ ] Unauthenticated access to a protected route redirects to `/auth`
- [ ] Wrong role accessing a route (e.g. client hitting `/admin`) is blocked
- [ ] `/healthz` returns 200
- [ ] Availability conflict check returns the right response when booking overlaps
- [ ] Form validation errors display correctly (required fields, wrong formats)
- [ ] Network error / API down — graceful error state shown

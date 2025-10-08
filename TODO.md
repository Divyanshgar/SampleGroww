# Fix OTP Verification Error Display Issue

## Overview
The issue is that invalid or expired OTP errors are shown on the registration page instead of the OTP page. This happens because the /verify-otp-only endpoint marks the OTP as verified, causing the /verify-otp endpoint to fail.

## Tasks

### Backend Fixes
- [ ] Modify `utils/otp_store.go`:
  - Rename `VerifyOTP` to `CheckOTP` and remove the line that marks OTP as verified.
  - Add new `VerifyAndMarkOTP` function that calls `CheckOTP` and marks as verified if valid.
- [ ] Update `handlers/notification_handler.go`:
  - In `VerifyOTP` handler, use `CheckOTP` instead of `VerifyOTP`.
  - In `VerifyOTPAndRegister` handler, use `VerifyAndMarkOTP` instead of `VerifyOTP`.

### Testing
- [ ] Test invalid OTP on OTP page: should show error on OTP page, not proceed to registration.
- [ ] Test valid OTP: should proceed to registration page.
- [ ] Test registration with valid OTP: should succeed.
- [ ] Test expired OTP on registration: should show error on registration page (acceptable).

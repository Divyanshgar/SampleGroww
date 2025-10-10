# TODO: Add PDF Attachment to User Registration Email

## Steps to Complete:

1. **Update `handlers/notification_handler.go`**:
   - In the `VerifyOTPAndRegister` function, after updating user details:
     - Generate PDF using `h.pdfService.GenerateUserProfilePDF(&user)`.
     - Replace `SendUserProfileEmail` with `SendUserProfileEmailWithAttachment(user, pdfBytes)`.
     - Add error handling for PDF generation and email sending (log errors, do not block registration).

2. **Test the changes**:
   - Run the server and test the full registration flow (request OTP, verify OTP, complete registration).
   - Verify that the welcome email includes the PDF attachment.
   - Check server logs for any errors.

3. **Verify and complete**:
   - Confirm the PDF is generated correctly and attached to the email.
   - Update this TODO with completion status.

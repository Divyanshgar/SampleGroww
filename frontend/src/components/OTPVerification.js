import React, { useState, useRef, useEffect } from 'react';

function OTPVerification({ email, onVerified, onBack }) {
  const [otp, setOtp] = useState(['', '', '', '', '', '']);
  const [error, setError] = useState('');
  const inputRefs = useRef([]);

  useEffect(() => {
    // Focus on first input when component mounts
    inputRefs.current[0]?.focus();
  }, []);

  const handleChange = (index, value) => {
    // Only allow numbers
    if (value && !/^\d$/.test(value)) return;

    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);

    // Move to next input if value is entered
    if (value && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (index, e) => {
    // Move to previous input on backspace if current input is empty
    if (e.key === 'Backspace' && !otp[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
  };

  const handlePaste = (e) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData('text').slice(0, 6).split('');
    const newOtp = [...otp];
    
    pastedData.forEach((char, index) => {
      if (/^\d$/.test(char) && index < 6) {
        newOtp[index] = char;
      }
    });
    
    setOtp(newOtp);
    
    // Focus on the last filled input or the next empty one
    const nextEmptyIndex = newOtp.findIndex(val => !val);
    const focusIndex = nextEmptyIndex === -1 ? 5 : nextEmptyIndex;
    inputRefs.current[focusIndex]?.focus();
  };

  const isOTPComplete = otp.every(digit => digit !== '');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (isOTPComplete) {
      const otpString = otp.join('');
      // Store OTP in sessionStorage to use in registration
      sessionStorage.setItem('otp', otpString);
      onVerified();
    }
  };

  return (
    <div className="otp-verification">
      <h2 style={{ textAlign: 'center', marginBottom: '10px', color: '#333' }}>
        Verify Your Email
      </h2>
      <p style={{ textAlign: 'center', color: '#666', marginBottom: '30px', fontSize: '14px' }}>
        We've sent a 6-digit code to<br />
        <strong>{email}</strong>
      </p>

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="otp-inputs">
          {otp.map((digit, index) => (
            <input
              key={index}
              ref={(el) => (inputRefs.current[index] = el)}
              type="text"
              maxLength="1"
              value={digit}
              onChange={(e) => handleChange(index, e.target.value)}
              onKeyDown={(e) => handleKeyDown(index, e)}
              onPaste={handlePaste}
              className="otp-input"
            />
          ))}
        </div>

        <button
          type="submit"
          className="btn btn-primary"
          disabled={!isOTPComplete}
        >
          Verify OTP
        </button>

        <button
          type="button"
          className="btn btn-secondary"
          onClick={onBack}
        >
          Change Email
        </button>
      </form>

      <div className="info-message" style={{ marginTop: '20px' }}>
        💡 Tip: Check your spam folder if you don't see the email
      </div>
    </div>
  );
}

export default OTPVerification;

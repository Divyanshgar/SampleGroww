 import React, { useState } from 'react';
import axios from 'axios';

function RegistrationForm({ email, onSuccess, onBack }) {
  const [formData, setFormData] = useState({
    first_name: '',
    last_name: '',
    phone: '',
    address: '',
    city: '',
    country: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [resendLoading, setResendLoading] = useState(false);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const otp = sessionStorage.getItem('otp');
      const response = await axios.post('/api/v1/auth/verify-otp', {
        email,
        otp,
        ...formData,
      });

      sessionStorage.removeItem('otp');
      onSuccess(response.data);
    } catch (err) {
      setError(
        err.response?.data?.error || 'Registration failed. Please try again.'
      );
    } finally {
      setLoading(false);
    }
  };

  const handleResend = async () => {
    setError('');
    setResendLoading(true);
    try {
      await axios.post('/api/v1/auth/request-otp', { email });
      // Go back to OTP verification step
      onBack();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to resend OTP');
    } finally {
      setResendLoading(false);
    }
  };

  return (
    <div className="registration-form">
      <h2 style={{ textAlign: 'center', marginBottom: '10px', color: '#333' }}>
        Complete Your Profile
      </h2>
      <p style={{ textAlign: 'center', color: '#666', marginBottom: '30px', fontSize: '14px' }}>
        Just a few more details to get you started
      </p>

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="first_name">First Name *</label>
          <input
            type="text"
            id="first_name"
            name="first_name"
            value={formData.first_name}
            onChange={handleChange}
            placeholder="Enter your first name"
            required
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="last_name">Last Name *</label>
          <input
            type="text"
            id="last_name"
            name="last_name"
            value={formData.last_name}
            onChange={handleChange}
            placeholder="Enter your last name"
            required
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="phone">Phone Number</label>
          <input
            type="tel"
            id="phone"
            name="phone"
            value={formData.phone}
            onChange={handleChange}
            placeholder="+1234567890"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="address">Address</label>
          <input
            type="text"
            id="address"
            name="address"
            value={formData.address}
            onChange={handleChange}
            placeholder="Enter your address"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="city">City</label>
          <input
            type="text"
            id="city"
            name="city"
            value={formData.city}
            onChange={handleChange}
            placeholder="Enter your city"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="country">Country</label>
          <input
            type="text"
            id="country"
            name="country"
            value={formData.country}
            onChange={handleChange}
            placeholder="Enter your country"
            disabled={loading}
          />
        </div>

        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? (
            <>
              <span className="loading"></span>
              Registering...
            </>
          ) : (
            'Complete Registration'
          )}
        </button>

        <button
          type="button"
          className="btn btn-secondary"
          onClick={onBack}
          disabled={loading}
        >
          Back
        </button>

        {error === 'Invalid or expired OTP' && (
          <button
            type="button"
            className="btn btn-link"
            onClick={handleResend}
            disabled={resendLoading}
            style={{ marginTop: '10px', fontSize: '14px' }}
          >
            {resendLoading ? 'Resending...' : 'Resend OTP'}
          </button>
        )}
      </form>

      <div className="info-message" style={{ marginTop: '20px', fontSize: '12px' }}>
        📧 You'll receive a welcome email at <strong>{email}</strong> after registration
      </div>
    </div>
  );
}

export default RegistrationForm;

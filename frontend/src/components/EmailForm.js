import React, { useState } from 'react';
import axios from 'axios';

function EmailForm({ onSubmit }) {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await axios.post('/api/v1/auth/request-otp', { email });
      console.log('OTP requested:', response.data);
      onSubmit(email);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to send OTP. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="email-form">
      <div style={{ textAlign: 'center', marginBottom: '20px' }}>
        <img
          src="/static/logo.jpg"
          alt="Centricity Financial Distribution Private Limited Logo"
          style={{
            maxWidth: '150px',
            height: 'auto',
            marginBottom: '10px'
          }}
          onError={(e) => {
            e.target.style.display = 'none';
          }}
        />
        <h2 style={{ textAlign: 'center', marginBottom: '10px', color: '#333' }}>
          Enter Your Email
        </h2>
        <p style={{ textAlign: 'center', color: '#666', marginBottom: '30px', fontSize: '14px' }}>
          We'll send you a verification code to confirm your email address
        </p>
      </div>

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="email">Email Address *</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter your email"
            required
            disabled={loading}
          />
        </div>

        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? (
            <>
              <span className="loading"></span>
              Sending OTP...
            </>
          ) : (
            'Send OTP'
          )}
        </button>
      </form>
    </div>
  );
}

export default EmailForm;

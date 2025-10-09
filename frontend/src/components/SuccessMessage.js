import React from 'react';

function SuccessMessage({ userData, onStartOver }) {
  const handleGeneratePDF = () => {
    if (userData?.user_id) {
      const url = `http://localhost:8443/api/v1/pdf/generate-user-profile/${userData.user_id}`;
      window.open(url, '_blank');
    }
  };

  return (
    <div className="success-container">
      <div style={{ textAlign: 'center', marginBottom: '30px' }}>
        <div style={{ fontSize: '64px', marginBottom: '20px' }}>🎉</div>
        <h2 style={{ color: '#667eea', marginBottom: '10px' }}>
          Registration Successful!
        </h2>
        <p style={{ color: '#666', fontSize: '14px' }}>
          Your account has been created successfully
        </p>
      </div>

      <div className="success-message">
        ✅ {userData?.message || 'User registered successfully'}
      </div>

      <div className="user-info">
        <h3>Your Account Details</h3>
        <div className="user-info-item">
          <label>User ID:</label>
          <span>{userData?.user_id}</span>
        </div>
        <div className="user-info-item">
          <label>Email:</label>
          <span>{userData?.email}</span>
        </div>
      </div>

      <div className="info-message">
        📧 A welcome email has been sent to your inbox with your complete profile details.
      </div>

      <button
        type="button"
        className="btn btn-primary"
        onClick={onStartOver}
        style={{ marginTop: '20px', marginRight: '10px' }}
      >
        Register Another User
      </button>

      <button
        type="button"
        className="btn btn-secondary"
        onClick={handleGeneratePDF}
        style={{ marginTop: '20px' }}
      >
        Generate PDF
      </button>

      <div style={{ textAlign: 'center', marginTop: '30px', color: '#999', fontSize: '12px' }}>
        <p>Thank you for registering with our service!</p>
      </div>
    </div>
  );
}

export default SuccessMessage;

import React from 'react';

function SuccessMessage({ userData, onStartOver }) {
const handleGeneratePDF = async () => {
  if (userData?.user_id) {
    try {
      const response = await fetch(`/api/v1/pdf/generate-user-profile/${userData.user_id}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'user_profile.pdf';
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      console.error('Error downloading PDF:', error);
      alert('Failed to download PDF. Please try again.');
    }
  }
};

const handleGenerateExcel = async () => {
  if (userData?.user_id) {
    try {
      const response = await fetch(`/api/v1/excel/generate-user-profile/${userData.user_id}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'user_profile.xlsx';
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      console.error('Error downloading Excel:', error);
      alert('Failed to download Excel. Please try again.');
    }
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
        style={{ marginTop: '20px', marginRight: '10px' }}
      >
        📄 Generate PDF
      </button>

      <button
        type="button"
        className="btn btn-secondary"
        onClick={handleGenerateExcel}
        style={{ marginTop: '20px' }}
      >
        📊 Download Excel
      </button>

      <div style={{ textAlign: 'center', marginTop: '30px', color: '#999', fontSize: '12px' }}>
        <p>Thank you for registering with our service!</p>
      </div>
    </div>
  );
}

export default SuccessMessage;

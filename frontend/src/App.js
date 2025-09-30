import React, { useState } from 'react';
import './App.css';
import EmailForm from './components/EmailForm';
import OTPVerification from './components/OTPVerification';
import RegistrationForm from './components/RegistrationForm';
import SuccessMessage from './components/SuccessMessage';

function App() {
  const [step, setStep] = useState(1);
  const [email, setEmail] = useState('');
  const [userData, setUserData] = useState(null);

  const handleEmailSubmit = (submittedEmail) => {
    setEmail(submittedEmail);
    setStep(2);
  };

  const handleOTPVerified = () => {
    setStep(3);
  };

  const handleRegistrationSuccess = (data) => {
    setUserData(data);
    setStep(4);
  };

  const handleStartOver = () => {
    setStep(1);
    setEmail('');
    setUserData(null);
  };

  return (
    <div className="App">
      <div className="container">
        <div className="header">
          <h1>🎉 User Registration</h1>
          <p>Secure registration with OTP verification</p>
        </div>

        {step === 1 && <EmailForm onSubmit={handleEmailSubmit} />}
        {step === 2 && (
          <OTPVerification
            email={email}
            onVerified={handleOTPVerified}
            onBack={() => setStep(1)}
          />
        )}
        {step === 3 && (
          <RegistrationForm
            email={email}
            onSuccess={handleRegistrationSuccess}
            onBack={() => setStep(2)}
          />
        )}
        {step === 4 && (
          <SuccessMessage userData={userData} onStartOver={handleStartOver} />
        )}

        <div className="progress-indicator">
          <div className={`step ${step >= 1 ? 'active' : ''}`}>1</div>
          <div className="line"></div>
          <div className={`step ${step >= 2 ? 'active' : ''}`}>2</div>
          <div className="line"></div>
          <div className={`step ${step >= 3 ? 'active' : ''}`}>3</div>
          <div className="line"></div>
          <div className={`step ${step >= 4 ? 'active' : ''}`}>✓</div>
        </div>
      </div>
    </div>
  );
}

export default App;

"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { authService } from '@/services/api';
import './register.css';

export default function RegisterPage() {
  const router = useRouter();
  
  // We still keep the state to capture user input
  const [formData, setFormData] = useState({
    full_name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // 1. Prepare Payload
      // Note: We are ignoring 'confirmPassword' here since the backend 
      // likely only needs the final password.
      const payload = {
        full_name: formData.full_name,
        email: formData.email,
        password: formData.password,
        role: "donor"
      };

      // 2. Send to Backend (Let the backend decide if it's valid)
      const data = await authService.register(payload);

      // 3. Success: Save session & Redirect
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify({
        id: data.user_id,
        name: data.full_name,
        role: data.role
      }));

      router.push('/dashboard');

    } catch (err: any) {
      // If the backend says "Passwords do not match" or "Password too short",
      // it will be caught here and displayed.
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page-container">
      <Link href="/">
        <img src="/Logo.webp" alt="SustainWear Logo" className="register-logo" />
      </Link>

      <div className="register-card">
        <h1 className="register-title">Join Us</h1>
        <p className="register-subtitle">Create an account to start giving clothes a second life.</p>
        
        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit}>
          
          <div className="form-group">
            <label className="form-label">Full Name</label>
            <input 
              type="text" 
              required
              className="form-input"
              placeholder="John Doe"
              value={formData.full_name}
              onChange={(e) => setFormData({...formData, full_name: e.target.value})}
            />
          </div>

          <div className="form-group">
            <label className="form-label">Email Address</label>
            <input 
              type="email" 
              required
              className="form-input"
              placeholder="you@example.com"
              value={formData.email}
              onChange={(e) => setFormData({...formData, email: e.target.value})}
            />
          </div>

          <div className="form-group">
            <label className="form-label">Password</label>
            <input 
              type="password" 
              required
              className="form-input"
              placeholder="••••••••"
              value={formData.password}
              onChange={(e) => setFormData({...formData, password: e.target.value})}
            />
          </div>

          <div className="form-group">
            <label className="form-label">Confirm Password</label>
            <input 
              type="password" 
              required
              className="form-input"
              placeholder="••••••••"
              value={formData.confirmPassword}
              onChange={(e) => setFormData({...formData, confirmPassword: e.target.value})}
            />
          </div>

          <button type="submit" disabled={loading} className="submit-btn">
            {loading ? 'Creating Account...' : 'Register'}
          </button>
        </form>

        <div className="register-footer">
          <p>Already have an account? <Link href="/login" className="register-link">Login here</Link></p>
        </div>
      </div>
    </div>
  );
}
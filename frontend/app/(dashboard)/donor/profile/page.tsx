"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Sidebar from '@/components/dashboard/Sidebar';
import { userService } from '@/services/api';
import './profile.css';

interface UserProfile {
    id: number;
    full_name: string;
    email: string;
    role: string;
    created_at: string;
}

export default function ProfilePage() {
    const router = useRouter();
    const [loading, setLoading] = useState(true);
    const [isEditing, setIsEditing] = useState(false);
    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [formData, setFormData] = useState({ full_name: '' });
    const [message, setMessage] = useState({ type: '', text: '' });

    useEffect(() => {
        fetchProfile();
    }, []);

    const fetchProfile = async () => {
        try {
            const data = await userService.getProfile();
            setProfile(data);
            setFormData({ full_name: data.full_name });
        } catch (err) {
            console.error(err);
            router.push('/login');
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault();
        setMessage({ type: '', text: '' });

        try {
            const updated = await userService.updateProfile(formData);
            setProfile(updated);
            setIsEditing(false);
            setMessage({ type: 'success', text: 'Profile updated successfully' });

            const storedUser = localStorage.getItem('user');
            if (storedUser) {
                const parsed = JSON.parse(storedUser);
                localStorage.setItem('user', JSON.stringify({ ...parsed, name: updated.full_name }));
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message || 'Failed to update profile' });
        }
    };

    if (loading) return null;

    return (
        <div className="dashboard-container">
            <Sidebar role="donor" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>My Profile</h1>
                    <p className="header-date">Manage your account settings and personal details.</p>
                </header>

                <div className="profile-container">
                    <div className="profile-card">
                        <div className="profile-header">
                            <div className="avatar-circle">
                                {profile?.full_name?.charAt(0).toUpperCase()}
                            </div>
                            <div className="profile-title">
                                <h2>{profile?.full_name}</h2>
                                <span className="role-badge">{profile?.role}</span>
                            </div>
                        </div>

                        {message.text && (
                            <div className={`message-alert ${message.type}`}>
                                {message.text}
                            </div>
                        )}

                        <form onSubmit={handleSave} className="profile-form">
                            <div className="form-group">
                                <label>Full Name</label>
                                <input
                                    type="text"
                                    className={`form-input ${isEditing ? 'editable' : ''}`}
                                    value={formData.full_name}
                                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                                    disabled={!isEditing}
                                    required
                                />
                            </div>

                            <div className="form-group">
                                <label>Email Address</label>
                                <input
                                    type="email"
                                    className="form-input"
                                    value={profile?.email || ''}
                                    disabled
                                    title="Email cannot be changed"
                                />
                                <span className="input-hint">Email cannot be changed</span>
                            </div>

                            <div className="form-group">
                                <label>Member Since</label>
                                <input
                                    type="text"
                                    className="form-input"
                                    value={profile?.created_at ? new Date(profile.created_at).toLocaleDateString() : ''}
                                    disabled
                                />
                            </div>

                            <div className="form-actions">
                                {isEditing ? (
                                    <>
                                        <button type="button" className="cancel-btn" onClick={() => {
                                            setIsEditing(false);
                                            setFormData({ full_name: profile?.full_name || '' });
                                            setMessage({ type: '', text: '' });
                                        }}>
                                            Cancel
                                        </button>
                                        <button type="submit" className="save-btn">
                                            Save Changes
                                        </button>
                                    </>
                                ) : (
                                    <button type="button" className="edit-btn" onClick={() => setIsEditing(true)}>
                                        Edit Profile
                                    </button>
                                )}
                            </div>
                        </form>
                    </div>
                </div>
            </main>
        </div>
    );
}
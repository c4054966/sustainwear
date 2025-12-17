"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { userService, organisationService, authService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './profile.css';

interface UserProfile {
    id: number;
    email: string;
    full_name: string;
    role: string;
    org_id: number | null;
}

interface OrgDetails {
    id: number;
    name: string;
    email: string;
    phone: string;
    address: string;
}

export default function StaffProfile() {
    const router = useRouter();
    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [org, setOrg] = useState<OrgDetails | null>(null);
    const [loading, setLoading] = useState(true);
    const [isEditing, setIsEditing] = useState(false);
    const [newName, setNewName] = useState("");
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        fetchProfileData();
    }, []);

    const fetchProfileData = async () => {
        try {
            const user = await userService.getProfile();
            setProfile(user);
            setNewName(user.full_name);

            if (user.org_id) {
                const orgData = await organisationService.getById(user.org_id);
                setOrg(orgData);
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        if (!newName.trim()) return;
        setSaving(true);
        try {
            await userService.updateProfile({ full_name: newName });
            if (profile) {
                setProfile({ ...profile, full_name: newName });
            }
            setIsEditing(false);
        } catch (error) {
            alert("Failed to update profile");
        } finally {
            setSaving(false);
        }
    };

    const handleLogout = () => {
        authService.logout();
        router.push('/login');
    };

    if (loading) return <div className="loading-screen">Loading profile...</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="staff" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>My Profile</h1>
                    <p className="header-date">Manage your account settings</p>
                </header>

                <div className="profile-layout">
                    <div className="profile-card user-details">
                        <div className="card-header">
                            <h2>Personal Information</h2>
                            {!isEditing && (
                                <button className="edit-btn" onClick={() => setIsEditing(true)}>Edit</button>
                            )}
                        </div>

                        <div className="form-group">
                            <label>Full Name</label>
                            {isEditing ? (
                                <input
                                    type="text"
                                    value={newName}
                                    onChange={(e) => setNewName(e.target.value)}
                                    className="profile-input"
                                />
                            ) : (
                                <p className="read-only-text">{profile?.full_name}</p>
                            )}
                        </div>

                        <div className="form-group">
                            <label>Email Address</label>
                            <p className="read-only-text">{profile?.email}</p>
                        </div>

                        <div className="form-group">
                            <label>Role</label>
                            <span className="role-badge">{profile?.role}</span>
                        </div>

                        {isEditing && (
                            <div className="edit-actions">
                                <button className="cancel-btn" onClick={() => {
                                    setIsEditing(false);
                                    setNewName(profile?.full_name || "");
                                }}>Cancel</button>
                                <button className="save-btn" onClick={handleSave} disabled={saving}>
                                    {saving ? 'Saving...' : 'Save Changes'}
                                </button>
                            </div>
                        )}
                    </div>

                    {org && (
                        <div className="profile-card org-details">
                            <h2>Organisation Details</h2>

                            <div className="info-row">
                                <span className="label">Organisation Name</span>
                                <span className="value fw-bold">{org.name}</span>
                            </div>

                            <div className="info-row">
                                <span className="label">Contact Email</span>
                                <span className="value">{org.email}</span>
                            </div>

                            <div className="info-row">
                                <span className="label">Phone</span>
                                <span className="value">{org.phone}</span>
                            </div>

                            <div className="info-row">
                                <span className="label">Address</span>
                                <span className="value">{org.address}</span>
                            </div>
                        </div>
                    )}

                    <div className="logout-section">
                        <button className="logout-btn-large" onClick={handleLogout}>
                            Sign Out
                        </button>
                    </div>
                </div>
            </main>
        </div>
    );
}
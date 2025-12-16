"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { authService, donationService } from '@/services/api';
import './donor_dashboard.css'; // <--- Import the specific CSS for this folder

export default function DonorDashboard() {
    const router = useRouter();
    const [user, setUser] = useState<any>(null);
    const [donations, setDonations] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // 1. Check if user is logged in
        const token = localStorage.getItem('token');
        const userData = localStorage.getItem('user');

        if (!token || !userData) {
            router.push('/login');
            return;
        }

        const parsedUser = JSON.parse(userData);

        // 2. Security Check: Ensure this user is actually a Donor
        if (parsedUser.role !== 'donor') {
            // If an admin tries to access this page, bounce them to their own dashboard
            if (parsedUser.role === 'admin') router.push('/admin');
            else if (parsedUser.role === 'staff') router.push('/staff');
            return;
        }

        setUser(parsedUser);
        fetchDashboardData();
    }, [router]);

    const fetchDashboardData = async () => {
        try {
            // Fetch recent donations
            const data = await donationService.getMyDonations();
            setDonations(data.data || data || []);
        } catch (err) {
            console.error("Failed to load dashboard data", err);
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        authService.logout();
        router.push('/login');
    };

    if (loading) return <div className="loading-screen">Loading...</div>;

    return (
        <div className="dashboard-container">
            {/* Sidebar */}
            <div className="dashboard-sidebar">
                <div className="sidebar-header">
                    <img src="/Logo.webp" alt="SustainWear" className="sidebar-logo" />
                </div>

                <nav className="sidebar-nav">
                    <a href="/donor" className="nav-item active">Overview</a>
                    <a href="/donor/donate" className="nav-item">Donate Clothes</a>
                    <a href="/donor/history" className="nav-item">History</a>
                    <a href="/donor/profile" className="nav-item">Profile</a>
                </nav>

                <button onClick={handleLogout} className="logout-btn">Sign out</button>
            </div>

            {/* Main Content */}
            <div className="dashboard-content">
                <header className="content-header">
                    <h1>Welcome back, {user?.name?.split(' ')[0]}!</h1>
                    <p className="header-date">{new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long' })}</p>
                </header>

                {/* Stats */}
                <section className="stats-grid">
                    <div className="stat-card blue">
                        <h3>Total Donations</h3>
                        <p className="stat-number">{donations.length}</p>
                        <span className="stat-label">Items contributed</span>
                    </div>
                    <div className="stat-card green">
                        <h3>Impact Status</h3>
                        <p className="stat-number">Active</p>
                        <span className="stat-label">Donor Level</span>
                    </div>
                </section>

                {/* Recent Activity */}
                <section className="recent-activity">
                    <div className="section-header">
                        <h2>Recent Donations</h2>
                    </div>

                    <div className="table-container">
                        {donations.length === 0 ? (
                            <div className="empty-state">
                                <p>No donations yet.</p>
                            </div>
                        ) : (
                            <table className="data-table">
                                <thead>
                                    <tr>
                                        <th>Item</th>
                                        <th>Date</th>
                                        <th>Status</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {donations.map((d: any) => (
                                        <tr key={d.id}>
                                            <td className="fw-bold">{d.item_name}</td>
                                            <td>{new Date(d.created_at).toLocaleDateString()}</td>
                                            <td>
                                                <span className={`status-badge ${d.status}`}>{d.status}</span>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        )}
                    </div>
                </section>
            </div>
        </div>
    );
}
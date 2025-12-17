"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { authService, donationService, analyticsService } from '@/services/api';
import './donor_dashboard.css';
import Sidebar from '@/components/dashboard/Sidebar';

export default function DonorDashboard() {
    const router = useRouter();
    const [user, setUser] = useState<any>(null);
    const [donations, setDonations] = useState<any[]>([]);
    const [impact, setImpact] = useState<any>({ total_donations: 0, status: 'New Donor' });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const userData = localStorage.getItem('user');

        if (!token || !userData) {
            router.push('/login');
            return;
        }

        setUser(JSON.parse(userData));
        fetchDashboardData();
    }, [router]);

    const fetchDashboardData = async () => {
        try {
            const [donationRes, impactRes] = await Promise.all([
                donationService.getMyDonations(),
                analyticsService.getDonorImpact()
            ]);

            setDonations(donationRes.data || []);
            setImpact(impactRes || { total_donations: 0, status: 'Active' });

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

    if (loading) return <div className="loading-screen">Loading your impact...</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="donor" />
            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Welcome back, {user?.name?.split(' ')[0]}!</h1>
                    <p className="header-date">
                        {new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long' })}
                    </p>
                </header>

                <section className="stats-grid">
                    <div className="stat-card blue">
                        <h3>Total Donations</h3>
                        <p className="stat-number">{impact.total_donations || donations.length}</p>
                        <span className="stat-label">Items contributed</span>
                    </div>

                    <div className="stat-card green">
                        <h3>Carbon Saved</h3>
                        <p className="stat-number">{impact.carbon_saved || '0'} kg</p>
                        <span className="stat-label">Estimated CO2 offset</span>
                    </div>

                    <div className="stat-card purple">
                        <h3>Pending</h3>
                        <p className="stat-number">
                            {donations.filter((d: any) => d.status === 'pending').length}
                        </p>
                        <span className="stat-label">Awaiting approval</span>
                    </div>
                </section>

                <section className="recent-activity">
                    <div className="section-header">
                        <h2>Recent Donations</h2>
                        <button className="view-all-btn" onClick={() => router.push('/donor/history')}>
                            View All
                        </button>
                    </div>

                    <div className="table-container">
                        {donations.length === 0 ? (
                            <div className="empty-state">
                                <p>You haven't made any donations yet.</p>
                                <button onClick={() => router.push('/donor/donate')} className="cta-btn">
                                    Make your first donation
                                </button>
                            </div>
                        ) : (
                            <table className="data-table">
                                <thead>
                                    <tr>
                                        <th>Item</th>
                                        <th>Category</th>
                                        <th>Date</th>
                                        <th>Status</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {donations.slice(0, 5).map((d: any) => (
                                        <tr key={d.id}>
                                            <td className="fw-bold">{d.item_name}</td>
                                            <td>{d.category}</td>
                                            <td>{new Date(d.created_at).toLocaleDateString()}</td>
                                            <td>
                                                <span className={`status-badge ${d.status}`}>
                                                    {d.status}
                                                </span>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        )}
                    </div>
                </section>
            </main>
        </div>
    );
}
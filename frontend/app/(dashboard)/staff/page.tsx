"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { userService, donationService, inventoryService, analyticsService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './staff.css';

interface Donation {
    id: number;
    item_name: string;
    quantity: number;
    created_at: string;
    status: string;
}

interface DashboardData {
    pendingCount: number;
    totalInventory: number;
    todaysIntake: number;
    pendingRequests: Donation[];
}

export default function StaffDashboard() {
    const router = useRouter();
    const [user, setUser] = useState<any>(null);
    const [data, setData] = useState<DashboardData>({
        pendingCount: 0,
        totalInventory: 0,
        todaysIntake: 0,
        pendingRequests: []
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }

        const fetchDashboardData = async () => {
            try {
                const profile = await userService.getProfile();
                setUser(profile);

                const orgId = profile.org_id;
                if (!orgId) {
                    setLoading(false);
                    return;
                }

                const today = new Date();
                const dateStr = `${String(today.getDate()).padStart(2, '0')}-${String(today.getMonth() + 1).padStart(2, '0')}-${today.getFullYear()}`;

                const [donationsData, inventoryStats, trends] = await Promise.all([
                    donationService.list({ status: 'pending', org_id: orgId, page_size: 5 }),
                    inventoryService.getStats(orgId),
                    analyticsService.getDonationTrends(orgId, 'daily', dateStr, dateStr)
                ]);

                const pendingReqs = donationsData.data || [];
                const todaysItems = trends && trends.length > 0 ? trends[0].total_items : 0;

                setData({
                    pendingCount: pendingReqs.length,
                    totalInventory: inventoryStats.total_items || 0,
                    todaysIntake: todaysItems,
                    pendingRequests: pendingReqs
                });

            } catch (error) {
                console.error(error);
            } finally {
                setLoading(false);
            }
        };

        fetchDashboardData();
    }, [router]);

    const handleReview = (id: number) => {
        router.push(`/staff/approvals/${id}`);
    };

    if (loading) return <div className="loading-screen">Loading dashboard...</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="staff" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Welcome back, {user?.full_name?.split(' ')[0]}!</h1>
                    <p className="header-date">
                        {new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long' })}
                    </p>
                </header>

                <section className="stats-grid">
                    <div className="stat-card blue">
                        <h3>Total Inventory</h3>
                        <p className="stat-number">{data.totalInventory}</p>
                        <span className="stat-label">Items in stock</span>
                    </div>

                    <div className="stat-card green">
                        <h3>Today's Intake</h3>
                        <p className="stat-number">{data.todaysIntake}</p>
                        <span className="stat-label">Processed today</span>
                    </div>

                    <div className="stat-card purple">
                        <h3>Pending Reviews</h3>
                        <p className="stat-number">{data.pendingCount}</p>
                        <span className="stat-label">Requires attention</span>
                    </div>
                </section>

                <section className="recent-activity">
                    <div className="section-header">
                        <h2>Pending Approvals</h2>
                        <button className="view-all-btn" onClick={() => router.push('/staff/approvals')}>
                            View All
                        </button>
                    </div>

                    <div className="table-container">
                        {data.pendingRequests.length === 0 ? (
                            <div className="empty-state">
                                <p>No pending approvals found.</p>
                            </div>
                        ) : (
                            <table className="data-table">
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>Item</th>
                                        <th>Qty</th>
                                        <th>Date</th>
                                        <th>Status</th>
                                        <th>Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {data.pendingRequests.map((req) => (
                                        <tr key={req.id}>
                                            <td className="fw-bold">#{req.id}</td>
                                            <td>{req.item_name}</td>
                                            <td>{req.quantity}</td>
                                            <td>{new Date(req.created_at).toLocaleDateString('en-GB')}</td>
                                            <td>
                                                <span className={`status-badge ${req.status}`}>
                                                    {req.status}
                                                </span>
                                            </td>
                                            <td>
                                                <button
                                                    className="view-all-btn"
                                                    onClick={() => handleReview(req.id)}
                                                >
                                                    Review
                                                </button>
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
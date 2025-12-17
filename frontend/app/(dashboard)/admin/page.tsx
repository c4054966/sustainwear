"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { analyticsService, organisationService, userService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './admin.css';

interface SystemOverview {
    total_organisations: number;
    total_donors: number;
    total_donations: number;
    total_items_processed: number;
    co2_saved_kg: number;
    landfill_reduction_kg: number;
}

interface Organisation {
    id: number;
    name: string;
    type: string;
    city: string;
    status: string;
}

export default function AdminDashboard() {
    const router = useRouter();
    const [stats, setStats] = useState<SystemOverview | null>(null);
    const [recentOrgs, setRecentOrgs] = useState<Organisation[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchAdminData();
    }, [router]);

    const fetchAdminData = async () => {
        try {
            const [overview, orgsResponse] = await Promise.all([
                analyticsService.getSystemOverview(),
                organisationService.list({ page_size: 5 })
            ]);

            setStats(overview);
            setRecentOrgs(orgsResponse.data || []);
        } catch (error) {
            console.error("Failed to load admin data", error);
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteOrg = async (id: number) => {
        if (!confirm("Are you sure you want to delete this organisation? This action cannot be undone.")) return;
        try {
            await organisationService.delete(id);
            setRecentOrgs(prev => prev.filter(org => org.id !== id));
        } catch (error) {
            alert("Failed to delete organisation");
        }
    };

    if (loading) return <div className="loading-screen">Loading system data...</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="admin" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>System Overview</h1>
                    <p className="header-date">
                        {new Date().toLocaleDateString('en-GB', {
                            weekday: 'long',
                            day: 'numeric',
                            month: 'long'
                        })}
                    </p>
                </header>

                <section className="stats-grid">
                    <div className="stat-card blue">
                        <h3>Organisations</h3>
                        <p className="stat-number">{stats?.total_organisations || 0}</p>
                        <span className="stat-label">Active partners</span>
                    </div>

                    <div className="stat-card purple">
                        <h3>Total Donors</h3>
                        <p className="stat-number">{stats?.total_donors || 0}</p>
                        <span className="stat-label">Registered users</span>
                    </div>

                    <div className="stat-card green">
                        <h3>Donations</h3>
                        <p className="stat-number">{stats?.total_donations || 0}</p>
                        <span className="stat-label">Total contributions</span>
                    </div>

                    <div className="stat-card teal">
                        <h3>Environmental Impact</h3>
                        <p className="stat-number">
                            {(stats?.co2_saved_kg || 0).toFixed(0)} kg
                        </p>
                        <span className="stat-label">CO2 Saved</span>
                    </div>
                </section>

                <section className="recent-activity">
                    <div className="section-header">
                        <h2>Recent Organisations</h2>
                        <div className="header-actions">
                            <button className="cta-btn small" onClick={() => router.push('/admin/organisations/create')}>
                                + Add New
                            </button>
                            <button className="view-all-btn" onClick={() => router.push('/admin/organisations')}>
                                View All
                            </button>
                        </div>
                    </div>

                    <div className="table-container">
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Name</th>
                                    <th>Type</th>
                                    <th>Location</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {recentOrgs.length > 0 ? (
                                    recentOrgs.map((org) => (
                                        <tr key={org.id}>
                                            <td className="fw-bold">#{org.id}</td>
                                            <td>{org.name}</td>
                                            <td className="capitalize">{org.type}</td>
                                            <td>{org.city}</td>
                                            <td>
                                                <span className={`status-badge ${org.status}`}>
                                                    {org.status}
                                                </span>
                                            </td>
                                            <td>
                                                <button
                                                    className="action-link delete"
                                                    onClick={() => handleDeleteOrg(org.id)}
                                                >
                                                    Delete
                                                </button>
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={6} className="empty-state">No organisations found.</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </section>
            </main>
        </div>
    );
}
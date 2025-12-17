"use client";

import { useEffect, useState } from 'react';
import { analyticsService } from '@/services/api';
// Ensure this path is correct for your project
import Sidebar from '@/components/dashboard/Sidebar'; 
import './analytics.css';

interface SystemOverview {
    total_organisations: number;
    total_donors: number;
    total_donations: number;
    total_items_processed: number;
    co2_saved_kg: number;
    landfill_reduction_kg: number;
    beneficiaries_helped: number;
}

export default function SystemAnalytics() {
    const [stats, setStats] = useState<SystemOverview | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // AUTH CHECK REMOVED
        // The page will now load immediately.
        // The API call below might still fail with 401 if the backend requires a token.
        fetchSystemData();
    }, []);

    const fetchSystemData = async () => {
        try {
            const data = await analyticsService.getSystemOverview();
            setStats(data);
        } catch (error) {
            console.error("Failed to fetch system stats", error);
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="loading-screen">Loading system data...</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="admin" />
            
            <main className="dashboard-content">
                <header className="content-header">
                    <h1>System Analytics</h1>
                    <p className="header-date">Global performance and impact metrics</p>
                </header>

                <div className="analytics-grid">
                    
                    <div className="card-section full-width">
                        <div className="card-header-row">
                            <h2>Environmental Impact</h2>
                            <div className="badge-group">Global</div>
                        </div>
                        <div className="metrics-row">
                            <div className="metric-box green">
                                <span className="metric-label">Total CO2 Saved</span>
                                <span className="metric-value">{(stats?.co2_saved_kg || 0).toFixed(1)} kg</span>
                            </div>
                            <div className="metric-box blue">
                                <span className="metric-label">Landfill Diverted</span>
                                <span className="metric-value">{(stats?.landfill_reduction_kg || 0).toFixed(1)} kg</span>
                            </div>
                        </div>
                    </div>

                    <div className="card-section">
                        <h2>Platform Growth</h2>
                        <div className="performance-stats">
                            <div className="perf-row">
                                <span>Active Organisations</span>
                                <strong>{stats?.total_organisations || 0}</strong>
                            </div>
                            <div className="perf-row">
                                <span>Registered Donors</span>
                                <strong>{stats?.total_donors || 0}</strong>
                            </div>
                            <div className="perf-row highlight">
                                <span>Total Donations</span>
                                <strong>{stats?.total_donations || 0}</strong>
                            </div>
                        </div>
                    </div>

                    <div className="card-section">
                        <h2>Throughput</h2>
                        <div className="performance-stats">
                            <div className="perf-row">
                                <span>Items Processed</span>
                                <strong>{stats?.total_items_processed || 0}</strong>
                            </div>
                            <div className="perf-row">
                                <span>Beneficiaries Helped</span>
                                <strong>{stats?.beneficiaries_helped || 0}</strong>
                            </div>
                            <div className="perf-row highlight">
                                <span>Avg Items / Donation</span>
                                <strong>
                                    {stats?.total_donations ? (stats.total_items_processed / stats.total_donations).toFixed(1) : 0}
                                </strong>
                            </div>
                        </div>
                    </div>

                </div>
            </main>
        </div>
    );
}
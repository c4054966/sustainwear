"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { userService, donationService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './approvals.css';

interface DonationRequest {
    id: number;
    item_name: string;
    category: string;
    quantity: number;
    created_at: string;
    status: string;
    donor_id: number;
}

export default function ApprovalsList() {
    const router = useRouter();
    const [requests, setRequests] = useState<DonationRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchRequests();
    }, [page]);

    const fetchRequests = async () => {
        try {
            setLoading(true);
            const profile = await userService.getProfile();
            const orgId = profile.org_id;

            if (!orgId) return;

            const response = await donationService.list({
                status: 'pending',
                org_id: orgId,
                page: page,
                page_size: 15
            });

            setRequests(response.data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleReview = (id: number) => {
        router.push(`/staff/approvals/${id}`);
    };

    return (
        <div className="dashboard-container">
            <Sidebar role="staff" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Donation Reviews</h1>
                    <p className="header-date">
                        Review and process incoming donation requests
                    </p>
                </header>

                <div className="controls-bar right-align">
                    <button className="cta-btn secondary" onClick={() => fetchRequests()}>Refresh List</button>
                </div>

                <div className="table-container">
                    {loading ? (
                        <div className="loading-state">Loading requests...</div>
                    ) : requests.length === 0 ? (
                        <div className="empty-state">
                            <p>No pending donations found.</p>
                        </div>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>Request ID</th>
                                    <th>Item Name</th>
                                    <th>Category</th>
                                    <th>Qty</th>
                                    <th>Date Received</th>
                                    <th>Status</th>
                                    <th>Action</th>
                                </tr>
                            </thead>
                            <tbody>
                                {requests.map((req) => (
                                    <tr key={req.id}>
                                        <td className="fw-bold">#{req.id}</td>
                                        <td>{req.item_name}</td>
                                        <td className="capitalize">{req.category}</td>
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

                <div className="pagination-controls">
                    <button
                        className="page-btn"
                        disabled={page === 1}
                        onClick={() => setPage(p => p - 1)}
                    >
                        Previous
                    </button>
                    <span className="page-info">Page {page}</span>
                    <button
                        className="page-btn"
                        disabled={requests.length < 15}
                        onClick={() => setPage(p => p + 1)}
                    >
                        Next
                    </button>
                </div>
            </main>
        </div>
    );
}
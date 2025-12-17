"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { userService, donationService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './approvals.css'; // Ensure this file exists as per your previous setup

// Matches the backend JSON response structure
interface DonationRequest {
    id: number;
    item_name: string;
    category: string;
    quantity: number;
    created_at: string;
    status: string;
    donor_id: number;
    // Optional: Add image if you want to show a thumbnail
    images?: string[]; 
}

export default function ApprovalsList() {
    const router = useRouter();
    const [requests, setRequests] = useState<DonationRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    
    // Debug state to show on screen
    const [currentOrgId, setCurrentOrgId] = useState<number | null>(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchRequests();
    }, [page]); // Re-run when page changes

    const fetchRequests = async () => {
        try {
            setLoading(true);
            
            // 1. Get User Profile to find Org ID
            const profile = await userService.getProfile();
            console.log("DEBUG: Logged in Profile:", profile);

            const orgId = profile.org_id;
            setCurrentOrgId(orgId); // Save for UI display

            // 2. Safety Check: If user isn't linked to an org, stop.
            if (!orgId) {
                console.warn("DEBUG: User has no Organization ID assigned!");
                setLoading(false);
                return;
            }

            console.log(`DEBUG: Fetching pending donations for Org ID: ${orgId}`);

            // 3. Fetch Donations for this specific Org
            const response = await donationService.list({
                status: 'pending',
                org_id: orgId,
                page: page,
                page_size: 15
            });

            console.log("DEBUG: API Response:", response);

            // 4. Handle response format (Backend might return Array directly or { data: [] })
            // This safely handles both cases to prevent crashes
            const donationList = Array.isArray(response) ? response : (response.data || []);
            
            setRequests(donationList);

        } catch (error) {
            console.error("DEBUG: Error fetching requests:", error);
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
                        {currentOrgId && <span style={{ marginLeft: '10px', fontSize: '0.8em', color: '#666' }}>(Org ID: {currentOrgId})</span>}
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
                            {/* Helpful hint for development */}
                            <small style={{ display: 'block', marginTop: '10px', color: '#888' }}>
                                (If you just created a donation, check if it was assigned to Org ID {currentOrgId})
                            </small>
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

                {/* Only show pagination if we have data */}
                {requests.length > 0 && (
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
                            // If we have fewer items than page_size, we are on the last page
                            disabled={requests.length < 15} 
                            onClick={() => setPage(p => p + 1)}
                        >
                            Next
                        </button>
                    </div>
                )}
            </main>
        </div>
    );
}
"use client";

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { donationService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import '../approvals.css';

interface DonationDetail {
    id: number;
    item_name: string;
    description: string;
    category: string;
    size: string;
    gender: string;
    condition: string;
    quantity: number;
    images: string;
    status: string;
    created_at: string;
    donor_id: number;
}

export default function ReviewDetails() {
    const router = useRouter();
    const params = useParams();
    const [donation, setDonation] = useState<DonationDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [parsedImages, setParsedImages] = useState<string[]>([]);

    // Modal States
    const [rejectReason, setRejectReason] = useState("");
    const [showRejectModal, setShowRejectModal] = useState(false);
    const [showStatusModal, setShowStatusModal] = useState(false); // NEW
    const [newStatus, setNewStatus] = useState(""); // NEW

    useEffect(() => {
        const id = Number(params.id);
        if (id) fetchDonation(id);
    }, [params.id]);

    const fetchDonation = async (id: number) => {
        try {
            const data = await donationService.getById(id);
            setDonation(data);
            setNewStatus(data.status); // Initialize with current status

            if (data.images) {
                try {
                    const parsed = JSON.parse(data.images);
                    setParsedImages(Array.isArray(parsed) ? parsed : []);
                } catch {
                    setParsedImages([]);
                }
            }
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleApprove = async () => {
        if (!donation) return;
        try {
            await donationService.approve(donation.id);
            router.push('/staff/approvals');
        } catch (error) {
            alert('Failed to approve donation');
        }
    };

    const handleReject = async () => {
        if (!donation) return;
        try {
            await donationService.reject(donation.id, rejectReason);
            setShowRejectModal(false);
            router.push('/staff/approvals');
        } catch (error) {
            alert('Failed to reject donation');
        }
    };

    // NEW: Handle Manual Status Update
    const handleStatusUpdate = async () => {
        if (!donation) return;
        try {
            // This calls the missing PUT /status endpoint
            await donationService.updateStatus(donation.id, newStatus);
            alert("Status updated successfully");
            setShowStatusModal(false);
            fetchDonation(donation.id); // Refresh data
        } catch (error) {
            alert('Failed to update status');
        }
    };

    if (loading) return <div className="loading-screen">Loading details...</div>;
    if (!donation) return <div className="error-screen">Donation not found</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="staff" />

            <main className="dashboard-content">
                <div className="back-nav">
                    <button onClick={() => router.back()} className="back-btn">Back to List</button>
                </div>

                <div className="review-layout">
                    <div className="review-card main-info">
                        <header className="review-header">
                            <div>
                                <h1>{donation.item_name}</h1>
                                <span className="id-badge">ID: #{donation.id}</span>
                            </div>
                            <div style={{ textAlign: 'right' }}>
                                <span className={`status-badge ${donation.status}`}>{donation.status}</span>
                                {/* NEW: Edit Status Link */}
                                <button
                                    className="btn-link"
                                    onClick={() => setShowStatusModal(true)}
                                    style={{ display: 'block', fontSize: '0.8rem', marginTop: '5px', color: '#666', textDecoration: 'underline', border: 'none', background: 'none', cursor: 'pointer' }}
                                >
                                    Edit Status
                                </button>
                            </div>
                        </header>

                        <div className="info-grid">
                            <div className="info-item">
                                <label>Category</label>
                                <p className="capitalize">{donation.category}</p>
                            </div>
                            <div className="info-item">
                                <label>Quantity</label>
                                <p>{donation.quantity}</p>
                            </div>
                            <div className="info-item">
                                <label>Condition</label>
                                <p className="capitalize">{donation.condition}</p>
                            </div>
                            <div className="info-item">
                                <label>Size / Gender</label>
                                <p>{donation.size} / {donation.gender}</p>
                            </div>
                        </div>

                        <div className="description-section">
                            <label>Description</label>
                            <p>{donation.description}</p>
                        </div>

                        <div className="date-info">
                            <small>Submitted on {new Date(donation.created_at).toLocaleString()}</small>
                        </div>

                        {/* Standard Workflow Buttons */}
                        {donation.status === 'pending' && (
                            <div className="action-buttons">
                                <button className="cta-btn reject-btn" onClick={() => setShowRejectModal(true)}>
                                    Reject
                                </button>
                                <button className="cta-btn approve-btn" onClick={handleApprove}>
                                    Approve Donation
                                </button>
                            </div>
                        )}
                    </div>

                    <div className="review-card images-section">
                        <h3>Item Images</h3>
                        {parsedImages.length > 0 ? (
                            <div className="image-grid">
                                {parsedImages.map((img, idx) => (
                                    <div key={idx} className="image-wrapper">
                                        <img src={`http://localhost:8080/${img}`} alt={`Donation ${idx + 1}`} />
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <p className="no-images">No images provided</p>
                        )}
                    </div>
                </div>

                {/* REJECT MODAL */}
                {showRejectModal && (
                    <div className="modal-overlay">
                        <div className="modal-content">
                            <h3>Reject Donation</h3>
                            <p>Please provide a reason for rejection:</p>
                            <textarea
                                value={rejectReason}
                                onChange={(e) => setRejectReason(e.target.value)}
                                placeholder="E.g., Item heavily damaged..."
                                className="form-textarea"
                            />
                            <div className="modal-actions">
                                <button className="cancel-btn" onClick={() => setShowRejectModal(false)}>Cancel</button>
                                <button
                                    className="confirm-reject-btn"
                                    onClick={handleReject}
                                    disabled={!rejectReason.trim()}
                                >
                                    Confirm Rejection
                                </button>
                            </div>
                        </div>
                    </div>
                )}

                {/* NEW: STATUS UPDATE MODAL */}
                {showStatusModal && (
                    <div className="modal-overlay">
                        <div className="modal-content">
                            <h3>Manually Update Status</h3>
                            <p>Change status without triggering inventory logic.</p>
                            <select
                                className="form-select"
                                value={newStatus}
                                onChange={(e) => setNewStatus(e.target.value)}
                                style={{ marginBottom: '20px' }}
                            >
                                <option value="pending">Pending</option>
                                <option value="approved">Approved</option>
                                <option value="rejected">Rejected</option>
                                <option value="received">Received</option>
                                <option value="in_transit">In Transit</option>
                            </select>
                            <div className="modal-actions">
                                <button className="cancel-btn" onClick={() => setShowStatusModal(false)}>Cancel</button>
                                <button className="submit-btn" onClick={handleStatusUpdate}>Save Change</button>
                            </div>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
}
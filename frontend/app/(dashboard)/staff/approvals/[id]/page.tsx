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
    const [rejectReason, setRejectReason] = useState("");
    const [showRejectModal, setShowRejectModal] = useState(false);

    useEffect(() => {
        const id = Number(params.id);
        if (id) fetchDonation(id);
    }, [params.id]);

    const fetchDonation = async (id: number) => {
        try {
            const data = await donationService.getById(id);
            setDonation(data);
            
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

    if (loading) return <div className="loading-screen">Loading details...</div>;
    if (!donation) return <div className="error-screen">Donation not found</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="staff" />
            
            <main className="dashboard-content">
                <div className="back-nav">
                    <button onClick={() => router.back()} className="back-btn">← Back to List</button>
                </div>

                <div className="review-layout">
                    <div className="review-card main-info">
                        <header className="review-header">
                            <div>
                                <h1>{donation.item_name}</h1>
                                <span className="id-badge">ID: #{donation.id}</span>
                            </div>
                            <span className={`status-badge ${donation.status}`}>{donation.status}</span>
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

                {showRejectModal && (
                    <div className="modal-overlay">
                        <div className="modal-content">
                            <h3>Reject Donation</h3>
                            <p>Please provide a reason for rejection:</p>
                            <textarea 
                                value={rejectReason}
                                onChange={(e) => setRejectReason(e.target.value)}
                                placeholder="E.g., Item heavily damaged, Does not meet hygiene standards..."
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
            </main>
        </div>
    );
}
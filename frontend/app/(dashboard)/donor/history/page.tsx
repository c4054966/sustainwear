"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Sidebar from '@/components/dashboard/Sidebar';
import { donationService } from '@/services/api';
import './history.css';

export default function HistoryPage() {
    const router = useRouter();
    const [donations, setDonations] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchHistory();
    }, [page, router]);

    const fetchHistory = async () => {
        setLoading(true);
        try {
            const response = await donationService.getMyDonations(page, 10);
            const newData = response.data || [];

            setDonations(newData);
            setHasMore(newData.length === 10);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="dashboard-container">
            <Sidebar role="donor" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Donation History</h1>
                    <p className="header-date">Track the status and impact of your contributions.</p>
                </header>

                <div className="history-card">
                    {loading ? (
                        <div className="loading-state">Loading records...</div>
                    ) : donations.length === 0 ? (
                        <div className="empty-state">
                            <p>No donation history found.</p>
                            <button onClick={() => router.push('/donor/donate')} className="cta-btn">
                                Make a Donation
                            </button>
                        </div>
                    ) : (
                        <>
                            <div className="table-responsive">
                                <table className="history-table">
                                    <thead>
                                        <tr>
                                            <th>Image</th>
                                            <th>Item Name</th>
                                            <th>Category</th>
                                            <th>Condition</th>
                                            <th>Date</th>
                                            <th>Status</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {donations.map((item) => (
                                            <tr key={item.id}>
                                                <td>
                                                    {item.image_url ? (
                                                        <img
                                                            src={`${process.env.NEXT_PUBLIC_API_URL}/${item.image_url}`}
                                                            alt={item.item_name}
                                                            className="item-thumbnail"
                                                        />
                                                    ) : (
                                                        <div className="no-image">No Img</div>
                                                    )}
                                                </td>
                                                <td className="fw-bold">{item.item_name}</td>
                                                <td>{item.category}</td>
                                                <td>{item.condition}</td>
                                                <td>{new Date(item.created_at).toLocaleDateString()}</td>
                                                <td>
                                                    <span className={`status-badge ${item.status}`}>
                                                        {item.status}
                                                    </span>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>

                            <div className="pagination-controls">
                                <button
                                    disabled={page === 1}
                                    onClick={() => setPage(p => p - 1)}
                                    className="page-btn"
                                >
                                    Previous
                                </button>
                                <span className="page-info">Page {page}</span>
                                <button
                                    disabled={!hasMore}
                                    onClick={() => setPage(p => p + 1)}
                                    className="page-btn"
                                >
                                    Next
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </main>
        </div>
    );
}
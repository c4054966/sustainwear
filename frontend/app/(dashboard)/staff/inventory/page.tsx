"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { userService, inventoryService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './inventory.css';

interface InventoryItem {
    id: number;
    item_name: string;
    category: string;
    quantity: number;
    status: string;
    condition: string;
    location: string;
    updated_at: string;
}

export default function StaffInventory() {
    const router = useRouter();
    const [items, setItems] = useState<InventoryItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [filters, setFilters] = useState({
        category: '',
        status: ''
    });

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchInventory();
    }, [page, filters]);

    const fetchInventory = async () => {
        try {
            setLoading(true);
            const profile = await userService.getProfile();
            const orgId = profile.org_id;

            if (!orgId) return;

            const response = await inventoryService.list({
                org_id: orgId,
                page: page,
                page_size: 15,
                category: filters.category,
                status: filters.status
            });

            setItems(response.data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const { name, value } = e.target;
        setFilters(prev => ({ ...prev, [name]: value }));
        setPage(1);
    };

    return (
        <div className="dashboard-container">
            <Sidebar role="staff" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Inventory Management</h1>
                    <p className="header-date">
                        Manage your organisation's stock levels and item details
                    </p>
                </header>

                <div className="controls-bar">
                    <div className="filter-group">
                        <select name="category" value={filters.category} onChange={handleFilterChange} className="filter-select">
                            <option value="">All Categories</option>
                            <option value="outerwear">Outerwear</option>
                            <option value="tops">Tops</option>
                            <option value="bottoms">Bottoms</option>
                            <option value="footwear">Footwear</option>
                        </select>

                        <select name="status" value={filters.status} onChange={handleFilterChange} className="filter-select">
                            <option value="">All Statuses</option>
                            <option value="available">Available</option>
                            <option value="allocated">Allocated</option>
                            <option value="distributed">Distributed</option>
                        </select>
                    </div>

                    <div className="action-group">
                        <button className="cta-btn secondary" onClick={() => fetchInventory()}>Refresh</button>
                    </div>
                </div>

                <div className="table-container">
                    {loading ? (
                        <div className="loading-state">Loading inventory...</div>
                    ) : items.length === 0 ? (
                        <div className="empty-state">
                            <p>No inventory items found matching your filters.</p>
                        </div>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>Item Name</th>
                                    <th>Category</th>
                                    <th>Condition</th>
                                    <th>Qty</th>
                                    <th>Location</th>
                                    <th>Status</th>
                                    <th>Last Updated</th>
                                </tr>
                            </thead>
                            <tbody>
                                {items.map((item) => (
                                    <tr key={item.id}>
                                        <td className="fw-bold">{item.item_name}</td>
                                        <td className="capitalize">{item.category}</td>
                                        <td className="capitalize">{item.condition}</td>
                                        <td>{item.quantity}</td>
                                        <td>{item.location || '-'}</td>
                                        <td>
                                            <span className={`status-badge ${item.status}`}>
                                                {item.status}
                                            </span>
                                        </td>
                                        <td>{new Date(item.updated_at).toLocaleDateString('en-GB')}</td>
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
                        disabled={items.length < 15}
                        onClick={() => setPage(p => p + 1)}
                    >
                        Next
                    </button>
                </div>
            </main>
        </div>
    );
}
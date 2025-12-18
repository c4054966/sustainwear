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
    allocated_quantity: number; // Needed to calculate available stock
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

    // --- NEW: State for Actions & Modals ---
    const [modalOpen, setModalOpen] = useState(false);
    const [actionType, setActionType] = useState<'create' | 'allocate' | 'distribute' | 'deallocate' | null>(null);
    const [selectedItem, setSelectedItem] = useState<InventoryItem | null>(null);

    // Form Data for all actions
    const [formData, setFormData] = useState({
        quantity: 1,
        reason: '',     // For Allocation
        recipient: '',  // For Distribution
        item_name: '',  // For Manual Create
        category: 'Clothing',
        condition: 'Good',
        location: 'Warehouse A'
    });

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

            // Ensure we handle the response format correctly
            setItems(response.data || response || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    // --- NEW: Action Handlers ---

    const openModal = (type: 'create' | 'allocate' | 'distribute' | 'deallocate', item?: InventoryItem) => {
        setActionType(type);
        setSelectedItem(item || null);
        setFormData({
            quantity: 1,
            reason: '',
            recipient: '',
            item_name: '',
            category: 'Clothing',
            condition: 'Good',
            location: 'Warehouse A'
        });
        setModalOpen(true);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!actionType) return;

        try {
            const profile = await userService.getProfile();

            if (actionType === 'create') {
                // Manual Create
                await inventoryService.create({
                    org_id: profile.org_id,
                    item_name: formData.item_name,
                    category: formData.category,
                    quantity: Number(formData.quantity),
                    condition: formData.condition,
                    location: formData.location
                });
            } else if (selectedItem) {
                // Inventory Movements
                if (actionType === 'allocate') {
                    await inventoryService.allocate(selectedItem.id, Number(formData.quantity), formData.reason);
                } else if (actionType === 'distribute') {
                    await inventoryService.distribute(selectedItem.id, Number(formData.quantity), formData.recipient);
                } else if (actionType === 'deallocate') {
                    await inventoryService.deallocate(selectedItem.id, Number(formData.quantity));
                }
            }

            alert('Action successful!');
            setModalOpen(false);
            fetchInventory(); // Refresh list to show changes
        } catch (error: any) {
            alert(error.message || 'Operation failed');
            console.error(error);
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
                    <h1>Inventory Manager</h1>
                    <p className="header-date">
                        Real-time stock control and distribution
                    </p>
                </header>

                <div className="controls-bar">
                    <div className="filter-group">
                        <select name="category" value={filters.category} onChange={handleFilterChange} className="filter-select">
                            <option value="">All Categories</option>
                            <option value="Clothing">Clothing</option>
                            <option value="Food">Food</option>
                            <option value="Medical">Medical</option>
                        </select>
                        <select name="status" value={filters.status} onChange={handleFilterChange} className="filter-select">
                            <option value="">All Statuses</option>
                            <option value="available">Available</option>
                            <option value="low_stock">Low Stock</option>
                        </select>
                    </div>

                    <div className="action-group">
                        <button className="cta-btn secondary" onClick={() => fetchInventory()}>Refresh</button>
                        {/* NEW: Manual Create Button */}
                        <button className="cta-btn primary" onClick={() => openModal('create')}>+ Add New Item</button>
                    </div>
                </div>

                <div className="table-container">
                    {loading ? (
                        <div className="loading-state">Loading inventory...</div>
                    ) : items.length === 0 ? (
                        <div className="empty-state">
                            <p>No inventory items found.</p>
                        </div>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>Item</th>
                                    <th>Category</th>
                                    <th>Total</th>
                                    <th>Allocated</th>
                                    <th>Available</th>
                                    <th>Actions</th> {/* NEW COLUMN */}
                                </tr>
                            </thead>
                            <tbody>
                                {items.map((item) => {
                                    const available = item.quantity - (item.allocated_quantity || 0);
                                    return (
                                        <tr key={item.id}>
                                            <td className="fw-bold">{item.item_name}</td>
                                            <td className="capitalize">{item.category}</td>
                                            <td>{item.quantity}</td>
                                            <td>{item.allocated_quantity || 0}</td>
                                            <td style={{ color: available > 0 ? '#2e7d32' : '#d32f2f', fontWeight: 'bold' }}>
                                                {available}
                                            </td>
                                            <td>
                                                <div className="action-buttons-row">
                                                    <button
                                                        className="btn-text"
                                                        onClick={() => openModal('allocate', item)}
                                                        disabled={available <= 0}
                                                        title="Reserve stock"
                                                    >Allocate</button>

                                                    <button
                                                        className="btn-text"
                                                        onClick={() => openModal('distribute', item)}
                                                        disabled={available <= 0}
                                                        title="Send stock out"
                                                    >Distribute</button>

                                                    {(item.allocated_quantity || 0) > 0 && (
                                                        <button
                                                            className="btn-text danger"
                                                            onClick={() => openModal('deallocate', item)}
                                                            title="Return to available"
                                                        >Return</button>
                                                    )}
                                                </div>
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    )}
                </div>

                {/* --- NEW: THE ACTION MODAL --- */}
                {modalOpen && (
                    <div className="modal-overlay">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h3>
                                    {actionType === 'create' ? 'Add New Item' :
                                        actionType === 'allocate' ? 'Allocate Stock' :
                                            actionType === 'distribute' ? 'Distribute Stock' : 'Deallocate Stock'}
                                </h3>
                                <button className="close-btn" onClick={() => setModalOpen(false)}>×</button>
                            </div>

                            <form onSubmit={handleSubmit} className="modal-form">

                                {/* DYNAMIC FIELDS BASED ON ACTION TYPE */}
                                {actionType === 'create' ? (
                                    <>
                                        <div className="form-group">
                                            <label>Item Name</label>
                                            <input required className="form-input"
                                                value={formData.item_name} onChange={e => setFormData({ ...formData, item_name: e.target.value })} />
                                        </div>
                                        <div className="form-group">
                                            <label>Category</label>
                                            <select className="form-select" value={formData.category} onChange={e => setFormData({ ...formData, category: e.target.value })}>
                                                <option value="Clothing">Clothing</option>
                                                <option value="Food">Food</option>
                                                <option value="Medical">Medical</option>
                                            </select>
                                        </div>
                                    </>
                                ) : (
                                    <p className="modal-subtitle">Item: <strong>{selectedItem?.item_name}</strong></p>
                                )}

                                <div className="form-group">
                                    <label>Quantity</label>
                                    <input type="number" min="1" required className="form-input"
                                        value={formData.quantity} onChange={e => setFormData({ ...formData, quantity: parseInt(e.target.value) })} />
                                </div>

                                {actionType === 'allocate' && (
                                    <div className="form-group">
                                        <label>Reason / Project</label>
                                        <input type="text" required placeholder="e.g. Winter Shelter Program" className="form-input"
                                            value={formData.reason} onChange={e => setFormData({ ...formData, reason: e.target.value })} />
                                    </div>
                                )}

                                {actionType === 'distribute' && (
                                    <div className="form-group">
                                        <label>Recipient Name</label>
                                        <input type="text" required placeholder="e.g. John Doe or Shelter A" className="form-input"
                                            value={formData.recipient} onChange={e => setFormData({ ...formData, recipient: e.target.value })} />
                                    </div>
                                )}

                                <div className="modal-actions">
                                    <button type="button" className="cancel-btn" onClick={() => setModalOpen(false)}>Cancel</button>
                                    <button type="submit" className="submit-btn">Confirm</button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Pagination Controls (Existing) */}
                <div className="pagination-controls">
                    <button className="page-btn" disabled={page === 1} onClick={() => setPage(p => p - 1)}>Previous</button>
                    <span className="page-info">Page {page}</span>
                    <button className="page-btn" disabled={items.length < 15} onClick={() => setPage(p => p + 1)}>Next</button>
                </div>
            </main>
        </div>
    );
}
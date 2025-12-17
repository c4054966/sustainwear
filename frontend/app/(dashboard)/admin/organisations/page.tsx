"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { organisationService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './organisations.css';

interface Organisation {
    id: number;
    name: string;
    type: string;
    city: string;
    email: string;
    status: string;
}

export default function OrganisationsList() {
    const router = useRouter();
    const [orgs, setOrgs] = useState<Organisation[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [filters, setFilters] = useState({
        type: '',
        status: ''
    });

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchOrganisations();
    }, [page, filters]);

    const fetchOrganisations = async () => {
        try {
            setLoading(true);
            const response = await organisationService.list({
                page,
                page_size: 15,
                type: filters.type,
                status: filters.status
            });
            setOrgs(response.data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure? This will delete the organisation and all associated data.")) return;
        try {
            await organisationService.delete(id);
            setOrgs(prev => prev.filter(o => o.id !== id));
        } catch (error) {
            alert("Failed to delete organisation");
        }
    };

    const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setFilters(prev => ({ ...prev, [e.target.name]: e.target.value }));
        setPage(1);
    };

    return (
        <div className="dashboard-container">
            <Sidebar role="admin" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Organisations</h1>
                    <p className="header-date">Manage partner charities and NGOs</p>
                </header>

                <div className="controls-bar">
                    <div className="filter-group">
                        <select name="type" value={filters.type} onChange={handleFilterChange} className="filter-select">
                            <option value="">All Types</option>
                            <option value="charity">Charity</option>
                            <option value="ngo">NGO</option>
                            <option value="community">Community</option>
                        </select>
                        <select name="status" value={filters.status} onChange={handleFilterChange} className="filter-select">
                            <option value="">All Statuses</option>
                            <option value="active">Active</option>
                            <option value="inactive">Inactive</option>
                            <option value="pending">Pending</option>
                        </select>
                    </div>
                    <button className="cta-btn" onClick={() => router.push('/admin/organisations/create')}>
                        + Add Organisation
                    </button>
                </div>

                <div className="table-container">
                    {loading ? (
                        <div className="loading-state">Loading organisations...</div>
                    ) : orgs.length === 0 ? (
                        <div className="empty-state">No organisations found.</div>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Name</th>
                                    <th>Type</th>
                                    <th>Location</th>
                                    <th>Email</th>
                                    <th>Status</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {orgs.map((org) => (
                                    <tr key={org.id}>
                                        <td className="fw-bold">#{org.id}</td>
                                        <td>{org.name}</td>
                                        <td className="capitalize">{org.type}</td>
                                        <td>{org.city}</td>
                                        <td>{org.email}</td>
                                        <td><span className={`status-badge ${org.status}`}>{org.status}</span></td>
                                        <td>
                                            <div className="action-buttons">
                                                <button className="action-link" onClick={() => router.push(`/admin/organisations/${org.id}`)}>Edit</button>
                                                <button className="action-link delete" onClick={() => handleDelete(org.id)}>Delete</button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>

                <div className="pagination-controls">
                    <button className="page-btn" disabled={page === 1} onClick={() => setPage(p => p - 1)}>Previous</button>
                    <span className="page-info">Page {page}</span>
                    <button className="page-btn" disabled={orgs.length < 15} onClick={() => setPage(p => p + 1)}>Next</button>
                </div>
            </main>
        </div>
    );
}
"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { userService } from '@/services/api';
import Sidebar from '@/components/dashboard/Sidebar';
import './users.css';

interface User {
    id: number;
    email: string;
    full_name: string;
    role: string;
    org_id?: number;
    created_at: string;
}

export default function UserManagement() {
    const router = useRouter();
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const pageSize = 15;

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }
        fetchUsers();
    }, [page]);

    const fetchUsers = async () => {
        try {
            setLoading(true);
            const response = await userService.list(page, pageSize);
            setUsers(response.data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure you want to delete this user? This cannot be undone.")) return;
        try {
            await userService.delete(id);
            setUsers(prev => prev.filter(u => u.id !== id));
        } catch (error) {
            alert("Failed to delete user");
        }
    };

    return (
        <div className="dashboard-container">
            <Sidebar role="admin" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>User Management</h1>
                    <p className="header-date">View and manage registered users</p>
                </header>

                <div className="controls-bar right-align">
                    <button className="cta-btn secondary" onClick={() => fetchUsers()}>Refresh List</button>
                </div>

                <div className="table-container">
                    {loading ? (
                        <div className="loading-state">Loading users...</div>
                    ) : users.length === 0 ? (
                        <div className="empty-state">
                            <p>No users found.</p>
                        </div>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Full Name</th>
                                    <th>Email</th>
                                    <th>Role</th>
                                    <th>Joined Date</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map((user) => (
                                    <tr key={user.id}>
                                        <td className="fw-bold">#{user.id}</td>
                                        <td>{user.full_name}</td>
                                        <td>{user.email}</td>
                                        <td>
                                            <span className={`role-badge ${user.role}`}>
                                                {user.role}
                                            </span>
                                        </td>
                                        <td>{new Date(user.created_at).toLocaleDateString('en-GB')}</td>
                                        <td>
                                            <button
                                                className="action-link delete"
                                                onClick={() => handleDelete(user.id)}
                                            >
                                                Delete
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
                        disabled={users.length < pageSize}
                        onClick={() => setPage(p => p + 1)}
                    >
                        Next
                    </button>
                </div>
            </main>
        </div>
    );
}
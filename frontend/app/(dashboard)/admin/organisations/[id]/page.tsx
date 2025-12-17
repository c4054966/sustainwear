"use client";

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { organisationService } from '@/services/api';
import Sidebar from '@/components//dashboard/Sidebar';
import '../organisations.css';

export default function OrganisationEditor() {
    const router = useRouter();
    const params = useParams();
    const isNew = params.id === 'create';
    
    const [formData, setFormData] = useState({
        name: '', description: '', type: 'charity', email: '', phone: '',
        address: '', city: '', county: '', postcode: '', country: 'United Kingdom', website: '', status: 'active'
    });
    const [loading, setLoading] = useState(!isNew);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (!isNew) {
            organisationService.getById(Number(params.id))
                .then(data => setFormData(data))
                .catch(() => router.push('/admin/organisations'))
                .finally(() => setLoading(false));
        }
    }, [isNew, params.id, router]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
        setFormData(prev => ({ ...prev, [e.target.name]: e.target.value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (isNew) {
                await organisationService.create(formData);
            } else {
                await organisationService.update(Number(params.id), formData);
            }
            router.push('/admin/organisations');
        } catch (error) {
            alert("Failed to save organisation");
        } finally {
            setSaving(false);
        }
    };

    if (loading) return <div className="loading-screen">Loading details...</div>;

    return (
        <div className="dashboard-container">
            <Sidebar role="admin" />
            <main className="dashboard-content">
                <header className="content-header">
                    <h1>{isNew ? 'Create Organisation' : 'Edit Organisation'}</h1>
                    <button className="back-btn" onClick={() => router.back()}>← Back</button>
                </header>

                <div className="form-container">
                    <form onSubmit={handleSubmit}>
                        <div className="form-section">
                            <h3>Basic Details</h3>
                            <div className="form-grid">
                                <div className="form-group">
                                    <label>Name</label>
                                    <input required name="name" value={formData.name} onChange={handleChange} className="form-input" />
                                </div>
                                <div className="form-group">
                                    <label>Type</label>
                                    <select name="type" value={formData.type} onChange={handleChange} className="form-input">
                                        <option value="charity">Charity</option>
                                        <option value="ngo">NGO</option>
                                        <option value="community">Community</option>
                                    </select>
                                </div>
                                <div className="form-group full-width">
                                    <label>Description</label>
                                    <textarea name="description" value={formData.description} onChange={handleChange} className="form-input" rows={3} />
                                </div>
                            </div>
                        </div>

                        <div className="form-section">
                            <h3>Contact Information</h3>
                            <div className="form-grid">
                                <div className="form-group">
                                    <label>Email</label>
                                    <input required type="email" name="email" value={formData.email} onChange={handleChange} className="form-input" />
                                </div>
                                <div className="form-group">
                                    <label>Phone</label>
                                    <input name="phone" value={formData.phone} onChange={handleChange} className="form-input" />
                                </div>
                                <div className="form-group">
                                    <label>Website</label>
                                    <input name="website" value={formData.website} onChange={handleChange} className="form-input" />
                                </div>
                            </div>
                        </div>

                        <div className="form-section">
                            <h3>Address</h3>
                            <div className="form-grid">
                                <div className="form-group full-width">
                                    <label>Street Address</label>
                                    <input name="address" value={formData.address} onChange={handleChange} className="form-input" />
                                </div>
                                <div className="form-group">
                                    <label>City</label>
                                    <input name="city" value={formData.city} onChange={handleChange} className="form-input" />
                                </div>
                                <div className="form-group">
                                    <label>County</label>
                                    <input name="county" value={formData.county} onChange={handleChange} className="form-input" />
                                </div>
                                <div className="form-group">
                                    <label>Postcode</label>
                                    <input name="postcode" value={formData.postcode} onChange={handleChange} className="form-input" />
                                </div>
                            </div>
                        </div>

                        <div className="form-actions">
                            <button type="button" className="cancel-btn" onClick={() => router.back()}>Cancel</button>
                            <button type="submit" className="save-btn" disabled={saving}>
                                {saving ? 'Saving...' : (isNew ? 'Create Organisation' : 'Save Changes')}
                            </button>
                        </div>
                    </form>
                </div>
            </main>
        </div>
    );
}
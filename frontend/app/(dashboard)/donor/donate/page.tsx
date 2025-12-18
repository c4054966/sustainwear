"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { donationService, uploadService, authService } from '@/services/api';
import './donate.css';
import Sidebar from '@/components/dashboard/Sidebar';

export default function DonatePage() {
    const router = useRouter();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const [formData, setFormData] = useState({
        item_name: '',
        category: 'Clothing',
        condition: 'Good',
        description: '',
        quantity: 1,
        size: '',
        gender: 'Unisex',
    });

    const [imageFile, setImageFile] = useState<File | null>(null);
    const [previewUrl, setPreviewUrl] = useState<string | null>(null);

    const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setImageFile(file);
            setPreviewUrl(URL.createObjectURL(file));
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            let imageUrl = "";

            if (imageFile) {
                const uploadResponse = await uploadService.uploadImage(imageFile);

                if (uploadResponse.paths && uploadResponse.paths.length > 0) {
                    imageUrl = uploadResponse.paths[0];
                }
            }

            const payload = {
                item_name: formData.item_name,
                category: formData.category,
                condition: formData.condition,
                description: formData.description,
                quantity: Number(formData.quantity),
                size: formData.size,                 
                gender: formData.gender,             
                images: imageUrl ? [imageUrl] : []
            };

            await donationService.create(payload);

            router.push('/donor/history');

        } catch (err: any) {
            console.error(err);
            setError("Failed to submit donation. Please try again.");
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        authService.logout();
        router.push('/login');
    };

    return (
        <div className="dashboard-container">
            <Sidebar role="donor" />

            <main className="dashboard-content">
                <header className="content-header">
                    <h1>Donate an Item</h1>
                    <p className="header-date">Give your clothes a second life.</p>
                </header>

                <div className="donate-card">
                    {error && <div className="error-message">{error}</div>}

                    <form onSubmit={handleSubmit} className="donate-form">

                        <div className="form-section">
                            {/* --- EXISTING ITEM NAME --- */}
                            <div className="form-group">
                                <label className="form-label">Item Name</label>
                                <input
                                    type="text"
                                    required
                                    className="form-input"
                                    placeholder="e.g. Vintage Denim Jacket"
                                    value={formData.item_name}
                                    onChange={(e) => setFormData({ ...formData, item_name: e.target.value })}
                                />
                            </div>

                            {/* --- NEW ROW: QUANTITY & SIZE --- */}
                            <div className="form-row">
                                <div className="form-group half">
                                    <label className="form-label">Quantity</label>
                                    <input
                                        type="number"
                                        min="1"
                                        required
                                        className="form-input"
                                        value={formData.quantity}
                                        onChange={(e) => setFormData({ ...formData, quantity: parseInt(e.target.value) || 1 })}
                                    />
                                </div>

                                <div className="form-group half">
                                    <label className="form-label">Size (Optional)</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        placeholder="e.g. M, L, 10, 42"
                                        value={formData.size}
                                        onChange={(e) => setFormData({ ...formData, size: e.target.value })}
                                    />
                                </div>
                            </div>

                            {/* --- NEW ROW: GENDER & CATEGORY --- */}
                            <div className="form-row">
                                <div className="form-group half">
                                    <label className="form-label">Gender / Target</label>
                                    <select
                                        className="form-select"
                                        value={formData.gender}
                                        onChange={(e) => setFormData({ ...formData, gender: e.target.value })}
                                    >
                                        <option value="Unisex">Unisex</option>
                                        <option value="Men">Men</option>
                                        <option value="Women">Women</option>
                                        <option value="Kids">Kids</option>
                                    </select>
                                </div>

                                <div className="form-group half">
                                    <label className="form-label">Category</label>
                                    <select
                                        className="form-select"
                                        value={formData.category}
                                        onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                                    >
                                        <option value="Clothing">Clothing</option>
                                        <option value="Footwear">Footwear</option>
                                        <option value="Accessories">Accessories</option>
                                    </select>
                                </div>

                                <div className="form-group half">
                                    <label className="form-label">Condition</label>
                                    <select
                                        className="form-select"
                                        value={formData.condition}
                                        onChange={(e) => setFormData({ ...formData, condition: e.target.value })}
                                    >
                                        <option value="New">Brand New</option>
                                        <option value="Good">Good Condition</option>
                                        <option value="Fair">Fair / Worn</option>
                                    </select>
                                </div>
                            </div>

                            <div className="form-group">
                                <label className="form-label">Description</label>
                                <textarea
                                    className="form-textarea"
                                    rows={4}
                                    value={formData.description}
                                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                />
                            </div>
                        </div>

                        <div className="image-section">
                            <label className="form-label">Upload Photo</label>
                            <div className="image-upload-box">
                                {previewUrl ? (
                                    <div className="image-preview">
                                        <img src={previewUrl} alt="Preview" />
                                        <button type="button" className="remove-img-btn" onClick={() => {
                                            setPreviewUrl(null);
                                            setImageFile(null);
                                        }}>Remove</button>
                                    </div>
                                ) : (
                                    <label className="upload-placeholder">
                                        <input
                                            type="file"
                                            accept="image/*"
                                            onChange={handleImageChange}
                                            hidden
                                        />
                                        <span>Click to upload image</span>
                                    </label>
                                )}
                            </div>
                        </div>

                        <div className="form-actions">
                            <button type="submit" disabled={loading} className="submit-btn">
                                {loading ? 'Submitting...' : 'Submit Donation'}
                            </button>
                        </div>

                    </form>
                </div>
            </main>
        </div>
    );
}
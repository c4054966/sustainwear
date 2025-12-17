"use client";

import { useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import Link from 'next/link';
import { authService } from '@/services/api';
import './sidebar.css';

interface SidebarProps {
    role?: 'donor' | 'admin' | 'staff';
}

export default function Sidebar({ role = 'donor' }: SidebarProps) {
    const router = useRouter();
    const pathname = usePathname();
    const [menuOpen, setMenuOpen] = useState(false);

    const toggleMenu = () => {
        setMenuOpen(!menuOpen);
    };

    const handleLogout = () => {
        authService.logout();
        router.push('/login');
    };

    const links = {
        donor: [
            { label: 'Overview', href: '/donor' },
            { label: 'Donate Clothes', href: '/donor/donate' },
            { label: 'History', href: '/donor/history' },
            { label: 'Profile', href: '/donor/profile' },
        ],
        admin: [
            { label: 'Dashboard', href: '/admin' },
            { label: 'Manage Users', href: '/admin/users' },
            { label: 'Organisations', href: '/admin/organisations' },
            { label: 'System Analytics', href: '/admin/analytics' },
        ],
        staff: [
            { label: 'Dashboard', href: '/staff' },
            { label: 'Inventory', href: '/staff/inventory' },
            { label: 'Donation Reviews', href: '/staff/approvals' },
            { label: 'Profile', href: '/staff/profile' },
        ]
    };

    const currentLinks = links[role] || links.donor;

    return (
        <aside className="dashboard-sidebar">
            <div className="sidebar-header">
                <img src="/Logo.webp" alt="SustainWear" className="sidebar-logo" />
                <img src="/icons/menu.webp" id="menubtn" alt="Menu" onClick={toggleMenu} />
            </div>

            <div className={`sidebar-content-wrapper ${menuOpen ? 'open' : ''}`}>
                <nav className="sidebar-nav">
                    {currentLinks.map((link) => (
                        <Link
                            key={link.href}
                            href={link.href}
                            className={`nav-item ${pathname === link.href ? 'active' : ''}`}
                            onClick={() => setMenuOpen(false)}
                        >
                            {link.label}
                        </Link>
                    ))}
                </nav>

                <button onClick={handleLogout} className="logout-btn">
                    Sign Out
                </button>
            </div>
        </aside>
    );
}
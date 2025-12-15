'use client';

import React, { useState } from 'react';
import dynamic from 'next/dynamic';
import styles from './FindUs.module.css'; // <--- Ensure this import is here
import Header_main from "@/components/layout/Header_main";
import Footer_main from "@/components/layout/Footer-main";

// Dynamic Import for Map
const MapWithNoSSR = dynamic(() => import('@/components/map'), { 
    ssr: false,
    loading: () => <div style={{textAlign:'center', padding:'20px'}}>Loading Map...</div>
});

interface LocationItem {
    id: number;
    name: string;
    city: string;
    lat: number;
    lng: number;
    address: string;
    distance?: number;
}

const locationsData: LocationItem[] = [
    { id: 1, name: "Sustain Wear Hub", city: "Sheffield", lat: 53.3811, lng: -1.4701, address: "12 High St, Sheffield City Centre, S1 2GA" },
    { id: 2, name: "Meadowhall Drop-Off", city: "Sheffield", lat: 53.4143, lng: -1.4116, address: "Meadowhall Shopping Centre, Sheffield, S9 1EP" },
    { id: 3, name: "Ecclesall Road Shop", city: "Sheffield", lat: 53.3664, lng: -1.5039, address: "Ecclesall Road, Sheffield, S11 8NX" },
    { id: 4, name: "Leeds Charity Shop", city: "Leeds", lat: 53.7965, lng: -1.5478, address: "Albion Place, Leeds, LS1 6JL" },
    { id: 5, name: "Headingley Drop Point", city: "Leeds", lat: 53.8203, lng: -1.5760, address: "Otley Rd, Leeds, LS6 3AD" },
];

const calculateDistance = (lat1: number, lon1: number, lat2: number, lon2: number) => {
    const R = 3958.8; 
    const dLat = (lat2 - lat1) * (Math.PI / 180);
    const dLon = (lon2 - lon1) * (Math.PI / 180);
    const a = 
        Math.sin(dLat / 2) * Math.sin(dLat / 2) +
        Math.cos(lat1 * (Math.PI / 180)) * Math.cos(lat2 * (Math.PI / 180)) * Math.sin(dLon / 2) * Math.sin(dLon / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c; 
};

const FindUs = () => {
    const [search, setSearch] = useState("");
    const [sortedLocations, setSortedLocations] = useState<LocationItem[]>(locationsData);
    const [userLocation, setUserLocation] = useState<[number, number] | null>(null);
    const [loading, setLoading] = useState(false);

    const handleSearch = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${search}&countrycodes=gb`);
            const data = await response.json();

            if (data && data.length > 0) {
                const searchLat = parseFloat(data[0].lat);
                const searchLon = parseFloat(data[0].lon);
                
                setUserLocation([searchLat, searchLon]);

                const sorted = locationsData.map(loc => ({
                    ...loc,
                    distance: calculateDistance(searchLat, searchLon, loc.lat, loc.lng)
                })).sort((a, b) => (a.distance || 0) - (b.distance || 0));

                const cityMatch = locationsData.find(l => search.toLowerCase().includes(l.city.toLowerCase()));
                if (cityMatch) {
                    setSortedLocations(sorted.filter(l => l.city.toLowerCase() === cityMatch.city.toLowerCase()));
                } else {
                    setSortedLocations(sorted);
                }
            } else {
                alert("Location not found.");
            }
        } catch (error) {
            console.error("Search error:", error);
        }
        setLoading(false);
    };

    const handleLocationFound = (lat: number, lng: number) => {
        setUserLocation([lat, lng]);
        const sorted = locationsData.map(loc => ({
            ...loc,
            distance: calculateDistance(lat, lng, loc.lat, loc.lng)
        })).sort((a, b) => (a.distance || 0) - (b.distance || 0));
        setSortedLocations(sorted);
    };

    return (
        // 1. Apply styles.pageContainer
        <div className={styles.pageContainer}>
            <Header_main />
            
            {/* 2. Apply styles.searchContainer */}
            <div className={styles.searchContainer}>
                <h1 className={styles.searchTitle}>Find a Donation Point</h1>
                <p className={styles.searchSubtitle}>Enter a Postcode (e.g. S1 2GA) or City.</p>
                
                {/* 3. Apply styles.searchForm */}
                <form onSubmit={handleSearch} className={styles.searchForm}>
                    <input 
                        type="text" 
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        placeholder="Enter Postcode or City..."
                        className={styles.searchInput}
                    />
                    <button type="submit" disabled={loading} className={styles.searchButton}>
                        {loading ? "Searching..." : "Search"}
                    </button>
                </form>
            </div>

            {/* 4. THIS IS THE FIX: styles.mapWrapper */}
            <div className={styles.mapWrapper}>
                <MapWithNoSSR 
                    locations={sortedLocations} 
                    userLocation={userLocation}
                    onLocationFound={handleLocationFound} 
                />
            </div>

            {/* 5. Apply styles.listContainer */}
            <div className={styles.listContainer}>
                <h3 className={styles.listTitle}>Nearest Locations:</h3>
                <ul className={styles.list}>
                    {sortedLocations.map((loc, index) => (
                        <li 
                            key={loc.id} 
                            // 6. Dynamic Classes (styles.listItem + optional styles.closestItem)
                            className={`${styles.listItem} ${index === 0 && loc.distance ? styles.closestItem : ''}`}
                        >
                            <div>
                                <strong className={styles.locationName}>{loc.name}</strong>
                                <span className={styles.locationAddress}>{loc.address}</span>
                            </div>
                            
                            {loc.distance !== undefined && (
                                <div className={styles.distanceBadge}>
                                    {loc.distance.toFixed(1)} miles away
                                </div>
                            )}
                        </li>
                    ))}
                </ul>
            </div>
            
            <Footer_main />
        </div>
    );
};

export default FindUs;
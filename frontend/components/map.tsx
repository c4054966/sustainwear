'use client';

import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { useEffect } from 'react';

// --- Red Icon Definition ---
const redIcon = new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});

interface MapProps {
    locations: any[];
    userLocation: [number, number] | null;
    onLocationFound: (lat: number, lng: number) => void;
}

const RecenterMap = ({ center }: { center: [number, number] | null }) => {
    const map = useMap();
    useEffect(() => {
        if (center) {
            map.flyTo(center, 13, { duration: 2 });
        }
    }, [center, map]);
    return null;
};

const MapComponent = ({ locations, userLocation, onLocationFound }: MapProps) => {

    useEffect(() => {
        // Fix Leaflet Default Icon
        // @ts-ignore
        delete (L.Icon.Default.prototype as any)._getIconUrl;
        L.Icon.Default.mergeOptions({
            iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
            iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
            shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
        });
    }, []);

    const handleMyLocation = () => {
        if (!navigator.geolocation) {
            alert("Geolocation is not supported");
            return;
        }
        navigator.geolocation.getCurrentPosition(
            (position) => {
                onLocationFound(position.coords.latitude, position.coords.longitude);
            },
            () => alert("Unable to retrieve location. Please check browser permissions."),
            { enableHighAccuracy: true }
        );
    };

    return (
        <div style={{ width: "100%", height: "100%", position: "relative" }}>

            {/* My Location Button */}
            <button
                onClick={handleMyLocation}
                style={{
                    position: "absolute",
                    top: "15px",
                    right: "15px",
                    zIndex: 1000,
                    padding: "10px 15px",
                    backgroundColor: "#e74c3c",
                    color: "white",
                    border: "none",
                    borderRadius: "8px",
                    cursor: "pointer",
                    boxShadow: "0 4px 6px rgba(0,0,0,0.2)",
                    fontWeight: "bold",
                    fontSize: "0.9rem"
                }}
            >
                Use My Location
            </button>

            <MapContainer
                center={[53.3811, -1.4701]}
                zoom={10}
                style={{ height: "100%", width: "100%" }}
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; OpenStreetMap contributors'
                />

                <RecenterMap center={userLocation} />

                {locations.map(loc => (
                    <Marker key={loc.id} position={[loc.lat, loc.lng]}>
                        <Popup>
                            <strong>{loc.name}</strong><br />
                            {loc.address}<br />
                            {loc.distance && <em>{loc.distance.toFixed(1)} miles away</em>}
                        </Popup>
                    </Marker>
                ))}

                {userLocation && (
                    <Marker position={userLocation} icon={redIcon}>
                        <Popup><strong>You are here!</strong></Popup>
                    </Marker>
                )}
            </MapContainer>
        </div>
    );
};

export default MapComponent;
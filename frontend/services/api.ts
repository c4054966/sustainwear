const API_URL = process.env.NEXT_PUBLIC_API_URL || '/api';

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    const url = `${API_URL}${cleanEndpoint}`;

    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>),
    };

    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
    }

    const res = await fetch(url, { ...options, headers });

    if (res.status === 204) {
        return {} as T;
    }

    const text = await res.text();
    let data;
    try {
        data = JSON.parse(text);
    } catch {
        data = { message: text };
    }

    if (!res.ok) {
        throw new Error(data.details || data.error || data.message || 'Something went wrong');
    }

    return data;
}


export const authService = {
    login: (credentials: { email: string; password: string }) =>
        request<any>('/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        }),

    register: (userData: any) =>
        request<any>('/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData),
        }),

    getProfile: () =>
        request<any>('/users/profile', {
            method: 'GET',
        }),

    logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
    }
};


export const donationService = {
    getMyDonations: () =>
        request<any>('/donations/my?page=1&page_size=5', { // Fetch top 5 recent
            method: 'GET',
        }),

    getStats: () =>
        request<any>('/analytics/donor-impact', {
            method: 'GET',
        }),

    create: (data: FormData) =>
        request<any>('/donations', {
            method: 'POST',
            body: data, // Note: This might need different handling if sending files/images
        }),
};
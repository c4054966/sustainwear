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
    if (res.status === 204) return {} as T;
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
    login: (credentials: any) => request<any>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(credentials),
    }),

    register: (userData: any) => request<any>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(userData),
    }),

    logout: () => {
        if (typeof window !== 'undefined') {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
        }
    }
};

export const uploadService = {
    uploadImage: async (file: File) => {
        const formData = new FormData();
        formData.append('images', file);

        const token = localStorage.getItem('token');

        const res = await fetch(`${API_URL}/upload`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            body: formData,
        });

        if (!res.ok) {
            const text = await res.text();
            throw new Error(text || 'Image upload failed');
        }

        return res.json();
    }
};

export const analyticsService = {
  getDonorImpact: () => 
    request<any>('/analytics/donor-impact', {
      method: 'GET',
    }),
};

export const donationService = {
  getMyDonations: (page = 1, pageSize = 10) => 
    request<any>(`/donations/my?page=${page}&page_size=${pageSize}`, {
      method: 'GET',
    }),

  create: (data: any) => 
    request<any>('/donations', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
};

export const userService = {
  getProfile: () => 
    request<any>('/users/profile', {
      method: 'GET',
    }),

  updateProfile: (data: { full_name: string }) => 
    request<any>('/users/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
};

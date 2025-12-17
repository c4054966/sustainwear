const API_URL = process.env.NEXT_PUBLIC_API_URL || '/api';

function buildQueryString(params: Record<string, any>) {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      query.append(key, String(value));
    }
  });
  return query.toString();
}

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

    const res = await fetch(`${API_URL}/uploads/images`, {
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

export const userService = {
  getProfile: () =>
    request<any>('/users/profile', { method: 'GET' }),

  updateProfile: (data: { full_name: string }) =>
    request<any>('/users/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  list: (page = 1, pageSize = 20) =>
    request<any>(`/users?page=${page}&page_size=${pageSize}`, { method: 'GET' }),

  delete: (id: number) =>
    request<any>(`/users/${id}`, { method: 'DELETE' }),

  getById: (id: number) =>
    request<any>(`/users/${id}`, { method: 'GET' }),
};

export const analyticsService = {
  getDonorImpact: () =>
    request<any>('/analytics/donor-impact', { method: 'GET' }),

  getDonationTrends: (orgId: number, period = 'daily', startDate: string, endDate: string) => {
    const query = buildQueryString({ org_id: orgId, period, start_date: startDate, end_date: endDate });
    return request<any>(`/analytics/trends?${query}`, { method: 'GET' });
  },

  getCategoryBreakdown: (orgId: number) =>
    request<any>(`/analytics/categories?org_id=${orgId}`, { method: 'GET' }),

  getSystemOverview: () =>
    request<any>('/analytics/system-overview', { method: 'GET' }),

  getSustainabilityMetrics: (orgId: number, period = 'all_time') => {
    const query = buildQueryString({ org_id: orgId, period });
    return request<any>(`/analytics/sustainability?${query}`, { method: 'GET' });
  },

  getOrgPerformance: (orgId: number) =>
    request<any>(`/analytics/org-performance?org_id=${orgId}`, { method: 'GET' }),
};

export const donationService = {
  getMyDonations: (page = 1, pageSize = 10) =>
    request<any>(`/donations/my?page=${page}&page_size=${pageSize}`, { method: 'GET' }),

  create: (data: any) =>
    request<any>('/donations', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  list: (filters: { status?: string; org_id?: number; page?: number; page_size?: number } = {}) => {
    const query = buildQueryString(filters);
    return request<any>(`/donations?${query}`, { method: 'GET' });
  },

  getById: (id: number) =>
    request<any>(`/donations/${id}`, { method: 'GET' }),

  approve: (id: number) =>
    request<any>(`/donations/${id}/approve`, { method: 'POST' }),

  reject: (id: number, reason: string) =>
    request<any>(`/donations/${id}/reject`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    }),
};

export const inventoryService = {
  getStats: (orgId: number) =>
    request<any>(`/inventory/stats?org_id=${orgId}`, { method: 'GET' }),

  list: (filters: { org_id?: number; page?: number; page_size?: number; category?: string; status?: string } = {}) => {
    const query = buildQueryString(filters);
    return request<any>(`/inventory?${query}`, { method: 'GET' });
  },

  create: (data: any) =>
    request<any>('/inventory', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: number, data: any) =>
    request<any>(`/inventory/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<any>(`/inventory/${id}`, { method: 'DELETE' }),
};

export const organisationService = {
  list: (filters: { type?: string; status?: string; city?: string; page?: number; page_size?: number } = {}) => {
    const query = buildQueryString(filters);
    return request<any>(`/organisations?${query}`, { method: 'GET' });
  },

  getById: (id: number) =>
    request<any>(`/organisations/${id}`, { method: 'GET' }),

  create: (data: any) =>
    request<any>('/organisations', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: number, data: any) =>
    request<any>(`/organisations/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<any>(`/organisations/${id}`, { method: 'DELETE' }),

  getStats: (id: number) =>
    request<any>(`/organisations/${id}/stats`, { method: 'GET' }),
};
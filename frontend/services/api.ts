const API_URL = process.env.NEXT_PUBLIC_API_URL || '/api';

const getHeaders = (isMultipart = false) => {
    const headers: HeadersInit = {};
    if (!isMultipart) {
        headers['Content-Type'] = 'application/json';
    }
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
    }
    return headers;
};

const request = async <T>(endpoint: string, options: RequestInit = {}): Promise<T> => {
    const url = `${API_URL}${endpoint}`;
    const headers = getHeaders(options.body instanceof FormData);

    const config = {
        ...options,
        headers: {
            ...headers,
            ...options.headers,
        },
    };

    const response = await fetch(url, config);
    let data;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        data = await response.json();
    } else {
        data = await response.text();
    }

    if (!response.ok) {
        const errorMessage = typeof data === 'object' && data.error ? data.error : (typeof data === 'string' ? data : 'Something went wrong');
        throw new Error(errorMessage);
    }

    return data as T;
};

const buildQueryString = (params: Record<string, any>) => {
    const query = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
            query.append(key, String(value));
        }
    });
    return query.toString();
};

export const authService = {
    register: (data: any) => 
        request<any>('/auth/register', { method: 'POST', body: JSON.stringify(data) }),
    
    login: (data: any) => 
        request<any>('/auth/login', { method: 'POST', body: JSON.stringify(data) }),
    
    logout: () => 
        request<any>('/auth/logout', { method: 'POST' }),
    
    refreshToken: () => 
        request<any>('/auth/refresh', { method: 'POST' }),
};

export const userService = {
    getProfile: () => 
        request<any>('/users/profile', { method: 'GET' }),
    
    updateProfile: (data: any) => 
        request<any>('/users/profile', { method: 'PUT', body: JSON.stringify(data) }),
    
    list: (page = 1, pageSize = 10) => 
        request<any>(`/users?page=${page}&page_size=${pageSize}`, { method: 'GET' }),
    
    getById: (id: number) => 
        request<any>(`/users/${id}`, { method: 'GET' }),
    
    delete: (id: number) => 
        request<any>(`/users/${id}`, { method: 'DELETE' }),
};

export const donationService = {
    create: (data: any) => 
        request<any>('/donations', { method: 'POST', body: JSON.stringify(data) }),
    
    getMyDonations: (page = 1, pageSize = 10) => 
        request<any>(`/donations/my?page=${page}&page_size=${pageSize}`, { method: 'GET' }),
    
    list: (filters: { status?: string; org_id?: number; page?: number; page_size?: number } = {}) => {
        const query = buildQueryString(filters);
        return request<any>(`/donations?${query}`, { method: 'GET' });
    },
    
    getById: (id: number) => 
        request<any>(`/donations/${id}`, { method: 'GET' }),
    
    updateStatus: (id: number, status: string) => 
        request<any>(`/donations/${id}/status`, { method: 'PUT', body: JSON.stringify({ status }) }),
    
    approve: (id: number) => 
        request<any>(`/donations/${id}/approve`, { method: 'POST' }),
    
    reject: (id: number, reason: string) => 
        request<any>(`/donations/${id}/reject`, { method: 'POST', body: JSON.stringify({ reason }) }),
    
    delete: (id: number) => 
        request<any>(`/donations/${id}`, { method: 'DELETE' }),
};

export const inventoryService = {
    list: (filters: { org_id?: number; page?: number; page_size?: number; category?: string; status?: string } = {}) => {
        const query = buildQueryString(filters);
        return request<any>(`/inventory?${query}`, { method: 'GET' });
    },
    
    getById: (id: number) => 
        request<any>(`/inventory/${id}`, { method: 'GET' }),
    
    create: (data: any) => 
        request<any>('/inventory', { method: 'POST', body: JSON.stringify(data) }),
    
    update: (id: number, data: any) => 
        request<any>(`/inventory/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    
    allocate: (id: number, quantity: number, reason: string) => 
        request<any>(`/inventory/${id}/allocate`, { method: 'POST', body: JSON.stringify({ quantity, reason }) }),
    
    distribute: (id: number, quantity: number, recipient: string) => 
        request<any>(`/inventory/${id}/distribute`, { method: 'POST', body: JSON.stringify({ quantity, recipient }) }),
    
    deallocate: (id: number, quantity: number) => 
        request<any>(`/inventory/${id}/deallocate`, { method: 'POST', body: JSON.stringify({ quantity }) }),
    
    getStats: (orgId: number) => 
        request<any>(`/inventory/stats?org_id=${orgId}`, { method: 'GET' }),
    
    delete: (id: number) => 
        request<any>(`/inventory/${id}`, { method: 'DELETE' }),
};

export const organisationService = {
    list: () => 
        request<any>('/organisations', { method: 'GET' }),
    
    getById: (id: number) => 
        request<any>(`/organisations/${id}`, { method: 'GET' }),
    
    getByEmail: (email: string) => 
        request<any>(`/organisations/email/${email}`, { method: 'GET' }),
    
    getStats: (id: number) => 
        request<any>(`/organisations/${id}/stats`, { method: 'GET' }),
    
    create: (data: any) => 
        request<any>('/organisations', { method: 'POST', body: JSON.stringify(data) }),
    
    update: (id: number, data: any) => 
        request<any>(`/organisations/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    
    delete: (id: number) => 
        request<any>(`/organisations/${id}`, { method: 'DELETE' }),
};

export const analyticsService = {
    getDonorImpact: () => 
        request<any>('/analytics/donor-impact', { method: 'GET' }),
    
    getDonationTrends: (orgId: number) => 
        request<any>(`/analytics/trends?org_id=${orgId}`, { method: 'GET' }),
    
    getCategoryBreakdown: (orgId: number) => 
        request<any>(`/analytics/categories?org_id=${orgId}`, { method: 'GET' }),
    
    getSustainabilityMetrics: (orgId: number) => 
        request<any>(`/analytics/sustainability?org_id=${orgId}`, { method: 'GET' }),
    
    getOrgPerformance: (orgId: number) => 
        request<any>(`/analytics/org-performance?org_id=${orgId}`, { method: 'GET' }),
    
    getSystemOverview: () => 
        request<any>('/analytics/system-overview', { method: 'GET' }),
};

export const uploadService = {
    uploadImages: async (files: File[]) => {
        const formData = new FormData();
        files.forEach(file => {
            formData.append('images', file);
        });
        return request<any>('/uploads/images', { method: 'POST', body: formData });
    }
};
export interface ApiResponse<T = any> {
    data?: T;
    message?: string;
    error?: ApiError;
}

export interface ApiError {
    code: string;
    message: string;
    details?: string[];
}

export interface PaginatedResponse<T> {
    data: T[];
    pagination: {
        page: number;
        page_size: number;
        total: number;
        total_pages: number;
    };
}
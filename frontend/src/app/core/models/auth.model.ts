import { User } from './user.model';

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    username: string;
    password: string;
}

export interface AuthResponse {
    user: User;
    token: string;
    refresh_token: string;
    expires_at: Date;
}

export interface RefreshTokenRequest {
    refresh_token: string;
}

export interface RefreshTokenResponse {
    token: string;
    refresh_token: string;
    expires_at: Date;
}
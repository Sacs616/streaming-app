import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, throwError, timer } from 'rxjs';
import { tap, catchError, switchMap } from 'rxjs/operators';
import { ApiService } from './api.service';
import { StorageService } from './storage.service';
import { environment } from '../../../environments/environment';
import {
    User,
    LoginRequest,
    RegisterRequest,
    AuthResponse,
    RefreshTokenRequest,
    RefreshTokenResponse
} from '../models';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private currentUserSubject: BehaviorSubject<User | null>;
    public currentUser$: Observable<User | null>;
    private tokenRefreshTimer?: any;

    constructor(
        private api: ApiService,
        private storage: StorageService,
        private router: Router
    ) {
        const storedUser = this.storage.getItem<User>(environment.userKey);
        this.currentUserSubject = new BehaviorSubject<User | null>(storedUser);
        this.currentUser$ = this.currentUserSubject.asObservable();

        // Start token refresh timer if user is logged in
        if (storedUser && this.getToken()) {
            this.scheduleTokenRefresh();
        }
    }

    public get currentUserValue(): User | null {
        return this.currentUserSubject.value;
    }

    /**
     * Register a new user
     */
    register(request: RegisterRequest): Observable<AuthResponse> {
        return this.api.post<AuthResponse>('/auth/register', request).pipe(
            tap(response => this.handleAuthSuccess(response)),
            catchError(error => {
                console.error('Registration failed:', error);
                return throwError(() => error);
            })
        );
    }

    /**
     * Login user
     */
    login(request: LoginRequest): Observable<AuthResponse> {
        return this.api.post<AuthResponse>('/auth/login', request).pipe(
            tap(response => this.handleAuthSuccess(response)),
            catchError(error => {
                console.error('Login failed:', error);
                return throwError(() => error);
            })
        );
    }

    /**
     * Logout user
     */
    logout(): Observable<any> {
        return this.api.post('/logout', {}).pipe(
            tap(() => this.handleLogout()),
            catchError(error => {
                // Logout locally even if server request fails
                this.handleLogout();
                return throwError(() => error);
            })
        );
    }

    /**
     * Refresh access token
     */
    refreshToken(): Observable<RefreshTokenResponse> {
        const refreshToken = this.getRefreshToken();

        if (!refreshToken) {
            this.handleLogout();
            return throwError(() => ({ message: 'No refresh token available' }));
        }

        const request: RefreshTokenRequest = { refresh_token: refreshToken };

        return this.api.post<RefreshTokenResponse>('/auth/refresh', request).pipe(
            tap(response => {
                this.setToken(response.token);
                this.setRefreshToken(response.refresh_token);
                this.scheduleTokenRefresh(response.expires_at);
            }),
            catchError(error => {
                console.error('Token refresh failed:', error);
                this.handleLogout();
                return throwError(() => error);
            })
        );
    }

    /**
     * Get current user from server
     */
    getCurrentUser(): Observable<{ user: User }> {
        return this.api.get<{ user: User }>('/me').pipe(
            tap(response => {
                this.storage.setItem(environment.userKey, response.user);
                this.currentUserSubject.next(response.user);
            })
        );
    }

    /**
     * Check if user is authenticated
     */
    isAuthenticated(): boolean {
        const token = this.getToken();
        if (!token) {
            return false;
        }

        // Check if token is expired
        try {
            const tokenPayload = this.decodeToken(token);
            const expirationDate = new Date(tokenPayload.exp * 1000);
            return expirationDate > new Date();
        } catch {
            return false;
        }
    }

    /**
     * Check if user has admin role
     */
    isAdmin(): boolean {
        return this.currentUserValue?.role === 'admin';
    }

    /**
     * Get access token
     */
    getToken(): string | null {
        return this.storage.getItem<string>(environment.tokenKey);
    }

    /**
     * Get refresh token
     */
    getRefreshToken(): string | null {
        return this.storage.getItem<string>(environment.refreshTokenKey);
    }

    /**
     * Set access token
     */
    private setToken(token: string): void {
        this.storage.setItem(environment.tokenKey, token);
    }

    /**
     * Set refresh token
     */
    private setRefreshToken(refreshToken: string): void {
        this.storage.setItem(environment.refreshTokenKey, refreshToken);
    }

    /**
     * Handle successful authentication
     */
    private handleAuthSuccess(response: AuthResponse): void {
        this.setToken(response.token);
        this.setRefreshToken(response.refresh_token);
        this.storage.setItem(environment.userKey, response.user);
        this.currentUserSubject.next(response.user);
        this.scheduleTokenRefresh(response.expires_at);
    }

    /**
     * Handle logout
     */
    private handleLogout(): void {
        this.storage.removeItem(environment.tokenKey);
        this.storage.removeItem(environment.refreshTokenKey);
        this.storage.removeItem(environment.userKey);
        this.currentUserSubject.next(null);

        if (this.tokenRefreshTimer) {
            clearTimeout(this.tokenRefreshTimer);
        }

        this.router.navigate(['/login']);
    }

    /**
     * Schedule automatic token refresh
     */
    private scheduleTokenRefresh(expiresAt?: Date): void {
        if (this.tokenRefreshTimer) {
            clearTimeout(this.tokenRefreshTimer);
        }

        const token = this.getToken();
        if (!token) {
            return;
        }

        try {
            const tokenPayload = this.decodeToken(token);
            const expirationDate = expiresAt ? new Date(expiresAt) : new Date(tokenPayload.exp * 1000);
            const now = new Date();

            // Refresh token 5 minutes before expiration
            const refreshTime = expirationDate.getTime() - now.getTime() - (5 * 60 * 1000);

            if (refreshTime > 0) {
                this.tokenRefreshTimer = setTimeout(() => {
                    this.refreshToken().subscribe({
                        error: (error) => console.error('Auto token refresh failed:', error)
                    });
                }, refreshTime);
            } else {
                // Token already expired or about to expire, refresh immediately
                this.refreshToken().subscribe();
            }
        } catch (error) {
            console.error('Error scheduling token refresh:', error);
        }
    }

    /**
     * Decode JWT token
     */
    private decodeToken(token: string): any {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(
                atob(base64)
                    .split('')
                    .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                    .join('')
            );
            return JSON.parse(jsonPayload);
        } catch (error) {
            throw new Error('Invalid token');
        }
    }
}
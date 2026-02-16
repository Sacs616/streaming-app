import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable, throwError, TimeoutError } from 'rxjs';
import { catchError, timeout } from 'rxjs/operators';
import { environment } from '../../../environments/environment';

export interface RequestOptions {
    headers?: HttpHeaders | { [header: string]: string | string[] };
    params?: HttpParams | { [param: string]: string | string[] };
    responseType?: 'json' | 'text' | 'blob' | 'arraybuffer';
    withCredentials?: boolean;
}

@Injectable({
    providedIn: 'root'
})
export class ApiService {
    private readonly baseUrl = environment.apiUrl;
    private readonly defaultTimeout = environment.apiTimeout;

    constructor(private http: HttpClient) { }

    get<T>(endpoint: string, options?: any): Observable<T> {
        return (this.http.get<T>(`${this.baseUrl}${endpoint}`, options) as Observable<T>).pipe(
            timeout({ first: this.defaultTimeout }),
            catchError(this.handleError)
        ) as Observable<T>;
    }

    post<T>(endpoint: string, data: any, options?: any): Observable<T> {
        return (this.http.post<T>(`${this.baseUrl}${endpoint}`, data, options) as Observable<T>).pipe(
            timeout({ first: this.defaultTimeout }),
            catchError(this.handleError)
        ) as Observable<T>;
    }

    put<T>(endpoint: string, data: any, options?: any): Observable<T> {
        return (this.http.put<T>(`${this.baseUrl}${endpoint}`, data, options) as Observable<T>).pipe(
            timeout({ first: this.defaultTimeout }),
            catchError(this.handleError)
        ) as Observable<T>;
    }

    patch<T>(endpoint: string, data: any, options?: any): Observable<T> {
        return (this.http.patch<T>(`${this.baseUrl}${endpoint}`, data, options) as Observable<T>).pipe(
            timeout({ first: this.defaultTimeout }),
            catchError(this.handleError)
        ) as Observable<T>;
    }

    delete<T>(endpoint: string, options?: any): Observable<T> {
        return (this.http.delete<T>(`${this.baseUrl}${endpoint}`, options) as Observable<T>).pipe(
            timeout({ first: this.defaultTimeout }),
            catchError(this.handleError)
        ) as Observable<T>;
    }

    private handleError(error: any): Observable<never> {
        let errorMessage = 'An unknown error occurred';

        if (error instanceof TimeoutError) {
            errorMessage = 'Request timed out. Please try again.';
        } else if (error.error instanceof ErrorEvent) {
            // Client-side error
            errorMessage = `Error: ${error.error.message}`;
        } else if (error.status) {
            // Server-side error
            errorMessage = error.error?.message || error.error?.error || `Error: ${error.status} ${error.statusText}`;
        }

        console.error('API Error:', error);
        return throwError(() => ({ message: errorMessage, originalError: error }));
    }
}
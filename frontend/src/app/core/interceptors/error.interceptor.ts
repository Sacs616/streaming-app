import { Injectable } from '@angular/core';
import {
    HttpRequest,
    HttpHandler,
    HttpEvent,
    HttpInterceptor,
    HttpErrorResponse
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {

    intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
        return next.handle(request).pipe(
            catchError((error: HttpErrorResponse) => {
                let errorMessage = 'An unknown error occurred';

                if (error.error instanceof ErrorEvent) {
                    // Client-side error
                    errorMessage = `Error: ${error.error.message}`;
                } else {
                    // Server-side error
                    if (error.error?.error) {
                        errorMessage = error.error.error;
                    } else if (error.error?.message) {
                        errorMessage = error.error.message;
                    } else if (error.message) {
                        errorMessage = error.message;
                    } else {
                        errorMessage = `Error Code: ${error.status}\nMessage: ${error.statusText}`;
                    }
                }

                console.error('HTTP Error:', errorMessage);

                return throwError(() => ({
                    message: errorMessage,
                    status: error.status,
                    originalError: error
                }));
            })
        );
    }
}
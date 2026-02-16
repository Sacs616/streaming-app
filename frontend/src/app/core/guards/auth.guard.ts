import { Injectable } from '@angular/core';
import {
    ActivatedRouteSnapshot,
    RouterStateSnapshot,
    Router,
    UrlTree
} from '@angular/router';
import { Observable } from 'rxjs';
import { AuthService } from '../services/auth.service';

@Injectable({
    providedIn: 'root'
})
export class AuthGuard {
    constructor(
        private authService: AuthService,
        private router: Router
    ) { }

    canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {

        if (this.authService.isAuthenticated()) {
            // Check if route requires admin
            if (route.data['requiresAdmin'] && !this.authService.isAdmin()) {
                return this.router.createUrlTree(['/']);
            }
            return true;
        }

        // Not logged in, redirect to login with return URL
        return this.router.createUrlTree(['/login'], {
            queryParams: { returnUrl: state.url }
        });
    }

    canActivateChild(
        childRoute: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
        return this.canActivate(childRoute, state);
    }
}
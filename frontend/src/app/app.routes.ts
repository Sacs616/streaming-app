import { Routes } from '@angular/router';
import { AuthGuard } from './core/guards/auth.guard';
import { GuestGuard } from './core/guards/guest.guard';

export const routes: Routes = [
    {
        path: '',
        redirectTo: '/browse',
        pathMatch: 'full'
    },
    {
        path: 'browse',
        canActivate: [AuthGuard],
        loadChildren: () => import('./features/home/home.module').then(m => m.HomeModule)
    },
    {
        path: '',
        loadChildren: () => import('./features/auth/auth.module').then(m => m.AuthModule)
    },
    {
        path: '**',
        redirectTo: '/browse'
    }
];

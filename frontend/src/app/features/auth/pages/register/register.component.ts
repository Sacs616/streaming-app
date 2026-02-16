import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, AbstractControl, ValidationErrors } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../../core/services/auth.service';

@Component({
    selector: 'app-register',
    templateUrl: './register.component.html',
    styleUrls: ['./register.component.scss'],
    standalone: false
})
export class RegisterComponent implements OnInit {
    registerForm!: FormGroup;
    loading = false;
    error = '';
    showPassword = false;

    constructor(
        private formBuilder: FormBuilder,
        private authService: AuthService,
        private router: Router
    ) { }

    ngOnInit(): void {
        this.registerForm = this.formBuilder.group({
            email: ['', [Validators.required, Validators.email]],
            username: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(20)]],
            password: ['', [Validators.required, Validators.minLength(6), this.passwordStrengthValidator]],
            confirmPassword: ['', [Validators.required]]
        }, {
            validators: this.passwordMatchValidator
        });
    }

    get f() {
        return this.registerForm.controls;
    }

    onSubmit(): void {
        if (this.registerForm.invalid) {
            this.markFormGroupTouched(this.registerForm);
            return;
        }

        this.loading = true;
        this.error = '';

        const { email, username, password } = this.registerForm.value;

        this.authService.register({ email, username, password }).subscribe({
            next: () => {
                this.router.navigate(['/browse']);
            },
            error: (error) => {
                this.error = error.message || 'Registration failed. Please try again.';
                this.loading = false;
            }
        });
    }

    togglePasswordVisibility(): void {
        this.showPassword = !this.showPassword;
    }

    private passwordMatchValidator(control: AbstractControl): ValidationErrors | null {
        const password = control.get('password');
        const confirmPassword = control.get('confirmPassword');

        if (!password || !confirmPassword) {
            return null;
        }

        return password.value === confirmPassword.value ? null : { passwordMismatch: true };
    }

    private passwordStrengthValidator(control: AbstractControl): ValidationErrors | null {
        const value = control.value;

        if (!value) {
            return null;
        }

        const hasNumber = /[0-9]/.test(value);
        const hasLetter = /[a-zA-Z]/.test(value);

        const valid = hasNumber && hasLetter;

        return valid ? null : { passwordStrength: true };
    }

    private markFormGroupTouched(formGroup: FormGroup): void {
        Object.keys(formGroup.controls).forEach(key => {
            const control = formGroup.get(key);
            control?.markAsTouched();

            if (control instanceof FormGroup) {
                this.markFormGroupTouched(control);
            }
        });
    }
}
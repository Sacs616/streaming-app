export interface User {
    id: number;
    email: string;
    username: string;
    role: 'user' | 'admin';
    created_at: Date;
    updated_at: Date;
}

export interface UserProfile extends User {
    avatar_url?: string;
    bio?: string;
}
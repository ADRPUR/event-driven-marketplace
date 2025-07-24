export type UserRole = "user" | "admin";

export interface UserDetails {
    firstName: string;
    lastName: string;
    dateOfBirth?: string; // ISO date
    phone?: string;
    address?: string; // Poate fi și JSON sau obiect separat
    photoPath?: string; // doar path relativ/absolut spre imagine
    thumbnailPath?: string;
}

export interface User {
    id: string;
    email: string;
    role: UserRole;
    details: UserDetails;
}

export interface UsersListResponse {
    users: User[];
}

export interface UpdateProfileRequest {
    firstName?: string;
    lastName?: string;
    dateOfBirth?: string;
    phone?: string;
    address?: string;
    photoPath?: string;
    thumbnailPath?: string;
}


export interface UploadPhotoResponse {
    photoPath: string;
    thumbnailPath?: string;
}


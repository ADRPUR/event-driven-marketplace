import type {User} from "./user";

export interface LoginResponse {
    accessToken: string;
    refreshToken?: string;
    sessionToken?: string;
    expiresAt?: number;    // timestamp unix
    user: User;
}

export interface RegisterResponse {
    id: string;
}

export interface MeResponse {
    user: User;
}

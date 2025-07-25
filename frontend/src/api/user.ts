import axios from "axios";
import type {User, UsersListResponse} from "../types/user";

const api = axios.create({ baseURL: "http://localhost:8090" });

export async function getProfile(token: string): Promise<{ user: User }> {
    const res = await api.get<{ user: User }>("/me", {
        headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
}

export async function updateProfile(token: string, data: Partial<User["details"]>): Promise<{ user: User }> {
    const res = await api.put<{ user: User }>("/me", data, {
        headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
}

export async function uploadPhoto(token: string, file: File): Promise<{ photoPath: string; thumbnailPath?: string }> {
    const form = new FormData();
    form.append("photo", file);
    const res = await api.post<{ photoPath: string; thumbnailPath?: string }>("/me/photo", form, {
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
        },
    });
    return res.data;
}

// Obține toți userii (admin)
export async function getAllUsers(token: string): Promise<UsersListResponse> {
    const res = await api.get<UsersListResponse>("/users", {
        headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
}

// Update user (admin poate edita pe oricine)
export async function updateUser(token: string, id: string, data: Partial<User["details"]>): Promise<{ user: User }> {
    const res = await api.put<{ user: User }>(`/users/${id}`, data, {
        headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
}

// Șterge user (admin)
export async function deleteUser(token: string, id: string): Promise<void> {
    await api.delete(`/users/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
    });
}


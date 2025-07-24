import axios from "axios";
import type { UploadPhotoResponse,  UsersListResponse, UpdateProfileRequest } from "../types/user";
import  type { UserDetails } from "../types/user";

const api = axios.create({ baseURL: "http://localhost:8090" });

export async function uploadPhoto(token: string, file: File): Promise<UploadPhotoResponse> {
    const form = new FormData();
    form.append("photo", file);
    const res = await api.post<UploadPhotoResponse>("/profile/photo", form, {
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
        },
    });
    return res.data;
}

export async function getAllUsers(token: string): Promise<UsersListResponse> {
    const res = await api.get<UsersListResponse>("/users", {
        headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
}

export async function updateProfile(token: string, data: UpdateProfileRequest): Promise<UserDetails> {
    const res = await api.put<UserDetails>("/profile", data, {
        headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
}

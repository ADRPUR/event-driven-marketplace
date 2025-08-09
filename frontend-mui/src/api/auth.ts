import axios from "axios";
import type { AuthResponse } from "../types/user";

export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8090";

const api = axios.create({ baseURL: API_URL });

export async function login(email: string, password: string): Promise<AuthResponse> {
    const res = await api.post("/login", { email, password });
    return res.data;
}

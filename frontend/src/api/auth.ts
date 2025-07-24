import axios from "axios";
import type {LoginResponse, RegisterResponse, MeResponse} from "../types/auth";

// Set your API base URL here!
const api = axios.create({ baseURL: "http://localhost:8090" });

export async function login(email: string, password: string): Promise<LoginResponse> {
    const res = await api.post<LoginResponse>("/login", { email, password });
    return res.data;
}

export async function register(data: {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
}): Promise<RegisterResponse> {
    const res = await api.post<RegisterResponse>("/register", data);
    return res.data;
}

export async function getProfile(token: string): Promise<MeResponse> {
    const res = await api.get<MeResponse>("/me", { headers: { Authorization: `Bearer ${token}` } });
    return res.data;
}


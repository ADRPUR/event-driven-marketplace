import { create } from "zustand";
import type { User } from "../types/user";

type AuthState = {
  user: User | null;
  token: string | null;
  login: (user: User, token: string) => void;
  logout: () => void;
  init: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  login: (user, token) => {
    set({ user, token });
    localStorage.setItem("user", JSON.stringify(user));
    localStorage.setItem("token", token);
  },
  logout: () => {
    set({ user: null, token: null });
    localStorage.removeItem("user");
    localStorage.removeItem("token");
  },
  init: () => {
    const token = localStorage.getItem("token");
    const userStr = localStorage.getItem("user");
    set({
      user: userStr ? JSON.parse(userStr) : null,
      token: token ?? null,
    });
  },
}));

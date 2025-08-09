#!/bin/bash

set -e

mkdir -p src/api
mkdir -p src/components
mkdir -p src/features/auth
mkdir -p src/features/profile
mkdir -p src/features/users
mkdir -p src/features/products
mkdir -p src/hooks
mkdir -p src/layouts
mkdir -p src/pages
mkdir -p src/router
mkdir -p src/store
mkdir -p src/theme
mkdir -p src/types

touch src/main.tsx

# Pagini cu template basic React component
pages=("LoginPage" "RegisterPage" "ProductsPage" "ProfilePage" "UsersPage")
for page in "${pages[@]}"; do
  cat <<EOF > src/pages/${page}.tsx
import React from "react";
export default function ${page}() {
  return <div>${page} works!</div>;
}
EOF
done

# Basic layout files
cat <<EOF > src/layouts/DashboardLayout.tsx
import React from "react";
export default function DashboardLayout({ children }: { children?: React.ReactNode }) {
  return <div>{children}</div>;
}
EOF

cat <<EOF > src/layouts/AppLayout.tsx
import React from "react";
export default function AppLayout({ children }: { children?: React.ReactNode }) {
  return <div>{children}</div>;
}
EOF

# Store basic authStore
cat <<EOF > src/store/authStore.ts
import { create } from "zustand";
import type { User } from "../types/user";
type AuthState = {
  user: User | null;
  token: string | null;
  login: (user: User, token: string) => void;
  logout: () => void;
};
export const useAuthStore = create<AuthState>((set) => ({
  user: null, token: null,
  login: (user, token) => set({ user, token }),
  logout: () => set({ user: null, token: null }),
}));
EOF

# Types example
cat <<EOF > src/types/user.ts
export type UserRole = "user" | "admin";
export interface User {
  id: string;
  email: string;
  role: UserRole;
  firstName?: string;
  lastName?: string;
  photo?: string;
  // etc...
}
EOF

echo "✅ Folders and base files created!"


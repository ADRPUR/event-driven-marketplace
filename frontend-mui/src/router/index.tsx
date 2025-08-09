import { lazy, Suspense } from "react";
import { Routes, Route, Navigate, Outlet } from "react-router-dom";
import DashboardLayout from "../layouts/DashboardLayout";
import { useAuthStore } from "../store/authStore";
import { CircularProgress, Box } from "@mui/material";

// Lazy imports:
const LoginPage = lazy(() => import("../pages/LoginPage"));
const RegisterPage = lazy(() => import("../pages/RegisterPage"));
const ProductsPage = lazy(() => import("../pages/ProductsPage"));
const ProfilePage = lazy(() => import("../pages/ProfilePage"));
const UsersPage = lazy(() => import("../pages/UsersPage"));

// Spinner wrapper:
function Loader() {
    return (
        <Box minHeight="60vh" display="flex" alignItems="center" justifyContent="center">
            <CircularProgress />
        </Box>
    );
}

// Guards:
function ProtectedRoute() {
    const { token } = useAuthStore();
    if (!token) return <Navigate to="/login" replace />;
    return <Outlet />;
}
function RoleRoute({ role }: { role: string }) {
    const { user } = useAuthStore();
    if (!user || user.role !== role) return <Navigate to="/products" replace />;
    return <Outlet />;
}
function GuestRoute() {
    const { token } = useAuthStore();
    if (token) return <Navigate to="/products" replace />;
    return <Outlet />;
}

export default function Router() {
    return (
        <Suspense fallback={<Loader />}>
            <Routes>
                {/* Only for guests (not logged in): */}
                <Route element={<GuestRoute />}>
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/register" element={<RegisterPage />} />
                </Route>

                {/* Protected routes */}
                <Route element={<ProtectedRoute />}>
                    <Route element={<DashboardLayout />}>
                        <Route path="/products" element={<ProductsPage />} />
                        <Route path="/profile" element={<ProfilePage />} />
                        {/* Only for admin */}
                        <Route element={<RoleRoute role="admin" />}>
                            <Route path="/users" element={<UsersPage />} />
                        </Route>
                        <Route path="/" element={<Navigate to="/products" replace />} />
                        <Route path="*" element={<Navigate to="/products" replace />} />
                    </Route>
                </Route>
            </Routes>
        </Suspense>
    );
}

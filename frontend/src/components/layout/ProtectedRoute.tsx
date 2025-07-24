import { useAuthStore } from "../../store/authStore";
import { Navigate, Outlet } from "react-router-dom";

export default function ProtectedRoute() {
    const { user } = useAuthStore();
    return user ? <Outlet /> : <Navigate to="/login" replace />;
}

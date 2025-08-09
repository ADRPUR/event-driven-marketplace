import { useState } from "react";
import { Card, Typography, TextField, Button, Box, Alert } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { login as apiLogin } from "../api/auth";
import { useAuthStore } from "../store/authStore";
import axios from "axios";

export default function LoginPage() {
  const [form, setForm] = useState({ email: "", password: "" });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const { login: setAuth } = useAuthStore();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res = await apiLogin(form.email, form.password); // <- AuthResponse
      setAuth(res.user, res.accessToken);
      navigate("/products", { replace: true });
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Server unavailable");
      } else {
        setError("Unknown error");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
      <Box minHeight="100vh" display="flex" alignItems="center" justifyContent="center" bgcolor="#f4f6f8">
        <Card sx={{ p: 4, width: 360 }}>
          <Typography variant="h5" gutterBottom>Sign in</Typography>
          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
          <Box component="form" autoComplete="off" onSubmit={handleSubmit}>
            <TextField
                label="Email"
                type="email"
                margin="normal"
                fullWidth
                required
                value={form.email}
                onChange={e => setForm(f => ({ ...f, email: e.target.value }))}
                autoFocus
            />
            <TextField
                label="Password"
                type="password"
                margin="normal"
                fullWidth
                required
                value={form.password}
                onChange={e => setForm(f => ({ ...f, password: e.target.value }))}
                onKeyDown={e => e.key === "Enter" && handleSubmit(e as any)}
            />
            <Button
                type="submit"
                variant="contained"
                fullWidth
                sx={{ mt: 2 }}
                disabled={loading}
            >
              {loading ? "Authenticating..." : "Login"}
            </Button>
          </Box>
        </Card>
      </Box>
  );
}

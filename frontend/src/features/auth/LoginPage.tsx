import { useState } from "react";
import { Form, Input, Button, Typography, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { useNavigate, Link } from "react-router-dom";
import { useAuthStore } from "../../store/authStore";
import { login } from "../../api/auth";
import axios from "axios";

const { Title } = Typography;

type LoginForm = {
    email: string;
    password: string;
};


export default function LoginPage() {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const { login: setAuth } = useAuthStore();

    async function onFinish(values: LoginForm) {
        setLoading(true);
        try {
            const res = await login(values.email, values.password);
            setAuth(res.user, res.accessToken);
            message.success("Login success!");
            navigate("/products");
        } catch (err) {
            if (axios.isAxiosError(err)) {
                message.error(err.response?.data?.error || "Login failed");
            } else {
                message.error("Login failed");
            }
        } finally {
            setLoading(false);
        }
    }

    return (
        <div style={{ maxWidth: 380, margin: "70px auto" }}>
            <Title level={2}>Login</Title>
            <Form<LoginForm> onFinish={onFinish} layout="vertical">
                <Form.Item name="email" label="Email" rules={[{ required: true }, { type: "email" }]}>
                    <Input prefix={<UserOutlined />} autoFocus />
                </Form.Item>
                <Form.Item name="password" label="Password" rules={[{ required: true }]}>
                    <Input.Password prefix={<LockOutlined />} />
                </Form.Item>
                <Form.Item>
                    <Button block type="primary" htmlType="submit" loading={loading}>
                        Login
                    </Button>
                </Form.Item>
                <div style={{ textAlign: "center" }}>
                    <Link to="/register">Don't have an account? Register</Link>
                </div>
            </Form>
        </div>
    );
}

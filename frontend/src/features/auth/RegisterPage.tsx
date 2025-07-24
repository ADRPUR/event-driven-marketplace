import { useState } from "react";
import { Form, Input, Button, Typography, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { useNavigate, Link } from "react-router-dom";
import { register } from "../../api/auth";
import axios from "axios";

const { Title } = Typography;

type RegisterForm = {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
};

export default function RegisterPage() {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    async function onFinish(values: RegisterForm) {
        setLoading(true);
        try {
            await register(values);
            message.success("Account created!");
            navigate("/login");
        } catch (err) {
            if (axios.isAxiosError(err)) {
                message.error(err.response?.data?.error || "Register failed");
            } else {
                message.error("Register failed");
            }
        } finally {
            setLoading(false);
        }
    }

    return (
        <div style={{ maxWidth: 420, margin: "70px auto" }}>
            <Title level={2}>Register</Title>
            <Form<RegisterForm> onFinish={onFinish} layout="vertical">
                <Form.Item name="email" label="Email" rules={[{ required: true }, { type: "email" }]}>
                    <Input prefix={<UserOutlined />} />
                </Form.Item>
                <Form.Item name="password" label="Password" rules={[{ required: true }, { min: 6 }]}>
                    <Input.Password prefix={<LockOutlined />} />
                </Form.Item>
                <Form.Item name="firstName" label="First Name" rules={[{ required: true }]}>
                    <Input />
                </Form.Item>
                <Form.Item name="lastName" label="Last Name" rules={[{ required: true }]}>
                    <Input />
                </Form.Item>
                <Form.Item>
                    <Button block type="primary" htmlType="submit" loading={loading}>
                        Register
                    </Button>
                </Form.Item>
                <div style={{ textAlign: "center" }}>
                    <Link to="/login">Already have an account? Login</Link>
                </div>
            </Form>
        </div>
    );
}

import { useEffect, useState } from "react";
import { Card, Typography, Avatar, Button, Form, Input, Upload, message } from "antd";
import { UserOutlined, UploadOutlined } from "@ant-design/icons";
import { useAuthStore } from "../../store/authStore";
import { getProfile } from "../../api/auth";
import type {User} from "../../types/user.ts";
import axios from "axios";

const { Title } = Typography;

type ProfileForm = {
    firstName: string;
    lastName: string;
    phone?: string;
    address?: string;
    // alte câmpuri dacă vrei
};

export default function ProfilePage() {
    const { user, token } = useAuthStore();
    const [loading, setLoading] = useState(false);
    const [editing, setEditing] = useState(false);
    const [profile, setProfile] = useState<User | null>(user);

    useEffect(() => {
        if (token) {
            getProfile(token).then(({ user }) => setProfile(user));
        }
    }, [token]);

    async function onFinish() {
        setLoading(true);
        try {
            // TODO: POST /profile update API call
            message.success("Profile updated!");
            setEditing(false);
        } catch (err) {
            if (axios.isAxiosError(err)) {
                message.error(err.response?.data?.error || "Update failed");
            } else {
                message.error("Update failed");
            }
        } finally {
            setLoading(false);
        }
    }

    return (
        <Card style={{ maxWidth: 450, margin: "50px auto" }}>
            <div style={{ textAlign: "center" }}>
                <Avatar
                    size={80}
                    src={profile?.details?.photoPath || undefined}
                    icon={<UserOutlined />}
                />
                <Title level={4} style={{ margin: 12 }}>
                    {profile?.details?.firstName} {profile?.details?.lastName}
                </Title>
            </div>
            <Form<ProfileForm>
                layout="vertical"
                initialValues={{
                    firstName: profile?.details?.firstName,
                    lastName: profile?.details?.lastName,
                    phone: profile?.details?.phone,
                    address: profile?.details?.address,
                }}
                onFinish={onFinish}
                disabled={!editing}
            >
                <Form.Item name="firstName" label="First Name">
                    <Input />
                </Form.Item>
                <Form.Item name="lastName" label="Last Name">
                    <Input />
                </Form.Item>
                <Form.Item name="phone" label="Phone">
                    <Input />
                </Form.Item>
                <Form.Item name="address" label="Address">
                    <Input />
                </Form.Item>
                <Form.Item label="Photo">
                    <Upload beforeUpload={() => false}>
                        <Button icon={<UploadOutlined />}>Upload Photo</Button>
                    </Upload>
                </Form.Item>
                <Form.Item>
                    {editing ? (
                        <Button htmlType="submit" type="primary" loading={loading}>
                            Save
                        </Button>
                    ) : (
                        <Button onClick={() => setEditing(true)}>Edit</Button>
                    )}
                </Form.Item>
            </Form>
        </Card>
    );
}

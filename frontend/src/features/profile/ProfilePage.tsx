import { useEffect, useState } from "react";
import {
    Card, Typography, Avatar, Button, Form, Input, Upload, message, DatePicker, Tabs, Row, Col, Menu, Divider
} from "antd";
import {UserOutlined, UploadOutlined, InfoCircleOutlined, TeamOutlined, BellOutlined} from "@ant-design/icons";
import { useAuthStore } from "../../store/authStore";
import { getProfile, updateProfile, uploadPhoto } from "../../api/user";
import dayjs from "dayjs";
import type { User } from "../../types/user";
import type { UploadRequestOption as RcCustomRequestOptions } from "rc-upload/lib/interface";
import "./ProfilePage.css";

const { Title, Text } = Typography;

type ProfileForm = {
    firstName: string;
    lastName: string;
    phone?: string;
    address?: string;
    dateOfBirth?: dayjs.Dayjs;
};

const sidebarItems = [
    { key: "overview", label: "Personal Info", icon: <UserOutlined /> },
    { key: "edit", label: "Edit Profile", icon: <UploadOutlined /> },
    { key: "password", label: "Change Password", icon: <UserOutlined /> },
    { key: "information", label: "Information", icon: <InfoCircleOutlined /> },
    { key: "social", label: "Social", icon: <TeamOutlined /> },
    { key: "notification", label: "Notification", icon: <BellOutlined /> },
];

export default function ProfilePage() {
    const { token, login } = useAuthStore();
    const [profile, setProfile] = useState<User | null>(null);
    const [tab, setTab] = useState("overview");
    const [loading, setLoading] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [form] = Form.useForm<ProfileForm>();

    useEffect(() => {
        if (token) {
            getProfile(token).then((resp ) => {
                setProfile(resp.user);
                console.log(resp.user);
                form.setFieldsValue({
                    firstName: resp.user.details?.firstName,
                    lastName: resp.user.details?.lastName,
                    phone: resp.user.details?.phone,
                    address: resp.user.details?.address,
                    dateOfBirth: resp.user.details?.dateOfBirth ? dayjs(resp.user.details.dateOfBirth).toISOString() : undefined,
                });
            });
        }
    }, [token, form]);

    async function handleUpload(options: RcCustomRequestOptions) {
        if (!token) return;
        setUploading(true);
        try {
            const file = options.file as File;
            const res = await uploadPhoto(token, file);
            message.success("Photo uploaded!");
            setProfile((old) =>
                old
                    ? {
                        ...old,
                        details: {
                            ...old.details,
                            photoPath: res.photoPath,
                            thumbnailPath: res.thumbnailPath,
                        },
                    }
                    : old
            );
        } catch (err) {
            message.error("Upload failed: " + err);
        } finally {
            setUploading(false);
        }
    }

    async function onFinish(values: ProfileForm) {
        if (!token) return;
        setLoading(true);
        try {
            const submitValues = {
                ...values,
                dateOfBirth: values.dateOfBirth
                    ? values.dateOfBirth.format("YYYY-MM-DD")
                    : undefined,
            };
            const res = await updateProfile(token, submitValues);
            setProfile(res.user);
            message.success("Profile updated!");
            login(res.user, token);
            form.setFieldsValue({
                ...submitValues,
                dateOfBirth: submitValues.dateOfBirth ? dayjs(submitValues.dateOfBirth) : undefined,
            });
            setTab("overview");
        } catch (err) {
            message.error("Update failed: " + (err instanceof Error ? err.message : "Unknown error"));
        } finally {
            setLoading(false);
        }
    }


    return (
        <Row style={{ minHeight: "calc(100vh - 64px)", background: "#f5f7fa" }}>
            <Col flex="260px" style={{ background: "#fff", borderRight: "1px solid #f0f0f0" }}>
                <Menu
                    mode="inline"
                    selectedKeys={[tab]}
                    style={{ height: "100%", paddingTop: 40 }}
                    items={sidebarItems}
                    onClick={({ key }) => setTab(key)}
                />
            </Col>
            <Col flex="auto" style={{ padding: "40px 0", display: "flex", justifyContent: "center" }}>
                <div style={{ width: "100%", maxWidth: 860 }}>
                    {/* Header card cu banner și avatar */}
                    <Card
                        className="profile-banner-card"
                        cover={
                            <div className="profile-banner" />
                        }
                        style={{ marginBottom: 40 }}
                        bodyStyle={{ paddingTop: 50, paddingBottom: 20, textAlign: "center" }}
                    >
                        <Avatar
                            size={96}
                            src={profile?.details?.photoPath || undefined}
                            icon={<UserOutlined />}
                            className="profile-avatar"
                        />
                        <div style={{ marginTop: 8 }}>
                            <Upload
                                accept="image/*"
                                showUploadList={false}
                                customRequest={async (options) => {
                                    try {
                                        await handleUpload(options);
                                        options.onSuccess?.({}, options.file);
                                    } catch {
                                        options.onError?.(new Error("Upload error"));
                                    }
                                }}
                                disabled={uploading || tab !== "edit"}
                            >
                                <Button size="small" icon={<UploadOutlined />} loading={uploading} disabled={tab !== "edit"}>
                                    Upload Photo
                                </Button>
                            </Upload>
                        </div>
                        <Title level={3} style={{ margin: 10 }}>{profile?.details?.firstName} {profile?.details?.lastName}</Title>
                        <Text type="secondary">{profile?.email}</Text>
                    </Card>

                    <Card>
                        <Tabs activeKey={tab} onChange={setTab} tabBarStyle={{ display: "none" }}>
                            <Tabs.TabPane key="overview">
                                <div style={{ padding: 16 }}>
                                    <Row gutter={32}>
                                        <Col span={12}>
                                            <Title level={5}>Personal Info</Title>
                                            <div style={{ marginBottom: 10 }}><b>Name:</b> {profile?.details?.firstName} {profile?.details?.lastName}</div>
                                            <div style={{ marginBottom: 10 }}><b>Date of Birth:</b> {profile?.details?.dateOfBirth}</div>
                                            <div style={{ marginBottom: 10 }}><b>Phone:</b> {profile?.details?.phone}</div>
                                            <div style={{ marginBottom: 10 }}><b>Address:</b> {profile?.details?.address}</div>
                                        </Col>
                                        <Col span={12}>
                                            <Title level={5}>Account</Title>
                                            <div style={{ marginBottom: 10 }}><b>Email:</b> {profile?.email}</div>
                                            <div style={{ marginBottom: 10 }}><b>Role:</b> {profile?.role}</div>
                                        </Col>
                                    </Row>
                                </div>
                            </Tabs.TabPane>
                            <Tabs.TabPane key="edit">
                                <Form<ProfileForm>
                                    layout="vertical"
                                    form={form}
                                    onFinish={onFinish}
                                    disabled={loading}
                                    style={{ maxWidth: 500, margin: "0 auto" }}
                                >
                                    <Form.Item name="firstName" label="First Name" rules={[{ required: true }]}>
                                        <Input />
                                    </Form.Item>
                                    <Form.Item name="lastName" label="Last Name" rules={[{ required: true }]}>
                                        <Input />
                                    </Form.Item>
                                    <Form.Item name="phone" label="Phone">
                                        <Input />
                                    </Form.Item>
                                    <Form.Item name="address" label="Address">
                                        <Input />
                                    </Form.Item>
                                    <Form.Item name="dateOfBirth" label="Date of Birth">
                                        <DatePicker
                                            format="YYYY-MM-DD"
                                            style={{ width: "100%" }}
                                            allowClear
                                            onChange={(_date, dateString) => {
                                                form.setFieldsValue({
                                                    dateOfBirth: Array.isArray(dateString) ? dateString[0] : dateString || undefined,
                                                });
                                            }}
                                        />
                                    </Form.Item>
                                    <Form.Item>
                                        <Button htmlType="submit" type="primary" loading={loading}>Save Changes</Button>
                                        <Button style={{ marginLeft: 8 }} onClick={() => setTab("overview")}>Cancel</Button>
                                    </Form.Item>
                                </Form>
                            </Tabs.TabPane>
                            <Tabs.TabPane key="password">
                                {/* Change password tab – poți implementa separat */}
                                <Divider />
                                <Title level={5}>Change Password</Title>
                                <p>Functionality coming soon…</p>
                            </Tabs.TabPane>
                            <Tabs.TabPane key="information">
                                <Divider />
                                <Title level={5}>Information</Title>
                                <p>Functionality coming soon…</p>
                            </Tabs.TabPane>
                            <Tabs.TabPane key="social">
                                <Divider />
                                <Title level={5}>Social</Title>
                                <p>Functionality coming soon…</p>
                            </Tabs.TabPane>
                            <Tabs.TabPane key="notification">
                                <Divider />
                                <Title level={5}>Notification</Title>
                                <p>Functionality coming soon…</p>
                            </Tabs.TabPane>
                        </Tabs>
                    </Card>
                </div>
            </Col>
        </Row>
    );
}

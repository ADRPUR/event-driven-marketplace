import {Layout, Avatar, Dropdown, Menu, Typography, type MenuProps} from "antd";
import { UserOutlined } from "@ant-design/icons";
import { Link, Outlet, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../store/authStore";

const { Header, Content } = Layout;
const { Title } = Typography;

export default function AppLayout() {
    const { user, logout } = useAuthStore();
    const navigate = useNavigate();

    const menuItems: MenuProps["items"] = [
        { key: "profile", label: <Link to="/profile">Profile</Link> },
        ...(user?.role === "admin"
            ? [{ key: "users", label: <Link to="/users">Users</Link> }]
            : []),
        { type: "divider" as const },
        {
            key: "logout",
            label: <span onClick={() => { logout(); navigate("/login"); }}>Logout</span>,
            danger: true,
        },
    ];

    const menu = <Menu items={menuItems} />;

    return (
        <Layout style={{ minHeight: "100vh" }}>
            <Header style={{ display: "flex", justifyContent: "space-between", alignItems: "center", color: "#fff" }}>
                <Title level={4} style={{ color: "#fff", margin: 0 }}>
                    Marketplace
                </Title>
                <Dropdown overlay={menu} placement="bottomRight">
                    <Avatar src={user?.details?.photoPath || undefined} icon={<UserOutlined />} style={{ cursor: "pointer" }} />
                </Dropdown>
            </Header>
            <Content style={{ padding: 24, background: "#fff" }}>
                <Outlet />
            </Content>
        </Layout>
    );
}

import {
  AppBar,
  Toolbar,
  Typography,
  Avatar,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemButton,
  useTheme,
} from "@mui/material";
import { useState } from "react";
import { Outlet, Link, useLocation, useNavigate } from "react-router-dom";
import { Logout, Person, List as ListIcon, Group } from "@mui/icons-material";
import { useAuthStore } from "../store/authStore";
import {API_URL} from "../api/auth.ts";

// Helper for user initials
function getInitials(name?: string, fallback = "?") {
  if (!name) return fallback;
  return name.split(" ").map(part => part[0]).join("").slice(0, 2).toUpperCase();
}

export default function DashboardLayout() {
  const { user, logout } = useAuthStore();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();

  const menuItems = [
    { key: "products", label: "Products", icon: <ListIcon />, to: "/products" },
    { key: "profile", label: "Profile", icon: <Person />, to: "/profile" },
    ...(user?.role === "admin"
        ? [{ key: "users", label: "Users", icon: <Group />, to: "/users" }]
        : []),
  ];

  return (
      <Box sx={{ display: "flex", minHeight: "100vh" }}>
        {/* Sidebar */}
        <Drawer
            variant="permanent"
            sx={{
              width: 220,
              flexShrink: 0,
              [`& .MuiDrawer-paper`]: {
                width: 220,
                boxSizing: "border-box",
                bgColor: "#F4F6F8",
                borderRight: `1px solid ${theme.palette.divider}`,
              },
            }}
        >
          <Toolbar>
            <Typography variant="h6" color="primary" sx={{ ml: 1 }}>
              {/* Small logo or title */}
              MKT
            </Typography>
          </Toolbar>
          <List>
            {menuItems.map(item => (
                <ListItem
                    key={item.key}
                    disablePadding
                    sx={{ borderRadius: 2, mb: 0.5 }}
                >
                  <ListItemButton
                      component={Link}
                      to={item.to}
                      selected={location.pathname.startsWith(item.to)}
                      sx={{
                        borderRadius: 2,
                        ...(location.pathname.startsWith(item.to)
                            ? {
                              bgColor: "#e0e7ff",
                              color: theme.palette.primary.main,
                              fontWeight: 600,
                            }
                            : {}),
                      }}
                  >
                    <ListItemIcon
                        sx={{
                          color: location.pathname.startsWith(item.to)
                              ? theme.palette.primary.main
                              : "inherit",
                        }}
                    >
                      {item.icon}
                    </ListItemIcon>
                    <ListItemText primary={item.label} />
                  </ListItemButton>
                </ListItem>
            ))}
          </List>
        </Drawer>

        {/* Main Content */}
        <Box sx={{ flexGrow: 1, bgColor: "#F4F6F8" }}>
          {/* Topbar */}
          <AppBar position="static" color="inherit" elevation={0}>
            <Toolbar sx={{ justifyContent: "space-between" }}>
              <Typography variant="h5" color="primary" sx={{ fontWeight: 700 }}>
                Marketplace
              </Typography>
              <IconButton onClick={e => setAnchorEl(e.currentTarget)} sx={{ ml: 2 }}>
                <Avatar
                    src={
                      user?.photo
                          ? `${API_URL}${user.photo}?t=${Date.now()}`
                          : undefined
                    }
                    sx={{ bgColor: "#6366f1" }}
                >
                  {getInitials(`${user?.firstName} ${user?.lastName}`)}
                </Avatar>
              </IconButton>
              <Menu
                  open={!!anchorEl}
                  anchorEl={anchorEl}
                  onClose={() => setAnchorEl(null)}
                  anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
              >
                <MenuItem
                    onClick={() => {
                      navigate("/profile");
                      setAnchorEl(null);
                    }}
                >
                  Profile
                </MenuItem>
                <MenuItem
                    onClick={() => {
                      logout();
                      navigate("/login");
                      setAnchorEl(null);
                    }}
                >
                  Logout <Logout fontSize="small" style={{ marginLeft: 8 }} />
                </MenuItem>
              </Menu>
            </Toolbar>
          </AppBar>
          {/* Page content */}
          <Box sx={{ p: { xs: 1, sm: 3 } }}>
            <Outlet />
          </Box>
        </Box>
      </Box>
  );
}

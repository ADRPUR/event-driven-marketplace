import { createTheme } from "@mui/material/styles";

const theme = createTheme({
    palette: {
        mode: "light",
        primary: { main: "#6366F1" },
        secondary: { main: "#F59E42" },
        background: { default: "#F4F6F8", paper: "#fff" },
    },
    shape: { borderRadius: 8 },
    components: {
        MuiButton: { styleOverrides: { root: { textTransform: "none" } } },
    },
});
export default theme;

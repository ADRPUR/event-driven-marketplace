import { BrowserRouter } from "react-router-dom";
import Router from "./router";
import { CssBaseline } from "@mui/material";

export default function App() {
    return (
        <BrowserRouter>
            <CssBaseline />
            <Router />
        </BrowserRouter>
    );
}
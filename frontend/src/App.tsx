import { BrowserRouter } from "react-router-dom";
import AppRouter from "./routes/AppRouter.tsx";


export default function App() {
    return (
        <BrowserRouter>
            <AppRouter/>
        </BrowserRouter>
    );
}

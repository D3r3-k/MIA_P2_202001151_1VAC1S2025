import { MiaContext } from "@/contexts/MiaContext";
import { useContext } from "react";

export const useMia = () => {
    const context = useContext(MiaContext);
    if (!context) throw new Error("useAuth must be used within an AuthProvider");
    return context;
}
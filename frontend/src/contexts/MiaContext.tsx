"use client";

import { MiaContextType, UserData } from "@/types/AuthTypes";
import { redirect } from "next/navigation";
import { createContext, useEffect, useState } from "react";

export const MiaContext = createContext<MiaContextType | undefined>(undefined);

export const MiaProvider = ({ children }: { children: React.ReactNode }) => {
    // Hooks
    // States
    const [systemState, setSystemState] = useState(false)
    const [loading, setLoading] = useState<boolean>(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);
    // ? Drivers
    // ? Partitions
    // ? Files
    // ? Loading 
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [userData, setUserData] = useState<UserData | null>(null);
    // Effects
    useEffect(() => {
        const fetchSystemStatus = async () => {
            try {
                const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/status`);
                if (response.ok) {
                    const data = await response.json();
                    setSystemState(data.status);
                } else {
                    throw new Error("Error al obtener el estado del sistema.");
                }
            } catch (error: any) {
                setErrorMsg(error.message);
            }
        };
        fetchSystemStatus();
    }, []);

    // Handlers
    // Functions
    const executeCommand = async (command: string) => {
        setLoading(true);
        setErrorMsg(null);
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/execute`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ command }),
            });
            if (response.ok) {
                const data = await response.json();
                return data.output;
            } else {
                throw new Error("Error al ejecutar el comando.");
            }
        } catch (error: any) {
            setErrorMsg(`Error del servidor. ${error.message}`);
            return error.message;
        } finally {
            setLoading(false);
        }
    }

    const login = async ({ partition_id, username, password }: { partition_id: string; username: string; password: string }) => {
        setLoading(true);
        setErrorMsg(null);
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ partition_id, username, password }),
            });
            if (response.ok) {
                const data = await response.json();
                document.cookie = `authToken=${data.token};`;
                setUserData(data.userData);
                setIsAuthenticated(true);
            } else {
                setIsAuthenticated(false);
                setUserData(null);
                throw new Error("Error al iniciar sesión. Verifica tus credenciales.");
            }
        } catch (error: any) {
            setErrorMsg(error.message);
        } finally {
            setLoading(false);
        }
    };
    const logout = () => {
        setIsAuthenticated(false);
        document.cookie = "authToken=; path=/; max-age=0";
        redirect("/");
    }
    // Renders
    return (
        <MiaContext.Provider value={{
            systemState,
            isAuthenticated,
            userData,
            login,
            logout,
            loading,
            errorMsg,
            executeCommand
        }}>
            {children}
        </MiaContext.Provider>
    );
}
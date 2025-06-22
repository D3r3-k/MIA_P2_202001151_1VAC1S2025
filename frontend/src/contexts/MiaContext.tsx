"use client";

import { MiaContextType, UserData } from "@/types/MiaTypes";
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
        const cookies = document.cookie.split(";").reduce((acc: Record<string, string>, cookie) => {
            const [key, value] = cookie.trim().split("=");
            acc[key] = value;
            return acc;
        }, {});
        if (cookies.authToken) {
            const userDataLocal = localStorage.getItem("userData");
            if (userDataLocal) {
                try {
                    setUserData(JSON.parse(userDataLocal));
                    setIsAuthenticated(true);
                } catch (error) {
                    console.error("Error al parsear userData desde localStorage:", error);
                    setUserData(null);
                    setIsAuthenticated(false);
                }
            }
        }
    }, []);

    useEffect(() => {
        const fetchSystemStatus = async () => {
            try {
                const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/status`);
                if (response.ok) {
                    const data = await response.json();
                    setSystemState(data.status);

                    if (data.authToken?.username) {
                        setIsAuthenticated(true);
                        if (data.userData?.username) setUserData(data.authToken);
                        document.cookie = `authToken=true; path=/;`;
                    } else {
                        document.cookie = "authToken=false; path=/; max-age=0";
                        localStorage.removeItem("userData");
                        setIsAuthenticated(false);
                        setUserData(null);
                    }
                } else {
                    setSystemState(false);
                    localStorage.removeItem("userData");
                    document.cookie = "authToken=false; path=/; max-age=0";
                    setIsAuthenticated(false);
                    setUserData(null);
                }
            } catch (error) {
                console.error("Error al obtener el estado del sistema:", error);
                setSystemState(false);
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
                throw new Error(await response.text() || "Error al ejecutar el comando.");
            }
        } catch (error) {
            setErrorMsg(`Error del servidor. ${error}`);
            return error instanceof Error ? error.message : "Error desconocido al ejecutar el comando.";
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
                if (data.user_data) {
                    localStorage.setItem("userData", JSON.stringify(data.user_data));
                }
                setUserData(data.user_data);
                setIsAuthenticated(true);
                if (typeof window !== "undefined") {
                    window.location.href = "/drives";
                }

            } else {
                setIsAuthenticated(false);
                setUserData(null);
                throw new Error("Error al iniciar sesión. Verifica tus credenciales.");
            }
        } catch (error) {
            console.error("Error al iniciar sesión:", error);
        } finally {
            setLoading(false);
        }
        return true;
    };
    const logout = () => {
        fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/logout`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
        }).catch((error) => {
            console.error("Error al cerrar sesión:", error);
        });
        setIsAuthenticated(false);
        setUserData(null);
        document.cookie = "authToken=; path=/; max-age=0";
        localStorage.removeItem("userData");
        if (typeof window !== "undefined") {
            window.location.href = "/";
        }
        return Promise.resolve();
    }

    // Renders
    return (
        <MiaContext.Provider value={{
            loading,
            setLoading,
            systemState,
            isAuthenticated,
            userData,
            executeCommand,
            login,
            logout,
            errorMsg,
        }}>
            {children}
        </MiaContext.Provider>
    );
}
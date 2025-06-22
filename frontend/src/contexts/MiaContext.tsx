"use client";

import { MiaContextType, SingleToastData, UserData } from "@/types/MiaTypes";
import { createContext, useEffect, useState } from "react";

export const MiaContext = createContext<MiaContextType | undefined>(undefined);

export const MiaProvider = ({ children }: { children: React.ReactNode }) => {
    // Hooks
    // States
    const [systemState, setSystemState] = useState(false);
    const [loading, setLoading] = useState<boolean>(false);
    // ? Toast
    const [toast, setToast] = useState<SingleToastData>({
        message: "",
        subtitle: "",
        type: "info",
        duration: 6000,
        visible: false,
    });
    // ? Partitions
    // ? Files
    // ? Loading 
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [userData, setUserData] = useState<UserData | null>(null);

    // Effects
    useEffect(() => {
        try {
            const cookies = document.cookie.split(";").reduce((acc: Record<string, string>, cookie) => {
                const [key, value] = cookie.trim().split("=");
                if (key && value) acc[key] = value;
                return acc;
            }, {});

            if (cookies.authToken) {
                const userDataLocal = localStorage.getItem("userData");
                if (userDataLocal) {
                    const parsedUserData = JSON.parse(userDataLocal);
                    setUserData(parsedUserData);
                    setIsAuthenticated(true);
                }
            }
        } catch (error) {
            const err = error as Error;
            activateToast("error", err.message || "Error al cargar los datos del usuario desde localStorage.", "Por favor, recarga la página.");
            setUserData(null);
            setIsAuthenticated(false);
        }
    }, []);

    useEffect(() => {
        const fetchSystemStatus = async () => {
            try {
                const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/status`);
                if (!response.ok) throw new Error("Fallo al obtener el estado del sistema");

                const data = await response.json();
                setSystemState(data.status ?? false);

                const hasValidAuth = data.authToken?.username;

                if (hasValidAuth) {
                    setIsAuthenticated(true);
                    setUserData(data.authToken);
                    document.cookie = `authToken=true; path=/;`;
                } else {
                    handleSessionReset();
                }
            } catch (error) {
                const err = error as Error;
                activateToast("error", err.message || "Error al obtener el estado del sistema.", "Por favor, recarga la página.");
                setSystemState(false);
                handleSessionReset();
            }
        };

        fetchSystemStatus();
    }, []);

    // Utils
    const handleSessionReset = () => {
        localStorage.removeItem("userData");
        document.cookie = "authToken=false; path=/; max-age=0";
        setIsAuthenticated(false);
        setUserData(null);
    };

    // Handlers
    const executeCommand = async (command: string) => {
        setLoading(true);
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/execute`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ command }),
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || "Error al ejecutar el comando.");
            }
            const data = await response.json();
            return data.output;
        } catch (error) {
            const err = error as Error;
            activateToast("error", "Error al ejecutar el comando.", err.message || "Error desconocido.");
            return "Error al ejecutar el comando.";
        } finally {
            setLoading(false);
        }
    };

    const login = async ({
        partition_id,
        username,
        password,
    }: {
        partition_id: string;
        username: string;
        password: string;
    }) => {
        setLoading(true);
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ partition_id, username, password }),
            });

            if (!response.ok) {
                throw new Error("Error al iniciar sesión. Verifica tus credenciales.");
            }

            const data = await response.json();
            document.cookie = `authToken=${data.token}; path=/;`;
            if (data.user_data) {
                localStorage.setItem("userData", JSON.stringify(data.user_data));
                setUserData(data.user_data);
                setIsAuthenticated(true);
            }

            if (typeof window !== "undefined") {
                window.location.href = "/drives";
            }

            return true;
        } catch (error) {
            const err = error as Error;
            activateToast("error", "Error al iniciar sesión.", err.message || "Error desconocido.");
            setIsAuthenticated(false);
            setUserData(null);
            return false;
        } finally {
            setLoading(false);
        }
    };

    const logout = async () => {
        setLoading(true);
        try {
            await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/logout`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });
        } catch (error) {
            const err = error as Error;
            activateToast("error", "Error al cerrar sesión.", err.message || "Por favor, recarga la página.");
        } finally {
            setLoading(false);
            handleSessionReset();
            if (typeof window !== "undefined") {
                window.location.href = "/";
            }
        }
    };

    const activateToast = (
        type: "info" | "success" | "error",
        message: string,
        subtitle?: string,
        duration: number = 6000
    ) => {
        if (toast.visible) {
            setToast((prev) => ({ ...prev, visible: false }));
            setTimeout(() => {
                setToast({ type, message, subtitle, duration, visible: true });
            }, 1000);
        } else {
            setToast({ type, message, subtitle, duration, visible: true });
        }
    };

    const handleClose = () => {
        setToast((prev) => ({ ...prev, visible: false }));
    };

    // Renders
    return (
        <MiaContext.Provider
            value={{
                loading,
                setLoading,
                systemState,
                isAuthenticated,
                userData,
                executeCommand,
                login,
                logout,
                activateToast,
                toast,
                handleClose,
            }}
        >
            {children}
        </MiaContext.Provider>
    );
};

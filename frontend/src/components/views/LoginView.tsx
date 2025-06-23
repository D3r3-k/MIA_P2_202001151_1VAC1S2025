import { useState, useEffect } from "react";
import { User } from "lucide-react";
import { useSearchParams } from "next/navigation";
import { useMia } from "@/hooks/useMia";
import Head from "next/head";

export default function LoginView({ onLoginSuccess }: { onLoginSuccess: () => void }) {
    const { login, userData } = useMia();
    const params = useSearchParams();

    const [loginData, setLoginData] = useState({
        partitionId: "",
        username: "",
        password: "",
    });

    useEffect(() => {
        if (userData?.partition_id) {
            onLoginSuccess();
        }
    }, [userData, onLoginSuccess]);

    useEffect(() => {
        const partitionId = params.get("partition_id") || "";
        setLoginData((prev) => ({
            ...prev,
            partitionId,
        }));
    }, [params]);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        const success = await login({
            partition_id: loginData.partitionId,
            username: loginData.username,
            password: loginData.password,
        });
        if (success) {
            onLoginSuccess();
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        const newValue = name === "partitionId" ? value.toUpperCase() : value;
        setLoginData((prev) => ({ ...prev, [name]: newValue }));
    };

    return (
        <>
            <Head>
                <title>Iniciar Sesión - F2 MIA</title>
                <meta name="description" content="Inicia sesión para acceder al sistema de discos" />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
            </Head>

            <main className="flex-1 p-6 ml-72 bg-gray-950 min-h-screen flex items-center justify-center">
                <div className="space-y-6">
                    <div className="flex flex-col items-center shadow-md rounded-lg p-6 mx-auto bg-gradient-to-br from-gray-700 to-gray-800">
                        <div className="relative flex items-center justify-center bg-gradient-to-br from-red-500 to-red-600 rounded-full w-16 h-16 mb-4 shadow-lg shadow-red-500/25">
                            <User scale={48} className="text-white" />
                        </div>
                        <h1 className="text-3xl font-bold text-white mb-2">Iniciar Sesión</h1>
                        <p className="text-white text-sm mb-4">
                            Ingresa tus credenciales para acceder al sistema
                        </p>
                        <form
                            className="flex flex-col space-y-4 w-72"
                            autoComplete="off"
                            onSubmit={handleLogin}
                        >
                            <input
                                type="text"
                                name="partitionId"
                                value={loginData.partitionId}
                                onChange={handleChange}
                                required
                                placeholder="ID de Partición"
                                className="bg-gray-800 border border-gray-600 rounded-md p-2 text-white"
                            />
                            <input
                                type="text"
                                name="username"
                                value={loginData.username}
                                onChange={handleChange}
                                required
                                placeholder="Usuario"
                                className="bg-gray-800 border border-gray-600 rounded-md p-2 text-white"
                            />
                            <input
                                type="password"
                                name="password"
                                value={loginData.password}
                                onChange={handleChange}
                                required
                                placeholder="Contraseña"
                                className="bg-gray-800 border border-gray-600 rounded-md p-2 text-white"
                            />
                            <button
                                type="submit"
                                className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-md"
                            >
                                Iniciar Sesión
                            </button>
                        </form>
                    </div>
                </div>
            </main>
        </>
    );
}

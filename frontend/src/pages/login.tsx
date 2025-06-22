import { useMia } from '@/hooks/useMia';
import { User } from 'lucide-react'
import Head from 'next/head';
import { useSearchParams } from 'next/navigation';
import { useRouter } from 'next/router';
import React, { useEffect, useState } from 'react'

export default function Login() {
    // Hooks
    const { login } = useMia();
    const params = useSearchParams();
    const router = useRouter();
    // States
    const [loginData, setLoginData] = useState({
        partitionId: "",
        username: "",
        password: "",
    });

    // Effects
    useEffect(() => {
        const partitionId = params.get("partition_id") || "";
        setLoginData((prev) => ({
            ...prev,
            partitionId,
        }));
    }, [params]);

    useEffect(() => {
        const searchParams = new URLSearchParams();
        if (loginData.partitionId) {
            searchParams.set("partition_id", loginData.partitionId);
        }
    }, [loginData.partitionId]);

    // Handlers
    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        const loginres = await login({
            partition_id: loginData.partitionId,
            username: loginData.username,
            password: loginData.password,
        });
        if (loginres) {
            router.push("/drives");
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        let newValue = value;
        if (name === "partitionId") {
            newValue = value.toUpperCase();
            const searchParams = new URLSearchParams();
            searchParams.set("partition_id", newValue);
            window.history.replaceState({}, "", `?${searchParams.toString()}`);
        }
        setLoginData((prev) => ({ ...prev, [name]: newValue }));
    };
    return (
        <>
            <Head>
                <title>Iniciar Sesión - F2 MIA</title>
                <meta name="description" content="Inicia sesión para acceder a la gestión de discos y particiones." />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <main className="flex-1 p-6 ml-72 bg-gray-950 min-h-screen items-center flex justify-center">
                <div className="space-y-6">
                    <div className="flex flex-col items-center shadow-md rounded-lg p-6 mx-auto bg-gradient-to-br from-gray-700 to-gray-800">
                        <div className="relative flex items-center justify-center bg-gradient-to-br from-red-500 to-red-600 rounded-full w-16 h-16 mb-4 shadow-lg shadow-red-500/25">
                            <User scale={48} className="text-white" />
                        </div>
                        <h1 className="text-3xl font-bold text-white mb-2">Iniciar Sesión</h1>
                        <p className="text-white text-sm mb-4">
                            Ingresa tus credenciales para acceder al sistema
                        </p>
                        <div className="flex flex-col items-center border-t border-gray-600 pt-4 w-full">
                            <form className="flex flex-col space-y-4 w-full" autoComplete="off" onSubmit={handleLogin}>
                                <input
                                    type="text"
                                    placeholder="ID de Partición"
                                    className="bg-gray-800 border border-gray-600 rounded-md p-2 text-white"
                                    autoComplete="off"
                                    name="partitionId"
                                    value={loginData.partitionId}
                                    onChange={handleChange}
                                    required
                                />
                                <input
                                    type="text"
                                    placeholder="Usuario"
                                    className="bg-gray-800 border border-gray-600 rounded-md p-2 text-white"
                                    autoComplete="off"
                                    name="username"
                                    value={loginData.username}
                                    onChange={handleChange}
                                    required
                                />
                                <input
                                    type="password"
                                    placeholder="Contraseña"
                                    className="bg-gray-800 border border-gray-600 rounded-md p-2 text-white"
                                    autoComplete="off"
                                    name="password"
                                    value={loginData.password}
                                    onChange={handleChange}
                                    required
                                />
                                <button className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-md">
                                    Iniciar Sesión
                                </button>
                            </form>
                        </div>
                    </div>
                </div>
            </main>
        </>
    )
}

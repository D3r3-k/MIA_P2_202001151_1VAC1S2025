"use client";
import { AlertTriangle, Compass, FileQuestion, Search, Terminal } from 'lucide-react'
import { useRouter } from 'next/router';
import React from 'react'

export default function Page404() {
    const router = useRouter();
    const handleGoHome = () => {
        router.push('/');
    };
    return (
        <main className="min-h-screen bg-gray-950 flex items-center justify-center p-6">
            <div className="max-w-4xl w-full space-y-8">
                {/* Header Section */}
                <div className="text-center space-y-6">
                    {/* 404 Animation */}
                    <div className="relative">
                        <div className="text-8xl md:text-9xl font-bold text-transparent bg-gradient-to-r from-red-500 via-red-400 to-red-600 bg-clip-text animate-pulse">
                            404
                        </div>
                        <div className="absolute inset-0 text-8xl md:text-9xl font-bold text-red-500/10 blur-sm">
                            404
                        </div>
                    </div>

                    {/* Error Icon */}
                    <div className="flex justify-center">
                        <div className="relative">
                            <div className="w-24 h-24 bg-gradient-to-br from-red-500/20 to-red-600/10 rounded-full flex items-center justify-center border border-red-500/30 backdrop-blur-sm">
                                <FileQuestion className="w-12 h-12 text-red-400" />
                            </div>
                            <div className="absolute -top-2 -right-2 w-8 h-8 bg-yellow-500/20 rounded-full flex items-center justify-center border border-yellow-500/30">
                                <AlertTriangle className="w-4 h-4 text-yellow-400" />
                            </div>
                        </div>
                    </div>

                    {/* Error Message */}
                    <div className="space-y-3">
                        <h1 className="text-3xl md:text-4xl font-bold text-white">
                            Página No Encontrada
                        </h1>
                        <p className="text-lg text-gray-400 max-w-2xl mx-auto">
                            Lo sentimos, la página que estás buscando no existe en el sistema de archivos MIA F2.
                            Puede que haya sido movida, eliminada o la URL sea incorrecta.
                        </p>
                    </div>
                </div>

                {/* Action buttons */}
                <div className="flex flex-wrap justify-center gap-4">
                    <button
                        onClick={handleGoHome}
                        className="bg-gradient-to-r from-corinto-600 to-corinto-700 hover:from-corinto-700 hover:to-corinto-800 text-white px-6 py-3 text-lg font-semibold shadow-lg shadow-corinto-600/25 hover:shadow-corinto-600/40 transition-all duration-300 flex gap-2 items-center rounded-lg cursor-pointer"
                    >
                        <Terminal className="w-5 h-5 mr-2" />
                        Regresar a la terminar
                    </button>
                </div>

                {/* Search Suggestion */}
                <div className="bg-gray-800/30 border-gray-700/50">
                    <div className="p-6">
                        <div className="flex items-center justify-center space-x-4">
                            <div className="p-3 rounded-xl bg-gray-700/50 border border-gray-600">
                                <Search className="w-6 h-6 text-gray-400" />
                            </div>
                            <div className="text-center">
                                <h3 className="font-semibold text-white text-lg mb-1">
                                    ¿Buscas algo específico?
                                </h3>
                                <p className="text-gray-400 text-sm">
                                    Asegúrate de que la URL sea correcta o intenta buscar en el sistema.
                                </p>
                            </div>
                            <div className="p-3 rounded-xl bg-gray-700/50 border border-gray-600">
                                <Compass className="w-6 h-6 text-gray-400" />
                            </div>
                        </div>
                    </div>
                </div>

                {/* Footer Info */}
                <div className="text-center text-gray-500 text-sm">
                    <p>Sistema de Archivos MIA F2 • Versión 1.0</p>
                    <p className="mt-1">
                        Si el problema persiste, contacta al administrador del sistema
                    </p>
                </div>
            </div>
        </main>
    )
}

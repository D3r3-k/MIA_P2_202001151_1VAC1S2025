"use client";

import { useMia } from '@/hooks/useMia';
import { HardDrive, LogIn, Power, Terminal, User } from 'lucide-react';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import React from 'react';

const menuItems = [
    { icon: Terminal, label: 'Consola', path: '/' },
    { icon: HardDrive, label: 'Discos', path: '/drives' },
    { icon: LogIn, label: 'Iniciar Sesión', path: '/login' }
];

export default function Sidebar() {
    const { systemState, isAuthenticated, logout, userData } = useMia();
    const pathname = usePathname();

    const isActive = (path: string): boolean => {
        if (!pathname) return false;
        const currentPartitionPath = userData?.partition_id
            ? `/drives/${userData.partition_id[0]}/${userData.partition_id}`
            : '';

        if (path === '/') return pathname === '/';
        if (path === '/login') return pathname.startsWith('/login');
        if (path === '/drives') return pathname.startsWith('/drives');
        if (path === currentPartitionPath) return pathname === currentPartitionPath;

        return pathname === path;
    };


    const activeClass = "bg-gradient-to-r from-corinto-600 to-corinto-700 text-white shadow-lg shadow-corinto-600/25";
    const defaultClass = "text-gray-300 hover:text-white hover:bg-gray-700/60";

    return (
        <aside className="fixed left-0 top-0 h-screen w-72 bg-gradient-to-b from-gray-900 via-gray-900 to-gray-800 border-r border-gray-700/50 backdrop-blur-xl flex flex-col z-50">
            {/* Logo / Encabezado */}
            <div className="p-5 border-b border-gray-700/30 flex-shrink-0">
                <div className="flex items-center space-x-3 mb-4">
                    <div className="relative">
                        <div className="w-12 h-12 bg-gradient-to-br from-corinto-500 via-corinto-600 to-corinto-700 rounded-xl flex items-center justify-center shadow-lg shadow-corinto-500/25">
                            <HardDrive className="w-6 h-6 text-white" />
                        </div>
                        <div className="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-gray-900 animate-pulse" />
                    </div>
                    <div className="flex-1 min-w-0">
                        <h1 className="text-white font-bold text-lg tracking-tight truncate">Sistema de Archivos</h1>
                        <p className="text-gray-400 text-sm font-medium">MIA Fase 2</p>
                    </div>
                </div>
            </div>

            {/* Menú de navegación */}
            <div className="flex-1 px-4 py-4">
                <div className="mb-6">
                    <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3 px-2">
                        Navegación
                    </h3>
                    <nav className="space-y-2">
                        {menuItems.map((item, index) => {
                            if (item.path === '/login' && isAuthenticated) return null;

                            const active = isActive(item.path);

                            return (
                                <div key={index}>
                                    <Link
                                        href={item.path}
                                        className={`w-full max-w-full group flex items-center px-3 py-2.5 rounded-xl text-left transition-all duration-200 relative cursor-pointer ${active ? activeClass : defaultClass}`}
                                    >
                                        <div className="flex items-center space-x-3 flex-1">
                                            <item.icon className="w-4 h-4 flex-shrink-0" />
                                            <span className="font-medium text-sm truncate">{item.label}</span>
                                        </div>
                                        {active && (
                                            <div className="w-2 h-2 bg-white rounded-full opacity-90 flex-shrink-0" />
                                        )}
                                    </Link>

                                    {/* Subenlace a partición si aplica */}
                                    {item.path === "/drives" && userData?.partition_id && (
                                        <div className="ml-5 mt-1">
                                            <span className="ml-1 text-xs text-gray-400">Sesión iniciada en</span>
                                            <div className="flex border-l border-gray-500 items-center justify-center ml-1">
                                                <Link
                                                    href={`/drives/${userData.partition_id[0]}/${userData.partition_id}`}
                                                    className={`flex items-center w-full text-left text-gray-300 hover:text-white hover:bg-gray-700/60 px-3 py-2 rounded-md text-sm ml-3 my-1 ${isActive(`/drives/${userData.partition_id[0]}/${userData.partition_id}`) ? activeClass : defaultClass}`}
                                                >
                                                    <div className="flex items-center space-x-3 flex-1">
                                                        {`Partición: ${userData.partition_id}`}
                                                    </div>
                                                    {isActive(`/drives/${userData.partition_id[0]}/${userData.partition_id}`) && (
                                                        <div className="w-2 h-2 bg-white rounded-full opacity-90 flex-shrink-0" />
                                                    )}
                                                </Link>
                                            </div>
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </nav>
                </div>
            </div>

            {/* Usuario autenticado */}
            {isAuthenticated && (
                <>
                    <div className="p-5 border-t border-gray-700/30 flex-shrink-0">
                        <div className="space-y-3">
                            <div className="bg-gradient-to-r from-gray-800/50 to-gray-700/30 rounded-xl p-4 border border-gray-700/30 relative">
                                <div className="flex items-start space-x-3">
                                    <User size={24} className="text-gray-300" />
                                    <div>
                                        <p className="text-xs text-gray-400">Partición</p>
                                        <p className="text-sm font-semibold text-white mb-1">{userData?.partition_id}</p>
                                        <div className="flex items-center space-x-4">
                                            <div>
                                                <span className="text-xs text-gray-400">Usuario</span>
                                                <span className="block text-xs font-medium text-gray-200">{userData?.username}</span>
                                            </div>
                                            <div>
                                                <span className="text-xs text-gray-400">Grupo</span>
                                                <span className="block text-xs font-medium text-gray-200">{userData?.group}</span>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="w-2 h-2 bg-green-400 rounded-full opacity-90 absolute right-0 mx-3" />
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="px-4 py-2 border-t border-gray-700/30">
                        <button
                            onClick={logout}
                            className="w-full flex items-center space-x-3 px-3 py-2.5 rounded-md text-gray-300 hover:text-white hover:bg-corinto-700/60 transition-all duration-200 relative cursor-pointer shadow-md bg-gradient-to-br from-corinto-500 to-corinto-700"
                        >
                            <Power className="w-4 h-4 flex-shrink-0" />
                            <span className="font-medium text-sm">Cerrar Sesión</span>
                        </button>
                    </div>
                </>
            )}

            {/* Estado del sistema */}
            <div className="p-5 border-t border-gray-700/30 flex-shrink-0">
                <div className="bg-gradient-to-r from-gray-800/50 to-gray-700/30 rounded-xl p-4 border border-gray-700/30">
                    <h4 className="text-sm font-semibold text-white mb-3">Estado del Sistema</h4>
                    <div className="space-y-2.5">
                        <div className="flex items-center justify-between">
                            <span className="text-xs text-gray-400">Backend</span>
                            <div className='flex items-center space-x-2'>
                                {systemState
                                    ? <span className="text-xs font-medium text-green-400">Activo</span>
                                    : <span className="text-xs font-medium text-corinto-400">Inactivo</span>}
                                <div className={`w-2 h-2 rounded-full ${systemState ? 'bg-green-400' : 'bg-corinto-400'}`} />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </aside>
    );
}

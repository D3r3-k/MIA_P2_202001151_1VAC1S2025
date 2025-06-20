"use client";

import { useMia } from '@/hooks/useMia';
import { HardDrive, LogIn, Power, Terminal, User } from 'lucide-react'
import { redirect, usePathname } from 'next/navigation';
import React from 'react'

const menuItems = [
    {
        icon: Terminal,
        label: 'Consola',
        path: '/'
    },
    {
        icon: HardDrive,
        label: 'Discos',
        path: '/drives'
    },
    {
        icon: LogIn,
        label: 'Iniciar Sesión',
        path: '/login'
    }
];

export default function Sidebar() {
    // Hooks
    const { systemState, isAuthenticated, logout, userData } = useMia();
    const pathname = usePathname();
    // States
    // Effects
    // Handlers
    // Functions
    // Renders
    const isActive = (path: string) => {
        return pathname === path ? "bg-gradient-to-r from-red-600 to-red-700 text-white shadow-lg shadow-red-600/25" : "text-gray-300 hover:text-white hover:bg-gray-700/60";
    };
    const isActiveIcon = (path: string) => {
        return pathname === path && <div className="w-2 h-2 bg-white rounded-full opacity-90 flex-shrink-0" />;
    };
    return (
        <aside className="fixed left-0 top-0 h-screen w-72 bg-gradient-to-b from-gray-900 via-gray-900 to-gray-800 border-r border-gray-700/50 backdrop-blur-xl flex flex-col z-50">
            <div className="p-5 border-b border-gray-700/30 flex-shrink-0">
                <div className="flex items-center space-x-3 mb-4">
                    <div className="relative">
                        <div className="w-12 h-12 bg-gradient-to-br from-red-500 via-red-600 to-red-700 rounded-xl flex items-center justify-center shadow-lg shadow-red-500/25">
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
            <div className="flex-1 px-4 py-4">
                <div className="mb-6">
                    <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3 px-2">
                        Navegación
                    </h3>
                    <nav className="space-y-2">
                        {menuItems.map((item, index) => (
                            <button
                                onClick={() => redirect(item.path)}
                                key={index}
                                className={`w-full group flex items-center px-3 py-2.5 rounded-xl text-left transition-all duration-200 relative cursor-pointer
                                    ${isActive(item.path)}
                                `}
                            >
                                <div className="flex items-center space-x-3 flex-1">
                                    <item.icon className="w-4 h-4 flex-shrink-0" />
                                    <span className="font-medium text-sm truncate">{item.label}</span>
                                </div>
                                {isActiveIcon(item.path)}
                            </button>
                        ))}
                    </nav>
                </div>
            </div>

            {
                isAuthenticated && (
                    <div className="p-5 border-t border-gray-700/30 flex-shrink-0">
                        <div className="space-y-3">
                            <div className="bg-gradient-to-r from-gray-800/50 to-gray-700/30 rounded-xl p-4 border border-gray-700/30 relative">
                                <div className="flex items-start space-x-3">
                                    <div className="relative">
                                        <User size={24} className="text-gray-300" />
                                    </div>
                                    <div>
                                        <p className="text-xs text-gray-400">Usuario</p>
                                        <p className="text-sm font-semibold text-white mb-1">{userData?.username}</p>
                                        <div className="flex items-center space-x-4">
                                            <div>
                                                <span className="text-xs text-gray-400">Grupo</span>
                                                <span className="block text-xs font-medium text-gray-200">{userData?.group}</span>
                                            </div>
                                            <div>
                                                <span className="text-xs text-gray-400">Permisos</span>
                                                <span className="block text-xs font-medium text-gray-200">{userData?.permissions}</span>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="w-2 h-2 bg-green-400 rounded-full opacity-90 absolute right-0 mx-3" />
                                </div>
                            </div>
                        </div>
                    </div>
                )
            }
            {
                isAuthenticated && (
                    <div className="px-4 py-2 border-t border-gray-700/30">
                        <button
                            onClick={logout}
                            className="w-full flex items-center space-x-3 px-3 py-2.5 rounded-md text-gray-300 hover:text-white hover:bg-red-700/60 transition-all duration-200 relative cursor-pointer shadow-md bg-gradient-to-br from-red-500 to-red-700"
                        >
                            <Power className="w-4 h-4 flex-shrink-0" />
                            <span className="font-medium text-sm">Cerrar Sesión</span>
                        </button>
                    </div>
                )
            }
            <div className="p-5 border-t border-gray-700/30 flex-shrink-0">
                <div className="bg-gradient-to-r from-gray-800/50 to-gray-700/30 rounded-xl p-4 border border-gray-700/30">
                    <h4 className="text-sm font-semibold text-white mb-3">Estado del Sistema</h4>
                    <div className="space-y-2.5">
                        <div className="flex items-center justify-between">
                            <span className="text-xs text-gray-400">Backend</span>
                            <div className='flex items-center space-x-2'>
                                {
                                    systemState
                                        ? <span className="text-xs font-medium text-green-400">Activo</span>
                                        : <span className="text-xs font-medium text-red-400">Inactivo</span>
                                }
                                {
                                    systemState
                                        ? <div className="w-2 h-2 bg-green-400 rounded-full" />
                                        : <div className="w-2 h-2 bg-red-400 rounded-full" />
                                }
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </aside>
    )
}

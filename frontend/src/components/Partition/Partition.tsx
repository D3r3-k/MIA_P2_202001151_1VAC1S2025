"use client";

import { useMia } from '@/hooks/useMia';
import { HardDrive, Database, Folder, Activity } from 'lucide-react';
import React, { useEffect, useState } from 'react';

type PartitionProps = {
    name: string;
    driveletter: string;
    size: string;
    type: string;
    filesystem: string;
    mountPoint: string;
    status: string;
    createDate: string;
    id: string;
    signature: string;
    onSelect?: (partitionId: string) => void;
};

export function Partition({
    name,
    driveletter: disk,
    size,
    type,
    filesystem,
    mountPoint,
    status,
    createDate: lastCheck,
    id,
    signature,
    onSelect,
}: PartitionProps) {
    const { isAuthenticated, userData, activateToast } = useMia();
    const [mounted, setMounted] = useState(false);

    // Asegurar render en cliente
    useEffect(() => {
        setMounted(true);
    }, []);

    const handleClick = () => {
        if (!mounted) return;
        if (status !== "Montada") {
            activateToast("error", "No se pudo acceder", "La partición no está montada.");
            return;
        }
        if (!isAuthenticated) {
            activateToast("error", "Inicia sesión", "Debes iniciar sesión para acceder.");
            return;
        }
        if (isAuthenticated && userData?.partition_id !== id) {
            activateToast("error", "Sesión activa en otra partición", "Debes cerrar sesión antes.");
            return;
        }
        if (userData?.partition_id === id && onSelect) {
            onSelect(id);
        }
    };

    const getStatusConfig = (status: string) => {
        switch (status) {
            case "Montada":
                return { color: "text-green-400", bg: "bg-green-500/20", border: "border-green-500/30", text: "Montada" };
            case "Desmontada":
                return { color: "text-yellow-400", bg: "bg-yellow-500/20", border: "border-yellow-500/30", text: "No Montada" };
            case "error":
                return { color: "text-red-400", bg: "bg-red-500/20", border: "border-red-500/30", text: "Error" };
            default:
                return { color: "text-gray-400", bg: "bg-gray-500/20", border: "border-gray-500/30", text: "Desconocido" };
        }
    };

    const getTypeConfig = (type: string) => {
        switch (type.toLowerCase()) {
            case "primaria":
                return { icon: HardDrive, color: "text-blue-400", bg: "bg-blue-500/20", border: "border-blue-500/30" };
            case "extendida":
                return { icon: Database, color: "text-purple-400", bg: "bg-purple-500/20", border: "border-purple-500/30" };
            case "lógica":
                return { icon: Folder, color: "text-green-400", bg: "bg-green-500/20", border: "border-green-500/30" };
            case "swap":
                return { icon: Activity, color: "text-orange-400", bg: "bg-orange-500/20", border: "border-orange-500/30" };
            default:
                return { icon: Database, color: "text-gray-400", bg: "bg-gray-500/20", border: "border-gray-500/30" };
        }
    };

    const statusConfig = getStatusConfig(status);
    const typeConfig = getTypeConfig(type);
    const TypeIcon = typeConfig.icon;

    return (
        <div
            onClick={handleClick}
            className="bg-gray-800/30 border border-gray-700/50 rounded-lg p-6 hover:bg-gray-800/50 transition-all duration-200 cursor-pointer"
        >
            <div className="grid grid-cols-1 lg:grid-cols-8 gap-6 items-center">
                {/* Información principal */}
                <div className="lg:col-span-5">
                    <div className="flex items-center space-x-4">
                        <div className={`p-3 rounded-xl border ${typeConfig.bg} ${typeConfig.border}`}>
                            <TypeIcon className={`w-6 h-6 ${typeConfig.color}`} />
                        </div>
                        <div className="flex-1">
                            <h3 className="text-lg font-semibold text-white">{name}</h3>
                            <p className="text-sm text-gray-400">
                                {disk} • {size}
                            </p>
                            <div className="flex gap-2 mt-1">
                                <span className={`text-xs px-2 py-1 rounded ${typeConfig.bg} ${typeConfig.border} ${typeConfig.color}`}>
                                    {type}
                                </span>
                                <span className="text-xs px-2 py-1 rounded text-gray-400 bg-gray-700/30 border border-gray-600">
                                    {filesystem}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Estado */}
                <div className="lg:col-span-1 space-y-2">
                    <div className={`flex items-center px-3 py-2 rounded-lg border ${statusConfig.bg} ${statusConfig.border}`}>
                        <span className={`text-sm font-medium ${statusConfig.color}`}>{statusConfig.text}</span>
                    </div>
                </div>

                {/* Punto de montaje */}
                <div className="lg:col-span-2 text-center">
                    <p className="text-xs text-gray-400 mb-1">Punto de montaje</p>
                    <code className="text-sm text-white bg-gray-900/50 px-2 py-1 rounded">{mountPoint}</code>
                </div>
            </div>

            {/* Footer */}
            <div className="mt-4 pt-4 border-t border-gray-700/30 text-xs text-gray-500 flex justify-between">
                <span>Última verificación: {lastCheck}</span>
                <div className="flex gap-4">
                    <span>ID: {id}</span>
                    <span>Firma: {signature}</span>
                </div>
            </div>
        </div>
    );
}

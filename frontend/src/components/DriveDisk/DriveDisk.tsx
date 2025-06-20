"use client";

import { DriveDiskProps } from "@/types/GlobalTypes";
import { HardDrive } from "lucide-react";

export default function DriveDisk({ name = 'Disco Duro',
    size = '500 GB',
    type = 'HDD',
    status = 'active',
    partitions = 1
}: DriveDiskProps) {
    // Hooks
    // States
    // Effects
    // Handlers
    // Functions
    const getStatusColor = (status: string) => {
        switch (status) {
            case 'active': return 'text-green-600 bg-green-100';
            case 'inactive': return 'text-yellow-600 bg-yellow-100';
            case 'error': return 'text-red-600 bg-red-100';
            default: return 'text-gray-600 bg-gray-100';
        }
    };

    const getStatusText = (status: "active" | "inactive" | "error") => {
        switch (status) {
            case 'active': return 'Activo';
            case 'inactive': return 'Inactivo';
            case 'error': return 'Error';
        }
    };
    // Renders
    return (
        <div className="disk-item group flex flex-col items-center bg-white rounded-xl shadow-md p-6 transition-all duration-300 hover:shadow-lg w-64 border border-gray-200 hover:border-corinto-300 cursor-pointer">
            <div className="relative mb-4">
                <div className="w-16 h-16 bg-corinto-100 rounded-full flex items-center justify-center group-hover:bg-corinto-200 transition-colors duration-300">
                    <HardDrive className="w-8 h-8 text-corinto-800" />
                </div>
                <div
                    className={`absolute -top-1 -right-1 w-4 h-4 rounded-full border-2 border-white ${status === "active"
                        ? "bg-green-500"
                        : status === "inactive"
                            ? "bg-yellow-500"
                            : "bg-red-500"
                        }`}
                ></div>
            </div>

            <div className="text-center space-y-1 w-full">
                <h3 className="font-semibold text-gray-900 truncate">{name}</h3>
                <p className="text-sm text-gray-600">{size} • {type}</p>
                <div className="flex items-center justify-center space-x-2">
                    <span className={`text-xs px-2 py-1 rounded-full ${getStatusColor(status)}`}>
                        {getStatusText(status)}
                    </span>
                </div>
                <p className="text-xs text-gray-500">
                    {partitions} partición{partitions !== 1 ? "es" : ""}
                </p>
            </div>

            <div className="w-full pt-4 mt-4 border-t border-gray-100 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                <button className="w-full text-sm text-corinto-800 hover:text-corinto-900 font-medium">
                    Ver detalles
                </button>
            </div>
        </div>
    )
}

"use client";

import { HardDrive } from "lucide-react";
import Link from "next/link";

interface DriveDiskProps {
    name: string;
    size: string;
    path: string;
    fit: 'FF' | 'BF' | 'WF';
    partitions: number;
}

export default function DriveDisk({
    name = 'Disco Duro',
    size = '500 GB',
    path = name[name.length - 1] + ".dsk",
    fit = 'FF',
    partitions = 1
}: DriveDiskProps) {
    // Hooks
    // States
    // Effects
    // Handlers
    // Functions
    // Renders
    return (
        <Link
            href={`/drives/${path}`}
            className="bg-gray-800/50 border-gray-700 hover:bg-gray-800/70 hover:border-corinto-600/70 transition-all duration-200 hover:scale-105 group cursor-pointer border rounded-lg shadow-lg">
            <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                    <div className="relative">
                        <div className="w-14 h-14 bg-red-500/20 rounded-full flex items-center justify-center shadow">
                            <HardDrive className="w-7 h-7 text-red-400" />
                        </div>
                        <div className={`absolute top-0 right-0 w-4 h-4 rounded-full border-2 border-gray-800 bg-green-500`} />
                    </div>
                    <div className="text-right">
                        <p className="text-xs text-gray-400 uppercase tracking-wide font-semibold">
                            {path}
                        </p>
                    </div>
                </div>

                <div className="space-y-3">
                    <div>
                        <h3 className="font-semibold text-white text-lg group-hover:text-red-400 transition-colors">
                            {name}
                        </h3>
                        <p className="text-gray-400 text-sm">
                            {size} • {path}
                        </p>
                    </div>

                    <div className="flex justify-between items-center text-sm mt-2">
                        <span className="text-gray-400">Tipo de Fit:</span>
                        <span className={`font-medium capitalize text-green-400`}>
                            {
                                fit === 'FF'
                                    ? 'Primer Ajuste'
                                    : fit === 'BF'
                                        ? 'Mejor Ajuste'
                                        : fit === 'WF'
                                            ? 'Peor Ajuste'
                                            : 'Desconocido'
                            }
                        </span>
                    </div>

                    <div className="flex items-center justify-between pt-3 border-t border-gray-700 mt-3">
                        <div className="flex items-center space-x-2">
                            <div className={`w-2.5 h-2.5 rounded-full`} />
                            <span className="text-sm text-gray-400">Particiones</span>
                        </div>
                        <span className="text-sm text-white font-semibold">
                            {partitions}
                        </span>
                    </div>
                </div>
            </div>
        </Link>
    )
}

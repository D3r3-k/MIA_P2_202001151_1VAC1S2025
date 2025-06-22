"use client";

import { useMia } from "@/hooks/useMia";
import { HardDrive, Loader } from "lucide-react";

export default function Loading() {
    // Hooks
    const { loading } = useMia();
    // States
    // Effects
    // Handlers
    // Functions
    // Renders
    return (
        loading && (
            <div className="fixed inset-0 flex items-start justify-end z-50 m-5">
                <div className="relative flex flex-col items-center justify-center bg-gradient-to-r animate-pulse from-gray-500/80 to-gray-600/60 rounded-xl p-4 border border-gray-600/40 w-36">
                    <div className="relative">
                        <div className="absolute top-0 -left-2 transform translate-x-1/2 -translate-y-1/2">
                            <Loader size={24} className="text-red-500 animate-spin" />
                        </div>
                        <HardDrive size={32} className="text-gray-400" />
                    </div>
                    <h2 className={`text-base font-semibold text-gray-200 ${loading ? "animate-bounce" : ""}`}>
                        Cargando...
                    </h2>
                </div>
            </div>
        )
    )
}

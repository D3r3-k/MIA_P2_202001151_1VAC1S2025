"use client";

import { Check, Info, X, XCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { useMia } from "@/hooks/useMia";

export default function Toast() {
    const { toast, handleClose } = useMia();
    const [progress, setProgress] = useState(100);

    useEffect(() => {
        if (!toast.visible) return;

        setProgress(100);
        const start = Date.now();

        const interval = setInterval(() => {
            const elapsed = Date.now() - start;
            const newProgress = Math.max(100 - (elapsed / (toast.duration || 6000)) * 100, 0);
            setProgress(newProgress);
        }, 100);

        const timeout = setTimeout(() => {
            handleClose();
        }, toast.duration || 6000);

        return () => {
            clearInterval(interval);
            clearTimeout(timeout);
        };
    }, [toast]);

    return (
        <div
            className={`fixed bottom-3 right-3 p-5 border-t backdrop-blur-md w-80 md:w-96 rounded-t-lg shadow-lg border-gray-700/30 bg-gray-800/50 ${toast.visible ? "translate-x-0 opacity-100" : "translate-x-full opacity-0 pointer-events-none"
                } transition-all duration-300`}
        >
            <div className="space-y-3">
                <div className="relative">
                    <div className="flex items-start space-x-3">
                        <div className="mt-1">
                            {toast.type === "info" ? (
                                <Info size={24} className="text-blue-400" />
                            ) : toast.type === "success" ? (
                                <Check size={24} className="text-green-400" />
                            ) : (
                                <XCircle size={24} className="text-red-400" />
                            )}
                        </div>
                        <div>
                            <p className="text-md font-semibold text-white mb-1">{toast.message}</p>
                        </div>
                        <button
                            onClick={handleClose}
                            type="button"
                            aria-label="Cerrar"
                            className="absolute top-0 right-0 p-1 rounded-full hover:bg-gray-700 transition cursor-pointer"
                        >
                            <X size={24} className="text-gray-300" />
                        </button>
                    </div>
                </div>
                {toast.subtitle && (
                    <div className="flex items-center space-x-4">
                        <span className="text-sm text-gray-400">{toast.subtitle}</span>
                    </div>
                )}
            </div>
            <div className="my-4 h-1 w-full bg-gray-700/40 rounded-md overflow-hidden">
                <div
                    className="h-full bg-zinc-700 transition-all duration-100"
                    style={{ width: `${progress}%` }}
                />
            </div>
        </div>
    );
}

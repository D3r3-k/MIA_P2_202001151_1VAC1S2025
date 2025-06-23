"use client";

import { useState } from "react";
import LoginView from "./LoginView";
import DrivesView from "./DrivesView";
import DriveDetailView from "./DriveDetailView";
import PartitionView from "./PartitionView";
import ConsoleView from "./ConsoleView";
import { Route } from "@/types/GlobalTypes";


export default function AppView({
    route,
    setRoute,
}: {
    route: Route;
    setRoute: (route: Route) => void;
}) {
    const [activeDrive, setActiveDrive] = useState<string | null>(null);
    const [activePartition, setActivePartition] = useState<string | null>(null);

    switch (route) {
        case "/":
            return <ConsoleView />;

        case "login":
            return <LoginView onLoginSuccess={() => setRoute("drives")} />;

        case "drives":
            return (
                <DrivesView
                    onSelectDrive={(driveLetter) => {
                        setActiveDrive(driveLetter);
                        setRoute("drive-detail");
                    }}
                />
            );

        case "drive-detail":
            if (!activeDrive) return null;
            return (
                <DriveDetailView
                    driveLetter={activeDrive}
                    onSelectPartition={(partitionId) => {
                        setActivePartition(partitionId);
                        setRoute("partition");
                    }}
                    onBack={() => {
                        setActiveDrive(null);
                        setRoute("drives");
                    }}
                />
            );

        case "partition":
            if (!activeDrive || !activePartition) return null;
            return (
                <PartitionView
                    driveLetter={activeDrive}
                    partitionId={activePartition}
                    onBack={() => {
                        setActivePartition(null);
                        setRoute("drive-detail");
                    }}
                />
            );

        default:
            return <div className="p-6 text-white">Ruta desconocida.</div>;
    }
}

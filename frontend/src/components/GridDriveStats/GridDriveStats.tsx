"use client";

import { Database, HardDriveIcon } from "lucide-react";
import { StatsCards } from "../StatsCards/StatsCards";
import { useEffect, useState } from "react";
import { DriveDiskStatus } from "@/types/GlobalTypes";
import SkeletonStatsCard from "../StatsCards/SkeletonStatsCard";


export default function GridDriveStats() {
    const [driveData, setDriveData] = useState<DriveDiskStatus | null>(null);

    useEffect(() => {
        const fetchDriveData = async () => {
            await new Promise((res) => setTimeout(res, 1500));
            const data = {
                totalDisks: 3,
                totalPartitions: 6,
                totalSize: "3TB"
            };
            setDriveData(data);
        };

        fetchDriveData();
    }, []);

    if (!driveData) {
        return (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <SkeletonStatsCard />
                <SkeletonStatsCard />
                <SkeletonStatsCard />
            </div>
        );
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 group">
            <StatsCards
                title="Total de Discos"
                value={driveData.totalDisks}
                color={{
                    color: "text-blue-500",
                    bgColor: "from-blue-500/10 to-blue-500/30",
                    borderColor: "border-blue-500",
                    accentColor: "bg-blue-500/10"
                }}
                icon={HardDriveIcon}
            />
            <StatsCards
                title="Total de Particiones"
                value={driveData.totalPartitions}
                color={{
                    color: "text-green-500",
                    bgColor: "from-green-500/10 to-green-500/30",
                    borderColor: "border-green-500",
                    accentColor: "bg-green-500/10"
                }}
                icon={Database}
            />
            <StatsCards
                title="Tamaño Total"
                value={driveData.totalSize}
                color={{
                    color: "text-red-500",
                    bgColor: "from-red-500/10 to-red-500/30",
                    borderColor: "border-red-500",
                    accentColor: "bg-red-500/10"
                }}
                icon={HardDriveIcon}
            />
        </div>
    );
}

"use client";

import { Database, HardDriveIcon } from "lucide-react";
import { DriveStats } from "../DriveStats/DriveStats";
import { DriveDiskStatusType } from "@/types/GlobalTypes";
import useFetchs from "@/hooks/useFetchs";
import { useEffect, useState } from "react";

export default function GridDrivesStats() {
    const { getDriveStats } = useFetchs();

    const [data, setData] = useState<DriveDiskStatusType>({
        totalDisks: 0,
        totalPartitions: 0,
        totalSize: "0 B",
    });

    useEffect(() => {
        const fetchData = async () => {
            const result = await getDriveStats();
            setData(result);
        };
        fetchData();
    }, []);

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 group">
            <DriveStats
                title="Total de Discos"
                value={data.totalDisks}
                color={{
                    color: "text-blue-500",
                    bgColor: "from-blue-500/10 to-blue-500/30",
                    borderColor: "border-blue-500",
                    accentColor: "bg-blue-500/10"
                }}
                icon={HardDriveIcon}
            />
            <DriveStats
                title="Particiones montadas"
                value={data.totalPartitions}
                color={{
                    color: "text-green-500",
                    bgColor: "from-green-500/10 to-green-500/30",
                    borderColor: "border-green-500",
                    accentColor: "bg-green-500/10"
                }}
                icon={Database}
            />
            <DriveStats
                title="Tamaño Total"
                value={data.totalSize}
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

"use client";

import { Calculator, HardDriveDownload, HardDriveIcon, Route } from "lucide-react";
import { DriveStats } from "../DriveStats/DriveStats";
import { DriveDiskInfoType } from "@/types/GlobalTypes";
import useFetchs from "@/hooks/useFetchs";
import { useEffect, useState } from "react";

interface GridDriveStatsProps {
    driveLetter: string;
}

export default function GridDriveStats({ driveLetter }: GridDriveStatsProps) {
    const { getDriveInfo } = useFetchs();
    const [data, setData] = useState<DriveDiskInfoType>({
        Name: driveLetter,
        Path: "N/A",
        Size: "0 B",
        Fit: "N/A",
        Partitions: 0,
    });

    useEffect(() => {
        const fetchData = async () => {
            const info = await getDriveInfo(driveLetter);
            setData(info);
        };
        fetchData();
    }, [driveLetter]);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 group">
            <DriveStats
                title="Disco"
                value={data.Name}
                color={{
                    color: "text-blue-500",
                    bgColor: "from-blue-500/10 to-blue-500/30",
                    borderColor: "border-blue-500",
                    accentColor: "bg-blue-500/10"
                }}
                direction="vertical"
                icon={HardDriveIcon}
            />
            <DriveStats
                title="Ruta del Disco"
                value={data.Path}
                color={{
                    color: "text-cyan-500",
                    bgColor: "from-cyan-500/10 to-cyan-500/30",
                    borderColor: "border-cyan-500",
                    accentColor: "bg-cyan-500/10"
                }}
                direction="vertical"
                icon={Route}
            />
            <DriveStats
                title="Tamaño del Disco"
                value={data.Size}
                color={{
                    color: "text-purple-500",
                    bgColor: "from-purple-500/10 to-purple-500/30",
                    borderColor: "border-purple-500",
                    accentColor: "bg-purple-500/10"
                }}
                direction="vertical"
                icon={Calculator}
            />
            <DriveStats
                title="Fit"
                value={data.Fit}
                color={{
                    color: "text-green-500",
                    bgColor: "from-green-500/10 to-green-500/30",
                    borderColor: "border-green-500",
                    accentColor: "bg-green-500/10"
                }}
                direction="vertical"
                icon={HardDriveDownload}
            />
            <DriveStats
                title="Particiones Montadas"
                value={data.Partitions}
                color={{
                    color: "text-red-500",
                    bgColor: "from-red-500/10 to-red-500/30",
                    borderColor: "border-red-500",
                    accentColor: "bg-red-500/10"
                }}
                direction="vertical"
                icon={HardDriveIcon}
            />
        </div>
    );
}

"use client";

import { useEffect, useState } from "react";
import { DriveDiskType } from "@/types/GlobalTypes";
import DriveDisk from "../DriveDisk/DriveDisk";
import useFetchs from "@/hooks/useFetchs";

export default function GridDrives() {
    const { getDrives } = useFetchs();
    const [drives, setDrives] = useState<DriveDiskType[]>([]);

    useEffect(() => {
        const fetchDrives = async () => {
            const data = await getDrives();
            setDrives(data);
        };
        fetchDrives();
    }, []);

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {drives.length === 0 && (
                <div className="col-span-4 text-center text-gray-400">
                    No hay discos disponibles.
                </div>
            )}
            {drives.map((drive, index) => (
                <DriveDisk
                    key={index}
                    name={`Disco ${drive.Name}`}
                    partitions={drive.Partitions}
                    size={drive.Size}
                    fit={drive.Fit}
                    path={drive.Path}
                />
            ))}
        </div>
    );
}

"use client";

import { useEffect, useState } from "react";
import DriveDisk from "../DriveDisk/DriveDisk";
import SkeletonDriveDisk from "../DriveDisk/SkeletonDriveDisk";
import { DriveDiskType } from "@/types/GlobalTypes";

export default function GridDrives() {
    const [driveData, setDriveData] = useState<DriveDiskType[] | null>(null);

    useEffect(() => {
        const fetchDriveData = async () => {
            try {
                const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/drives`);
                if (!response.ok) {
                    throw new Error("Failed to fetch drive data");
                }
                const data = await response.json();
                setDriveData(data);
            } catch (error) {
                console.error("Error fetching drive data:", error);
            }
        };

        fetchDriveData();
    }, []);

    if (!driveData) {
        return (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <SkeletonDriveDisk />
                <SkeletonDriveDisk />
                <SkeletonDriveDisk />
            </div>
        );
    }
    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {driveData.map((drive, index) => (
                <DriveDisk key={index} name={`Disco ${drive.Name}`} partitions={drive.Partitions} size={drive.Size} fit={drive.Fit} path={drive.Path} />
            ))}
        </div>
    )
}

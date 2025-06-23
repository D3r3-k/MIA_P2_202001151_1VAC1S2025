"use client";

import DriveDisk from "../DriveDisk/DriveDisk";
import { DriveDiskType } from "@/types/GlobalTypes";

export default function GridDrives({ onSelectDrive, drive }: { onSelectDrive: (driveLetter: string) => void, drive: DriveDiskType }) {

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <DriveDisk
                name={`Disco ${drive.Name}`}
                partitions={drive.Partitions}
                size={drive.Size}
                fit={drive.Fit}
                path={drive.Path}
                onClick={() => onSelectDrive(drive.Name)}
            />
        </div>
    );
}

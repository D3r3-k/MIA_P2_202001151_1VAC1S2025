import { useEffect, useState, Suspense } from "react";
import Head from "next/head";
import SkeletonDriveDisk from "@/components/DriveDisk/SkeletonDriveDisk";
import SkeletonStatsCard from "@/components/DriveStats/SkeletonStatsCard";
import GridDrivesStats from "@/components/Grids/GridDrivesStats";
import { DriveDiskInfoType } from "@/types/GlobalTypes";
import DriveDisk from "../DriveDisk/DriveDisk";
import useFetchs from "@/hooks/useFetchs";

export default function DrivesView({
    onSelectDrive,
}: {
    onSelectDrive: (driveLetter: string) => void;
}) {
    const { getDrives } = useFetchs();
    const [drives, setDrives] = useState<DriveDiskInfoType[]>([]);

    useEffect(() => {
        const fetchDrives = async () => {
            const data = await getDrives();
            setDrives(data || []);
        };

        fetchDrives();
    }, []);

    const handleDriveClick = (driveLetter: string) => {
        onSelectDrive(driveLetter);
    };

    return (
        <>
            <Head>
                <title>Gestion de Discos - F2 MIA</title>
                <meta name="description" content="Explora y gestiona los discos y particiones del sistema." />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <main className="flex-1 p-6 ml-72">
                <div className="mb-8 grid grid-cols-2 gap-6">
                    <div className="flex flex-col justify-center">
                        <h1 className="text-3xl font-bold text-white mb-2">
                            Gestión de Discos
                        </h1>
                        <p className="text-gray-400">
                            Administra y monitorea todas las particiones del sistema
                        </p>
                    </div>
                    <Suspense
                        fallback={
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 group">
                                <SkeletonStatsCard />
                                <SkeletonStatsCard />
                                <SkeletonStatsCard />
                            </div>
                        }
                    >
                        <GridDrivesStats />
                    </Suspense>
                </div>
                <div className="space-y-6">
                    <h2 className="text-xl font-semibold text-white mb-6">Discos del Sistema</h2>

                    <Suspense
                        fallback={
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                                <SkeletonDriveDisk />
                                <SkeletonDriveDisk />
                                <SkeletonDriveDisk />
                            </div>
                        }
                    >
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
                                    onClick={() => handleDriveClick(drive.Name)}
                                />
                            ))}
                        </div>
                    </Suspense>
                </div>
            </main>
        </>
    );
}

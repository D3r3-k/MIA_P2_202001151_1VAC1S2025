import SkeletonStatsCard from "@/components/DriveStats/SkeletonStatsCard";
import GridDriveStats from "@/components/Grids/GridDriveStats";
import GridPartitions from "@/components/Grids/GridPartitions";
import { PartitionSkeleton } from "@/components/Partition/PartitionSkeleton";
import { Suspense } from "react";

interface DriveLetterPageProps {
    params: {
        driveletter: string;
    };
}
export default async function DriveLetterPage({ params }: DriveLetterPageProps) {
    const { driveletter } = await params;
    // Hooks
    // States
    // Effects
    // Handlers
    // Functions
    // Renders
    return (
        <main className="flex-1 p-6 ml-72">
            <div className="mb-8 grid grid-cols-2 gap-6">
                <div className="flex flex-col justify-center">
                    <h1 className="text-3xl font-bold text-white mb-2">
                        Gestión del Disco {driveletter.toUpperCase()}
                    </h1>
                    <p className="text-gray-400">
                        Administra y monitorea las particiones del disco {driveletter.toUpperCase()} del sistema
                    </p>
                </div>
                <Suspense fallback={
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 group">
                        <SkeletonStatsCard />
                        <SkeletonStatsCard />
                        <SkeletonStatsCard />
                    </div>
                }>
                    <GridDriveStats driveLetter={driveletter} />
                </Suspense>
            </div>
            <div className="space-y-6">
                <h2 className="text-xl font-semibold text-white mb-6">Particiones del Disco</h2>
                <Suspense fallback={
                    <div className="space-y-6">
                        <PartitionSkeleton />
                        <PartitionSkeleton />
                        <PartitionSkeleton />
                        <PartitionSkeleton />
                    </div>
                }>
                    <GridPartitions driveletter={driveletter} />
                </Suspense>
            </div>
        </main>
    );
}
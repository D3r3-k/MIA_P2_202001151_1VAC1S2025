import { Suspense } from "react";
import Head from "next/head";
import SkeletonStatsCard from "@/components/DriveStats/SkeletonStatsCard";
import { PartitionSkeleton } from "@/components/Partition/PartitionSkeleton";
import GridDriveStats from "@/components/Grids/GridDriveStats";
import GridPartitions from "../Grids/GridPartitions";
import { ChevronLeft } from "lucide-react";

export default function DriveDetailView({
    driveLetter,
    onSelectPartition,
    onBack,
}: {
    driveLetter: string;
    onSelectPartition: (partitionId: string) => void;
    onBack: () => void;
}) {
    return (
        <>
            <Head>
                <title>{`Disco ${driveLetter} - F2 MIA`}</title>
                <meta name="description" content={`Particiones del disco ${driveLetter}`} />
            </Head>

            <main className="flex-1 p-6 ml-72">
                <button
                    onClick={onBack}
                    className="mb-6 flex items-center gap-2 px-4 py-2 rounded-lg bg-gray-700 text-gray-200 hover:bg-gray-600 hover:text-white transition-colors shadow cursor-pointer"
                >
                    <ChevronLeft size={20} />
                    Regresar
                </button>
                {/* Header */}
                <div className="mb-8 grid grid-cols-2 gap-6">
                    <div className="flex flex-col justify-center">
                        <h1 className="text-3xl font-bold text-white mb-2">
                            Gestión del Disco {driveLetter.toUpperCase()}
                        </h1>
                        <p className="text-gray-400">
                            Administra y monitorea las particiones del disco {driveLetter.toUpperCase()}
                        </p>
                    </div>

                    <Suspense
                        fallback={
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 group">
                                <SkeletonStatsCard />
                                <SkeletonStatsCard />
                            </div>
                        }
                    >
                        <GridDriveStats driveLetter={driveLetter} />
                    </Suspense>
                </div>

                {/* Particiones */}
                <div className="space-y-6">
                    <h2 className="text-xl font-semibold text-white mb-6">Particiones del Disco</h2>

                    <Suspense
                        fallback={
                            <>
                                <PartitionSkeleton />
                                <PartitionSkeleton />
                            </>
                        }
                    >
                        <GridPartitions
                            driveletter={driveLetter}
                            onSelectPartition={onSelectPartition}
                        />
                    </Suspense>
                </div>
            </main>
        </>
    );
}

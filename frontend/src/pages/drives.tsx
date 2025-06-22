import SkeletonDriveDisk from '@/components/DriveDisk/SkeletonDriveDisk'
import SkeletonStatsCard from '@/components/DriveStats/SkeletonStatsCard'
import GridDrives from '@/components/Grids/GridDrives'
import GridDrivesStats from '@/components/Grids/GridDrivesStats'
import Head from 'next/head'
import React, { Suspense } from 'react'

export default function Drives() {
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
                        <GridDrives />
                    </Suspense>
                </div>
            </main>
        </>
    )
}

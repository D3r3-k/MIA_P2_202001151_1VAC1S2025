import ContentPartition from '@/components/ContentPartition/ContentPartition';
import { DriveDiskInfoType } from '@/types/GlobalTypes';
import Head from 'next/head';
import React from 'react';

export async function getStaticPaths() {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/drives`);
    if (!res.ok) {
        throw new Error('Failed to fetch drives data');
    }
    const drivesData: DriveDiskInfoType[] = await res.json();
    const paths: { params: { driveletter: string; partition_id: string } }[] = [];

    for (const drive of drivesData) {
        const driveLetter = drive?.Name;
        if (!driveLetter) continue;

        const partitionsRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/drives/${driveLetter}/partitions`);
        if (!partitionsRes.ok) {
            throw new Error(`Failed to fetch partitions for drive ${driveLetter}`);
        }
        const partitionsData = await partitionsRes.json();

        for (const partition of partitionsData) {
            const partitionId = partition?.ID;
            if (typeof partitionId === 'string' && partitionId.trim() !== '') {
                paths.push({
                    params: {
                        driveletter: driveLetter,
                        partition_id: partitionId,
                    },
                });
            }
        }
    }

    return {
        paths,
        fallback: false,
    };
}



export async function getStaticProps({ params }: { params: { driveletter: string, partition_id: string } }) {
    return {
        props: {
            driveletter: params.driveletter,
            partition_id: params.partition_id,
        },
    };
}

export default function ParticionId({ driveletter, partition_id }: { driveletter: string; partition_id: string }) {
    return (
        <>
            <Head>
                <title>{`Partición ${partition_id} - Disco ${driveletter} - F2 MIA`}</title>
                <meta name="description" content={`Explora el contenido de la partición ${partition_id} del disco ${driveletter}`} />
            </Head>
            <main className="flex-1 p-6 ml-72">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-white mb-2">Explorador de Archivos</h1>
                    <p className="text-gray-400">
                        Navega y administra el sistema de archivos del disco {driveletter}, partición {partition_id}.
                    </p>
                </div>
                <div className="space-y-6">
                    <ContentPartition />
                </div>
            </main>
        </>
    );
}

import Head from "next/head";
import ContentPartition from "@/components/ContentPartition/ContentPartition";
import { ChevronLeft } from "lucide-react";

export default function PartitionView({
    driveLetter,
    partitionId,
    onBack,
}: {
    driveLetter: string;
    partitionId: string;
    onBack: () => void;
}) {
    return (
        <>
            <Head>
                <title>{`Partición ${partitionId} - Disco ${driveLetter} - F2 MIA`}</title>
                <meta
                    name="description"
                    content={`Explora el contenido de la partición ${partitionId} del disco ${driveLetter}`}
                />
            </Head>
            <main className="flex-1 p-6 ml-72">
                <button
                    onClick={onBack}
                    className="mb-6 flex items-center gap-2 px-4 py-2 rounded-lg bg-gray-700 text-gray-200 hover:bg-gray-600 hover:text-white transition-colors shadow cursor-pointer"
                >
                    <ChevronLeft size={20} />
                    Regresar
                </button>

                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-white mb-2">Explorador de Archivos</h1>
                    <p className="text-gray-400">
                        Navega y administra el sistema de archivos del disco {driveLetter}, partición {partitionId}.
                    </p>
                </div>

                <div className="space-y-6">
                    <ContentPartition />
                </div>
            </main>
        </>
    );
}

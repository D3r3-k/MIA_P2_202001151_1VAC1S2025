import { Partition } from "@/components/Partition/Partition";

interface DriveLetterPageProps {
    params: {
        driveletter: string;
    };
}
export default function DriveLetterPage({ params: { driveletter } }: DriveLetterPageProps) {
    return (
        <main className="flex-1 p-6 ml-72">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-white mb-2">
                    Gestión de Particiones del Disco {driveletter.toUpperCase()}
                </h1>
                <p className="text-gray-400">
                    Visualiza y administra todos los discos y particiones del sistema
                </p>
            </div>
            <div className="space-y-6">
                <Partition partition={{
                    name: "Partición 1",
                    disk: driveletter,
                    size: "100 GB",
                    type: "Primaria",
                    filesystem: "NTFS",
                    mountPoint: "/mnt/part1",
                    usedSpace: 50,
                    freeSpace: "50 GB",
                    status: "mounted",
                    lastCheck: "2023-10-01",
                    id: "part1"
                }} />
            </div>
        </main>
    );
}
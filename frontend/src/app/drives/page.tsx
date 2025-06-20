import DriveDisk from "@/components/DriveDisk/DriveDisk";
import GridDriveStats from "@/components/GridDriveStats/GridDriveStats";

export default function DrivesPage() {
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
                        Gestión de Discos
                    </h1>
                    <p className="text-gray-400">
                        Visualiza y administra todos los discos y particiones del sistema
                    </p>
                </div>
                <GridDriveStats />
            </div>
            <div className="space-y-6">
                <h2 className="text-xl font-semibold text-white mb-6">Discos del Sistema</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <DriveDisk name="Disco A" partitions={1} size="100" fit="FF" path="A.dsk" />
                    <DriveDisk name="Disco B" partitions={3} size="100" fit="BF" path="B.dsk" />
                    <DriveDisk name="Disco C" partitions={2} size="100" fit="WF" path="C.dsk" />
                </div>
            </div>
        </main>
    );
}

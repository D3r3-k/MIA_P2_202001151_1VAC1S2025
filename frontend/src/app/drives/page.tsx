import GridDrives from "@/components/Grids/GridDrives";
import GridDriveStats from "@/components/Grids/GridDriveStats";

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
                        Administra y monitorea todas las particiones del sistema
                    </p>
                </div>
                <GridDriveStats />
            </div>
            <div className="space-y-6">
                <h2 className="text-xl font-semibold text-white mb-6">Discos del Sistema</h2>
                <GridDrives />
            </div>
        </main>
    );
}

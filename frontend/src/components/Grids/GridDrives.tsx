import { DriveDiskType } from "@/types/GlobalTypes";
import DriveDisk from "../DriveDisk/DriveDisk";
import useFetchs from "@/hooks/useFetchs";

export default async function GridDrives() {
    // Hooks
    const { getDrives } = useFetchs();
    // States
    // Effects
    const data: DriveDiskType[] = await getDrives();
    // Functions
    // Renders
    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {
                data.length === 0 && (
                    <div className="col-span-4 text-center text-gray-400">
                        No hay discos disponibles.
                    </div>
                )
            }
            {data.map((drive, index) => (
                <DriveDisk key={index} name={`Disco ${drive.Name}`} partitions={drive.Partitions} size={drive.Size} fit={drive.Fit} path={drive.Path} />
            ))}
        </div>
    )
}

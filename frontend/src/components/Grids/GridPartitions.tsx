import useFetchs from "@/hooks/useFetchs";
import { Partition } from "../Partition/Partition";

interface GridPartitionsProps {
    driveletter: string;
}

export default async function GridPartitions({ driveletter }: GridPartitionsProps) {
    // Hooks
    const { getPartitions } = useFetchs();
    // States
    const data = await getPartitions(driveletter);
    // Effects
    // Handlers
    // Functions
    // Renders
    return (
        <div className="space-y-6">
            {data.length > 0 ? (
                data.map((partition, index) => (
                    <Partition
                        key={index}
                        name={partition.Name !== "" ? partition.Name : `Libre`}
                        driveletter={driveletter}
                        size={partition.Size}
                        type={partition.Type}
                        filesystem={partition.Filesystem}
                        mountPoint={partition.Path}
                        status={partition.Status}
                        createDate={partition.Date}
                        id={partition.ID}
                        signature={partition.Signature}
                    />
                ))
            ) : (
                <div className="text-gray-500 text-center">
                    No hay particiones disponibles en este disco.
                </div>
            )}
        </div>
    )
}

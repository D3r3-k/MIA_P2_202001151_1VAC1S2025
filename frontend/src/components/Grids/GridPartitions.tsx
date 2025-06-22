"use client";

import { useEffect, useState } from "react";
import useFetchs from "@/hooks/useFetchs";
import { Partition } from "../Partition/Partition";
import { DrivePartitionType } from "@/types/GlobalTypes";

interface GridPartitionsProps {
    driveletter: string;
}

export default function GridPartitions({ driveletter }: GridPartitionsProps) {
    const { getPartitions } = useFetchs();
    const [partitions, setPartitions] = useState<DrivePartitionType[]>([]);

    useEffect(() => {
        const fetchPartitions = async () => {
            const data = await getPartitions(driveletter);
            setPartitions(data);
        };
        fetchPartitions();
    }, [driveletter]);

    return (
        <div className="space-y-6">
            {partitions.length > 0 ? (
                partitions.map((partition, index) => (
                    <Partition
                        key={index}
                        name={partition.Name !== "" ? partition.Name : "Libre"}
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
    );
}

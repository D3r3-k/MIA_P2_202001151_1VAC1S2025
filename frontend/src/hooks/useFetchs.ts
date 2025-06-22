"use client";

import {
  DriveDiskInfoType,
  DriveDiskStatusType,
  DriveDiskType,
  DrivePartitionType,
  FileSystemItemType,
} from "@/types/GlobalTypes";
import { useMia } from "./useMia";

const useFetchs = () => {
  const { activateToast } = useMia();
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;

  const fetchJson = async <T>(
    url: string,
    options?: RequestInit,
    defaultValue?: T,
    errorMessage?: string
  ): Promise<T> => {
    try {
      const response = await fetch(url, options);
      if (!response.ok)
        throw new Error(errorMessage || "Error en la solicitud");
      const data = await response.json();
      return data as T;
    } catch (error) {
      const err = error as Error;
      activateToast(
        "error",
        "Error en la solicitud",
        errorMessage || err.message
      );
      return defaultValue as T;
    }
  };

  // [Drives]
  const getDrivesStats = async (): Promise<DriveDiskStatusType> =>
    fetchJson<DriveDiskStatusType>(
      `${baseUrl}/api/drives/info`,
      undefined,
      {
        totalDisks: 0,
        totalPartitions: 0,
        totalSize: "0 B",
      },
      "No se pudo obtener el resumen de los discos."
    );

  const getDrives = async (): Promise<DriveDiskType[]> =>
    fetchJson<DriveDiskType[]>(
      `${baseUrl}/api/drives`,
      undefined,
      [],
      "No se pudo obtener la lista de discos."
    );

  // [Drive]
  const getDriveInfo = async (
    driveLetter: string
  ): Promise<DriveDiskInfoType> =>
    fetchJson<DriveDiskInfoType>(
      `${baseUrl}/api/drives/${driveLetter}`,
      undefined,
      {
        Name: driveLetter.toUpperCase(),
        Path: "N/A",
        Size: "0 B",
        Fit: "N/A",
        Partitions: 0,
      },
      `No se pudo obtener la información del disco ${driveLetter.toUpperCase()}.`
    );

  const getPartitions = async (
    driveLetter: string
  ): Promise<DrivePartitionType[]> =>
    fetchJson<DrivePartitionType[]>(
      `${baseUrl}/api/drives/${driveLetter}/partitions`,
      undefined,
      [],
      `No se pudieron obtener las particiones del disco ${driveLetter.toUpperCase()}.`
    );

  const getFileSystemItems = async (
    path: string
  ): Promise<FileSystemItemType | []> =>
    fetchJson<FileSystemItemType | []>(
      `${baseUrl}/api/find`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ path }),
      },
      [],
      `No se pudo obtener el contenido de la ruta "${path}".`
    );

  const getContentFile = async (path: string): Promise<string> =>
    fetchJson<string>(
      `${baseUrl}/api/cat`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ path }),
      },
      "",
      `No se pudo leer el archivo "${path}".`
    );

  return {
    getDrives,
    getDriveStats: getDrivesStats,
    getDriveInfo,
    getPartitions,
    getFileSystemItems,
    getContentFile,
  };
};

export default useFetchs;

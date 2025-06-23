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
  const { activateToast, setLoading } = useMia();
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;

  const fetchJson = async <T>(
    url: string,
    options?: RequestInit,
    defaultValue?: T,
    fallbackMessage?: string
  ): Promise<T> => {
    try {
      setLoading(true);
      const response = await fetch(url, options);
      const data = await response.json();

      if (!response.ok) {
        const backendError =
          data?.error || fallbackMessage || response.statusText;
        activateToast("error", "Error en la solicitud", backendError);
        return data as T;
      }

      return data as T;
    } catch (error) {
      const err = error as Error;
      activateToast(
        "error",
        "Error en la solicitud",
        err.message || fallbackMessage || "Ocurrió un error inesperado"
      );
      return defaultValue as T;
    } finally {
      setLoading(false);
    }
  };

  // [Execute]
  const executeCommand = async (commands: string): Promise<string> => {
    const body = JSON.stringify({ commands });

    const response = await fetchJson<{
      response: string;
      status: string;
      error?: string;
    }>(
      `${baseUrl}/api/execute`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body,
      },
      { response: "", status: "" },
      `No se pudo ejecutar el/los comando(s) "${commands}".`
    );

    return response.response;
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
    executeCommand,
    getDrives,
    getDriveStats: getDrivesStats,
    getDriveInfo,
    getPartitions,
    getFileSystemItems,
    getContentFile,
  };
};

export default useFetchs;

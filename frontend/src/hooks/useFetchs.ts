// custom hook to fetch data from an API

import {
  DriveDiskInfoType,
  DriveDiskStatusType,
  DriveDiskType,
  DrivePartitionType,
} from "@/types/GlobalTypes";

const useFetchs = () => {
  // [Drives]
  const getDrivesStats = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/drives/info`
      );
      if (!response.ok) {
        throw new Error("Failed to fetch drive data");
      }
      const data: DriveDiskStatusType = await response.json();
      return data;
    } catch (error) {
      return {
        totalDisks: 0,
        totalPartitions: 0,
        totalSize: "0 B",
      };
    }
  };
  const getDrives = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/drives`
      );
      if (!response.ok) {
        throw new Error("Failed to fetch drive data");
      }
      const data: DriveDiskType[] = await response.json();
      return data || [];
    } catch (error) {
      return [];
    }
  };
  // [Drive]
  const getDriveInfo = async (driveLetter: string) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/drives/${driveLetter}`
      );
      if (!response.ok) {
        throw new Error("Failed to fetch drive info");
      }
      const data: DriveDiskInfoType = await response.json();
      return data;
    } catch (error) {
      return {
        Name: driveLetter.toUpperCase(),
        Path: "N/A",
        Size: "0 B",
        Fit: "N/A",
        Partitions: 0,
      };
    }
  };
  const getPartitions = async (driveLetter: string) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/drives/${driveLetter}/partitions`
      );
      if (!response.ok) {
        throw new Error("Failed to fetch partitions data");
      }
      const data: DrivePartitionType[] = await response.json();
      return data || [];
    } catch (error) {
      return [];
    }
  };

  return {
    getDrives,
    getDriveStats: getDrivesStats,
    getDriveInfo,
    getPartitions,
  };
};

export default useFetchs;

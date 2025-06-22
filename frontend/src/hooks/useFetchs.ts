// custom hook to fetch data from an API

import {
  DriveDiskInfoType,
  DriveDiskStatusType,
  DriveDiskType,
  DrivePartitionType,
  FileSystemItemType,
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
      console.log("Error fetching drive stats:", error);
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
      console.log("Error fetching drives:", error);
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
      console.log("Error fetching drive info:", error);
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
      console.log("Error fetching partitions:", error);
      return [];
    }
  };

  const getFileSystemItems = async (
    path: string
  ): Promise<FileSystemItemType | []> => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/find`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ path }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to fetch partition find");
      }
      const data: FileSystemItemType = await response.json();
      return data;
    } catch (error) {
      console.log("Error fetching partition find:", error);
      return [];
    }
  };

  const getContentFile = async (path: string) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/cat`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ path }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to fetch file content");
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.log("Error fetching file content:", error);
      return "";
    }
  };

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

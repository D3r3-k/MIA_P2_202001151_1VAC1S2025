export type DriveDiskStatusType = {
  totalDisks: number;
  totalPartitions: number;
  totalSize: string;
};

export type DriveDiskType = {
  Name: string;
  Partitions: number;
  Size: string;
  Fit: string;
  Path: string;
};

export type DriveDiskInfoType = {
  Name: string;
  Path: string;
  Size: string;
  Fit: string;
  Partitions: number;
};

export type DrivePartitionType = {
  Status: string;
  Type: string;
  Fit: string;
  Size: string;
  Start: string;
  Name: string;
  ID: string;
  Path: string;
  Date: string;
  Filesystem: string;
  Signature: string;
};

export type FileSystemItemType = {
  ID: string;
  Name: string;
  Type: string;
  Path: string;
  Children: FileSystemItemType[];
  Size?: string;
  CreatedAt?: string;
  Owner?: string;
  Content?: string;
  Extension?: string;
  Permissions?: string;
};


export type Route =
    | "login"
    | "drives"
    | "drive-detail"
    | "partition"
    | "/";
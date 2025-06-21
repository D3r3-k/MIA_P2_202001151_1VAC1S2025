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

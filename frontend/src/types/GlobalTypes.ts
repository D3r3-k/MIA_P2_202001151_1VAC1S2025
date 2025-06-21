export type DriveDiskStatus = {
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
}
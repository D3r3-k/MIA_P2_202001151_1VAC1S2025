export interface DriveDiskProps {
  name: string;
  size: string;
  type: string;
  status: "active" | "inactive" | "error";
  partitions: number;
}

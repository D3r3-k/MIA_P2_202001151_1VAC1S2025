type LoginParams = {
  partition_id: string;
  username: string;
  password: string;
};

export interface UserData {
  username: string;
  group: string;
  partition_id: string;
  permissions: string;
}

export interface MiaContextType {
  executeCommand: (command: string) => Promise<string>;
  isAuthenticated: boolean;
  login: ({ partition_id, username, password }: LoginParams) => Promise<void>;
  logout: () => Promise<void>;
  userData: UserData | null;
  loading: boolean;
  errorMsg: string | null;
}

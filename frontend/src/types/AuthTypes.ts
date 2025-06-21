type LoginParams = {
  partition_id: string;
  username: string;
  password: string;
};

export interface UserData {
  username: string;
  group: string;
  partition_id: string;
}

export interface MiaContextType {
  systemState: boolean;
  executeCommand: (command: string) => Promise<string>;
  isAuthenticated: boolean;
  login: ({ partition_id, username, password }: LoginParams) => Promise<boolean>;
  logout: () => Promise<void>;
  userData: UserData | null;
  loading: boolean;
  errorMsg: string | null;
}

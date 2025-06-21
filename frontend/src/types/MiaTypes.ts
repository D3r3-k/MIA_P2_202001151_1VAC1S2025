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
  loading: boolean;
  systemState: boolean;
  isAuthenticated: boolean;
  userData: UserData | null;
  executeCommand: (command: string) => Promise<string>;
  login: ({ partition_id, username, password }: LoginParams) => Promise<boolean>;
  logout: () => Promise<void>;
  errorMsg: string | null;
}

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
  setLoading: (loading: boolean) => void;
  systemState: boolean;
  isAuthenticated: boolean;
  userData: UserData | null;
  executeCommand: (command: string) => Promise<string>;
  login: ({ partition_id, username, password }: LoginParams) => Promise<boolean>;
  logout: () => Promise<void>;
  activateToast: (
    type: "info" | "success" | "error",
    message: string,
    subtitle?: string,
    duration?: number
  ) => void;
  toast: SingleToastData;
  handleClose: () => void;
}


export interface SingleToastData {
    message: string;
    subtitle?: string;
    type?: "info" | "success" | "error";
    duration?: number;
    visible: boolean;
}

export interface ToastContextType {
    activate: (
        type: "info" | "success" | "error",
        message: string,
        subtitle?: string,
        duration?: number
    ) => void;
    toast: SingleToastData;
    handleClose: () => void;
}

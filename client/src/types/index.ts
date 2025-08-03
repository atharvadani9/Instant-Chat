export type UserAuthRequest = {
  username: string;
  password: string;
};

export type Message = {
    id: number;
    sender_id: number;
    receiver_id: number;
    content: string;
    created_at: string;
}

export type User = {
    id: number;
    username: string;
    created_at: string;
}

export type WSMessage = {
    type: string;
    sender_id?: number;
    receiver_id?: number;
    content?: string;
    error?: string;
    created_at?: string;
}

export type ConnectionState = "connecting" | "connected" | "disconnected";
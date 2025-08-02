import {useCallback, useEffect, useState} from "react";
import {useAuthContext} from "../contexts/AuthContext";
import type {WSMessage} from "../types";

export const useWebSocket = () => {
    const {userID} = useAuthContext();
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [messages, setMessages] = useState<WSMessage[]>([]);
    const [connectionState, setConnectionState] = useState<"connecting" | "connected" | "disconnected">("disconnected");
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!userID) {
            return;
        }
        setConnectionState("connecting");

        const newSocket = new WebSocket(`ws://localhost:8080/chat/ws?user_id=${userID}`);
        newSocket.onopen = () => {
            setConnectionState("connected");
            setSocket(newSocket);
        };
        newSocket.onmessage = (event) => {
            const message = JSON.parse(event.data) as WSMessage;

            if (message.type === "error") {
                setError(message.error || "Server error occurred");
                return;
            }

            setMessages((prevMessages) => [...prevMessages, message]);
            setError(null);
        };
        newSocket.onclose = () => {
            setConnectionState("disconnected");
            setSocket(null);
        };
        newSocket.onerror = () => {
            setError("WebSocket error occurred");
            setConnectionState("disconnected");
            setSocket(null);
        };

        return () => {
            newSocket.close();
        };
    }, [userID]);

    const sendMessage = useCallback((receiverID: number, content: string) => {
        if (!socket) {
            return;
        }

        const message: WSMessage = {
            type: "send_message",
            receiver_id: receiverID,
            content,
        };
        socket.send(JSON.stringify(message));
    }, [socket]);

    const getMessages = useCallback((receiverID: number) => {
        if (!socket) {
            return;
        }

        const message: WSMessage = {
            type: "get_history",
            receiver_id: receiverID,
        };
        socket.send(JSON.stringify(message));
    }, [socket]);

    const clearMessages = useCallback(() => {
        setMessages([]);
    }, []);

    return {
        sendMessage,
        getMessages,
        messages,
        connectionState,
        clearMessages,
        error,
    };
}
import {useCallback, useEffect, useRef, useState} from "react";
import {useAuthContext} from "../contexts/AuthContext";
import type {WSMessage, ConnectionState} from "../types";

export const useWebSocket = () => {
    const {userID} = useAuthContext();
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [messages, setMessages] = useState<WSMessage[]>([]);
    const [connectionState, setConnectionState] = useState<ConnectionState>("disconnected");
    const [error, setError] = useState<string | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const socketRef = useRef<WebSocket | null>(null);

    const connectWebSocket = useCallback(() => {
        if (!userID) {
            return;
        }

        // Clear any existing reconnect timeout
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }

        // Close existing socket if any
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.close();
        }

        setConnectionState("connecting");
        setError(null);

        const newSocket = new WebSocket(`ws://localhost:8080/chat/ws?user_id=${userID}`);
        socketRef.current = newSocket;

        newSocket.onopen = () => {
            console.log("WebSocket connected");
            setConnectionState("connected");
            setSocket(newSocket);
            setError(null);
        };

        newSocket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                console.log("Received message:", data);

                if (data.type === "error") {
                    setError(data.error || "Server error occurred");
                    return;
                }

                if (data.type === "messages_history") {
                    // Handle message history - replace current messages
                    const historyMessages = data.messages || [];
                    const formattedMessages = historyMessages.map((msg: any) => ({
                        type: "message_history",
                        sender_id: msg.sender_id,
                        receiver_id: msg.receiver_id,
                        content: msg.content,
                        created_at: msg.created_at
                    }));
                    setMessages(formattedMessages);
                } else if (data.type === "new_message") {
                    // Handle new real-time message (both received and sent messages)
                    setMessages((prevMessages) => [...prevMessages, data as WSMessage]);
                }

                setError(null);
            } catch (err) {
                console.error("Error parsing WebSocket message:", err);
                setError("Failed to parse message");
            }
        };

        newSocket.onclose = (event) => {
            console.log("WebSocket closed:", event.code, event.reason);
            setConnectionState("disconnected");
            setSocket(null);
            socketRef.current = null;

            // Only attempt reconnection if it wasn't a normal closure and user is still logged in
            if (event.code !== 1000 && userID) {
                console.log("Attempting to reconnect in 3 seconds...");
                reconnectTimeoutRef.current = setTimeout(() => {
                    connectWebSocket();
                }, 3000);
            }
        };

        newSocket.onerror = (event) => {
            console.error("WebSocket error:", event);
            setError("WebSocket connection error");
            setConnectionState("disconnected");
            setSocket(null);
        };
    }, [userID]);

    useEffect(() => {
        connectWebSocket();

        return () => {
            // Clear reconnect timeout
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
                reconnectTimeoutRef.current = null;
            }

            // Close socket
            if (socketRef.current) {
                socketRef.current.close(1000, "Component unmounting");
                socketRef.current = null;
            }
        };
    }, [connectWebSocket]);

    const sendMessage = useCallback((receiverID: number, content: string) => {
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            setError("WebSocket is not connected");
            return;
        }

        const message: WSMessage = {
            type: "send_message",
            receiver_id: receiverID,
            content,
        };

        try {
            socket.send(JSON.stringify(message));
            console.log("Message sent:", message);
        } catch (err) {
            console.error("Error sending message:", err);
            setError("Failed to send message");
        }
    }, [socket]);

    const getMessages = useCallback((receiverID: number) => {
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            console.log("WebSocket not ready for getting messages");
            return;
        }

        const message: WSMessage = {
            type: "get_history",
            receiver_id: receiverID,
        };

        try {
            socket.send(JSON.stringify(message));
            console.log("Requesting message history for user:", receiverID);
        } catch (err) {
            console.error("Error requesting messages:", err);
            setError("Failed to get messages");
        }
    }, [socket]);

    const clearMessages = useCallback(() => {
        setMessages([]);
    }, []);

    const reconnect = useCallback(() => {
        connectWebSocket();
    }, [connectWebSocket]);

    return {
        sendMessage,
        getMessages,
        messages,
        connectionState,
        clearMessages,
        error,
        reconnect,
    };
}
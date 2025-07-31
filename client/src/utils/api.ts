import axios from "axios";

const api = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});


export const postAPI = async (url: string, payload: any) => {
    try {
        const response = await api.post(url, payload);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            return error.response.data;
        } else {
            throw new Error("Server error occurred");
        }
    }
}

export const getAPI = async (url: string) => {
    try {
        const response = await api.get(url);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            return error.response.data;
        } else {
            throw new Error("Server error occurred");
        }
    }
}

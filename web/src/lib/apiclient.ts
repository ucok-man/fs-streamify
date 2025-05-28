import axios from "axios";

export const apiclient = axios.create({
  baseURL: `${import.meta.env.VITE_BASE_API_URL}/api/v1`,
  withCredentials: true,
});

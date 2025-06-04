import axios from "axios";

export const apiclient = axios.create({
  baseURL: `/api/v1`,
  withCredentials: true,
});

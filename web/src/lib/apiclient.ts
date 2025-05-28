import axios from "axios";

export const apiclient = axios.create({
  baseURL: `${process.env.BASE_API_URL}/api`,
  withCredentials: true,
});

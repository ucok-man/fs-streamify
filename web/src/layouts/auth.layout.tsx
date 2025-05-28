import { useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { ShipWheel } from "lucide-react";
import toast from "react-hot-toast";
import { Navigate, Outlet } from "react-router";
import ThreeDotLoader from "../components/three-dot-loader";
import { apiclient } from "../lib/apiclient";
import type { UserResponse } from "../types/user-response.type";

export default function AuthLayout() {
  const { data, isPending, isError, error } = useQuery({
    queryKey: ["auth:session"],
    queryFn: async () => {
      const { data } = await apiclient.get<{ user: UserResponse }>("/auth/me");
      return data;
    },
    retry: false, // prevent retry on 401 errors
  });

  // Show a loading state while the session is being fetched
  if (isPending) {
    return (
      <div
        data-theme="forest"
        className="flex h-screen w-full items-center justify-center rounded"
      >
        <div className="space-y-4 text-center">
          <ThreeDotLoader size="lg" />
          <div className="flex items-center gap-4">
            <ShipWheel className="size-9 text-primary" />
            <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text font-mono text-3xl font-bold tracking-wider text-transparent">
              Streamify
            </span>
          </div>
        </div>
      </div>
    );
  }

  if (isError) {
    const axiosError = error as AxiosError;
    if (axiosError.response?.status && axiosError.response.status >= 500) {
      toast.error("Sorry we have a problem on our server. Try again later!");
      console.error("Server error during session check:", axiosError);
    }
    return <Navigate to="/signin" replace />;
  }

  // Not authenticated if no user found
  if (!data?.user) {
    return <Navigate to="/signin" replace />;
  }

  return <Outlet />;
}

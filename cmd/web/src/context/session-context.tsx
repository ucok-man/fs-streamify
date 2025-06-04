import { useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { createContext, useContext, type PropsWithChildren } from "react";
import { apiclient } from "../lib/apiclient";
import type { UserResponse } from "../types/user-response.type";

export type SessionContextType = {
  data?: UserResponse;
  error?: AxiosError;
  isError: boolean;
  isPending: boolean;
};

const SessionContext = createContext<SessionContextType>({
  data: undefined,
  error: undefined,
  isError: false,
  isPending: true,
});

export default function SessionContextProvider({
  children,
}: PropsWithChildren) {
  const { data, isPending, isError, error } = useQuery({
    queryKey: ["auth:session"],
    queryFn: async () => {
      const { data } = await apiclient.get<{ user: UserResponse }>("/auth/me");
      return data;
    },
    retry: false, // prevent retry on 401 errors
  });

  return (
    <SessionContext.Provider
      value={{
        data: data?.user,
        error: (error as AxiosError) || undefined,
        isError: isError,
        isPending: isPending,
      }}
    >
      {children}
    </SessionContext.Provider>
  );
}

export function useSessionContext() {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error(
      "useSessionContext must be used within an SessionContextProvider"
    );
  }
  return context;
}

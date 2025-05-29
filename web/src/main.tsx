import { QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "react-hot-toast";
import App from "./App.tsx";
import SessionContextProvider from "./context/session-context.tsx";
import "./index.css";
import { queryclient } from "./lib/query-client.ts";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryclient}>
      <SessionContextProvider>
        <App />
        <Toaster />
      </SessionContextProvider>
    </QueryClientProvider>
  </StrictMode>
);

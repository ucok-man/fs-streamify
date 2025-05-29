import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import type { SessionContextType } from "../context/session-context";

type RouterContext = {
  session: Pick<SessionContextType, "data" | "error" | "isError">;
};

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
});

function RootLayout() {
  return (
    <main>
      <Outlet />
      <TanStackRouterDevtools />
    </main>
  );
}

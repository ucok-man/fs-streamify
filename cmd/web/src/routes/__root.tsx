import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import NotFound from "../components/not-found";
import type { SessionContextType } from "../context/session-context";

type RouterContext = {
  session: Pick<SessionContextType, "data" | "error" | "isError">;
};

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
  notFoundComponent: NotFound,
});

function RootLayout() {
  return (
    <main data-theme="forest">
      <Outlet />
      <TanStackRouterDevtools />
    </main>
  );
}

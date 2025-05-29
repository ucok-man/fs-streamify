import {
  createRootRouteWithContext,
  Link,
  Outlet,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
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

function NotFound() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-base-100 px-6">
      <div className="text-center">
        <h1 className="text-9xl font-bold text-primary">404</h1>
        <p className="mt-4 text-2xl font-semibold text-base-content md:text-3xl">
          Page Not Found
        </p>
        <p className="mt-2 text-base-content/70">
          Looks like you wandered off the path and into the woods 🌲
        </p>

        <div className="mt-6">
          <Link to="/" className="btn btn-primary">
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}

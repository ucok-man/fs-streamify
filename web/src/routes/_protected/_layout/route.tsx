import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_layout")({
  loader: ({ location }) => {
    return {
      showSidebar: !location.pathname.startsWith("/chat"),
    };
  },
  component: Layout,
});

function Layout() {
  const { showSidebar } = Route.useLoaderData();

  return (
    <div className="min-h-screen">
      <div className="flex">
        {showSidebar && <div>[SIDEBAR]</div>}
        <div className="flex flex-1 flex-col">
          [NAVBAR]
          <div className="flex-1 overflow-y-auto">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}

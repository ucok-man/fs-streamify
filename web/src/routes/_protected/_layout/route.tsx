import { createFileRoute, Outlet } from "@tanstack/react-router";
import Navbar from "../../../components/navbar";
import Sidebar from "../../../components/sidebar";

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
        {showSidebar && <Sidebar />}
        <div className="flex flex-1 flex-col">
          <Navbar />
          <div className="flex-1 overflow-y-auto">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}

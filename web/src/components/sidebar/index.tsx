import { Link, useRouteContext } from "@tanstack/react-router";
import { BellIcon, HomeIcon, ShipWheelIcon, UsersIcon } from "lucide-react";

const Sidebar = () => {
  const { session } = useRouteContext({
    from: "/_protected",
  });

  return (
    <aside className="sticky top-0 hidden h-screen w-64 flex-col border-r border-base-300 bg-base-200 lg:flex">
      <div className="border-b border-base-300 p-5">
        <Link to="/" className="flex items-center gap-2.5">
          <ShipWheelIcon className="size-9 text-primary" />
          <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text font-mono text-3xl font-bold tracking-wider  text-transparent">
            Streamify
          </span>
        </Link>
      </div>

      <nav className="flex-1 space-y-3 p-4">
        <Link
          to="/"
          className={`btn w-full justify-start gap-3 px-3 normal-case btn-ghost`}
          activeProps={{
            className: "btn-active bg-secondary text-accent-content",
          }}
        >
          <HomeIcon className="size-5 opacity-70" />
          <span className="text-lg">Home</span>
        </Link>

        <Link
          to="/friend"
          className={`btn w-full justify-start gap-3 px-3 normal-case btn-ghost`}
          activeProps={{
            className: "btn-active bg-secondary text-accent-content",
          }}
        >
          <UsersIcon className="size-5 opacity-70" />
          <span className="text-lg">Friends</span>
        </Link>

        <Link
          to="/notification"
          className={`btn w-full justify-start gap-3 px-3 normal-case btn-ghost`}
          activeProps={{
            className: "btn-active bg-secondary text-accent-content",
          }}
        >
          <BellIcon className="size-5 opacity-70" />
          <span className="text-lg">Notifications</span>
        </Link>
      </nav>

      {/* USER PROFILE SECTION */}
      <div className="mt-auto border-t border-base-300 p-4">
        <div className="flex items-center gap-3">
          <div className="avatar">
            <div className="w-10 rounded-full">
              <img src={session.data?.profile_pic} alt="User Avatar" />
            </div>
          </div>
          <div className="flex-1">
            <p className="text-sm font-semibold">{session.data?.full_name}</p>
            <p className="flex items-center gap-1 text-xs text-success">
              <span className="inline-block size-2 rounded-full bg-success" />
              Online
            </p>
          </div>
        </div>
      </div>
    </aside>
  );
};
export default Sidebar;

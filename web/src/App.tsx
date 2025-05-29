import { RouterProvider, createRouter } from "@tanstack/react-router";
import toast from "react-hot-toast";
import PageLoader from "./components/page-loader";
import { useSessionContext } from "./context/session-context";
import "./index.css";
import { routeTree } from "./routeTree.gen";

const router = createRouter({
  routeTree,
  context: {
    session: undefined!, // docs.
  },
});
// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

export default function App() {
  const session = useSessionContext();

  if (session.isPending) return <PageLoader />;

  if (!session.isPending && session.error) {
    if (!session.error.status || session.error.status >= 500) {
      toast.error("Sorry we have a problem on our server. Try again later!");
      console.error("Server error during session check:", session?.error);
    }

    return <PageLoader />;
  }

  return (
    <RouterProvider
      router={router}
      context={{
        session: {
          data: session.data,
          error: session.error,
          isError: session.isError,
        },
      }}
    />
  );
}

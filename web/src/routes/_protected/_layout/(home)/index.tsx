import { createFileRoute } from "@tanstack/react-router";
import Recommended from "../../../../components/pages/homes/recommended";

export const Route = createFileRoute("/_protected/_layout/(home)/")({
  validateSearch: (search: Record<string, unknown>) => {
    const query = search?.query;
    if (!query) return { query: undefined };
    if (Array.isArray(query)) return { query: query[0] as string };
    return { query: query as string };
  },
  component: HomePage,
});

function HomePage() {
  return (
    <div className="p-4 sm:p-6 lg:p-8">
      <div className="container mx-auto space-y-8">
        {/* <Friends /> */}
        <Recommended />
      </div>
    </div>
  );
}

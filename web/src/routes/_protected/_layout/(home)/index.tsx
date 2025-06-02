import { createFileRoute } from "@tanstack/react-router";
import Friends from "../../../../components/pages/homes/friends";
import Recommended from "../../../../components/pages/homes/recommended";

export const Route = createFileRoute("/_protected/_layout/(home)/")({
  component: HomePage,
});

function HomePage() {
  return (
    <div className="p-4 sm:p-6 lg:p-8">
      <div className="container mx-auto space-y-8">
        <Friends />
        <Recommended />
      </div>
    </div>
  );
}

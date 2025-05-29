import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_layout/")({
  component: HomePage,
});

function HomePage() {
  return <div>HomePage</div>;
}

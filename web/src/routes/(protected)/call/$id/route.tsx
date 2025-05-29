import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/call/$id")({
  component: CallPage,
});

function CallPage() {
  const { id } = Route.useParams();
  return <div>CallPage #ID: {id}</div>;
}

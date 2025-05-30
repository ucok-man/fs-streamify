import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_layout/chat/$id")({
  component: ChatPage,
});

function ChatPage() {
  const { id } = Route.useParams();
  return <div>ChatPage #ID: {id}</div>;
}

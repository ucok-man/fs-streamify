import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_layout/friend")({
  component: FriendPage,
});

function FriendPage() {
  return <div>Hello "/_protected/_layout/friend"!</div>;
}

import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_layout/notification")({
  component: NotificationPage,
});

function NotificationPage() {
  return <div>NotificationPage</div>;
}

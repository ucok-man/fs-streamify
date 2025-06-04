import { createFileRoute } from "@tanstack/react-router";
import IncomingFriendRequest from "./-incoming-friend-request";
import OutgoingFriendRequest from "./-outgoing-friend-request";

export const Route = createFileRoute("/_protected/_layout/notification")({
  component: NotificationPage,
});

function NotificationPage() {
  return (
    <div className="p-4 sm:p-6 lg:p-8">
      <div className="container mx-auto space-y-8">
        <h1 className="mb-6 text-2xl font-bold tracking-tight sm:text-3xl">
          Notifications
        </h1>

        <div className="space-y-8">
          <IncomingFriendRequest />
          <OutgoingFriendRequest />
        </div>
      </div>
    </div>
  );
}

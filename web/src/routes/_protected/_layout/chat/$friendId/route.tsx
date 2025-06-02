import { createFileRoute, notFound } from "@tanstack/react-router";
import { AxiosError } from "axios";
import NotFound from "../../../../../components/not-found";
import { apiclient } from "../../../../../lib/apiclient";
import type { UserResponse } from "../../../../../types/user-response.type";

export const Route = createFileRoute("/_protected/_layout/chat/$friendId")({
  loader: async ({ params }) => {
    const user = await apiclient
      .get<{ user: UserResponse }>(`/users/${params.friendId}`)
      .then((res) => res.data.user)
      .catch((error) => {
        if (
          error instanceof AxiosError &&
          (error.response?.status === 404 || error.response?.status === 400)
        )
          return null;
        throw error;
      });

    if (!user) throw notFound();

    const token = await apiclient
      .get<{ token: string }>(`/chat/token`)
      .then((res) => res.data.token)
      .catch((error) => {
        throw error;
      });

    return { token };
  },
  component: ChatPage,
  notFoundComponent: NotFound,
});

function ChatPage() {
  const { token } = Route.useLoaderData();
  return <div>ChatPage Token: {token}</div>;
}

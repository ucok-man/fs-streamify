import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected")({
  beforeLoad: ({ context }) => {
    const user = context.session.data;

    if (!user) {
      throw redirect({
        to: "/signin",
      });
    }
    if (user && !user.is_onboarded) {
      throw redirect({
        to: "/onboarding",
      });
    }
  },
});

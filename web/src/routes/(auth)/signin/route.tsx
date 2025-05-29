import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/(auth)/signin")({
  beforeLoad: ({ context, location }) => {
    const user = context.session.data;
    if (user && user.is_onboarded) {
      throw redirect({
        to: "/",
        search: {
          redirect: location.href,
        },
      });
    }

    if (user && !user.is_onboarded) {
      throw redirect({
        to: "/onboarding",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  component: SigninPage,
});

function SigninPage() {
  return <div>SigninPage</div>;
}

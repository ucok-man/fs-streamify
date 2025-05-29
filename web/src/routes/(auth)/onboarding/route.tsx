import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/(auth)/onboarding")({
  beforeLoad: ({ context, location }) => {
    const user = context.session.data;

    if (!user) {
      throw redirect({
        to: "/signin",
        search: {
          redirect: location.href,
        },
      });
    }

    if (user && user.is_onboarded) {
      throw redirect({
        to: "/",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  component: OnboardingPage,
});

function OnboardingPage() {
  return <div>OnboardingPage</div>;
}

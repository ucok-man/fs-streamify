import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import type { AxiosError } from "axios";
import { ShipWheel } from "lucide-react";
import { useTransition } from "react";
import { useForm } from "react-hook-form";
import toast from "react-hot-toast";
import SpinnerBtn from "../../../components/spinner-btn";
import { apiclient } from "../../../lib/apiclient";
import { refetchQuery } from "../../../lib/query-client";
import { parseApiError } from "../../../lib/utils";
import type { UserResponse } from "../../../types/user-response.type";
import { signupSchema, type SignupData } from "./-signup.schema";

export const Route = createFileRoute("/_auth/signup")({
  beforeLoad: ({ context }) => {
    const user = context.session.data;

    if (user && user.is_onboarded) {
      throw redirect({
        to: "/",
        search: { query: undefined },
      });
    }

    if (user && !user.is_onboarded) {
      throw redirect({
        to: "/onboarding",
      });
    }
  },
  component: SignupPage,
});

function SignupPage() {
  const navigate = Route.useNavigate();
  const [isRedirecting, startTransition] = useTransition();

  const form = useForm({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      fullname: "",
      email: "",
      password: "",
    },
  });

  const signup = useMutation({
    mutationFn: async (payload: SignupData) => {
      const { data } = await apiclient.post<{ user: UserResponse }>(
        "/auth/signup",
        payload
      );
      return data;
    },

    onSuccess: () => {
      form.reset();
      refetchQuery(["auth:session"]);
      startTransition(() => {
        navigate({
          to: "/onboarding",
        });
      });
    },
    onError: (err: AxiosError) => {
      if (err.status === 422) {
        parseApiError(err.response?.data, form);
        return;
      }
      toast.error(
        "Sorry our server encountered problem. Please try again later!"
      );
    },
  });

  const formerror = form.formState.errors;

  return (
    <div
      className="flex h-screen items-center justify-center p-4 sm:p-6 md:p-8"
      data-theme="forest"
    >
      <div className="mx-auto flex w-full max-w-5xl flex-col overflow-hidden rounded-xl border border-primary/25 bg-base-100 shadow-lg lg:flex-row">
        {/* Signup Form LEFT */}
        <div className="flex w-full flex-col p-4 sm:p-8 lg:w-1/2">
          {/* Logo Streamify */}
          <div className="mb-4 flex items-center justify-start gap-2">
            <ShipWheel className="size-9 text-primary" />
            <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text font-mono text-3xl font-bold tracking-wider text-transparent">
              Streamify
            </span>
          </div>

          {/* Form */}
          <div className="w-full">
            <form
              onSubmit={form.handleSubmit((payload) => signup.mutate(payload))}
            >
              <div className="space-y-4">
                <div>
                  <h2 className="text-xl font-semibold">Create an Account</h2>
                  <p className="text-sm opacity-70">
                    Join Streamify and start your language learning adventure!
                  </p>
                </div>

                <div className="space-y-4">
                  <div className="form-control w-full space-y-2">
                    <label htmlFor="fullname" className="label">
                      <span className="label-text px-1">Full Name</span>
                    </label>

                    <input
                      id="fullname"
                      type="text"
                      placeholder="John Doe"
                      className="input-bordered input w-full px-4"
                      {...form.register("fullname")}
                    />
                    <div className="px-1 text-xs text-red-500">
                      {formerror.fullname?.message}
                    </div>
                  </div>

                  <div className="form-control w-full space-y-2">
                    <label htmlFor="email" className="label">
                      <span className="label-text px-1">Email</span>
                    </label>

                    <input
                      id="email"
                      type="text"
                      placeholder="john@example.com"
                      className="input-bordered input w-full px-4"
                      {...form.register("email")}
                    />
                    <div className="px-1 text-xs text-red-500">
                      {formerror.email?.message}
                    </div>
                  </div>

                  <div className="form-control w-full space-y-2">
                    <label htmlFor="password" className="label">
                      <span className="label-text px-1">Password</span>
                    </label>

                    <input
                      id="password"
                      type="password"
                      placeholder="********"
                      className="input-bordered input w-full px-4"
                      {...form.register("password")}
                    />
                    <div className="px-1 text-xs text-red-500">
                      {formerror.password?.message}
                    </div>
                  </div>

                  <div className="form-control px-1 py-2">
                    <label
                      htmlFor="termService"
                      className="label cursor-pointer justify-start gap-2"
                    >
                      <input
                        type="checkbox"
                        className="checkbox checkbox-sm"
                        required
                      />
                      <span className="text-xs leading-tight">
                        I agree to the{" "}
                        <span className="text-primary hover:underline">
                          terms of service
                        </span>{" "}
                        and{" "}
                        <span className="text-primary hover:underline">
                          privacy policy
                        </span>
                      </span>
                    </label>
                  </div>
                </div>

                <button type="submit" className="btn w-full btn-primary">
                  {signup.isPending ? (
                    <SpinnerBtn msg="Signing up..." />
                  ) : isRedirecting ? (
                    <SpinnerBtn msg="Redirecting..." />
                  ) : (
                    "Create Account"
                  )}
                </button>

                <div className="mt-4 text-center">
                  <p className="text-sm">
                    Already have an account?{" "}
                    <Link
                      to={"/signin"}
                      className="text-primary underline-offset-4 hover:underline"
                    >
                      Sign in
                    </Link>
                  </p>
                </div>
              </div>
            </form>
          </div>
        </div>

        {/* Signup Form RIGHT */}
        <div className="hidden w-full items-center justify-center bg-primary/10 lg:flex lg:w-1/2">
          <div className="max-w-md p-8">
            <div className="relative mx-auto aspect-square max-w-sm">
              <img src="/hero.png" alt="Hero Image" className="h-full w-full" />
            </div>

            <div className="mt-6 space-y-3 text-center">
              <h2 className="text-xl font-semibold">
                Connect with language partners worldwide
              </h2>
              <p className="opacity-70">
                Practice conversations, make friends, and improve your language
                skills together
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

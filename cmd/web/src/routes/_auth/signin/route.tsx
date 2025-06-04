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
import { signinSchema, type SigninData } from "./-signin-schema";

export const Route = createFileRoute("/_auth/signin")({
  beforeLoad: ({ context }) => {
    const user = context.session.data;
    if (user && user.is_onboarded) {
      throw redirect({
        to: "/",
        search: {
          query: undefined,
        },
      });
    }

    if (user && !user.is_onboarded) {
      throw redirect({
        to: "/onboarding",
      });
    }
  },
  component: SigninPage,
});

function SigninPage() {
  const navigate = Route.useNavigate();
  const [isRedirecting, startTransition] = useTransition();

  const form = useForm({
    resolver: zodResolver(signinSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const signin = useMutation({
    mutationFn: async (payload: SigninData) => {
      const { data } = await apiclient.post<{ user: UserResponse }>(
        "/auth/signin",
        payload
      );
      return data;
    },

    onSuccess: () => {
      form.reset();
      refetchQuery(["auth:session"]);
      startTransition(() => {
        navigate({
          to: "/",
          search: {
            query: undefined,
          },
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
              onSubmit={form.handleSubmit((payload) => signin.mutate(payload))}
            >
              <div className="space-y-4">
                <div>
                  <h2 className="text-xl font-semibold">Welcome Back</h2>
                  <p className="text-sm opacity-70">
                    Sign in to your account to continue your language journey
                  </p>
                </div>

                <div className="space-y-4">
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
                </div>

                <button
                  disabled={signin.isPending || isRedirecting}
                  type="submit"
                  className="btn w-full btn-primary"
                >
                  {signin.isPending ? (
                    <SpinnerBtn msg="Signing up..." />
                  ) : isRedirecting ? (
                    <SpinnerBtn msg="Redirecting..." />
                  ) : (
                    "Sign In"
                  )}
                </button>

                <div className="mt-4 text-center">
                  <p className="text-sm">
                    Don't have an account?{" "}
                    <Link
                      to={"/signup"}
                      className="text-primary underline-offset-4 hover:underline"
                    >
                      Sign up
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

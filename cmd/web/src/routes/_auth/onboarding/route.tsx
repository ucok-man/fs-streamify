import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";
import type { AxiosError } from "axios";
import { CameraIcon, MapPin, ShipWheelIcon, ShuffleIcon } from "lucide-react";
import { useTransition } from "react";
import { useForm } from "react-hook-form";
import toast from "react-hot-toast";
import SpinnerBtn from "../../../components/spinner-btn";
import { LANGUAGES } from "../../../constants";
import { apiclient } from "../../../lib/apiclient";
import { refetchQuery } from "../../../lib/query-client";
import { generateAvatar, parseApiError } from "../../../lib/utils";
import type { UserResponse } from "../../../types/user-response.type";
import { onboardingSchema, type OnboardingData } from "./-onboarding.schema";

export const Route = createFileRoute("/_auth/onboarding")({
  beforeLoad: ({ context }) => {
    const user = context.session.data;

    if (!user) {
      throw redirect({
        to: "/signin",
      });
    }

    if (user && user.is_onboarded) {
      throw redirect({
        to: "/",
        search: { query: undefined },
      });
    }
  },
  component: OnboardingPage,
});

function OnboardingPage() {
  const { session } = Route.useRouteContext();
  const navigate = Route.useNavigate();
  const [isRedirecting, startTransition] = useTransition();

  const form = useForm({
    resolver: zodResolver(onboardingSchema),
    defaultValues: {
      fullname: session.data?.full_name || "",
      bio: "",
      learning_lng: "",
      native_lng: "",
      location: "",
      profile_pic: session.data?.profile_pic || "",
    },
  });

  const onboarding = useMutation({
    mutationFn: async (payload: OnboardingData) => {
      const { data } = await apiclient.post<{ user: UserResponse }>(
        "/auth/onboarding",
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
          reloadDocument: true,
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

  const handleRandomAvatar = () => {
    form.setValue("profile_pic", generateAvatar());
    toast.success("Random profile picture generated!");
  };

  form.watch("profile_pic");
  const formerror = form.formState.errors;

  return (
    <div className="flex min-h-screen items-center justify-center bg-base-100 p-4">
      <div className="card w-full max-w-3xl bg-base-200 shadow-xl">
        <div className="card-body p-6 sm:p-8">
          <h1 className="mb-6 text-center text-2xl font-bold sm:text-3xl">
            Complete Your Profile
          </h1>

          <form
            onSubmit={form.handleSubmit((payload) =>
              onboarding.mutate(payload)
            )}
            className="space-y-6"
          >
            {/* PROFILE PIC CONTAINER */}
            <div className="flex flex-col items-center justify-center space-y-4">
              {/* IMAGE PREVIEW */}
              <div className="size-32 overflow-hidden rounded-full bg-base-300">
                {form.getValues("profile_pic") ? (
                  <img
                    src={form.getValues("profile_pic")}
                    alt="Profile Preview"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full items-center justify-center">
                    <CameraIcon className="size-12 text-base-content opacity-40" />
                  </div>
                )}
              </div>

              {/* Generate Random Avatar BTN */}
              <div className="flex items-center gap-2">
                <button
                  type="button"
                  onClick={handleRandomAvatar}
                  className="btn btn-accent"
                >
                  <ShuffleIcon className="mr-2 size-4" />
                  Generate Random Avatar
                </button>
              </div>
            </div>

            {/* FULL NAME */}
            <div className="form-control w-full space-y-2">
              <label htmlFor="fullname" className="label">
                <span className="label-text px-1">Full Name</span>
              </label>

              <input
                id="fullname"
                type="text"
                placeholder="Your full name"
                className="input-bordered input w-full px-4"
                {...form.register("fullname")}
              />
              <div className="px-1 text-xs text-red-500">
                {formerror.fullname?.message}
              </div>
            </div>

            {/* BIO */}
            <div className="form-control w-full space-y-2">
              <label htmlFor="bio" className="label">
                <span className="label-text px-1">Bio</span>
              </label>
              <textarea
                id="bio"
                className="textarea-bordered textarea h-24 w-full resize-none px-4 py-4"
                placeholder="Tell others about yourself and your language learning goals"
                {...form.register("bio")}
              />
              <div className="px-1 text-xs text-red-500">
                {formerror.bio?.message}
              </div>
            </div>

            {/* LANGUAGES */}
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              {/* NATIVE LANGUAGE */}
              <div className="form-control w-full space-y-2">
                <label htmlFor="native_lng" className="label">
                  <span className="label-text px-1">Native Language</span>
                </label>
                <select
                  id="native_lng"
                  className="select-bordered select w-full"
                  {...form.register("native_lng")}
                >
                  <option disabled value="">
                    Select your native language
                  </option>
                  {LANGUAGES.map((lang) => (
                    <option key={`native-${lang}`} value={lang.toLowerCase()}>
                      {lang}
                    </option>
                  ))}
                </select>
                <div className="px-1 text-xs text-red-500">
                  {formerror.native_lng?.message}
                </div>
              </div>

              {/* LEARNING LANGUAGE */}
              <div className="form-control w-full space-y-2">
                <label htmlFor="learning_lng" className="label">
                  <span className="label-text px-1">Learning Language</span>
                </label>
                <select
                  id="learning_lng"
                  className="select-bordered select w-full"
                  {...form.register("learning_lng")}
                >
                  <option value="" disabled>
                    Select language you're learning
                  </option>
                  {LANGUAGES.map((lang) => (
                    <option key={`learning-${lang}`} value={lang.toLowerCase()}>
                      {lang}
                    </option>
                  ))}
                </select>
                <div className="px-1 text-xs text-red-500">
                  {formerror.learning_lng?.message}
                </div>
              </div>
            </div>

            {/* LOCATION */}
            <div className="form-control w-full space-y-2">
              <label htmlFor="location" className="label">
                <span className="label-text px-1">Location</span>
              </label>
              <div className="relative">
                <MapPin className="absolute top-1/2 left-3 z-10 size-5 -translate-y-1/2 transform text-white opacity-70" />
                <input
                  id="location"
                  type="text"
                  className="input-bordered input w-full pl-10"
                  placeholder="City, Country"
                  {...form.register("location")}
                />
              </div>
              <div className="px-1 text-xs text-red-500">
                {formerror.location?.message}
              </div>
            </div>

            {/* SUBMIT BUTTON */}

            <button
              className="btn w-full btn-primary"
              disabled={onboarding.isPending || isRedirecting}
              type="submit"
            >
              {onboarding.isPending ? (
                <SpinnerBtn msg="Onboarding..." />
              ) : isRedirecting ? (
                <SpinnerBtn msg="Redirecting..." />
              ) : (
                <>
                  <ShipWheelIcon className="mr-2 size-5" />
                  Complete Onboarding
                </>
              )}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

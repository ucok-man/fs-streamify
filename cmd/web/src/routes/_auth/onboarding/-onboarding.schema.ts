import { z } from "zod";

export const onboardingSchema = z.object({
  fullname: z
    .string()
    .min(2, { message: "Full name must be at least 2 characters long" })
    .max(100, { message: "Full name must be less than 100 characters" }),

  bio: z
    .string()
    .min(10, { message: "Bio must be at least 10 characters long" })
    .max(500, { message: "Bio must be less than 500 characters" }),

  native_lng: z
    .string()
    .trim()
    .min(1, { message: "Native language must be selected" }),

  learning_lng: z.string().trim().min(1, {
    message: "Learning language must be selected",
  }),

  location: z
    .string()
    .trim()
    .min(1, { message: "Location must be provided" })
    .regex(/^[^,]+,\s*[^,]+$/, {
      message: "Location must be in the format: City, Country",
    }),

  profile_pic: z.string().url(),
});

export type OnboardingData = z.infer<typeof onboardingSchema>;

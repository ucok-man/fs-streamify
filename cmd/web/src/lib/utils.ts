import { clsx, type ClassValue } from "clsx";
import type { useForm } from "react-hook-form";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/* eslint-disable @typescript-eslint/no-explicit-any */
// Parsing 422 Error ONLY from Backend
export function parseApiError(
  data: any,
  form: ReturnType<typeof useForm<any>>
) {
  const error = data.error;
  for (const key in error) {
    for (const validkey in form.getValues()) {
      if (key === validkey) {
        form.setError(key as any, {
          message: error[key][0],
          type: "onChange",
        });
      }
    }
  }
}

export function generateAvatar() {
  const idx = Math.floor(Math.random() * 100) + 1; // 1-100 included
  return `https://avatar.iran.liara.run/public/${idx}.png`;
}

export const capitialize = (str: string) =>
  str.charAt(0).toUpperCase() + str.slice(1);

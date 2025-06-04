import { QueryClient } from "@tanstack/react-query";

export const queryclient = new QueryClient();

export const refetchQuery = (querykey: string[]) => {
  querykey.forEach((key) => {
    queryclient.invalidateQueries({
      queryKey: [key],
    });
    queryclient.refetchQueries({
      queryKey: [key],
    });
  });
};

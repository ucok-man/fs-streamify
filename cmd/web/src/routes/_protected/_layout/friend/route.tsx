import { useInfiniteQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useMemo, type PropsWithChildren } from "react";
import toast from "react-hot-toast";
import { useIntersectionObserver } from "usehooks-ts";
import FriendCard from "../../../../components/friend-card";
import { apiclient } from "../../../../lib/apiclient";
import type { MetadataResponse } from "../../../../types/metadata-response.type";
import type { UserResponse } from "../../../../types/user-response.type";
import SearchBox from "./-search-box";

export const Route = createFileRoute("/_protected/_layout/friend")({
  validateSearch: (search: Record<string, unknown>) => {
    const query = search?.query;
    if (!query) return { query: undefined };
    if (Array.isArray(query)) return { query: query[0] as string };
    return { query: query as string };
  },
  component: FriendPage,
});

function FriendPage() {
  const searchParams = Route.useSearch();
  const {
    data,
    isPending,
    error,
    hasNextPage,
    isFetchingNextPage,
    fetchNextPage,
  } = useInfiniteQuery({
    queryKey: ["all:users", "friends", searchParams.query],
    queryFn: async ({ pageParam }) => {
      const { data } = await apiclient.get(
        `/users/friends-with-me?page_size=8&page=${pageParam}&query=${searchParams.query || ""}`
      );
      return data as {
        users: UserResponse[];
        metadata: MetadataResponse;
      };
    },
    initialPageParam: 1,
    getNextPageParam: (last) => {
      return last.metadata.current_page < last.metadata.last_page
        ? last.metadata.current_page + 1
        : undefined;
    },
  });

  const observer = useIntersectionObserver({
    threshold: 0,
    rootMargin: "260px",
  });

  useEffect(() => {
    if (observer.isIntersecting && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [observer.isIntersecting, hasNextPage, isFetchingNextPage, fetchNextPage]);

  const users = useMemo(
    () => data?.pages.flatMap((page) => page.users) ?? [],
    [data?.pages]
  );

  if (isPending) {
    return (
      <Wrapper>
        <div className="flex h-[calc(100vh-320px)] w-full items-center justify-center">
          <span className="loading loading-lg loading-spinner" />
        </div>
      </Wrapper>
    );
  }

  if (error) {
    toast.error("Sorry we have problem in our server. Please try again later!");
    return (
      <Wrapper>
        <div className="flex h-[calc(100vh-320px)] w-full items-center justify-center">
          <span className="loading loading-lg loading-spinner" />
        </div>
      </Wrapper>
    );
  }

  if (users.length === 0) {
    return (
      <Wrapper>
        <div className="card flex h-[calc(100vh-320px)] w-full items-center justify-center bg-base-200 p-5 text-center">
          <h3 className="mb-2 text-lg font-semibold">No friends yet</h3>
          <p className="text-base-content opacity-70">
            Connect with language partners to start practicing together!
          </p>
        </div>
      </Wrapper>
    );
  }

  return (
    <Wrapper>
      <div className="relative">
        <div className="relative grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {users.map((user, idx) => (
            <FriendCard key={idx} friend={user} />
          ))}
        </div>
        {/* Marker */}
        {hasNextPage && (
          <div
            ref={observer.ref}
            id="marker"
            className="pointer-events-none absolute bottom-0 -z-50"
          />
        )}
      </div>
    </Wrapper>
  );
}

function Wrapper({ children }: PropsWithChildren) {
  return (
    <div className="p-4 sm:p-6 lg:p-8">
      <div className="container mx-auto space-y-8">
        <section className="relative">
          {/* HEADER */}
          <div className="mb-6 space-y-6 sm:mb-8 sm:space-y-8">
            <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
              <div>
                <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">
                  Your Friends
                </h2>
                <p className="opacity-70">
                  Message your friend and start learning together
                </p>
              </div>
            </div>

            <div className="flex w-full items-center justify-start">
              <SearchBox />
            </div>
          </div>

          <div className="my-8">{children}</div>
        </section>
      </div>
    </div>
  );
}

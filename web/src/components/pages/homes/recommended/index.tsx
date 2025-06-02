import { useInfiniteQuery } from "@tanstack/react-query";
import { Link, useSearch } from "@tanstack/react-router";
import { UsersIcon } from "lucide-react";
import { useEffect, useMemo, type PropsWithChildren } from "react";
import toast from "react-hot-toast";
import { useIntersectionObserver } from "usehooks-ts";
import { apiclient } from "../../../../lib/apiclient";
import type { MetadataResponse } from "../../../../types/metadata-response.type";
import type { UserWithFriendRequestResponse } from "../../../../types/user-with-friend-request-response.type";
import RecommendCard from "../../../recommend-card";
import SearchBox from "./search-box";

export default function Recommended() {
  const searchParams = useSearch({
    from: "/_protected/_layout/(home)/",
  });

  const {
    data,
    isPending,
    error,
    hasNextPage,
    isFetchingNextPage,
    fetchNextPage,
  } = useInfiniteQuery({
    queryKey: ["all:users", "recommended", searchParams.query],
    queryFn: async ({ pageParam }) => {
      const { data } = await apiclient.get(
        `/users/recommended?page_size=8&page=${pageParam}&query=${searchParams.query || ""}`
      );
      return data as {
        users: UserWithFriendRequestResponse[];
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
    rootMargin: "380px",
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
        <div className="flex h-[280px] w-full items-center justify-center">
          <span className="loading loading-lg loading-spinner" />
        </div>
      </Wrapper>
    );
  }

  if (error) {
    toast.error("Sorry we have problem in our server. Please try again later!");
    return (
      <Wrapper>
        <div className="flex h-[280px] w-full items-center justify-center">
          <span className="loading loading-lg loading-spinner" />
        </div>
      </Wrapper>
    );
  }

  if (users.length === 0) {
    return (
      <Wrapper>
        <div className="card flex h-[280px] w-full items-center justify-center bg-base-200 p-5 text-center">
          <h3 className="mb-2 text-lg font-semibold">
            No recommendations available
          </h3>
          <p className="text-base-content opacity-70">
            Check back later for new language partners!
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
            <RecommendCard key={idx} user={user} />
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
    <section className="relative">
      <div className="mb-6 space-y-6 sm:mb-8 sm:space-y-8">
        <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">
              Meet New Learners
            </h2>
            <p className="opacity-70">
              Discover perfect language exchange partners based on your profile
            </p>
          </div>
        </div>

        <div className="flex w-full items-center justify-start">
          <SearchBox />
        </div>
      </div>

      <Link
        to="/notification"
        className="btn absolute top-1 right-0 btn-outline btn-sm"
      >
        <UsersIcon className="mr-2 size-4" />
        Friend Requests
      </Link>

      <div className="my-8">{children}</div>
    </section>
  );
}

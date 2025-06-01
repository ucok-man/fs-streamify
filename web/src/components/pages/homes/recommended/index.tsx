import { useQuery } from "@tanstack/react-query";
import type { PropsWithChildren } from "react";
import toast from "react-hot-toast";
import { apiclient } from "../../../../lib/apiclient";
import type { MetadataResponse } from "../../../../types/metadata-response.type";
import type { UserWithFriendRequestResponse } from "../../../../types/user-with-friend-request-response.type";
import RecommendCard from "../../../recommend-card";

export default function Recommended() {
  const { data, isPending, error } = useQuery({
    queryKey: ["all:users", "recommended"],
    queryFn: async () => {
      const { data } = await apiclient.get(
        "/users/recommended?page=1&page_size=9"
      );
      return data as {
        users: UserWithFriendRequestResponse[];
        metadata: MetadataResponse;
      };
    },
  });

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

  if (data.users.length === 0) {
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
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {data.users.map((user) => (
          <RecommendCard key={user.id} user={user} />
        ))}
      </div>
    </Wrapper>
  );
}

function Wrapper({ children }: PropsWithChildren) {
  return (
    <section>
      <div className="mb-6 sm:mb-8">
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
      </div>

      <div className="my-8">{children}</div>
    </section>
  );
}

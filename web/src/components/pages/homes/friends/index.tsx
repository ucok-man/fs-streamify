import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { UsersIcon } from "lucide-react";
import type { JSX, PropsWithChildren } from "react";

import toast from "react-hot-toast";
import { apiclient } from "../../../../lib/apiclient";
import type { MetadataResponse } from "../../../../types/metadata-response.type";
import type { UserResponse } from "../../../../types/user-response.type";
import FriendCard from "../../../friend-card";

export default function Friends(): JSX.Element {
  const { data, isPending, error } = useQuery({
    queryKey: ["all:users", "friends"],
    queryFn: async () => {
      const response = await apiclient.get(
        "/users/friends-with-me?page=1&page_size=5"
      );
      return response.data as {
        users: UserResponse[];
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

  if (data?.users.length === 0) {
    return (
      <Wrapper>
        <div className="card flex h-[280px] w-full items-center justify-center bg-base-200 p-5 text-center">
          <h3 className="mb-2 text-lg font-semibold">No friends yet</h3>
          <p className="text-base-content opacity-70">
            Connect with language partners below to start practicing together!
          </p>
        </div>
      </Wrapper>
    );
  }

  return (
    <Wrapper>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {data?.users.map((friend) => (
          <FriendCard key={friend.id} friend={friend} />
        ))}
      </div>
    </Wrapper>
  );
}

function Wrapper({ children }: PropsWithChildren) {
  return (
    <section>
      <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
        <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">
          Your Friends
        </h2>
        <Link to="/notification" className="btn btn-outline btn-sm">
          <UsersIcon className="mr-2 size-4" />
          Friend Requests
        </Link>
      </div>

      <div className="my-8">{children}</div>
    </section>
  );
}

import { useQuery } from "@tanstack/react-query";
import { UserCheckIcon } from "lucide-react";
import type { PropsWithChildren } from "react";
import toast from "react-hot-toast";
import { apiclient } from "../../../../lib/apiclient";
import type { FriendRequestWithSenderResponse } from "../../../../types/friend-request-with-sender-response.type";
import type { MetadataResponse } from "../../../../types/metadata-response.type";

export default function IncomingFriendRequest() {
  const { data, isPending, error } = useQuery({
    queryKey: ["incoming:friend:request"],
    queryFn: async () => {
      const { data } = await apiclient.get("/users/friends-request/from");
      return data as {
        friend_requests: FriendRequestWithSenderResponse[];
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

  if (data?.friend_requests.length === 0) {
    return (
      <Wrapper>
        <div className="card flex h-[280px] w-full items-center justify-center bg-base-200 p-5 text-center">
          <h3 className="mb-2 text-lg font-semibold">
            No incoming friends request yet
          </h3>
          {/* TODO: */}
          <p className="text-base-content opacity-70">
            Connect with language partners below to start practicing together!
          </p>
        </div>
      </Wrapper>
    );
  }

  console.log({ data });

  return <div>IncomingFriendRequest</div>;
}

function Wrapper({ children }: PropsWithChildren) {
  return (
    <section className="space-y-4">
      <h2 className="flex items-center gap-2 text-xl font-semibold">
        <UserCheckIcon className="h-5 w-5 text-primary" />
        Incoming Friend Requests
      </h2>

      {children}
    </section>
  );
}

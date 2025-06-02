import { useMutation } from "@tanstack/react-query";
import { CheckCircleIcon, MapPinIcon, UserPlusIcon } from "lucide-react";
import toast from "react-hot-toast";
import { apiclient } from "../../lib/apiclient";
import { refetchQuery } from "../../lib/query-client";
import { capitialize, cn } from "../../lib/utils";
import type { UserWithFriendRequestResponse } from "../../types/user-with-friend-request-response.type";
import LanguageFlag from "../language-flag";

type Props = {
  user: UserWithFriendRequestResponse;
};

export default function RecommendCard({ user }: Props) {
  const isHasSent = user.sent_friend_request.length > 0;
  const isHasFrom = user.from_friend_request.length > 0;

  const sentRequest = useMutation({
    mutationFn: async (receipentId: string) => {
      return await apiclient.post(
        `/users/friends-request/create/${receipentId}`
      );
    },
    onError: (err) => {
      console.log({ err });
      toast.error("Failed to sent friend request!");
    },
    onSuccess: () => {
      toast.success("Success sending friend request");
      refetchQuery(["all:users", "recommended"]);
      // TODO: refetch query
    },
  });

  return (
    <div className="card bg-base-200 transition-all duration-300 hover:shadow-lg">
      <div className="card-body space-y-4 p-5">
        <div className="flex items-center gap-3">
          <div className="avatar size-16 rounded-full">
            <img src={user.profile_pic} alt={user.full_name} />
          </div>

          <div>
            <h3 className="text-lg font-semibold">{user.full_name}</h3>
            {user.location && (
              <div className="mt-1 flex items-center text-xs opacity-70">
                <MapPinIcon className="mr-1 size-3" />
                {user.location}
              </div>
            )}
          </div>
        </div>

        {/* Languages with flags */}
        <div className="flex flex-wrap gap-1.5 space-y-1">
          <span className="badge badge-secondary">
            <LanguageFlag language={user.native_lng} />
            Native: {capitialize(user.native_lng)}
          </span>
          <span className="badge-outline badge">
            <LanguageFlag language={user.learning_lng} />
            Learning: {capitialize(user.learning_lng)}
          </span>
        </div>

        <p className="line-clamp-2 text-sm opacity-70">{user.bio}</p>

        {/* Action button */}
        <button
          className={cn(
            `btn mt-2 w-full`,
            isHasSent || isHasFrom ? "btn-disabled" : "btn-primary"
          )}
          onClick={() => sentRequest.mutate(user.id)}
          disabled={isHasSent || isHasFrom || sentRequest.isPending}
        >
          {isHasSent || isHasFrom ? (
            <>
              <CheckCircleIcon className="mr-2 size-4" />
              See Notification
            </>
          ) : (
            <>
              <UserPlusIcon className="mr-2 size-4" />
              Send Friend Request
            </>
          )}
        </button>
      </div>
    </div>
  );
}

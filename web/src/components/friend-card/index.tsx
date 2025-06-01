import { Link } from "@tanstack/react-router";
import { capitialize } from "../../lib/utils";
import type { UserResponse } from "../../types/user-response.type";
import LanguageFlag from "../language-flag";

type Props = {
  friend: UserResponse;
};

export default function FriendCard({ friend }: Props) {
  return (
    <div className="card bg-base-200 transition-shadow hover:shadow-md">
      <div className="card-body p-4">
        {/* USER INFO */}
        <div className="mb-3 flex items-center gap-3">
          <div className="avatar size-12">
            <img src={friend.profile_pic} alt={friend.full_name} />
          </div>
          <h3 className="truncate font-semibold">{friend.full_name}</h3>
        </div>

        <div className="mb-3 flex flex-wrap gap-1.5 space-y-1">
          <span className="badge text-xs badge-secondary">
            <LanguageFlag language={friend.native_lng} />
            Native: {capitialize(friend.native_lng)}
          </span>
          <span className="badge-outline badge text-xs">
            <LanguageFlag language={friend.learning_lng} />
            Learning: {capitialize(friend.learning_lng)}
          </span>
        </div>

        <Link
          to={`/chat/$id`}
          params={{
            id: friend.id,
          }}
          className="btn w-full btn-outline"
        >
          Message
        </Link>
      </div>
    </div>
  );
}

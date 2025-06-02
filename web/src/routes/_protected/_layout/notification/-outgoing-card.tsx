import type { FriendRequestWithRecipientResponse } from "../../../../types/friend-request-with-recipient-response.type";

type Props = {
  item: FriendRequestWithRecipientResponse;
};

export default function OutgoingCard({ item }: Props) {
  return (
    <div className="card bg-base-200 shadow-sm transition-shadow hover:shadow-md">
      <div className="card-body p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="avatar h-14 w-14 rounded-full bg-base-300">
              <img
                src={item.recipient.profile_pic}
                alt={item.recipient.full_name}
              />
            </div>
            <div>
              <h3 className="font-semibold">{item.recipient.full_name}</h3>
              <div className="mt-1 flex flex-wrap gap-1.5">
                <span className="badge badge-sm badge-secondary">
                  Native: {item.recipient.native_lng}
                </span>
                <span className="badge-outline badge badge-sm">
                  Learning: {item.recipient.learning_lng}
                </span>
              </div>
            </div>
          </div>

          <div className="badge badge-neutral">{item.status}</div>
        </div>
      </div>
    </div>
  );
}

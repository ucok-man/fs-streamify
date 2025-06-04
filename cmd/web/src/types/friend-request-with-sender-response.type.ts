export type FriendRequestWithSenderResponse = {
  id: string;
  sender_id: string;
  recipient_id: string;
  status: "Pending" | "Accepted";
  created_at: string; // ISO date string (time.Time in Go)
  updated_at: string;
  sender: {
    id: string;
    full_name: string;
    email: string;
    bio: string;
    profile_pic: string;
    native_lng: string;
    learning_lng: string;
    location: string;
    is_onboarded: boolean;
    friend_ids: string[];
    created_at: string;
    updated_at: string;
  };
};

import type { FriendRequestResponse } from "./friend-request-response.type";

export type UserWithFriendRequestResponse = {
  id: string;
  full_name: string;
  email: string;
  bio: string;
  profile_pic: string;
  native_lng: string;
  learning_lng: string;
  location: string;
  is_onboarded: true;
  friend_ids: string[];
  created_at: string;
  updated_at: string;
  sent_friend_request: FriendRequestResponse[];
  from_friend_request: FriendRequestResponse[];
};

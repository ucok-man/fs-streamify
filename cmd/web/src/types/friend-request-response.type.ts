export type FriendRequestResponse = {
  id: string;
  sender_id: string;
  recipient_id: string;
  status: "Accepted" | "Pending";
  created_at: string;
  uppdated_at: string;
};

package com.mardomsara.social.json;

public class J {
	public static class Chat {//oridnal: 0
		public String ChatKey;
		public String RoomKey;
		public int RoomTypeEnumId;
		public int UserId;
		public int PeerUserId;
		public int GroupId;
		public int CreatedSe;
		public int StartMessageIdFrom;
		public int LastSeenMessageId;
		public int UpdatedMs;
		public int LastMessageId;
		public int LastDeletedMessageId;
		public int LastSeqSeen;
		public int LastSeqDelete;
		public int CurrentSeq;
	}

	public static class Comments {//oridnal: 1
		public int Id;
		public int UserId;
		public int PostId;
		public String Text;
		public int CreatedTime;
	}

}
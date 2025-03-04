import { create } from "zustand";

interface Message {
  sender: string;
  receiver: string;
  text: string;
  gameID: string;
}

interface ChatState {
  messages: Message[];
  setMessages: (message: Message) => void;
}

const useChatStore = create<ChatState>()((set, get) => ({
  messages: [],
  setMessages: (message: Message) => {
    if (get().messages) {
      set((state) => ({ messages: [...state.messages, message] }));
    } else {
      set(() => ({ messages: [message] }));
    }
  },
}));

export default useChatStore;

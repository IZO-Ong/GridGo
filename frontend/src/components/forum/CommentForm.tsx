import { useState } from "react";

interface CommentFormProps {
  onSubmit: (content: string) => Promise<void>;
}

export default function CommentForm({ onSubmit }: CommentFormProps) {
  const [text, setText] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!text.trim() || submitting) return;
    setSubmitting(true);
    try {
      await onSubmit(text);
      setText("");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="APPEND_DATA_TO_THREAD..."
        className="w-full border-4 border-black p-4 font-bold text-sm focus:bg-zinc-50 outline-none h-32"
        required
      />
      <button
        disabled={submitting}
        className="px-8 py-3 bg-black text-white font-black uppercase italic border-4 border-black hover:bg-zinc-800 disabled:opacity-50 cursor-pointer transition-all active:translate-y-1"
      >
        {submitting ? "TRANSMITTING..." : "POST_COMMENT"}
      </button>
    </form>
  );
}

interface VoteSidebarProps {
  upvotes: number;
  userVote?: number;
  onVote: (val: number) => void;
  small?: boolean;
}

export default function VoteSidebar({
  upvotes,
  userVote,
  onVote,
  small,
}: VoteSidebarProps) {
  const size = small ? "18" : "28";
  const containerWidth = small ? "w-10" : "w-14";
  const textSize = small ? "text-[10px]" : "text-lg";

  return (
    <div
      className={`${containerWidth} bg-zinc-50 border-r-4 border-black flex flex-col items-center py-6 gap-2 shrink-0 cursor-default`}
    >
      <button
        onClick={() => onVote(1)}
        className={`cursor-pointer transition-all hover:scale-110 ${userVote === 1 ? "text-black opacity-100" : "opacity-20 hover:opacity-100"}`}
      >
        <svg
          width={size}
          height={size}
          viewBox="0 0 24 24"
          fill={userVote === 1 ? "black" : "none"}
          stroke="black"
          strokeWidth="4"
        >
          <path d="M18 15l-6-6-6 6" />
        </svg>
      </button>

      <span className={`font-black italic ${textSize} select-none`}>
        {upvotes}
      </span>

      <button
        onClick={() => onVote(-1)}
        className={`cursor-pointer transition-all hover:scale-110 ${userVote === -1 ? "text-black opacity-100" : "opacity-20 hover:opacity-100"}`}
      >
        <svg
          width={size}
          height={size}
          viewBox="0 0 24 24"
          fill={userVote === -1 ? "black" : "none"}
          stroke="black"
          strokeWidth="4"
        >
          <path d="M6 9l6 6 6-6" />
        </svg>
      </button>
    </div>
  );
}

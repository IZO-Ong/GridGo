import Link from "next/link";

interface CommentCardProps {
  comment: any;
  isOwner: boolean;
  onDelete: (id: string) => void;
}

export default function ProfileCommentCard({
  comment,
  isOwner,
  onDelete,
}: CommentCardProps) {
  return (
    <div className="relative group transition-all duration-200 hover:-translate-y-1">
      <div className="border-4 border-black p-6 bg-white space-y-3 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:shadow-[10px_10px_0px_0px_rgba(0,0,0,1)] transition-all relative">
        <div className="text-[10px] font-black uppercase opacity-40">
          On Thread:{" "}
          <Link
            href={`/forum/post/${comment.post_id}`}
            className="text-black underline"
          >
            {comment.post?.title || "Unknown_Thread"}
          </Link>
        </div>
        <p className="font-medium text-sm border-l-4 border-black pl-4">
          {comment.content}
        </p>

        {/* Sync'd Delete Button */}
        {isOwner && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onDelete(comment.id);
            }}
            className="absolute top-2 right-2 bg-black text-white p-1.5 border-2 border-black opacity-0 group-hover:opacity-100 transition-all hover:bg-red-600 z-20 cursor-pointer shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:shadow-none"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="4"
            >
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>
    </div>
  );
}

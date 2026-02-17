import Link from "next/link";

interface MazeCardProps {
  maze: any;
  isOwner: boolean;
  onDelete: (e: React.MouseEvent, id: string) => void;
}

export default function MazeRepositoryCard({
  maze,
  isOwner,
  onDelete,
}: MazeCardProps) {
  return (
    <div className="relative group hover:-translate-y-1 transition-all duration-200">
      <Link href={`/solve?id=${maze.id}`}>
        <div className="border-4 border-black p-4 bg-white shadow-[0px_0px_0px_0px_rgba(0,0,0,1)] group-hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] transition-all flex flex-col h-full">
          <div
            className="aspect-video bg-zinc-50 border-2 border-black mb-4 flex items-center justify-center overflow-hidden relative bg-cover bg-center"
            style={{ backgroundImage: `url(${maze.thumbnail})` }}
          >
            <div className="absolute inset-0 bg-white/80 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center backdrop-blur-[2px]">
              <span className="text-[10px] font-black uppercase tracking-widest z-10">
                VIEW_MATRIX_{maze.id}
              </span>
            </div>
          </div>

          <div className="mt-auto space-y-3">
            <div className="flex justify-between items-center">
              <span className="font-black text-sm italic tracking-tighter">
                {maze.id}
              </span>
              <span className="bg-black text-white text-[9px] font-black px-2 py-0.5">
                {maze.rows}x{maze.cols}
              </span>
            </div>

            <div className="flex justify-between items-center border-t border-black border-dotted pt-2">
              <span className="text-[9px] font-bold opacity-60">
                COMPLEXITY: {maze.complexity?.toFixed(2) ?? "0.00"}
              </span>
              <span className="text-[9px] font-bold opacity-60">
                NODES: {maze.rows * maze.cols}
              </span>
            </div>

            <div className="text-[8px] font-bold opacity-30 text-right">
              {new Date(maze.created_at).toLocaleDateString()}
            </div>
          </div>
        </div>
      </Link>

      {isOwner && (
        <button
          onClick={(e) => onDelete(e, maze.id)}
          className="absolute top-2 right-2 bg-black text-white p-1.5 border-2 border-black opacity-0 group-hover:opacity-100 transition-all hover:bg-red-600 z-20 cursor-pointer shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
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
  );
}

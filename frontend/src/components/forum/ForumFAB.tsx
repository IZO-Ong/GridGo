import Link from "next/link";

export default function ForumFAB() {
  return (
    <div className="fixed bottom-20 left-1/2 -translate-x-1/2 max-w-5xl w-full pointer-events-none z-50">
      <Link
        href="/forum/new"
        className="absolute right-8 w-16 h-16 bg-black text-white rounded-full flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,0.3)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all border-4 border-black group cursor-pointer pointer-events-auto"
        title="START_NEW_THREAD"
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="4"
          className="group-hover:rotate-90 transition-transform"
        >
          <path d="M12 5v14M5 12h14" />
        </svg>
      </Link>
    </div>
  );
}

interface AuthCardProps {
  title: string;
  error: string | null;
  children: React.ReactNode;
  footerAction?: () => void;
  footerLabel?: string;
}

export default function AuthCard({
  title,
  error,
  children,
  footerAction,
  footerLabel,
}: AuthCardProps) {
  return (
    <div className="flex flex-col items-center justify-start min-h-screen pt-12 pb-20 space-y-4">
      {error && (
        <div
          className={`w-full max-w-[380px] p-3 border-2 font-bold uppercase text-[10px] ${
            error.includes("SUCCESS")
              ? "border-green-600 bg-green-50 text-green-600"
              : "border-red-600 bg-red-50 text-red-600"
          }`}
        >
          {`>> ${error}`}
        </div>
      )}
      <div className="w-full max-w-[380px] border-4 border-black p-6 bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
        <h2 className="text-2xl font-black uppercase tracking-tighter mb-6">
          {title}
        </h2>
        {children}
        {footerAction && (
          <div className="mt-8 pt-4 border-t-2 border-black border-dotted">
            <button
              onClick={footerAction}
              className="text-[10px] font-black uppercase underline decoration-2 underline-offset-4 hover:text-zinc-500 cursor-pointer"
            >
              {footerLabel}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

"use client";

interface DeleteModalProps {
  type: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export default function DeleteModal({
  type,
  onConfirm,
  onCancel,
}: DeleteModalProps) {
  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in duration-200">
      {/* Background click to close */}
      <div className="absolute inset-0" onClick={onCancel} />

      <div className="relative w-full max-w-sm border-4 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] p-6 animate-in zoom-in-95 duration-200">
        {/* Header */}
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-black uppercase tracking-tight italic text-red-600">
            PURGE {type.toUpperCase()}?
          </h2>

          {/* Custom X Button: Black bg, White X, Red on Hover */}
          <button
            onClick={onCancel}
            className="flex items-center justify-center w-8 h-8 bg-black text-white hover:bg-[#ef4444] border-2 border-black transition-colors"
            aria-label="Close"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="4"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        {/* Content Section */}
        <div className="mb-8">
          <p className="text-xs font-bold opacity-60 leading-tight tracking-tighter">
            This action is irreversible. All associated data will be scrubbed.
          </p>
        </div>

        {/* Action Buttons with Neo-Brutalist Shadow and Click Animation */}
        <div className="flex flex-col gap-3">
          <button
            onClick={onConfirm}
            className="w-full bg-red-600 text-white border-4 border-black py-3 font-black uppercase tracking-widest shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:bg-red-700 active:shadow-none active:translate-x-1 active:translate-y-1 transition-all"
          >
            Confirm
          </button>

          <button
            onClick={onCancel}
            className="w-full bg-white text-black border-4 border-black py-3 font-black uppercase tracking-widest shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:bg-zinc-100 active:shadow-none active:translate-x-1 active:translate-y-1 transition-all"
          >
            Abort
          </button>
        </div>
      </div>
    </div>
  );
}

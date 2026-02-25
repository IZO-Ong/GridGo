"use client";
import { useState } from "react";

interface ShareModalProps {
  url: string;
  onClose: () => void;
}

export default function ShareModal({ url, onClose }: ShareModalProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // Shared encoded strings for social links
  const shareText = encodeURIComponent("Solve this maze on GridGo!");
  const shareUrl = encodeURIComponent(url);

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="absolute inset-0" onClick={onClose} />

      <div className="relative w-full max-w-sm border-4 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] p-6 animate-in zoom-in-95 duration-200">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-black uppercase tracking-tight">
            Share Maze
          </h2>
          <button
            onClick={onClose}
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

        <div className="mb-8">
          <p className="text-[10px] font-bold uppercase mb-2 opacity-60 tracking-wider">
            Maze Access Link
          </p>
          <div className="flex border-2 border-black bg-zinc-100 p-1">
            <input
              readOnly
              value={url}
              className="flex-1 bg-transparent px-3 py-2 text-xs font-medium outline-none truncate"
            />
            <button
              onClick={handleCopy}
              className={`w-12 flex items-center justify-center transition-colors border-l-2 border-black ${
                copied
                  ? "bg-green-400"
                  : "bg-black text-white hover:bg-zinc-800"
              }`}
              title={copied ? "Copied!" : "Copy Link"}
            >
              {copied ? (
                <svg
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="3"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <polyline points="20 6 9 17 4 12" />
                </svg>
              ) : (
                <svg
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                </svg>
              )}
            </button>
          </div>
        </div>

        {/* Social Circles: Now including Facebook, Email, and Pinterest */}
        <div className="flex flex-wrap justify-center gap-4 py-2">
          {/* WhatsApp */}
          <a
            href={`https://wa.me/?text=${shareText}%20${shareUrl}`}
            target="_blank"
            className="group"
          >
            <div className="w-12 h-12 rounded-full bg-[#25D366] border-2 border-black flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:translate-x-[2px] group-hover:translate-y-[2px] group-hover:shadow-none transition-all">
              <svg width="24" height="24" fill="white" viewBox="0 0 24 24">
                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.395 0 .01 5.388 0 12.044c0 2.129.559 4.207 1.626 6.07L0 24l6.148-1.613a11.771 11.771 0 005.9 1.594h.005c6.654 0 12.04-5.39 12.042-12.046a11.813 11.813 0 00-3.535-8.527z" />
              </svg>
            </div>
          </a>

          {/* Telegram */}
          <a
            href={`https://t.me/share/url?url=${shareUrl}&text=${shareText}`}
            target="_blank"
            className="group"
          >
            <div className="w-12 h-12 rounded-full bg-[#0088cc] border-2 border-black flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:translate-x-[2px] group-hover:translate-y-[2px] group-hover:shadow-none transition-all">
              <svg width="22" height="22" fill="white" viewBox="0 0 24 24">
                <path d="M11.944 0C5.347 0 0 5.347 0 11.944 0 18.541 5.347 23.888 11.944 23.888s11.944-5.347 11.944-11.944C23.888 5.347 18.541 0 11.944 0zm5.889 8.161l-1.996 9.417c-.151.67-.546.835-1.107.52l-3.041-2.241-1.467 1.412c-.162.162-.298.298-.612.298l.218-3.098 5.639-5.094c.245-.218-.054-.339-.381-.121l-6.968 4.387-3.001-.938c-.652-.204-.666-.652.136-.964l11.734-4.52c.544-.197 1.018.13 1.018.842z" />
              </svg>
            </div>
          </a>

          {/* Facebook */}
          <a
            href={`https://www.facebook.com/sharer/sharer.php?u=${shareUrl}`}
            target="_blank"
            className="group"
          >
            <div className="w-12 h-12 rounded-full bg-[#1877F2] border-2 border-black flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:translate-x-[2px] group-hover:translate-y-[2px] group-hover:shadow-none transition-all">
              <svg width="24" height="24" fill="white" viewBox="0 0 24 24">
                <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
              </svg>
            </div>
          </a>

          {/* Pinterest */}
          <a
            href={`https://pinterest.com/pin/create/button/?url=${shareUrl}&description=${shareText}`}
            target="_blank"
            className="group"
          >
            <div className="w-12 h-12 rounded-full bg-[#BD081C] border-2 border-black flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:translate-x-[2px] group-hover:translate-y-[2px] group-hover:shadow-none transition-all">
              <svg width="20" height="20" fill="white" viewBox="0 0 24 24">
                <path d="M12.017 0C5.396 0 .029 5.367.029 11.987c0 5.079 3.158 9.417 7.618 11.162-.105-.949-.199-2.403.041-3.439.219-.937 1.406-5.965 1.406-5.965s-.359-.718-.359-1.78c0-1.667.967-2.911 2.168-2.911 1.024 0 1.518.769 1.518 1.688 0 1.029-.653 2.567-.992 3.992-.285 1.193.6 2.165 1.775 2.165 2.128 0 3.768-2.245 3.768-5.487 0-2.861-2.063-4.869-5.008-4.869-3.41 0-5.409 2.562-5.409 5.199 0 1.033.394 2.143.889 2.741.099.12.112.225.085.345-.09.375-.293 1.199-.334 1.363-.053.211-.174.256-.402.151-1.504-.699-2.445-2.895-2.445-4.659 0-3.792 2.759-7.271 7.944-7.271 4.166 0 7.402 2.969 7.402 6.936 0 4.14-2.611 7.471-6.233 7.471-1.214 0-2.356-.632-2.746-1.378 0 0-.6 2.281-.746 2.847-.27 1.042-1.002 2.348-1.493 3.141 1.127.348 2.325.539 3.577.539 6.62 0 11.987-5.367 11.987-11.987C24 5.367 18.633 0 12.017 0z" />
              </svg>
            </div>
          </a>

          {/* Email */}
          <a
            href={`mailto:?subject=Check%20out%20this%20maze!&body=Solve%20this%20maze%20on%20GridGo:%20${url}`}
            className="group"
          >
            <div className="w-12 h-12 rounded-full bg-[#71717a] border-2 border-black flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:translate-x-[2px] group-hover:translate-y-[2px] group-hover:shadow-none transition-all">
              <svg
                width="22"
                height="22"
                fill="none"
                stroke="white"
                strokeWidth="2.5"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75"
                />
              </svg>
            </div>
          </a>
        </div>
      </div>
    </div>
  );
}

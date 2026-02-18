"use client";
const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function OAuthButton() {
  const handleOAuth = () => {
    window.location.href = `${BASE_URL}/api/auth/google`;
  };

  return (
    <>
      <div className="relative flex items-center my-8">
        <div className="flex-grow border-t-2 border-dotted border-black"></div>
        <span className="flex-shrink mx-4 text-[9px] font-black uppercase bg-white px-2 opacity-40">
          External Auth
        </span>
        <div className="flex-grow border-t-2 border-dotted border-black"></div>
      </div>

      <button
        onClick={handleOAuth}
        className="w-full border-2 border-black py-2 font-black uppercase text-[11px] flex items-center justify-center gap-3 hover:bg-black hover:text-white transition-all active:translate-y-0.5 group cursor-pointer"
      >
        <svg
          width="18"
          height="18"
          viewBox="0 0 24 24"
          className="transition-colors"
        >
          <path
            fill="currentColor"
            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09zM12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23zM5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84zM12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
          />
        </svg>
        Continue with Google
      </button>
    </>
  );
}

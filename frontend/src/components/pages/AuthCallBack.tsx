"use client";
import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

export default function AuthCallback() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth();

  useEffect(() => {
    const token = searchParams.get("token");
    const username = searchParams.get("username");

    if (token && username) {
      login(token, username);
      router.push("/");
    } else {
      router.push("/login?error=OAUTH_MISSING_DATA");
    }
  }, [router, searchParams, login]);

  return (
    <div className="flex flex-col items-center justify-center h-[calc(100vh-250px)]">
      <div className="border-4 border-black p-8 bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
        <h2 className="text-xl font-black uppercase italic animate-pulse">
          Finalizing_Identity_Sync...
        </h2>
      </div>
    </div>
  );
}

"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function AuthPage() {
  const { login } = useAuth();
  const [isLogin, setIsLogin] = useState(true);
  const [isVerifying, setIsVerifying] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [otp, setOtp] = useState("");
  const [formData, setFormData] = useState({
    email: "",
    username: "",
    password: "",
    verify: "",
  });

  const router = useRouter();

  const handleOAuth = () => {
    window.location.href = `${BASE_URL}/auth/google`;
  };

  const handleAction = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!isLogin && formData.password !== formData.verify) {
      return setError("Passwords do not match");
    }

    setLoading(true);
    try {
      const endpoint = isLogin ? "/login" : "/register";
      const payload = isLogin
        ? { username: formData.email, password: formData.password }
        : {
            email: formData.email,
            username: formData.username,
            password: formData.password,
          };

      const res = await fetch(`${BASE_URL}${endpoint}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(
          data.error ||
            (isLogin ? "Invalid credentials" : "Registration failed")
        );
      }

      if (isLogin) {
        const data = await res.json();
        login(data.token, data.username);
        router.push("/");
      } else {
        setIsVerifying(true);
        setError("SUCCESS: Verification code sent to email.");
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyOTP = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await fetch(`${BASE_URL}/verify`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: formData.email, code: otp }),
      });
      if (!res.ok) throw new Error("Invalid or expired code");

      setIsLogin(true);
      setIsVerifying(false);
      setError("SUCCESS: Account verified. Please login.");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

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

      <div className="w-full max-w-[380px] border-4 border-black p-6 bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] transition-all duration-300">
        <h2 className="text-2xl font-black uppercase tracking-tighter mb-6">
          {isVerifying ? "Verify Code" : isLogin ? "Login" : "Register"}
        </h2>

        {!isVerifying ? (
          <>
            <form className="space-y-4" onSubmit={handleAction}>
              <div className="space-y-1">
                <label className="text-[9px] font-black uppercase opacity-50">
                  {isLogin ? "Username or Email" : "Email Address"}
                </label>
                <input
                  type="text"
                  required
                  value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50 transition-colors"
                  placeholder={
                    isLogin ? "Enter handle or email" : "user@domain.com"
                  }
                />
              </div>

              {!isLogin && (
                <div className="space-y-1">
                  <label className="text-[9px] font-black uppercase opacity-50">
                    Public Username
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.username}
                    onChange={(e) =>
                      setFormData({ ...formData, username: e.target.value })
                    }
                    className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50 transition-colors"
                    placeholder="e.g. MazeMaster99"
                  />
                </div>
              )}

              <div className="space-y-1">
                <label className="text-[9px] font-black uppercase opacity-50">
                  Password
                </label>
                <input
                  type="password"
                  required
                  value={formData.password}
                  onChange={(e) =>
                    setFormData({ ...formData, password: e.target.value })
                  }
                  className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50 transition-colors"
                  placeholder="••••••••"
                />
              </div>

              {!isLogin && (
                <div className="space-y-1">
                  <label className="text-[9px] font-black uppercase opacity-50">
                    Confirm Password
                  </label>
                  <input
                    type="password"
                    required
                    value={formData.verify}
                    onChange={(e) =>
                      setFormData({ ...formData, verify: e.target.value })
                    }
                    className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50 transition-colors"
                    placeholder="••••••••"
                  />
                </div>
              )}

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-black text-white py-3 font-black uppercase italic hover:bg-zinc-800 cursor-pointer transition-all active:translate-y-1 active:shadow-none mt-4 text-sm border-2 border-black disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? "Processing..." : isLogin ? "Login" : "Register"}
              </button>
            </form>

            <div className="relative flex items-center my-8">
              <div className="flex-grow border-t-2 border-dotted border-black"></div>
              <span className="flex-shrink mx-4 text-[9px] font-black uppercase bg-white px-2 opacity-40">
                External Auth
              </span>
              <div className="flex-grow border-t-2 border-dotted border-black"></div>
            </div>

            <button
              onClick={handleOAuth}
              // Monochrome styling: Black border, Black text, Black SVG
              className="w-full border-2 border-black py-2 font-black uppercase text-[11px] flex items-center justify-center gap-3 hover:bg-black hover:text-white transition-all active:translate-y-0.5 group"
            >
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                className="transition-colors"
              >
                <path
                  fill="currentColor" // Inherits text color
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09zM12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23zM5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84zM12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                />
              </svg>
              Continue with Google
            </button>
          </>
        ) : (
          <form className="space-y-4" onSubmit={handleVerifyOTP}>
            <div className="space-y-1">
              <label className="text-[9px] font-black uppercase opacity-50 text-center block">
                Enter 6-Digit Code
              </label>
              <input
                type="text"
                maxLength={6}
                required
                value={otp}
                onChange={(e) => setOtp(e.target.value)}
                className="w-full border-2 border-black p-4 outline-none font-black text-center text-2xl tracking-[0.5em] focus:bg-zinc-50"
                placeholder="000000"
              />
            </div>
            <button
              type="submit"
              className="w-full bg-black text-white py-3 font-black uppercase italic hover:bg-zinc-800 cursor-pointer transition-all active:translate-y-1 text-sm border-2 border-black"
            >
              Verify Identity
            </button>
          </form>
        )}

        <div className="mt-8 pt-4 border-t-2 border-black border-dotted flex justify-between items-center">
          <button
            onClick={() => {
              setIsLogin(!isLogin);
              setIsVerifying(false);
              setError(null);
            }}
            className="text-[10px] font-black uppercase underline decoration-2 underline-offset-4 hover:text-zinc-500 cursor-pointer transition-colors"
          >
            {isLogin ? "Create an account" : "Back to Login"}
          </button>
        </div>
      </div>
    </div>
  );
}

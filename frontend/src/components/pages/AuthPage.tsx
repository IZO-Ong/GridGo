"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import AuthCard from "@/components/auth/AuthCard";
import LoginForm from "@/components/auth/LoginForm";
import RegisterForm from "@/components/auth/RegisterForm";
import VerifyForm from "@/components/auth/VerifyForm";
import OAuthButton from "@/components/auth/OAuthButton";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function AuthPage() {
  const { login } = useAuth();
  const router = useRouter();

  const [view, setView] = useState<"login" | "register" | "verify">("login");
  const [formData, setFormData] = useState({
    email: "",
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAuthSubmit = async (payload: any) => {
    setError(null);
    setLoading(true);

    try {
      if (view === "login") {
        // Payload: { username, password }
        const res = await fetch(`${BASE_URL}/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });

        if (!res.ok) {
          const data = await res.json().catch(() => ({}));
          throw new Error(data.error || "Invalid credentials");
        }

        const data = await res.json();
        login(data.token, data.username);
        router.push("/");
      } else if (view === "register") {
        setFormData({ email: payload.email });

        const res = await fetch(`${BASE_URL}/register`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });

        if (!res.ok) {
          const data = await res.json().catch(() => ({}));
          throw new Error(data.error || "Registration failed");
        }

        setView("verify");
        setError("SUCCESS: Verification code sent to email.");
      } else if (view === "verify") {
        // Payload: string (the 6-digit OTP code)
        const res = await fetch(`${BASE_URL}/verify`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email: formData.email, code: payload }),
        });

        if (!res.ok) throw new Error("Invalid or expired code");

        setView("login");
        setError("SUCCESS: Account verified. Please login.");
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthCard
      title={
        view === "login"
          ? "Login"
          : view === "register"
            ? "Register"
            : "Verify Code"
      }
      error={error}
      footerLabel={view === "login" ? "Create an account" : "Back to Login"}
      footerAction={() => {
        setError(null); // Clear errors when switching views
        setView(view === "login" ? "register" : "login");
      }}
    >
      {view === "login" && (
        <LoginForm onSubmit={handleAuthSubmit} loading={loading} />
      )}
      {view === "register" && (
        <RegisterForm
          onSubmit={handleAuthSubmit}
          loading={loading}
          setError={setError}
        />
      )}
      {view === "verify" && (
        <VerifyForm
          email={formData.email}
          onVerify={handleAuthSubmit}
          loading={loading}
        />
      )}

      {view !== "verify" && <OAuthButton />}
    </AuthCard>
  );
}

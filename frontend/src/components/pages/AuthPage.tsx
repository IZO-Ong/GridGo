"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import AuthCard from "@/components/auth/AuthCard";
import LoginForm from "@/components/auth/LoginForm";
import RegisterForm from "@/components/auth/RegisterForm";
import VerifyForm from "@/components/auth/VerifyForm";
import OAuthButton from "@/components/auth/OAuthButton";
import {
  login as apiLogin,
  register as apiRegister,
  verifyAccount,
} from "@/lib/api";

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
        const data = await apiLogin(payload);
        login(data.token, data.username);
        router.push("/");
      } else if (view === "register") {
        setFormData({ email: payload.email });
        await apiRegister(payload);
        setView("verify");
        setError("SUCCESS: Code sent to email.");
      } else if (view === "verify") {
        await verifyAccount(formData.email, payload);
        setView("login");
        setError("SUCCESS: Verified. Please login.");
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

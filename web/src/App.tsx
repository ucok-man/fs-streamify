import { Toaster } from "react-hot-toast";
import { Route, Routes } from "react-router";
import AuthLayout from "./layouts/auth.layout";
import CallPage from "./pages/call.page";
import ChatPage from "./pages/chat.page";
import HomePage from "./pages/home.page";
import NotificationPage from "./pages/notification.page";
import OnboardingPage from "./pages/onboarding.page";
import SigninPage from "./pages/signin.page";
import SignupPage from "./pages/signup.page";

export default function App() {
  return (
    <div data-theme="night" className="size-full min-h-screen">
      <Routes>
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/signin" element={<SigninPage />} />

        <Route element={<AuthLayout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/onboarding" element={<OnboardingPage />} />
          <Route path="/notification" element={<NotificationPage />} />
          <Route path="/call" element={<CallPage />} />
          <Route path="/chat" element={<ChatPage />} />
        </Route>
      </Routes>
      <Toaster />
    </div>
  );
}

import { useAppSelector } from "@/store/hooks";
import { Navigate, Outlet } from "react-router-dom";

const ProtectedRoute = () => {
  const user = useAppSelector((state) => state.user);

  if (!user.startupComplete) {
    return <div>Loading...</div>;
  }

  if (!user.isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
};

export default ProtectedRoute;

import { Navigate, Outlet } from "react-router-dom";

const ProtectedRoute = () => {
  const user = { isAuthenticated: true };

  if (!user.isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
};

export default ProtectedRoute;

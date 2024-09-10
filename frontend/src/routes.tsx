import { createBrowserRouter } from "react-router-dom";

import AppPage from "@/pages/app-page";
import AppShell from "@/pages/AppShell";
import LandingPage from "@/pages/landing-page";
import NotFoundPage from "@/pages/not-found-page";
import ProtectedRoute from "@/pages/ProtectedRoute";
import VaultItemPage from "@/pages/vault-item-page";
import VaultLogsPage from "@/pages/vault-logs-page";
import VaultPage from "@/pages/vault-page";
import VaultUsersPage from "@/pages/vault-users-page";

const router = createBrowserRouter([
  {
    index: true,
    element: <LandingPage />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: "/app",
        element: <AppShell />,
        children: [
          {
            index: true,
            element: <AppPage />,
          },
          {
            path: "vault/:vaultId",
            element: <VaultPage />,
          },
          {
            path: "vault/:vaultId/item/:vaultItemId",
            element: <VaultItemPage />,
          },
          {
            path: "vault/:vaultId/users",
            element: <VaultUsersPage />,
          },
          {
            path: "vault/:vaultId/logs",
            element: <VaultLogsPage />,
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
]);

export default router;

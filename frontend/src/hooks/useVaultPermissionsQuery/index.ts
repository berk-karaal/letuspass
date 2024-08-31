import { listMyVaultPermissions } from "@/api/letuspass";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

/**
 * Returns the query object for the vault permissions of the current user on given vault.
 */
export const useVaultPermissionsQuery = (vaultId: number) =>
  useQuery({
    queryKey: ["vaultPermissions", vaultId],
    queryFn: () => listMyVaultPermissions(Number(vaultId)),
    retry: (failureCount: number, error: Error) => {
      if (failureCount > 2) {
        return false;
      }
      if (axios.isAxiosError(error)) {
        if (
          (error.response?.status ?? 500 >= 400) &&
          (error.response?.status ?? 500 < 500)
        ) {
          return false;
        }
      }
      return true;
    },
  });

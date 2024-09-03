import { listMyVaultPermissions } from "@/api/letuspass";
import { restQueryRetryFunc } from "@/common/queryRetry";
import { useQuery } from "@tanstack/react-query";

/**
 * Returns the query object for the vault permissions of the current user on given vault.
 */
export const useVaultPermissionsQuery = (vaultId: number) =>
  useQuery({
    queryKey: ["vaultPermissions", vaultId],
    queryFn: () => listMyVaultPermissions(Number(vaultId)),
    retry: restQueryRetryFunc,
  });

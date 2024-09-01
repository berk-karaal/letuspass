import { ControllersHandleVaultsMyKeyVaultKeyResponse } from "@/api/letuspass.schemas";
import { AESService, ECService } from "@/services/letuscrypto";

export async function decryptVaultKey(
  keyResponse: ControllersHandleVaultsMyKeyVaultKeyResponse,
  userPrivateKey: string
) {
  if (keyResponse.inviter_user_id === keyResponse.key_owner_user_id) {
    return await AESService.decrypt(
      userPrivateKey,
      keyResponse.encryption_iv,
      keyResponse.encrypted_vault_key
    );
  }

  const sharedKey = await ECService.getSharedKey(
    userPrivateKey,
    keyResponse.inviter_user_public_key
  );
  return await AESService.decrypt(
    sharedKey,
    keyResponse.encryption_iv,
    keyResponse.encrypted_vault_key
  );
}

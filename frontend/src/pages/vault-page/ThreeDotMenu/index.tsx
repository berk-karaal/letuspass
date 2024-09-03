import { Menu } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconAbc,
  IconLogout2,
  IconTrash,
  IconUsersGroup,
} from "@tabler/icons-react";
import { Link } from "react-router-dom";
import DeleteVaultModal from "./DeleteVaultModal";
import LeaveVaultModal from "./LeaveVaultModal";
import RenameVaultModal from "./RenameVaultModal";

export default function ThreeDotMenu({
  vaultId,
  vaultName,
  target,
}: {
  vaultId: number;
  vaultName: string;
  target: React.ReactNode;
}) {
  const [deleteConfirmationModalOpened, deleteConfirmationModal] =
    useDisclosure(false);
  const [leaveConfirmationModalOpened, leaveConfirmationModal] =
    useDisclosure(false);
  const [vaultRenameModalOpened, vaultRenameModal] = useDisclosure(false);

  return (
    <>
      <Menu shadow="md" withArrow>
        <Menu.Target>{target}</Menu.Target>

        <Menu.Dropdown>
          <Menu.Item
            leftSection={<IconAbc size={"1.1rem"} />}
            onClick={vaultRenameModal.open}
          >
            Rename Vault
          </Menu.Item>
          <Menu.Item
            leftSection={<IconUsersGroup size={"1.1rem"} />}
            component={Link}
            to={`/app/vault/${vaultId}/users`}
          >
            Manage Vault Users
          </Menu.Item>
          <Menu.Divider />
          <Menu.Item
            leftSection={<IconLogout2 size={"1.1rem"} />}
            color="red"
            onClick={leaveConfirmationModal.open}
          >
            Leave Vault
          </Menu.Item>
          <Menu.Item
            leftSection={<IconTrash size={"1.1rem"} />}
            color="red"
            onClick={deleteConfirmationModal.open}
          >
            Delete Vault
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <RenameVaultModal
        vaultId={vaultId}
        currentName={vaultName}
        opened={vaultRenameModalOpened}
        close={vaultRenameModal.close}
      />

      <LeaveVaultModal
        vaultId={vaultId}
        opened={leaveConfirmationModalOpened}
        onClose={leaveConfirmationModal.close}
      />

      <DeleteVaultModal
        vaultId={vaultId}
        opened={deleteConfirmationModalOpened}
        onClose={deleteConfirmationModal.close}
      />
    </>
  );
}

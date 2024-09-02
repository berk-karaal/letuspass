import { Button, Group, Menu, Modal, Text } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconAbc,
  IconLogout2,
  IconTrash,
  IconUsersGroup,
} from "@tabler/icons-react";
import { Link } from "react-router-dom";

export default function ThreeDotMenu({
  vaultId,
  target,
}: {
  vaultId: number;
  target: React.ReactNode;
}) {
  const [deleteConfirmationModalOpened, deleteConfirmationModal] =
    useDisclosure(false);

  return (
    <>
      <Menu shadow="md" withArrow>
        <Menu.Target>{target}</Menu.Target>

        <Menu.Dropdown>
          <Menu.Item
            leftSection={<IconAbc size={"1.1rem"} />}
            onClick={() => null}
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
            onClick={() => null}
          >
            Leave Vault
          </Menu.Item>
          <Menu.Item
            leftSection={<IconTrash size={"1.1rem"} />}
            onClick={deleteConfirmationModal.open}
          >
            Delete Vault
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Modal
        opened={deleteConfirmationModalOpened}
        onClose={deleteConfirmationModal.close}
        centered
        title="Confirmation"
      >
        <Text ta={"center"} size="lg" mb={"md"}>
          Are you sure you want to delete the vault?
        </Text>
        <Group justify="space-evenly" my={"lg"}>
          <Button
            color="red"
            leftSection={<IconTrash size={"1.25rem"} />}
            onClick={() => null}
          >
            Delete
          </Button>
          <Button color="gray" onClick={deleteConfirmationModal.close}>
            Cancel
          </Button>
        </Group>
      </Modal>
    </>
  );
}

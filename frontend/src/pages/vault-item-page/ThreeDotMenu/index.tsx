import { deleteVaultItem } from "@/api/letuspass";
import { Button, Group, Menu, Modal, Text } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconTrash } from "@tabler/icons-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

export default function ThreeDotMenu({
  vaultId,
  vaultItemId,
  target,
}: {
  vaultId: number;
  vaultItemId: number;
  target: React.ReactNode;
}) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [deleteConfirmationModalOpened, deleteConfirmationModal] =
    useDisclosure(false);

  const deleteVaultItemMutation = useMutation({
    mutationFn: () => deleteVaultItem(vaultId, vaultItemId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaultItems", vaultId] });
      navigate(`/app/vault/${vaultId}`);
      notifications.show({
        title: "Vault item deleted successfully",
        message: "",
        color: "green",
      });
    },
    onError: () => {
      deleteConfirmationModal.close();
      notifications.show({
        title: "Failed to delete vault item",
        message: "Please try again later",
        color: "red",
      });
    },
  });

  return (
    <>
      <Menu shadow="md" withArrow>
        <Menu.Target>{target}</Menu.Target>

        <Menu.Dropdown>
          <Menu.Item
            color="red"
            leftSection={<IconTrash size={"1.1rem"} />}
            onClick={deleteConfirmationModal.open}
          >
            Delete item
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
          Are you sure you want to delete this item?
        </Text>
        <Group justify="space-evenly" my={"lg"}>
          <Button
            color="red"
            leftSection={<IconTrash size={"1.25rem"} />}
            onClick={() => deleteVaultItemMutation.mutate()}
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

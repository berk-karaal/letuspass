import { deleteVault } from "@/api/letuspass";
import { Button, Group, Modal, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconTrash } from "@tabler/icons-react";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

export default function DeleteVaultModal({
  vaultId,
  opened,
  onClose,
}: {
  vaultId: number;
  opened: boolean;
  onClose: () => void;
}) {
  const navigate = useNavigate();

  const deleteVaultMutation = useMutation({
    mutationFn: () => deleteVault(vaultId),
    onSuccess: () => {
      notifications.show({
        title: "Vault has been deleted",
        message: "",
        color: "green",
      });
      navigate("/app");
    },
    onError: (error) => {
      console.error(error);
      notifications.show({
        title: "Failed to delete vault",
        message: "Please try again later",
        color: "red",
      });
    },
  });

  const handleSubmit = () => {
    deleteVaultMutation.mutate();
  };

  return (
    <>
      <Modal opened={opened} onClose={onClose} centered title="Confirmation">
        <Text ta={"center"} size="lg" mb={"md"}>
          Are you sure you want to delete the vault?
        </Text>
        <Group justify="space-evenly" my={"lg"}>
          <Button
            color="red"
            leftSection={<IconTrash size={"1.25rem"} />}
            onClick={handleSubmit}
          >
            Delete
          </Button>
          <Button color="gray" onClick={onClose}>
            Cancel
          </Button>
        </Group>
      </Modal>
    </>
  );
}

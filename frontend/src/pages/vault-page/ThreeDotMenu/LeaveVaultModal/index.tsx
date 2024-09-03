import { leaveVault } from "@/api/letuspass";
import { Button, Group, Modal, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

export default function LeaveVaultModal({
  vaultId,
  opened,
  onClose,
}: {
  vaultId: number;
  opened: boolean;
  onClose: () => void;
}) {
  const navigate = useNavigate();

  const leaveVaultMutation = useMutation({
    mutationFn: () => leaveVault(vaultId),
    onSuccess: () => {
      notifications.show({
        title: "You have left the vault",
        message: "",
        color: "green",
      });
      navigate("/app");
    },
  });

  const handleSubmit = () => {
    leaveVaultMutation.mutate();
  };

  return (
    <>
      <Modal opened={opened} onClose={onClose} centered title="Confirmation">
        <Text ta={"center"} size="lg" mb={"md"}>
          Are you sure you want to leave the vault?
        </Text>
        <Group justify="space-evenly" my={"lg"}>
          <Button color="red" onClick={handleSubmit}>
            Leave
          </Button>
          <Button color="gray" onClick={onClose}>
            Cancel
          </Button>
        </Group>
      </Modal>
    </>
  );
}

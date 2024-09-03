import { removeUserFromVault } from "@/api/letuspass";
import {
  ActionIcon,
  Button,
  Group,
  Modal,
  Text,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconSquareX } from "@tabler/icons-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function RemoveUserButtonAndModal({
  vaultId,
  userId,
}: {
  vaultId: number;
  userId: number;
}) {
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const queryClient = useQueryClient();

  const [deleteConfirmationModalOpened, deleteConfirmationModal] =
    useDisclosure(false);

  const removeUserMutation = useMutation({
    mutationFn: () => removeUserFromVault(vaultId, { user_id: userId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaultUsers", vaultId] });
      notifications.show({
        title: "User removed",
        message: "",
        color: "green",
      });
      deleteConfirmationModal.close();
    },
    onError: (error) => {
      console.error(error);
      queryClient.invalidateQueries({ queryKey: ["vaultUsers", vaultId] });
      notifications.show({
        title: "Failed to remove user",
        message: "Please try again later",
        color: "red",
      });
      deleteConfirmationModal.close();
    },
  });

  return (
    <>
      <ActionIcon
        variant="transparent"
        onClick={deleteConfirmationModal.open}
        color={
          colorScheme == "light" ? theme.colors.red[7] : theme.colors.red[6]
        }
      >
        <IconSquareX size={"1.5rem"} />
      </ActionIcon>

      <Modal
        opened={deleteConfirmationModalOpened}
        onClose={deleteConfirmationModal.close}
        centered
        title="Confirmation"
      >
        <Text ta={"center"} size="lg" mb={"md"}>
          Are you sure you want to remove the user from vault?
        </Text>
        <Group justify="space-evenly" my={"lg"}>
          <Button
            color="red"
            onClick={() => {
              removeUserMutation.mutate();
            }}
            loading={removeUserMutation.isPending}
          >
            Remove
          </Button>
          <Button color="gray" onClick={deleteConfirmationModal.close}>
            Cancel
          </Button>
        </Group>
      </Modal>
    </>
  );
}

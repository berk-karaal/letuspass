import { renameVault } from "@/api/letuspass";
import { ControllersHandleVaultsManageRenameRenameVaultResponse } from "@/api/letuspass.schemas";
import { Button, Group, Modal, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function RenameVaultModal({
  vaultId,
  currentName,
  opened,
  close,
}: {
  vaultId: number;
  currentName: string;
  opened: boolean;
  close: () => void;
}) {
  const queryClient = useQueryClient();

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      name: currentName,
    },

    validate: {
      name: (value) => (value.length > 0 ? null : "Name is required"),
    },
  });

  const renameVaultMutation = useMutation({
    mutationFn: (
      data: ControllersHandleVaultsManageRenameRenameVaultResponse
    ) => renameVault(vaultId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vault", vaultId] });
      close();
      notifications.show({
        title: "Vault renamed successfully",
        message: "",
        color: "green",
      });
    },
    onError: (error) => {
      console.log(error);
      notifications.show({
        title: "Failed to rename vault",
        message: "Please try again later",
        color: "red",
      });
    },
  });

  const handleSubmit = (values: typeof form.values) => {
    renameVaultMutation.mutate({ name: values.name });
  };

  return (
    <>
      <Modal opened={opened} onClose={close} title="Rename Vault">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="Name"
            placeholder="Vault Name"
            key={form.key("name")}
            {...form.getInputProps("name")}
            disabled={false}
          />
          <Group justify="flex-end" mt="md">
            <Button type="submit" disabled={false} loading={false}>
              Rename
            </Button>
          </Group>
        </form>
      </Modal>
    </>
  );
}

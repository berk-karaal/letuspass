import { updateVaultItem } from "@/api/letuspass";
import {
  ControllersHandleVaultItemsUpdateVaultItemUpdateRequest,
  SchemasBadRequestResponse,
} from "@/api/letuspass.schemas";
import { AESService } from "@/services/letuscrypto";
import {
  ActionIcon,
  Button,
  Group,
  Modal,
  PasswordInput,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconEdit } from "@tabler/icons-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useState } from "react";

export default function EditItemButtonAndModal({
  vaultId,
  vaultItemId,
  vaultKey,
  vaultItemEncryptionIV,
  currentPlainValues,
}: {
  vaultId: number;
  vaultItemId: number;
  vaultKey: string;
  vaultItemEncryptionIV: string;
  currentPlainValues: {
    title: string;
    username: string;
    password: string;
    notes: string;
  };
}) {
  const [isUpdateModalOpened, updateModal] = useDisclosure(false);
  const queryClient = useQueryClient();
  const [errorText, setErrorText] = useState<string | null>(null);

  const updateVaultItemMutation = useMutation({
    mutationFn: (
      updateData: ControllersHandleVaultItemsUpdateVaultItemUpdateRequest
    ) => updateVaultItem(vaultId, vaultItemId, updateData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaultItem", vaultItemId] });
      updateModal.close();
      notifications.show({
        title: "Vault item updated successfully",
        message: "",
        color: "green",
      });
    },
    onError: (error) => {
      if (axios.isAxiosError(error)) {
        switch (error.response?.status) {
          case 400:
            const data = error.response.data as SchemasBadRequestResponse;
            setErrorText(data.error);
            break;
          default:
            setErrorText(
              "Failed to create vault item. Please try again later."
            );
            break;
        }
      } else {
        setErrorText("Failed to create vault item. Please try again later.");
      }
    },
  });

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      title: currentPlainValues.title,
      username: currentPlainValues.username,
      password: currentPlainValues.password,
      notes: currentPlainValues.notes,
    },

    validate: {
      title: (value) => (value.length > 0 ? null : "Title is required"),
    },
  });

  const handleSubmit = async (values: typeof form.values) => {
    updateVaultItemMutation.mutate({
      title: values.title,
      encrypted_username: await AESService.encrypt(
        vaultKey,
        vaultItemEncryptionIV,
        values.username
      ),
      encrypted_password: await AESService.encrypt(
        vaultKey,
        vaultItemEncryptionIV,
        values.password
      ),
      encrypted_note: await AESService.encrypt(
        vaultKey,
        vaultItemEncryptionIV,
        values.notes
      ),
    });
  };

  return (
    <>
      <Modal
        opened={isUpdateModalOpened}
        onClose={updateModal.close}
        title="Edit Item"
      >
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="Title"
            placeholder="Item Title"
            key={form.key("title")}
            {...form.getInputProps("title")}
            disabled={updateVaultItemMutation.isPending}
          />
          <TextInput
            label="Username"
            placeholder=""
            key={form.key("username")}
            {...form.getInputProps("username")}
            disabled={updateVaultItemMutation.isPending}
          />
          <PasswordInput
            label="Password"
            placeholder=""
            key={form.key("password")}
            {...form.getInputProps("password")}
            disabled={updateVaultItemMutation.isPending}
          />
          <Textarea
            label="Notes"
            placeholder=""
            key={form.key("notes")}
            {...form.getInputProps("notes")}
            disabled={updateVaultItemMutation.isPending}
            autosize
            minRows={3}
            maxRows={5}
          />

          <Text c={"red"} mt={"xs"} display={errorText ? "block" : "none"}>
            {errorText}
          </Text>

          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={updateVaultItemMutation.isPending}>
              Update
            </Button>
          </Group>
        </form>
      </Modal>
      <ActionIcon
        variant="transparent"
        color="dark"
        onClick={() => updateModal.open()}
        mx={"0.35rem"}
      >
        <IconEdit size={"1.5rem"} />
      </ActionIcon>
    </>
  );
}

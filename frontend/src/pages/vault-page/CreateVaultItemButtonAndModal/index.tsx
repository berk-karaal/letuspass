import { createVaultItem } from "@/api/letuspass";
import {
  ControllersHandleVaultItemsCreateVaultItemCreateRequest,
  SchemasBadRequestResponse,
} from "@/api/letuspass.schemas";
import {
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
import { IconPlus } from "@tabler/icons-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useState } from "react";

export default function CreateVaultItemButtonAndModal({
  vaultId,
}: {
  vaultId: number;
}) {
  const queryClient = useQueryClient();
  const [opened, { open, close }] = useDisclosure(false);

  const [errorText, setErrorText] = useState<string | null>(null);

  const createVaultMutation = useMutation({
    mutationFn: (
      newVaultItem: ControllersHandleVaultItemsCreateVaultItemCreateRequest
    ) => createVaultItem(vaultId, newVaultItem),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vault", vaultId, "items"] });
      close();
      form.reset();
      notifications.show({
        title: "Vault item created successfully",
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
      title: "",
      username: "",
      password: "",
      note: "",
    },

    validate: {
      title: (value) => (value.length > 0 ? null : "Title is required"),
    },
  });

  const handleSubmit = (values: typeof form.values) => {
    setErrorText(null);

    // TODO: make encryption
    createVaultMutation.mutate({
      title: values.title,
      encrypted_username: values.username,
      encrypted_password: values.password,
      encrypted_note: values.note,
    });
  };

  return (
    <>
      <Modal opened={opened} onClose={close} title="New Item">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="Title"
            placeholder="Item Title"
            key={form.key("title")}
            {...form.getInputProps("title")}
            disabled={createVaultMutation.isPending}
          />
          <TextInput
            label="Username"
            placeholder=""
            key={form.key("username")}
            {...form.getInputProps("username")}
            disabled={createVaultMutation.isPending}
          />
          <PasswordInput
            label="Password"
            placeholder=""
            key={form.key("password")}
            {...form.getInputProps("password")}
            disabled={createVaultMutation.isPending}
          />
          <Textarea
            label="Note"
            placeholder=""
            key={form.key("note")}
            {...form.getInputProps("note")}
            disabled={createVaultMutation.isPending}
            autosize
            minRows={3}
            maxRows={5}
          />

          <Text c={"red"} mt={"xs"} display={errorText ? "block" : "none"}>
            {errorText}
          </Text>

          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={createVaultMutation.isPending}>
              Create Vault
            </Button>
          </Group>
        </form>
      </Modal>

      <Button onClick={open} leftSection={<IconPlus size={"1.25rem"} />}>
        New Item
      </Button>
    </>
  );
}

import { createVault } from "@/api/letuspass";
import { SchemasBadRequestResponse } from "@/api/letuspass.schemas";
import { AESService } from "@/services/letuscrypto";
import { useAppSelector } from "@/store/hooks";
import { Button, Group, Modal, Text, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconPlus } from "@tabler/icons-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useState } from "react";

export default function CreateVaultButtonAndModal() {
  const user = useAppSelector((state) => state.user);

  const queryClient = useQueryClient();
  const [opened, { open, close }] = useDisclosure(false);

  const [errorText, setErrorText] = useState<string | null>(null);

  const createVaultMutation = useMutation({
    mutationFn: createVault,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaults"] });
      close();
      form.reset();
      notifications.show({
        title: "Vault created successfully",
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
            setErrorText("Failed to create vault. Please try again later.");
            break;
        }
      } else {
        setErrorText("Failed to create vault. Please try again later.");
      }
    },
  });

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      name: "",
    },

    validate: {
      name: (value) => (value.length > 0 ? null : "Name is required"),
    },
  });

  const handleSubmit = async (values: typeof form.values) => {
    setErrorText(null);
    const vaultKey = await AESService.generateRandomKey();
    const vaultKeyEncryptionIV = AESService.generateRandomIV();
    const encryptedVaultKey = await AESService.encrypt(
      user.privateKey,
      vaultKeyEncryptionIV,
      vaultKey
    );
    createVaultMutation.mutate({
      name: values.name,
      encryption_iv: vaultKeyEncryptionIV,
      encrypted_vault_key: encryptedVaultKey,
    });
  };

  return (
    <>
      <Modal opened={opened} onClose={close} title="New Vault">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="Vault Name"
            placeholder="Enter vault name"
            key={form.key("name")}
            {...form.getInputProps("name")}
            disabled={createVaultMutation.isPending}
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
        New Vault
      </Button>
    </>
  );
}

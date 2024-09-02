import {
  addUserToVault,
  getUserByEmail,
  retrieveMyVaultKey,
} from "@/api/letuspass";
import { SchemasBadRequestResponse } from "@/api/letuspass.schemas";
import { decryptVaultKey } from "@/common/vaultkey";
import { AESService, ECService } from "@/services/letuscrypto";
import { useAppSelector } from "@/store/hooks";
import {
  ActionIcon,
  Button,
  Checkbox,
  Group,
  Modal,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconUserPlus } from "@tabler/icons-react";
import axios from "axios";
import { useState } from "react";
import classes from "./styles.module.css";

const permissions = [
  {
    value: "manage_vault",
    description: "Can add and remove users and list users.",
  },
  {
    value: "delete_vault",
    description: "Can delete the vault.",
  },
  {
    value: "manage_items",
    description: "Can create, update and delete vault items.",
  },
];

export default function AddUserButtonAndModal({
  vaultId,
}: {
  vaultId: number;
}) {
  const [opened, { open, close }] = useDisclosure(false);
  const user = useAppSelector((state) => state.user);

  const [errorText, setErrorText] = useState<string | null>(null);

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      email: "",
      permissions: [],
    },

    validate: {
      email: (value) => (value.length > 0 ? null : "Email is required"),
    },
  });

  const handleSubmit = async (values: typeof form.values) => {
    setErrorText(null);

    if (values.email === user.email) {
      form.setFieldError("email", "You cannot add yourself to the vault.");
      return;
    }

    let newUserPublicKey: string;
    try {
      const response = await getUserByEmail({ email: values.email });
      newUserPublicKey = response.public_key;
    } catch (e) {
      if (axios.isAxiosError(e) && e.response?.status == 404) {
        form.setFieldError("email", "User with given email not found.");
        return;
      }
      setErrorText("An error occurred, please try again later.");
      return;
    }

    let rawVaultKey: string;
    try {
      const response = await retrieveMyVaultKey(vaultId);
      rawVaultKey = await decryptVaultKey(response, user.privateKey);
    } catch (e) {
      setErrorText("An error occurred, please try again later.");
      return;
    }

    const sharedKey = ECService.getSharedKey(user.privateKey, newUserPublicKey);
    console.log("Shared key", sharedKey);
    const encryptionIV = AESService.generateRandomIV();
    const encryptedVaultKey = await AESService.encrypt(
      sharedKey,
      encryptionIV,
      rawVaultKey
    );

    try {
      await addUserToVault(vaultId, {
        email: values.email,
        permissions: values.permissions,
        encrypted_vault_key: encryptedVaultKey,
        vault_key_encryption_iv: encryptionIV,
      });
    } catch (e) {
      if (axios.isAxiosError(e) && e.response?.status == 400) {
        let r = e.response.data as SchemasBadRequestResponse;
        setErrorText(r.error);
        return;
      }
      setErrorText("An error occurred, please try again later.");
      return;
    }

    notifications.show({
      title: "User added",
      message: `User '${values.email}' has been added to the vault successfully.`,
      color: "green",
    });
    form.reset();
    close();
  };

  const cards = permissions.map((item) => (
    <Checkbox.Card
      p={"xs"}
      className={classes.permissionCard}
      radius="md"
      value={item.value}
      key={item.value}
    >
      <Group wrap="nowrap" align="flex-start">
        <Checkbox.Indicator my={"auto"} />
        <div>
          <Text inline>
            <b>{item.value}</b>
          </Text>
          <Text size="sm">{item.description}</Text>
        </div>
      </Group>
    </Checkbox.Card>
  ));

  return (
    <>
      <Modal opened={opened} onClose={close} title="Add User to Vault">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="User Email"
            placeholder="user@email.com"
            key={form.key("email")}
            {...form.getInputProps("email")}
            disabled={false}
          />

          <Checkbox.Group
            label="Additional Permissions"
            description="Choose extra permissions for the user"
            key={form.key("permissions")}
            {...form.getInputProps("permissions")}
            mt={"xs"}
          >
            <Stack mt={"0.2rem"} gap="0.5rem">
              {cards}
            </Stack>
          </Checkbox.Group>

          <Text c={"red"} mt={"xs"} display={errorText ? "block" : "none"}>
            {errorText}
          </Text>

          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={false}>
              Add User
            </Button>
          </Group>
        </form>
      </Modal>

      <ActionIcon variant="transparent" color="dark" onClick={open}>
        <IconUserPlus size={"1.5rem"} />
      </ActionIcon>
    </>
  );
}

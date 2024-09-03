import {
  retrieveMyVaultKey,
  retrieveVault,
  retrieveVaultItem,
} from "@/api/letuspass";
import { useVaultPermissionsQuery } from "@/hooks/useVaultPermissionsQuery";
import {
  ActionIcon,
  Box,
  Group,
  Loader,
  Overlay,
  PasswordInput,
  Text,
  Textarea,
  TextInput,
  Title,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";

import { restQueryRetryFunc } from "@/common/queryRetry";
import { decryptVaultKey } from "@/common/vaultkey";
import { AESService } from "@/services/letuscrypto";
import { useAppSelector } from "@/store/hooks";
import {
  IconArrowLeft,
  IconBriefcase2,
  IconDotsVertical,
  IconEye,
  IconEyeOff,
  IconFileTime,
  IconKey,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import strftime from "strftime";
import EditItemButtonAndModal from "./EditItemButtonAndModal";
import ThreeDotMenu from "./ThreeDotMenu";

function VaultItemPage() {
  const { vaultId, vaultItemId } = useParams();
  const navigate = useNavigate();
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const user = useAppSelector((state) => state.user);

  const [vaultItemFieldsDecrypted, setVaultItemFieldsDecrypted] = useState<{
    username: string | null;
    password: string | null;
    note: string | null;
  }>({
    username: null,
    password: null,
    note: null,
  });

  const vaultKey = useRef<string | null>(null);

  const [isOverlayActive, setIsOverlayActive] = useState(true);

  const vaultQuery = useQuery({
    queryKey: ["vault", vaultId],
    queryFn: () => retrieveVault(Number(vaultId)),
    retry: restQueryRetryFunc,
  });

  const vaultPermissionsQuery = useVaultPermissionsQuery(Number(vaultId));

  const vaultItemQuery = useQuery({
    queryKey: ["vaultItem", Number(vaultItemId)],
    queryFn: () => retrieveVaultItem(Number(vaultId), Number(vaultItemId)),
    retry: restQueryRetryFunc,
  });

  const vaultKeyQuery = useQuery({
    queryKey: ["vaultKey", vaultId],
    queryFn: () => retrieveMyVaultKey(Number(vaultId)),
    retry: restQueryRetryFunc,
  });

  useEffect(() => {
    if (vaultItemQuery.isError) {
      let errorText = "Failed to retrieve vault. Please try again later.";
      if (axios.isAxiosError(vaultItemQuery.error)) {
        if (vaultItemQuery.error.response?.status === 404) {
          errorText = "Vault item not found.";
        }
      }
      notifications.show({
        title: errorText,
        message: "",
        color: "red",
      });
      return;
    }
  }, [vaultItemQuery.isError]);

  useEffect(() => {
    if (
      !vaultItemQuery.isFetching &&
      !vaultKeyQuery.isFetching &&
      vaultItemQuery.isSuccess &&
      vaultKeyQuery.isSuccess
    ) {
      const decryptFields = async () => {
        setVaultItemFieldsDecrypted({
          username: await decryptVaultItemField(
            vaultItemQuery.data.encrypted_username
          ),
          password: await decryptVaultItemField(
            vaultItemQuery.data.encrypted_password
          ),
          note: await decryptVaultItemField(vaultItemQuery.data.encrypted_note),
        });
      };
      decryptFields();
    }
  }, [vaultItemQuery.isFetching, vaultKeyQuery.isFetching]);

  const decryptVaultItemField = async (
    encryptedData: string
  ): Promise<string> => {
    if (!vaultKeyQuery.isSuccess) {
      console.error("Vault key query is not successful");
      return "";
    }
    if (!vaultItemQuery.isSuccess) {
      console.error("Vault item query is not successful");
      return "";
    }

    vaultKey.current = await decryptVaultKey(
      vaultKeyQuery.data,
      user.privateKey
    );

    if (encryptedData === "") {
      return "";
    }

    return await AESService.decrypt(
      vaultKey.current,
      vaultItemQuery.data.encryption_iv,
      encryptedData
    );
  };

  return (
    <>
      {/* Back button and vault name row */}
      <Group gap={"sm"} justify="flex-start">
        <ActionIcon
          variant="transparent"
          color="dark"
          onClick={() => navigate(-1)}
        >
          <IconArrowLeft size={"1.5rem"} />
        </ActionIcon>
        <IconBriefcase2 size={"1.75rem"} />
        <Title order={2} fw={"lighter"}>
          {vaultQuery.isSuccess ? (
            vaultQuery.data.name
          ) : (
            <Loader color="gray" />
          )}
        </Title>
      </Group>

      {/* Vault item title row */}
      <Group mt={"lg"} gap={"sm"} justify="flex-start">
        <IconKey size={"1.5rem"} />
        <Text fz={"h3"} fw={"lighter"}>
          {vaultItemQuery.isSuccess ? (
            vaultItemQuery.data.title
          ) : (
            <Loader color="gray" />
          )}
        </Text>
        <Box style={{ marginLeft: "auto" }}>
          <ActionIcon
            variant="transparent"
            color="dark"
            onClick={() => null}
            mx={"0.35rem"}
          >
            <IconFileTime size={"1.5rem"} />
          </ActionIcon>
          {vaultPermissionsQuery.isSuccess &&
            vaultPermissionsQuery.data.includes("manage_items") && (
              <>
                {vaultItemQuery.isSuccess &&
                  vaultKeyQuery.isSuccess &&
                  vaultItemFieldsDecrypted.username !== null && (
                    <EditItemButtonAndModal
                      vaultId={Number(vaultId)}
                      vaultItemId={Number(vaultItemId)}
                      vaultKey={vaultKey.current ?? ""}
                      vaultItemEncryptionIV={vaultItemQuery.data.encryption_iv}
                      currentPlainValues={{
                        title: vaultItemQuery.data.title,
                        username: vaultItemFieldsDecrypted.username ?? "",
                        password: vaultItemFieldsDecrypted.password ?? "",
                        notes: vaultItemFieldsDecrypted.note ?? "",
                      }}
                    />
                  )}
                <ThreeDotMenu
                  vaultId={Number(vaultId)}
                  vaultItemId={Number(vaultItemId)}
                  target={
                    <ActionIcon
                      variant="transparent"
                      color="dark"
                      onClick={() => null}
                      mx={"0.35rem"}
                    >
                      <IconDotsVertical size={"1.5rem"} />
                    </ActionIcon>
                  }
                />
              </>
            )}
        </Box>
      </Group>

      {/* Username */}
      <TextInput
        readOnly
        label="Username"
        value={vaultItemFieldsDecrypted.username ?? ""}
        placeholder={vaultItemFieldsDecrypted.username === "" ? "No data" : ""}
        mt={"sm"}
      />

      {/* Password */}
      <PasswordInput
        readOnly
        label="Password"
        value={vaultItemFieldsDecrypted.password ?? ""}
        placeholder={vaultItemFieldsDecrypted.password === "" ? "No data" : ""}
        mt={"xs"}
      />

      {/* Note */}
      <>
        <Group gap={"0.2rem"} mt={"sm"} justify="space-between">
          <Text size="sm">Notes</Text>
          <ActionIcon
            variant="transparent"
            color="dark"
            onClick={() => setIsOverlayActive(!isOverlayActive)}
          >
            {isOverlayActive ? (
              <IconEye
                color={"var(--mantine-color-gray-light-color)"}
                stroke={1.25}
                size={"1rem"}
              />
            ) : (
              <IconEyeOff
                color={"var(--mantine-color-gray-light-color)"}
                stroke={1.25}
                size={"1rem"}
              />
            )}
          </ActionIcon>
        </Group>
        <Box pos={"relative"}>
          <Textarea
            readOnly
            value={vaultItemFieldsDecrypted.note ?? ""}
            placeholder={vaultItemFieldsDecrypted.note === "" ? "No data" : ""}
            autosize
            minRows={5}
            maxRows={8}
            style={{ border: "none" }}
          />

          {vaultItemFieldsDecrypted.note !== "" && isOverlayActive && (
            <Overlay
              m={"1px"}
              color={
                colorScheme == "light"
                  ? theme.colors.gray[0]
                  : theme.colors.dark[6]
              }
              blur={2}
              backgroundOpacity={0.85}
            />
          )}
        </Box>
      </>

      {/* Last Update text */}
      {vaultItemQuery.isSuccess && (
        <Text size={"xs"} ta={"end"} mt={"md"}>
          Last Update:{" "}
          {strftime(
            "%d-%m-%Y %H:%M %z",
            new Date(vaultItemQuery.data.updated_at)
          )}
        </Text>
      )}
    </>
  );
}

export default VaultItemPage;

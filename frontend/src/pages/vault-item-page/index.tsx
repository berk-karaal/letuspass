import { retrieveVault, retrieveVaultItem } from "@/api/letuspass";
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

import {
  IconArrowLeft,
  IconBriefcase2,
  IconDotsVertical,
  IconEdit,
  IconEye,
  IconEyeOff,
  IconFileTime,
  IconKey,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import strftime from "strftime";

function VaultItemPage() {
  const { vaultId, vaultItemId } = useParams();
  const navigate = useNavigate();
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const [isOverlayActive, setIsOverlayActive] = useState(true);

  const vaultQuery = useQuery({
    queryKey: ["vault", vaultId],
    queryFn: () => retrieveVault(Number(vaultId)),
    retry: (failureCount: number, error: Error) => {
      if (failureCount > 2) {
        return false;
      }
      if (axios.isAxiosError(error)) {
        if (
          (error.response?.status ?? 500 >= 400) &&
          (error.response?.status ?? 500 < 500)
        ) {
          return false;
        }
      }
      return true;
    },
  });

  const vaultPermissionsQuery = useVaultPermissionsQuery(Number(vaultId));

  const vaultItemQuery = useQuery({
    queryKey: ["vaultItem", vaultItemId],
    queryFn: () => retrieveVaultItem(Number(vaultId), Number(vaultItemId)),
    retry: (failureCount: number, error: Error) => {
      if (failureCount > 2) {
        return false;
      }
      if (axios.isAxiosError(error)) {
        if (
          (error.response?.status ?? 500 >= 400) &&
          (error.response?.status ?? 500 < 500)
        ) {
          return false;
        }
      }
      return true;
    },
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
    }
  }, [vaultItemQuery]);

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
                <ActionIcon
                  variant="transparent"
                  color="dark"
                  onClick={() => null}
                  mx={"0.35rem"}
                >
                  <IconEdit size={"1.5rem"} />
                </ActionIcon>
                <ActionIcon
                  variant="transparent"
                  color="dark"
                  onClick={() => null}
                  mx={"0.35rem"}
                >
                  <IconDotsVertical size={"1.5rem"} />
                </ActionIcon>
              </>
            )}
        </Box>
      </Group>

      {/* Username */}
      <TextInput
        readOnly
        label="Username"
        value={
          vaultItemQuery.isSuccess ? vaultItemQuery.data.encrypted_username : ""
        }
        placeholder={
          vaultItemQuery.isSuccess &&
          vaultItemQuery.data.encrypted_username == ""
            ? "No data"
            : ""
        }
        mt={"sm"}
      />

      {/* Password */}
      <PasswordInput
        readOnly
        label="Password"
        value={
          vaultItemQuery.isSuccess ? vaultItemQuery.data.encrypted_password : ""
        }
        placeholder={
          vaultItemQuery.isSuccess &&
          vaultItemQuery.data.encrypted_password == ""
            ? "No data"
            : ""
        }
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
            value={
              vaultItemQuery.isSuccess ? vaultItemQuery.data.encrypted_note : ""
            }
            placeholder={
              vaultItemQuery.isSuccess &&
              vaultItemQuery.data.encrypted_note == ""
                ? "No data"
                : ""
            }
            autosize
            minRows={5}
            maxRows={8}
            style={{ border: "none" }}
          />

          {/* The first 2 conditions disable the overlay if there is no saved note. */}
          {vaultItemQuery.isSuccess &&
            vaultItemQuery.data.encrypted_note != "" &&
            isOverlayActive && (
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

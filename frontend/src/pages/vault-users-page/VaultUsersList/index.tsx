import { listVaultUsers } from "@/api/letuspass";
import { useAppSelector } from "@/store/hooks";
import {
  ActionIcon,
  Box,
  Group,
  Loader,
  Text,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { IconSquareX, IconUserFilled } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import classes from "./styles.module.css";

function VaultUserBox({
  vaultId,
  userEmail,
  userPermissions,
}: {
  vaultId: number;
  userEmail: string;
  userPermissions: string[];
}) {
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const user = useAppSelector((state) => state.user);

  return (
    <Box py={"sm"} px={"xs"} my={"sm"} className={classes.VaultUserBox}>
      <Group gap={"xs"}>
        <IconUserFilled size={"1.75rem"} />
        <Text size="lg">{userEmail}</Text>
        {user.email !== userEmail && (
          <Box ml={"auto"}>
            <ActionIcon
              variant="transparent"
              onClick={() => {}}
              color={
                colorScheme == "light"
                  ? theme.colors.red[7]
                  : theme.colors.red[6]
              }
            >
              <IconSquareX size={"1.5rem"} />
            </ActionIcon>
          </Box>
        )}
      </Group>
      <Text mt={"sm"}>
        <b>Permissions:</b> {userPermissions.join(", ")}
      </Text>
    </Box>
  );
}

export default function VaultUsersList({ vaultId }: { vaultId: number }) {
  const vaultUsersQuery = useQuery({
    queryKey: ["vaultUsers", vaultId],
    queryFn: () => listVaultUsers(Number(vaultId)),
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

  return (
    <>
      {!vaultUsersQuery.isSuccess ? (
        <Loader color="gray" />
      ) : (
        vaultUsersQuery.data.map((user) => (
          <VaultUserBox
            key={user.email}
            vaultId={Number(vaultId)}
            userEmail={user.email}
            userPermissions={user.permissions}
          />
        ))
      )}{" "}
    </>
  );
}
